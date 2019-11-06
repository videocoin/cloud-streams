package manager

import (
	"context"

	"github.com/opentracing/opentracing-go"
	v1 "github.com/videocoin/cloud-api/users/v1"
	tracer "github.com/videocoin/cloud-pkg/tracer"
)

func (m *Manager) GetUserByID(ctx context.Context, id string) (*v1.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.GetUserByID")
	defer span.Finish()

	user, err := m.ds.User.Get(ctx, id)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return user, nil
}
