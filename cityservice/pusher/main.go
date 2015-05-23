package main

import (
	"flag"
	"github.com/dimroc/urban-events/cityservice/cityrecorder"
	"github.com/dimroc/urban-events/cityservice/flagvalidator"
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

	pusher := cityrecorder.NewPusher(*pusherAppId, *pusherKey, *pusherSecret)
	pusher.Listen(os.Stdin)
}
