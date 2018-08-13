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
		&VFile{Name: "./2.mp4"},
		&VFile{Name: "./1.mp4"},
	}
	reader, _ := CreateReader(ch, chclose, files)
	// rtmp := "rtmp://live-prg.twitch.tv/live_129862765_H7988wWNq4m2kNaPPwnHkIxRKIsoDB"
	writer, _ := CreateWriter(ch, chclose, "./t.flv")
	// read
	log.Println("INFO: Reader Start Loop")
	go reader.StartLoop()
	// write
	log.Println("INFO: Writer Prepare")
	writer.Prepare()
	log.Println("INFO: Writer Start Loop")
	go writer.StartLoop()
	time.Sleep(time.Minute * 4)
}

func main() {
	asyncCopyPackets()
}
