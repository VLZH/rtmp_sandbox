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
		&VFile{Name: "/Users/vladimirzhdanov/go/src/redlinestudio/stream_sandbox/1.mp4"},
		&VFile{Name: "/Users/vladimirzhdanov/go/src/redlinestudio/stream_sandbox/2.mp4"},
		&VFile{Name: "/Users/vladimirzhdanov/go/src/redlinestudio/stream_sandbox/3.mp4"},
	}
	// reader
	reader, closeChan := CreateReader(ch, headCh, files)
	// writer
	// rtmp_server := "rtmp://live-hel.twitch.tv/app/live_129862765_H7988wWNq4m2kNaPPwnHkIxRKIsoDB"
	// rtmp_server := "rtmp://a.rtmp.youtube.com/live2/x090-5e9d-dwra-em1g"
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
