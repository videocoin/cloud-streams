package rpc

import (
	"github.com/jinzhu/copier"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	ds "github.com/videocoin/cloud-streams/datastore"
)

func toStreamProfile(stream *ds.Stream) (*v1.StreamProfile, error) {
	profile := new(v1.StreamProfile)
	if err := copier.Copy(profile, stream); err != nil {
		return nil, err
	}

	return profile, nil
}

func toStreamProfiles(streams []*ds.Stream) (*v1.StreamProfiles, error) {
	profiles := &v1.StreamProfiles{}
	if err := copier.Copy(&profiles.Items, &streams); err != nil {
		return nil, err
	}

	return profiles, nil
}
