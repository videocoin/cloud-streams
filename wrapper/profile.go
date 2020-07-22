package wrapper

import (
	"strings"

	v1 "github.com/videocoin/cloud-api/streams/v1"
	ds "github.com/videocoin/cloud-streams/datastore"
)

type Profile struct {
	*ds.Profile
}

func (p *Profile) Render(input, output string) string {
	built := []string{"ffmpeg"}

	for _, c := range p.Spec.Components {
		if c.Type == v1.ComponentTypeDemuxer {
			built = append(built, c.Render())
		}
	}

	built = append(built, "-i "+input)

	for _, c := range p.Spec.Components {
		if c.Type != v1.ComponentTypeDemuxer {
			built = append(built, c.Render())
		}
	}

	built = append(built, output)
	return strings.Join(built, " ")
}
