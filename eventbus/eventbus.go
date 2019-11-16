package eventbus

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	notificationsv1 "github.com/videocoin/cloud-api/notifications/v1"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/mqmux"
	tracerext "github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/cloud-streams/manager"
)

type Config struct {
	Logger  *logrus.Entry
	URI     string
	Name    string
	DM      *manager.Manager
	Users   usersv1.UserServiceClient
	Emitter emitterv1.EmitterServiceClient
}

type EventBus struct {
	logger  *logrus.Entry
	mq      *mqmux.WorkerMux
	dm      *manager.Manager
	users   usersv1.UserServiceClient
	emitter emitterv1.EmitterServiceClient
}

func New(c *Config) (*EventBus, error) {
	mq, err := mqmux.NewWorkerMux(c.URI, c.Name)
	if err != nil {
		return nil, err
	}
	if c.Logger != nil {
		mq.Logger = c.Logger
	}

	return &EventBus{
		logger:  c.Logger,
		mq:      mq,
		dm:      c.DM,
		users:   c.Users,
		emitter: c.Emitter,
	}, nil
}

func (e *EventBus) Start() error {
	err := e.registerPublishers()
	if err != nil {
		return err
	}

	err = e.registerConsumers()
	if err != nil {
		return err
	}

	return e.mq.Run()
}

func (e *EventBus) registerPublishers() error {
	if err := e.mq.Publisher("streams.events"); err != nil {
		return err
	}

	if err := e.mq.Publisher("notifications/send"); err != nil {
		return err
	}
	return nil
}

func (e *EventBus) registerConsumers() error {
	if err := e.mq.Consumer("streams.status", 1, false, e.handleStreamStatus); err != nil {
		return err
	}

	return nil
}

func (e *EventBus) Stop() error {
	return e.mq.Close()
}

func (e *EventBus) emitCUDStream(ctx context.Context, t privatev1.EventType, id string) error {
	headers := make(amqp.Table)

	span := opentracing.SpanFromContext(ctx)
	ext.SpanKindRPCServer.Set(span)
	ext.Component.Set(span, "streams")

	span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		mqmux.RMQHeaderCarrier(headers),
	)

	event := &privatev1.Event{
		Type:     t,
		StreamID: id,
	}
	err := e.mq.PublishX("streams.events", event, headers)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventBus) EmitCreateStream(ctx context.Context, id string) error {
	e.logger.Debugf("emitting create stream: %s", id)
	return e.emitCUDStream(ctx, privatev1.EventTypeCreate, id)
}

func (e *EventBus) EmitUpdateStream(ctx context.Context, id string) error {
	e.logger.Debugf("emitting update stream: %s", id)
	return e.emitCUDStream(ctx, privatev1.EventTypeUpdate, id)
}

func (e *EventBus) EmitDeleteStream(ctx context.Context, id string) error {
	e.logger.Debugf("emitting delete stream: %s", id)
	return e.emitCUDStream(ctx, privatev1.EventTypeDelete, id)
}

func (e *EventBus) SendNotification(span opentracing.Span, req *notificationsv1.Notification) error {
	headers := make(amqp.Table)
	ext.SpanKindRPCServer.Set(span)
	ext.Component.Set(span, "streams")

	span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		mqmux.RMQHeaderCarrier(headers),
	)

	return e.mq.PublishX("notifications/send", req, headers)
}

func (e *EventBus) handleStreamStatus(d amqp.Delivery) error {
	var span opentracing.Span
	tracer := opentracing.GlobalTracer()
	spanCtx, err := tracer.Extract(opentracing.TextMap, mqmux.RMQHeaderCarrier(d.Headers))

	e.logger.Debugf("handling body: %+v", string(d.Body))

	if err != nil {
		span = tracer.StartSpan("eventbus.handleStreamStatus")
	} else {
		span = tracer.StartSpan("eventbus.handleStreamStatus", ext.RPCServerOption(spanCtx))
	}

	defer span.Finish()

	req := new(privatev1.Event)
	err = json.Unmarshal(d.Body, req)
	if err != nil {
		tracerext.SpanLogError(span, err)
		return err
	}

	span.SetTag("stream_id", req.StreamID)
	span.SetTag("event_type", req.Type.String())
	span.SetTag("status", req.Status.String())

	e.logger.Infof("handling request %+v", req)

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	switch req.Type {
	case privatev1.EventTypeUpdateStatus:
		{
			logger := e.logger.WithFields(logrus.Fields{
				"status":    req.Status.String(),
				"stream_id": req.StreamID,
			})
			logger.Info("updating status")

			stream, err := e.dm.GetStreamByID(ctx, req.StreamID)
			if err != nil {
				logger.WithError(err).Error("failed to get stream")
				return nil
			}

			updates := map[string]interface{}{"status": req.Status}
			err = e.dm.UpdateStream(ctx, stream, updates)
			if err != nil {
				logger.WithError(err).Error("failed to update stream status")
				return nil
			}

			if req.Status == v1.StreamStatusFailed {
				_, err = e.emitter.EndStream(ctx, &emitterv1.EndStreamRequest{
					UserId:                stream.UserId,
					StreamContractId:      stream.StreamContractId,
					StreamContractAddress: stream.StreamContractAddress,
				})

				if err != nil {
					logger.WithError(err).Error("failed to end stream")
				}
			}

			if req.Status == v1.StreamStatusReady {
				user, err := e.users.GetById(ctx, &usersv1.UserRequest{
					Id: stream.UserId,
				})
				if err != nil {
					logger.WithError(err).Error("failed to get user by id")
					return nil
				}
				err = e.sendStreamPublished(ctx, user.Email, stream.OutputUrl)
				if err != nil {
					logger.WithError(err).Error("failed to send email notification")
					return nil
				}
			}
		}
	}

	return nil
}

func (e *EventBus) sendStreamPublished(ctx context.Context, by, url string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendStreamPublished")
	defer span.Finish()

	md := metautils.ExtractIncoming(ctx)

	params := map[string]string{
		"by":       by,
		"url":      url,
		"internal": "",
		"domain":   md.Get("x-forwarded-host"),
	}

	notification := &notificationsv1.Notification{
		Target:   notificationsv1.NotificationTarget_EMAIL,
		Template: "stream_published",
		Params:   params,
	}

	if err := e.SendNotification(span, notification); err != nil {
		return err
	}

	return nil
}
