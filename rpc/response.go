package rpc

import (
	"github.com/jinzhu/copier"
	v1 "github.com/videocoin/cloud-api/streams/v1"
)

func toStreamResponse(stream *v1.Stream) (*v1.StreamResponse, error) {
	response := new(v1.StreamResponse)
	if err := copier.Copy(response, stream); err != nil {
		return nil, err
	}

	return response, nil
}

func toStreamListResponse(streams []*v1.Stream) (*v1.StreamListResponse, error) {
	response := &v1.StreamListResponse{}
	if err := copier.Copy(&response.Items, &streams); err != nil {
		return nil, err
	}

	return response, nil
}
