package main

import (
	"log"
	"time"

	"github.com/imkira/go-libav/avcodec"
	"github.com/imkira/go-libav/avformat"
)

func init() {
	avformat.RegisterAll()
}

func asyncCopyPackets() {
	ch := make(chan *avcodec.Packet)
	files := []*VFile{
		&VFile{Name: "./1.mp4"},
		&VFile{Name: "./2.mp4"},
		&VFile{Name: "./3.mp4"},
	}
	reader, _ := CreateReader(ch, files)
	// rtmp := "rtmp://live-prg.twitch.tv/live_129862765_H7988wWNq4m2kNaPPwnHkIxRKIsoDB"
	writer, _ := CreateWriter(ch, "./t.flv")
	// read
	log.Println("Reader Start Loop")
	go reader.StartLoop()
	// write
	log.Println("Writer Prepare")
	writer.Prepare()
	log.Println("Writer Start Loop")
	go writer.StartLoop()
	time.Sleep(time.Minute * 4)
}

func main() {
	asyncCopyPackets()
}
