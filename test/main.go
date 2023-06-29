package main

import (
	"log"

	"gitea.heinrich.blue/PHI/webrtc-gst/pkg/webrtc"
	"github.com/tinyzimmer/go-glib/glib"
)

func main() {
	loop := glib.NewMainLoop(glib.MainContextDefault(), false)

	audioPipeline, audioCh, err := webrtc.CreateVideoPipeline("videotestsrc")
	if err != nil {
		panic(err)
	}

	go func() {
		for s := range audioCh {
			log.Println(s.Timestamp)
		}
	}()

	err = audioPipeline.Start()
	if err != nil {
		panic(err)
	}

	loop.Run()
}
