package cityrecorder

import (
	"bufio"
	"encoding/json"
	"github.com/pusher/pusher-http-go"
	"io"
	"log"
)

type Pusher struct {
	client *pusher.Client
}

func (p *Pusher) Write(g GeoEvent) {
	p.client.Trigger("nyc", "tweet", g)
}

// Blocking
func (p *Pusher) Listen(rd io.Reader) {
	reader := bufio.NewReader(rd)

	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			continue
		}

		g := GeoEvent{}
		err = json.Unmarshal(data, &g)
		p.Write(g)
	}
}

func NewPusher(appId string, key string, secret string) *Pusher {
	p := &Pusher{}
	p.client = &pusher.Client{
		AppId:  appId,
		Key:    key,
		Secret: secret,
	}

	return p
}
