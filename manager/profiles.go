package manager

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/videocoin/cloud-pkg/tracer"
	ds "github.com/videocoin/cloud-streams/datastore"
)

func (m *Manager) GetProfileByID(ctx context.Context, id string) (*ds.Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.GetProfileByID")
	defer span.Finish()

	profile, err := m.ds.Profile.Get(ctx, id)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return profile, nil
}

func (m *Manager) GetProfileByIDOrName(ctx context.Context, idOrName string) (*ds.Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.GetProfileByIDOrName")
	defer span.Finish()

	profile, err := m.ds.Profile.GetByIdOrName(ctx, idOrName)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return profile, nil
}

func (m *Manager) GetProfileByName(ctx context.Context, name string) (*ds.Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.GetProfileByName")
	defer span.Finish()

	profile, err := m.ds.Profile.GetByName(ctx, name)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return profile, nil
}

func (m *Manager) ListEnabledProfiles(ctx context.Context) ([]*ds.Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.ListEnabledProfiles")
	defer span.Finish()

	profiles, err := m.ds.Profile.ListEnabled(ctx)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return profiles, nil
}

func (m *Manager) ListAllProfiles(ctx context.Context) ([]*ds.Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.ListAllProfiles")
	defer span.Finish()

	profiles, err := m.ds.Profile.List(ctx)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return profiles, nil
}
