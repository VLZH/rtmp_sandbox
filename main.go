package main

import (
	"log"
	"time"

	"github.com/3d0c/gmf"
)

func asyncCopyPackets() {
	ch := make(chan *gmf.Packet)
	chclose := make(chan bool)
	files := []*VFile{
		&VFile{Path: "./1.mp4"},
		&VFile{Path: "./2.mp4"},
	}
	reader, _ := CreateReader(ch, chclose, files)
	// rtmp := "rtmp://live-fra.twitch.tv/app/_"
	rtmp := "rtmp://95.213.204.75:1935/stream/test"
	writer, _ := CreateWriter(ch, chclose, rtmp)
	// read
	log.Println("INFO: Reader Start Loop")
	go reader.StartLoop()
	// write
	log.Println("INFO: Writer Prepare")
	writer.Prepare()
	log.Println("INFO: Writer Start Loop")
	go writer.StartLoop()
	time.Sleep(time.Minute * 60)
}

func main() {
	asyncCopyPackets()
}
