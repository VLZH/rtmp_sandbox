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
	// output
	OutputCodecContext *gmf.CodecCtx
	OutputCodec        *gmf.Codec
	// meta
	Height int
	Width  int
}

// Prepare is
func (v *VFile) Prepare() error {
	var err error
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
	// height, width
	v.Width = v.InputCodecContext.Width()
	v.Height = v.InputCodecContext.Height()
	log.Printf("INFO: Input codec width: %v, height: %v \n", v.Width, v.Height)
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
	// output
	v.OutputCodec, err = gmf.FindEncoder(gmf.AV_CODEC_ID_FLV1)
	if err != nil {
		log.Println("ERROR: on finding flv codec", err.Error())
	}
	v.OutputCodecContext = gmf.NewCodecCtx(v.OutputCodec)
	v.OutputCodecContext.SetBitRate(v.InputCodecContext.BitRate()).
		SetWidth(v.InputCodecContext.Width()).
		SetHeight(v.InputCodecContext.Height()).
		SetPixFmt(v.InputCodecContext.PixFmt()).
		SetTimeBase(v.InputCodecContext.TimeBase().AVR()).
		SetFlag(gmf.SWS_BILINEAR)
	if err = v.OutputCodecContext.Open(nil); err != nil {
		log.Println("ERROR: on open codecContext", err.Error())
	}
	return err
}

func (v *VFile) free() {
	v.InputContext.CloseInputAndRelease()
	gmf.Release(v.InputVideoStream)
	gmf.Release(v.InputAudioStream)
	gmf.Release(v.InputCodecContext)
	gmf.Release(v.InputCodec)
	// output
	gmf.Release(v.OutputCodec)
	gmf.Release(v.OutputCodecContext)
}

func (v *VFile) ReadPacket() *gmf.Packet {
	var err error
	var op *gmf.Packet
	var currentPacketStream *gmf.Stream
	ip := v.InputContext.GetNextPacket()
	if ip == nil {
		return nil
	}
	defer gmf.Release(ip)
	if ip.StreamIndex() == v.InputVideoStream.Index() {
		// is video packet
		currentPacketStream = v.InputVideoStream
		f, err := ip.Frames(v.InputCodecContext)
		if err != nil {
			log.Println("ERROR: on getting frame from packet", err.Error())
		}
		defer gmf.Release(f)
		op, err = f.Encode(v.OutputCodecContext)
	} else if ip.StreamIndex() == v.InputAudioStream.Index() {
		// is audio
		currentPacketStream = v.InputAudioStream
		op = ip.Clone()
	}
	gmf.RescaleTs(op, currentPacketStream.TimeBase(), gmf.AVR{Num: 1, Den: 1000}.AVRational())
	log.Printf(`
New packet in chan:
Source packet:      | pts: %v; data size: %v; duration: %v;
Destination packet: | pts: %v; data size: %v; duration: %v;
		`, ip.Pts(), len(ip.Data()), ip.Duration(), op.Pts(), len(op.Data()), op.Duration())
	if err != nil {
		log.Println("ERROR: on encoding to flv", err.Error())
	}
	return op
}
