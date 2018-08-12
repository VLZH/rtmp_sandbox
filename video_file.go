package main

import (
	"log"

	"github.com/3d0c/gmf"
)

// VFile is
type VFile struct {
	Name string
	// input
	InputContext      *gmf.FmtCtx
	InputVideoStream  *gmf.Stream // video stream for detecting codec
	InputCodecContext *gmf.CodecCtx
	InputCodec        *gmf.Codec
	// output
	OutputCodecContext *gmf.CodecCtx
	OutputCodec        *gmf.Codec
	// meta
	Height int
	Width  int
}

// ReadPacket is
func (v *VFile) readPacket() (*gmf.Packet, error) {
	return v.getNextFlvPacket(), nil
}

// Prepare is
func (v *VFile) prepare() error {
	var err error
	// input
	v.InputContext, err = gmf.NewInputCtx(v.Name)
	if err != nil {
		log.Println("ERROR: on getting context for input", err.Error())
	}
	v.InputVideoStream, err = v.InputContext.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		log.Println("ERROR: on getting best stream from input context", err.Error())
	}
	v.InputCodecContext = v.InputVideoStream.CodecCtx()
	v.InputCodec = v.InputCodecContext.Codec()
	v.Width = v.InputCodecContext.Width()
	v.Height = v.InputCodecContext.Height()
	log.Printf("INFO: Input codec width: %v, height: %v \n", v.Width, v.Height)
	log.Println("INFO: Streams:")
	for i := 0; i < v.InputContext.StreamsCnt(); i++ {
		srcStream, err := v.InputContext.GetStream(i)
		if err != nil {
			log.Println("ERROR: on getting stream by index: ", i, err.Error())
		}
		log.Printf(
			"Stream #%v; Is audio: %v; Is video: %v; Codec: %v, Codec id: %v \n",
			srcStream.Index(), srcStream.IsAudio(), srcStream.IsVideo(),
			srcStream.CodecCtx().Codec().Name(), srcStream.CodecCtx().Codec().Id())
	}
	// output
	v.OutputCodec, err = gmf.FindEncoder(gmf.AV_CODEC_ID_FLV1)
	if err != nil {
		log.Println("ERROR: on finding flv codec", err.Error())
	}
	v.OutputCodecContext = gmf.NewCodecCtx(v.OutputCodec)
	v.OutputCodecContext.SetBitRate(40000).
		SetWidth(v.InputCodecContext.Width()).
		SetHeight(v.InputCodecContext.Width()).
		SetPixFmt(v.InputCodecContext.PixFmt()).
		SetTimeBase(gmf.AVR{Num: 1, Den: 25})
	if err = v.OutputCodecContext.Open(nil); err != nil {
		log.Println("ERROR: on open codecContext", err.Error())
	}
	return err
}

func (v *VFile) free() {
	v.InputContext.CloseInputAndRelease()
	gmf.Release(v.InputVideoStream)
	gmf.Release(v.InputCodecContext)
	gmf.Release(v.InputCodec)
	// output
	gmf.Release(v.OutputCodec)
	gmf.Release(v.OutputCodecContext)
}

// TODO: rename or move code about audio to func readPacket
func (v *VFile) getNextFlvPacket() *gmf.Packet {
	var err error
	ip := v.InputContext.GetNextPacket()
	defer gmf.Release(ip)
	if ip.StreamIndex() == v.InputVideoStream.Index() {
		op := ip.Clone()
		return op
	}
	f, err := ip.Frames(v.InputCodecContext)
	if err != nil {
		log.Println("ERROR: on getting frame from packet", err.Error())
	}
	defer gmf.Release(f)
	of := f.CloneNewFrame()
	defer gmf.Release(of)
	op, err := of.Encode(v.OutputCodecContext)
	if err != nil {
		log.Println("ERROR: on encoding to flv", err.Error())
	}
	return op
}
