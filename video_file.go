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
	OutputFrame        *gmf.Frame
	// sws
	SwsContext *gmf.SwsCtx
	// meta
	Height int
	Width  int
	//
	DestHeight int
	DestWidth  int
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
	v.OutputCodec, err = gmf.FindEncoder(gmf.AV_CODEC_ID_H264)
	if err != nil {
		log.Println("ERROR: on finding output codec", err.Error())
	}
	v.OutputCodecContext = gmf.NewCodecCtx(v.OutputCodec)
	v.OutputCodecContext.
		SetWidth(v.DestWidth).
		SetHeight(v.DestHeight).
		SetPixFmt(OutputPixFormat).
		SetTimeBase(v.InputCodecContext.TimeBase().AVR()).
		SetFlag(gmf.SWS_BILINEAR).
		SetProfile(gmf.FF_PROFILE_H264_BASELINE)
	if err = v.OutputCodecContext.Open(nil); err != nil {
		log.Println("ERROR: on open codecContext", err.Error())
	}
	v.OutputFrame = gmf.NewFrame(). // output frame
					SetWidth(v.DestWidth).
					SetHeight(v.DestHeight).
					SetFormat(OutputPixFormat)
	if err := v.OutputFrame.ImgAlloc(); err != nil {
		log.Fatal(err)
	}
	// sws context
	log.Printf(`
INFO: SWS CONTEXT:
Source:      Height:%v; Width:%v; PixFmt:%v;
Destination: Height:%v; Width:%v; PixFmt:%v;
	`, v.InputCodecContext.Height(), v.InputCodecContext.Width(), v.InputCodecContext.PixFmt(),
		v.OutputCodecContext.Height(), v.OutputCodecContext.Width(), v.OutputCodecContext.PixFmt())
	v.SwsContext = gmf.NewSwsCtx(v.InputCodecContext, v.OutputCodecContext, gmf.SWS_BILINEAR)
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
	gmf.Release(v.OutputFrame)
	// sws
	gmf.Release(v.SwsContext)
}

// ReadPacket is
func (v *VFile) ReadPacket() (*gmf.Packet, string) {
	var err error
	var op *gmf.Packet
	var f *gmf.Frame
	var currentPacketStream *gmf.Stream
	var packetType = ""
	//
	ip := v.InputContext.GetNextPacket()
	if ip == nil {
		return nil, ""
	}
	defer gmf.Release(ip)
	if ip.StreamIndex() == v.InputVideoStream.Index() {
		// is video packet
		packetType = IS_VIDEO
		currentPacketStream = v.InputVideoStream
		f, err = ip.Frames(v.InputCodecContext)
		if f == nil {
			return v.ReadPacket()
		}
		if err != nil {
			log.Println("ERROR: on getting frame from packet", err.Error())
		}
		v.SwsContext.Scale(f, v.OutputFrame)
		defer gmf.Release(f)
		op, err = v.OutputFrame.Encode(v.OutputCodecContext)
		if op == nil {
			return v.ReadPacket()
		}
		op.SetPts(ip.Pts())
		op.SetDts(ip.Dts())
		if err != nil {
			log.Println("ERROR: on encoding", err.Error())
		}
	} else if ip.StreamIndex() == v.InputAudioStream.Index() {
		// is audio
		packetType = IS_AUDIO
		currentPacketStream = v.InputAudioStream
		op = ip.Clone()
	}
	gmf.RescaleTs(op, currentPacketStream.TimeBase(), gmf.AVR{Num: 1, Den: 1000}.AVRational())
	log.Printf(`
New packet in chan:
Packet type:        | %v;
Source packet:      | pts: %v; data size: %v; duration: %v;
Destination packet: | pts: %v; data size: %v; duration: %v;`,
		sgreen(packetType),
		sgreen(ip.Pts()), sgreen(len(ip.Data())), sgreen(ip.Duration()),
		sgreen(op.Pts()), sgreen(len(op.Data())), sgreen(op.Duration()))
	return op, packetType
}
