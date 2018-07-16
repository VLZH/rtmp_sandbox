package main

import (
	"fmt"
	"log"

	"github.com/imkira/go-libav/avcodec"
	"github.com/imkira/go-libav/avformat"
	"github.com/imkira/go-libav/avutil"
)

type VFile struct {
	Name string
	// decode
	decFmt    *avformat.Context
	decStream *avformat.Stream // video stream
	codecCtx  *avcodec.Context
	decCodec  *avcodec.Context
	//
	decPkt   *avcodec.Packet
	decFrame *avutil.Frame
}

func (v *VFile) prepareContext() {
	var err error
	v.alloc()
	v.decFmt, err = avformat.NewContextForInput()
	if err != nil {
		log.Fatalf("Failed to open input context: %v\n", err)
	}
	options := avutil.NewDictionary()
	defer options.Free()
	// open
	if err = v.decFmt.OpenInput(v.Name, nil, options); err != nil {
		log.Fatalf("Failed to open input file: %v\n", err)
	}
	// find streams
	if err := v.decFmt.FindStreamInfo(nil); err != nil {
		log.Fatalf("Failed to find stream info: %v\n", err)
	}
	// dump streams to standard output
	v.decFmt.Dump(0, v.Name, false)

	//
	// PREPARE FIRST VIDEO STREAM FOR DECODING
	//

	// find first video stream
	if v.decStream = getFirstVideoStream(v.decFmt); v.decStream == nil {
		log.Fatalf("Could not find a video stream. Aborting...\n")
	}

	v.codecCtx = v.decStream.CodecContext()
	codec := avcodec.FindDecoderByID(v.codecCtx.CodecID())
	if codec == nil {
		log.Fatalf("Could not find decoder: %v\n", v.codecCtx.CodecID())
	}

	if v.decCodec, err = avcodec.NewContextWithCodec(codec); err != nil {
		log.Fatalf("Failed to create codec context: %v\n", err)
	}
	// if err := v.codecCtx.CopyTo(v.decCodec); err != nil {
	// 	log.Fatalf("Failed to copy codec context: %v\n", err)
	// }
	if err = v.decCodec.SetInt64Option("refcounted_frames", 1); err != nil {
		log.Fatalf("Failed to copy codec context: %v\n", err)
	}
	if err = v.decCodec.OpenWithCodec(codec, nil); err != nil {
		log.Fatalf("Failed to open codec: %v\n", err)
	}
}

func (v *VFile) readPacket() bool {
	// reading
	reading, err := v.decFmt.ReadFrame(v.decPkt)
	if err != nil {
		log.Fatalf("Failed to read packet: %v\n", err)
	}
	if !reading {
		return false
	}
	defer v.decPkt.Unref()
	// is video packet?
	if v.decPkt.StreamIndex() != v.decStream.Index() {
		return true
	}
	v.decPkt.RescaleTime(v.decStream.TimeBase(), v.decCodec.TimeBase())
	var decoded bool
	for v.decPkt.Size() > 0 {
		fmt.Println("Size:", v.decPkt.Size(), "Packet:", v.decPkt.PTS())
		decoded = true
	}
	return decoded
}

func (v *VFile) alloc() error {
	var err error
	if v.decPkt, err = avcodec.NewPacket(); err != nil {
		return err
	}
	if v.decFrame, err = avutil.NewFrame(); err != nil {
		return err
	}
	return nil
}

func (v *VFile) free() {
	v.decFmt.Free()
	v.codecCtx.Free()
	v.decCodec.Free()
	v.decPkt.Free()
	v.decFrame.Free()
}

func getFirstVideoStream(fmtCtx *avformat.Context) *avformat.Stream {
	for _, stream := range fmtCtx.Streams() {
		switch stream.CodecContext().CodecType() {
		case avutil.MediaTypeVideo:
			return stream
		}
	}
	return nil
}
