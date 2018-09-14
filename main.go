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

type SFrame struct {
	Frames      []*gmf.Frame
	StreamIndex int
	TimeBase    *gmf.AVRational
	OriginalPts int64
	Flush       int
}

func asyncCopyPackets() {
	ch := make(chan *SFrame, 100)
	chclose := make(chan bool)
	files := []*VFile{
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/0197072fdf4a4932427f16af754bfd34.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/04379f7432e53f85a3fd5cce00fb05b3.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/06da175b41884c53e985a950a6e5db4b.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/084691b8a639da050b5cf7330665df03.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/18786e742ec2c1d9298a33e5796166c7.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/1a4258ad35467687b5a04e88620aee40.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/1c5b7616b45a02074fcde2ea6e6820b0.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/1ccf32ecee68d4ad68946a3ef90b9b97.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/1d1a3225641a42fbbbcdf0e1f6416f71.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/23f062fa3000af7ac2d73eef4e8dd0b1.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/332d2a7c50c15e908668e6352989b4f8.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/3815ab29516820019c551817282bcc48.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/3931b879cc8112a0751a20120907dd32.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/4d3ebd2e6ef6742d280063c47ac18d57.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/53290712ae7f745d16208b5a016714a0.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/7fab65aafe672cf3433c483fbfd3a3be.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/858b9e10325283bb4116dee4589db37d.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/8d107a50304d4102b0fc575462f583c4.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/9510cb7550c3d7702ebed9eb09003714.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/99fba410ca62f5800688ecdb6f2bc638.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/9ebe67bd409fdd886b24268fff5f067a.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/a5ef36e31fe30f7cba36f511ef23cb79.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/a7ffeab6c339b1d911f62a0d04641057.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/a90e38db19520ada9f7fdb589f076075.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/bd5bc0218cb2f555c1397456de255cef.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/c993338943a6abba2581a1b4cb4a075b.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/da0149c7e60807109945bc424b65e6f9.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/eaaf1f3171184482df5559f8ca6fb8a9.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/ec2c9e3de8d7d041b6d4bd058d6a3b88.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/f6bb480b02b19b3b2b7396936d7eee86.mp4", DestHeight: 320, DestWidth: 640},
		&VFile{Path: "/Users/vladimirzhdanov/go/src/bss/bss_server/media/f8c5ec2702b8cb8042cc269d2bf9d94c.mp4", DestHeight: 320, DestWidth: 640},
	}
	reader, _ := CreateReader(ch, chclose, files)
	rtmp := "rtmp://bsslive.com:1935/stream/"
	// rtmp := "test.flv"
	writer, _ := CreateWriter(ch, chclose, rtmp)
	// write
	log.Println("INFO: Writer Prepare")
	writer.Prepare()
	log.Println("INFO: Writer Start Loop")
	go writer.StartLoop()
	// read
	log.Println("INFO: Reader Start Loop")
	go reader.StartLoop()
	time.Sleep(time.Minute * 60)
}

func main() {
	asyncCopyPackets()
}
