package main

import (
	"fmt"
	"log"
	"time"

	"github.com/3d0c/gmf"
)

// CreateWriter is
func CreateWriter(ch chan *SFrame, chclose chan bool, dst string) (Writer, error) {
	return Writer{
		Ch:          ch,
		Destination: dst,
	}, nil
}

// Writer is struct of writer to a muxer
type Writer struct {
	Ch          chan *SFrame
	ChS         chan *gmf.Stream
	Destination string
	//
	OutputContex *gmf.FmtCtx
	//
	OutputVideoCodecContext *gmf.CodecCtx
	OutputVideoCodec        *gmf.Codec
	OutputVideoStream       *gmf.Stream
	//
	OutputAudioCodecContext *gmf.CodecCtx
	OutputAudioCodec        *gmf.Codec
	OutputAudioStream       *gmf.Stream
}

// Prepare is
func (wr *Writer) Prepare() {
	var err error
	wr.OutputContex, err = gmf.NewOutputCtxWithFormatName(wr.Destination, "flv")
	if err != nil {
		log.Fatal("ERROR: on createing output context", err.Error())
	}
	if wr.OutputContex.IsGlobalHeader() {
		wr.OutputContex.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}
	wr.RegisterStreams()
	wr.writeHeader()
}

// RegisterStreams is
func (wr *Writer) RegisterStreams() {
	wr.RegisterVideoStream()
	wr.RegisterAudioStream()
}

// RegisterVideoStream is
func (wr *Writer) RegisterVideoStream() {
	var err error
	// codec
	wr.OutputVideoCodec, err = gmf.FindEncoder("libx264")
	if err != nil {
		log.Fatal("ERROR: cannot get video encoder")
	}
	// stream
	wr.OutputVideoStream = wr.OutputContex.NewStream(wr.OutputVideoCodec)
	// codec context
	wr.OutputVideoCodecContext = gmf.NewCodecCtx(wr.OutputVideoCodec).
		SetHeight(320).
		SetWidth(640).
		SetProfile(gmf.FF_PROFILE_H264_BASELINE).
		SetPixFmt(gmf.AV_PIX_FMT_YUV420P).
		SetTimeBase(gmf.AVR{Num: 1, Den: 25})
	defer gmf.Release(wr.OutputVideoStream)
	if wr.OutputVideoCodec.IsExperimental() {
		wr.OutputVideoCodecContext.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}
	if wr.OutputContex.IsGlobalHeader() {
		wr.OutputVideoCodecContext.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}
	//
	wr.OutputVideoStream.SetTimeBase(gmf.AVR{Num: 1, Den: 25})
	wr.OutputVideoStream.SetRFrameRate(gmf.AVR{Num: 25, Den: 1})
	// prepare from input stream
	if err := wr.OutputVideoCodecContext.Open(nil); err != nil {
		log.Fatal("ERROR: Cannot open OutputVideoCodecContext")
	}
	fmt.Printf("wr.OutputVideoStream time base: %v, Index: %v\n", wr.OutputVideoStream.TimeBase(), wr.OutputVideoStream.Index())
	wr.OutputVideoStream.SetCodecCtx(wr.OutputVideoCodecContext)
}

// RegisterAudioStream is
func (wr *Writer) RegisterAudioStream() {
	var err error
	// codec
	wr.OutputAudioCodec, err = gmf.FindEncoder("aac")
	if err != nil {
		log.Fatal("ERROR: cannot get audio encoder")
	}
	// stream
	wr.OutputAudioStream = wr.OutputContex.NewStream(wr.OutputAudioCodec)
	// codec context
	wr.OutputAudioCodecContext = gmf.NewCodecCtx(wr.OutputAudioCodec)
	defer gmf.Release(wr.OutputAudioCodecContext)
	//
	if wr.OutputAudioCodec.IsExperimental() {
		wr.OutputAudioCodecContext.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}
	if wr.OutputContex.IsGlobalHeader() {
		wr.OutputAudioCodecContext.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}
	//
	wr.OutputAudioCodecContext.
		SetTimeBase(gmf.AVR{Num: 1, Den: 44100}).
		SetSampleRate(44100).
		SetChannels(2).
		SetSampleFmt(gmf.AV_SAMPLE_FMT_FLTP)
		// SetChannelLayout(255)
	//
	if err := wr.OutputAudioCodecContext.Open(nil); err != nil {
		log.Fatal("ERROR: Cannot open OutputAudioCodecContext")
	}
	fmt.Printf("wr.OutputAudioStream time base: %v; Index: %v\n", wr.OutputAudioStream.TimeBase(), wr.OutputAudioStream.Index())
	wr.OutputAudioStream.SetCodecCtx(wr.OutputAudioCodecContext)
}

// StartLoop is function for starting listen chan of packets and write they to muxer
func (wr *Writer) StartLoop() {
	var (
		packets   []*gmf.Packet
		startTime int64
	)
	for {
		f := <-wr.Ch
		stream, err := wr.OutputContex.GetStream(f.StreamIndex)
		if err != nil {
			log.Fatal("ERROR: cannot get output stream by index")
		}
		for _, _f := range f.Frames {
			fmt.Printf("Frame pts: %v\n", _f.PktPts())
		}
		packets, err = stream.CodecCtx().Encode(f.Frames, f.Flush)
		if err != nil {
			log.Fatalf("ERROR: on getting packets from frames; error: %v", err.Error())
		}
		if startTime == 0 {
			startTime = time.Now().UnixNano()
		}
		for _, op := range packets {
			fmt.Printf("Packet pts: %v\n", op.Pts())
			gmf.RescaleTs(op, *f.TimeBase, stream.TimeBase())
			op.SetStreamIndex(f.StreamIndex)
			//
			err = wr.OutputContex.WritePacket(op)
			if err != nil {
				log.Fatalf("ERROR: on writing packet to output; %v", err.Error())
			}
			//
			if f.StreamIndex == 0 {
				diff := (time.Now().UnixNano() - startTime) / 1000000
				sleep_time := op.Pts() - diff
				log.Printf("Sleep: %v; St: %v; Pts: %v; TimeBase: %v; stream.TimeBase: %v ; Is video: %v\n",
					sleep_time, startTime, op.Pts(), *f.TimeBase, stream.TimeBase(), stream.IsVideo())
				time.Sleep(time.Millisecond * time.Duration(sleep_time))
			}
			op.Free()
		}
	}
}

func (wr *Writer) writeHeader() {
	log.Println("INFO: write header")
	if err := wr.OutputContex.WriteHeader(); err != nil {
		log.Println("ERROR: on write header to output; ", err.Error())
	}
}

func (wr *Writer) writeTrailer() {
	log.Println("INFO: write trailer")
	wr.OutputContex.WriteTrailer()
}

func (wr *Writer) free() {
	wr.OutputContex.CloseOutputAndRelease()
}
