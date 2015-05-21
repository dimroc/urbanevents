package cityrecorder

import (
	"bufio"
	"encoding/json"
	"github.com/pusher/pusher-http-go"
	"io"
	"log"
)

type Pusher struct {
	AppId  string
	Key    string
	Secret string
}

func (p *Pusher) Start(rd io.Reader) {
	client := pusher.Client{
		AppId:  p.AppId,
		Key:    p.Key,
		Secret: p.Secret,
	}

	reader := bufio.NewReader(rd)

	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			continue
		}

		var objmap *json.RawMessage
		err = json.Unmarshal(data, &objmap)

		if err != nil {
			log.Fatal(err)
		} else {
			client.Trigger("nyc", "tweet", objmap)
		}
	}
}
