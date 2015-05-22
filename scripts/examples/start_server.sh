export TWITTER_CONSUMER_KEY=
export TWITTER_CONSUMER_SECRET=
export TWITTER_TOKEN=
export TWITTER_TOKEN_SECRET=

#Clients can be instantiated from a specially-crafted Pusher URL. For example:
	#client := pusher.ClientFromURL("http://key:secret@api.pusherapp.com/apps/app_id")

export PUSHER_URL=

cd ~/go_workspace
export GOPATH=`pwd`;
cd /Users/dimroc/go_workspace/src/github.com/dimroc/urban-events/cityrecorder
go run main.go


