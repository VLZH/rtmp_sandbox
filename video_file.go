package main

import (
	"log"

	"github.com/3d0c/gmf"
)

// VFile is
type VFile struct {
	Path string
	// input
	InputContext      *gmf.FmtCtx
	InputVideoStream  *gmf.Stream // video stream for detecting codec
	InputAudioStream  *gmf.Stream // video stream for detecting codec
	InputCodecContext *gmf.CodecCtx
	InputCodec        *gmf.Codec
	// sws
	SwsContext *gmf.SwsCtx
	//
	DestHeight int
	DestWidth  int
	//
	CurrentPacket       *gmf.Packet
	CurrentPacketStream *gmf.Stream
	//
	Flush int
}

// Prepare is
func (v *VFile) Prepare() error {
	var err error
	v.Flush = -1
	// input
	v.InputContext, err = gmf.NewInputCtx(v.Path)
	if err != nil {
		log.Println("ERROR: on getting context for input", err.Error())
	}
	v.InputVideoStream, err = v.InputContext.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	v.InputAudioStream, err = v.InputContext.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
	if err != nil {
		log.Println("ERROR: on getting best stream from input context", err.Error())
	}
	v.InputCodecContext = v.InputVideoStream.CodecCtx()
	v.InputCodec = v.InputCodecContext.Codec()

	return err
}

// LogStreams is
func (v *VFile) LogStreams() {
	// height, width
	log.Printf("INFO: Input codec width: %v, height: %v \n", v.InputCodecContext.Width(), v.InputCodecContext.Height())
	log.Println("INFO: Streams:")
	// print info about input file streams
	for i := 0; i < v.InputContext.StreamsCnt(); i++ {
		srcStream, err := v.InputContext.GetStream(i)
		if err != nil {
			log.Println("ERROR: on getting stream by index: ", i, err.Error())
		}
		if srcStream.IsVideo() {
			log.Printf(
				"Stream #%v; VIDEO; Codec: %v, Codec profile: %v, Codec id: %v, Timebase: %+v\n",
				srcStream.Index(),
				srcStream.CodecCtx().Codec().Name(), srcStream.CodecCtx().Profile(), srcStream.CodecCtx().Codec().Id(),
				srcStream.TimeBase())
		} else if srcStream.IsAudio() {
			log.Printf(
				"Stream #%v; AUDIO; Codec: %v, Codec profile: %v, Codec id: %v, Timebase: %+v, sample fmt: %v, sample rate: %v\n",
				srcStream.Index(),
				srcStream.CodecCtx().Codec().Name(), srcStream.CodecCtx().Profile(), srcStream.CodecCtx().Codec().Id(),
				srcStream.TimeBase(), srcStream.CodecCtx().SampleFmt(), srcStream.CodecCtx().SampleRate())
		}
	}
}

func (v *VFile) PrepareSws() {
	// 	v.OutputCodecContext.
	// 		SetWidth(v.DestWidth).
	// 		SetHeight(v.DestHeight).
	// 		SetPixFmt(OutputPixFormat).
	// 		SetTimeBase(v.InputCodecContext.TimeBase().AVR()).
	// 		SetFlag(gmf.SWS_BILINEAR).
	// 		SetProfile(gmf.FF_PROFILE_H264_BASELINE)
	// 	if err = v.OutputCodecContext.Open(nil); err != nil {
	// 		log.Println("ERROR: on open codecContext", err.Error())
	// 	}
	// 	v.OutputFrame = gmf.NewFrame(). // output frame
	// 					SetWidth(v.DestWidth).
	// 					SetHeight(v.DestHeight).
	// 					SetFormat(OutputPixFormat)
	// 	if err := v.OutputFrame.ImgAlloc(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// sws context
	// 	log.Printf(`
	// INFO: SWS CONTEXT:
	// Source:      Height:%v; Width:%v; PixFmt:%v;
	// Destination: Height:%v; Width:%v; PixFmt:%v;
	// `, v.InputCodecContext.Height(), v.InputCodecContext.Width(), v.InputCodecContext.PixFmt(),
	// 		v.OutputCodecContext.Height(), v.OutputCodecContext.Width(), v.OutputCodecContext.PixFmt())
	// 	v.SwsContext = gmf.NewSwsCtx(v.InputCodecContext, v.OutputCodecContext, gmf.SWS_BILINEAR)
}

func (v *VFile) free() {
	v.InputContext.CloseInputAndRelease()
	gmf.Release(v.InputVideoStream)
	gmf.Release(v.InputAudioStream)
	gmf.Release(v.InputCodecContext)
	gmf.Release(v.InputCodec)
	// output
	if v.SwsContext != nil {
		gmf.Release(v.SwsContext)
	}
}

// ReadFrame is
func (v *VFile) ReadFrames() *SFrame {
	var (
		err               error
		outputStreamIndex = 0
		frames            []*gmf.Frame
		streamIndex       int
	)
	if v.Flush < 0 {
		if v.CurrentPacket != nil {
			v.CurrentPacket.Release()
		}
		v.CurrentPacket, err = v.InputContext.GetNextPacket()
		if err != nil && err.Error() != "End of file" {
			if v.CurrentPacket != nil {
				v.CurrentPacket.Free()
			}
			log.Fatalf("error getting next packet - %s", err)
		} else if err != nil && v.CurrentPacket == nil {
			log.Printf("=== Flushing \n")
			v.Flush++
			return nil
		}
	}
	if v.Flush == 2 {
		return nil
	}
	if v.Flush < 0 {
		streamIndex = v.CurrentPacket.StreamIndex()
	} else {
		streamIndex = v.Flush
		v.Flush++
	}
	v.CurrentPacketStream, err = v.InputContext.GetStream(streamIndex)
	if err != nil {
		log.Println("ERROR: on getting stream by id", err.Error())
	}
	// decode
	frames, err = v.CurrentPacketStream.CodecCtx().Decode(v.CurrentPacket)
	if err != nil {
		log.Println("ERROR: on getting frame from packet", err.Error())
	}
	// get outputStreamIndex
	if v.CurrentPacket.StreamIndex() == v.InputVideoStream.Index() {
		outputStreamIndex = 0
	} else if v.CurrentPacket.StreamIndex() == v.InputAudioStream.Index() {
		outputStreamIndex = 1
	}
	tb := v.CurrentPacketStream.TimeBase()
	return &SFrame{Frames: frames, StreamIndex: outputStreamIndex, TimeBase: &tb, Flush: v.Flush}
}
