package main

import (
	"log"
	"time"

	"github.com/3d0c/gmf"
)

// CreateWriter is
func CreateWriter(ch chan *gmf.Packet, chclose chan bool, dst string) (Writer, error) {
	return Writer{
		Ch:          ch,
		Destination: dst,
	}, nil
}

// Writer is struct of writer to a muxer
type Writer struct {
	Ch          chan *gmf.Packet
	Destination string
	//
	OutputContex *gmf.FmtCtx
}

// Prepare is
func (wr *Writer) Prepare() {
	var err error
	wr.OutputContex, err = gmf.NewOutputCtxWithFormatName(wr.Destination, "flv")
	if err != nil {
		log.Println("ERROR: on createing output context", err.Error())
	}
	// video
	vc, err := gmf.FindEncoder(gmf.AV_CODEC_ID_H264)
	if err != nil {
		log.Println("ERROR: on finding video encoder in writer", err.Error())
	}
	log.Printf("INFO: Output video codec: %v, id: %v, full: %v \n", vc.Name(), vc.Id(), vc.LongName())
	vcc := gmf.NewCodecCtx(vc).
		SetHeight(320).
		SetWidth(620).
		SetProfile(gmf.FF_PROFILE_H264_BASELINE).
		SetPixFmt(gmf.AV_PIX_FMT_RGB32)
	sv, _ := wr.OutputContex.AddStreamWithCodeCtx(vcc)
	// audio
	ac, err := gmf.FindDecoder("aac")
	if err != nil {
		log.Println("ERROR: on finding audio encoder in writer", err.Error())
	}
	log.Printf("INFO: Output audio codec: %v, id: %v, full: %v \n", ac.Name(), ac.Id(), ac.LongName())
	acc := gmf.NewCodecCtx(ac).SetSampleFmt(int32(8)).SetSampleRate(44100)
	sa, _ := wr.OutputContex.AddStreamWithCodeCtx(acc)
	log.Printf("INFO: Output video stream index: %v, audio stream index: %v, streams count: %v \n", sv.Index(), sa.Index(), wr.OutputContex.StreamsCnt())
	wr.writeHeader()
}

// StartLoop is function for starting listen chan of packets and write they to muxer
func (wr *Writer) StartLoop() {
	var err error
	var prevPacketTime int64
	var currentPacketTime int64
	for {
		pkt := <-wr.Ch
		if pkt != nil {
			err = wr.OutputContex.WritePacketNoBuffer(pkt)
			gmf.Release(pkt)
			if err != nil {
				log.Println("ERROR: on writing packet to output context", err.Error())
			}
			currentPacketTime = time.Now().UnixNano()
			if prevPacketTime != 0 {
				log.Printf(
					"Time since previous packet: %v\n",
					sred((currentPacketTime-prevPacketTime)/1000000),
				)
			}
			prevPacketTime = currentPacketTime
		} else {
			// wr.writeTrailer()
			// wr.free()
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
