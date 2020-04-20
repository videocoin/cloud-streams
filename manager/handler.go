package manager

import (
	"context"

	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
)

func (m *Manager) handleStreamStatus(ctx context.Context, event *privatev1.Event) error {
	switch event.Type {
	case privatev1.EventTypeUpdateStatus:
		{
			logger := m.logger.WithFields(logrus.Fields{
				"status":    event.Status.String(),
				"stream_id": event.StreamID,
			})
			logger.Info("updating status")

			stream, err := m.GetStreamByID(ctx, event.StreamID)
			if err != nil {
				logger.WithError(err).Error("failed to get stream")
				return nil
			}

			if event.Status == streamsv1.StreamStatusFailed ||
				event.Status == streamsv1.StreamStatusReady ||
				event.Status == streamsv1.StreamStatusCancelled ||
				event.Status == streamsv1.StreamStatusDeleted {
				if stream.Status == streamsv1.StreamStatusFailed ||
					stream.Status == streamsv1.StreamStatusReady ||
					stream.Status == streamsv1.StreamStatusCancelled ||
					stream.Status == streamsv1.StreamStatusDeleted {
					return nil
				}
			}

			updates := map[string]interface{}{"status": event.Status}
			err = m.UpdateStream(ctx, stream, updates)
			if err != nil {
				logger.WithError(err).Error("failed to update stream status")
				return nil
			}

			if event.Status == v1.StreamStatusFailed {
				_, err = m.emitter.EndStream(ctx, &emitterv1.EndStreamRequest{
					UserId:                stream.UserID,
					StreamContractId:      stream.StreamContractID,
					StreamContractAddress: stream.StreamContractAddress,
				})

				if err != nil {
					logger.WithError(err).Error("failed to end stream")
				}
			}

			if event.Status == v1.StreamStatusReady {
				user, err := m.users.GetById(ctx, &usersv1.UserRequest{
					Id: stream.UserID,
				})
				if err != nil {
					logger.WithError(err).Error("failed to get user by id")
					return nil
				}
				err = m.eb.EmitStreamPublished(ctx, user.Email, stream.OutputURL)
				if err != nil {
					logger.WithError(err).Error("failed to send email notification")
					return nil
				}
			}
		}
	}

	return nil
}
