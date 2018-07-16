package main

import (
	"github.com/imkira/go-libav/avcodec"
	"github.com/imkira/go-libav/avformat"
)

func init() {
	avformat.RegisterAll()
}

func asyncCopyPackets() {
	ch := make(chan avcodec.Packet)
	files := []*VFile{
		&VFile{Name: "./1.mp4"},
		&VFile{Name: "./2.mp4"},
		&VFile{Name: "./3.mp4"},
	}
	reader, _ := CreateReader(ch, files)
	reader.StartLoop()
}

func main() {
	asyncCopyPackets()
}
