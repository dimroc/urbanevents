package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/dimroc/urban-events/flagvalidator"
	"github.com/pusher/pusher-http-go"
	"log"
	"os"
)

var (
	pusherAppId  = flag.String("appid", "", "Pusher App Id")
	pusherKey    = flag.String("key", "", "Pusher Key")
	pusherSecret = flag.String("secret", "", "Pusher Secret")
)

func main() {
	flag.Parse()
	flags := []string{"appid", "key", "secret"}
	flagvalidator.ValidateFlags(flags)

	client := pusher.Client{
		AppId:  *pusherAppId,
		Key:    *pusherKey,
		Secret: *pusherSecret,
	}

	reader := bufio.NewReader(os.Stdin)
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
