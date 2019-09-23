package eventbus

import (
	"context"

	"github.com/sirupsen/logrus"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/cloud-pkg/mqmux"
)

type Config struct {
	Logger *logrus.Entry
	URI    string
	Name   string
}

type EventBus struct {
	logger *logrus.Entry
	mq     *mqmux.WorkerMux
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
		logger: c.Logger,
		mq:     mq,
	}, nil
}

func (e *EventBus) Start() error {
	err := e.mq.Publisher("streams.events")
	if err != nil {
		return err
	}
	return e.mq.Run()
}

func (e *EventBus) Stop() error {
	return e.mq.Close()
}

func (e *EventBus) emitCUDStream(t privatev1.EventType, id string) error {
	event := &privatev1.Event{
		Type:     t,
		StreamID: id,
	}
	err := e.mq.Publish("streams.events", event)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventBus) EmitCreateStream(ctx context.Context, id string) error {
	e.logger.Debugf("emitting create stream: %s", id)
	return e.emitCUDStream(privatev1.EventTypeCreate, id)
}

func (e *EventBus) EmitUpdateStream(ctx context.Context, id string) error {
	e.logger.Debugf("emitting update stream: %s", id)
	return e.emitCUDStream(privatev1.EventTypeUpdate, id)
}

func (e *EventBus) EmitDeleteStream(ctx context.Context, id string) error {
	e.logger.Debugf("emitting delete stream: %s", id)
	return e.emitCUDStream(privatev1.EventTypeDelete, id)
}
