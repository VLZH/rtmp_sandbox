package main

import (
	"log"
	"time"

	"github.com/3d0c/gmf"
	"github.com/fatih/color"
)

var OutputPixFormat = gmf.AV_PIX_FMT_YUV420P
var IS_VIDEO = "IS_VIDEO"
var IS_AUDIO = "IS_AUDIO"

var sgreen = color.New(color.FgGreen).SprintFunc()
var sblue = color.New(color.FgBlue).SprintFunc()
var sred = color.New(color.FgRed).SprintFunc()

func asyncCopyPackets() {
	ch := make(chan *gmf.Packet)
	chclose := make(chan bool)
	files := []*VFile{
		&VFile{Path: "./1.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "./2.mp4", DestHeight: 320, DestWidth: 640},
	}
	reader, _ := CreateReader(ch, chclose, files)
	rtmp := "rtmp://95.213.204.75:1935/stream/test"
	writer, _ := CreateWriter(ch, chclose, rtmp)
	// write
	log.Println("INFO: Writer Prepare")
	writer.Prepare()
	log.Println("INFO: Writer Start Loop")
	go writer.StartLoop()
	// TEST!
	time.Sleep(time.Second * 3)
	// read
	log.Println("INFO: Reader Start Loop")
	go reader.StartLoop()
	time.Sleep(time.Minute * 60)
}

func main() {
	asyncCopyPackets()
}
