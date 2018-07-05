package main

import (
	"fmt"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"
)

func init() {
	rtmp.Debug = true
	format.RegisterAll()
}

func asyncCopyPackets() {
	ch := make(chan av.Packet)
	headCh := make(chan []av.CodecData)
	files := []*VFile{
		&VFile{Name: "./1.mp4"},
		&VFile{Name: "./2.mp4"},
		&VFile{Name: "./3.mp4"},
	}
	// reader
	reader, closeChan := CreateReader(ch, headCh, files)
	// writer
	// rtmp_server := "rtmp://live-hel.twitch.tv/app/<key>"
	// rtmp_server := "rtmp://a.rtmp.youtube.com/live2/<key>"
	rtmp_server := "rtmp://127.0.0.1:1935/live/test"
	conn, err := rtmp.Dial(rtmp_server)
	if err != nil {
		fmt.Println("Error on open rtmp connection", err)
		return
	}
	writer := CreateWriter(ch, headCh, conn)
	// start gorutines
	go reader.StartLoop()
	go writer.StartLoop()
	<-closeChan
}

func main() {
	asyncCopyPackets()
}
