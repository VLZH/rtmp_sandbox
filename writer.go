package main

import (
	"log"

	"github.com/3d0c/gmf"
)

// CreateWriter is
func CreateWriter(ch chan *gmf.Packet, dst string) (Writer, error) {
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
	vc, err := gmf.FindEncoder(gmf.AV_CODEC_ID_FLV1)
	if err != nil {
		log.Println("ERROR: on finding encoder for flv in writer", err.Error())
	}
	log.Printf("INFO: Output video codec: %v \n", vc.Name())
	vcc := gmf.NewCodecCtx(vc)
	// audio
	ac, err := gmf.FindEncoder("libmp3lame")
	if err != nil {
		log.Println("ERROR: on finding encoder for mp3 in writer", err.Error())
	}
	log.Printf("INFO: Output audio codec: %v \n", ac.Name())
	acc := gmf.NewCodecCtx(ac)
	// add Streams
	sv, _ := wr.OutputContex.AddStreamWithCodeCtx(vcc)
	sa, _ := wr.OutputContex.AddStreamWithCodeCtx(acc)
	log.Printf("INFO: output video stream index: %v, audio stream index: %v, streams count: %v \n", sv.Index(), sa.Index(), wr.OutputContex.StreamsCnt())
	wr.writeHeader()
}

// StartLoop is function for starting listen chan of packets and write they to muxer
func (wr *Writer) StartLoop() {
	var err error
	for {
		pkt := <-wr.Ch
		if pkt != nil {
			log.Printf("INFO: Packet pts: %v, data size: %v stream index: %v \n", pkt.Pts(), len(pkt.Data()), pkt.StreamIndex())
			err = wr.OutputContex.WritePacket(pkt)
			gmf.Release(pkt)
			if err != nil {
				log.Println("ERROR: on writing packet to output context", err.Error())
			}
		} else {
		}
	}
}

func (wr *Writer) writeHeader() {
	if err := wr.OutputContex.WriteHeader(); err != nil {
		log.Println("ERROR: on write header to output; ", err.Error())
	}
}

func (wr *Writer) writeTrailer() {
	wr.OutputContex.WriteTrailer()
}

func (wr *Writer) free() {
	wr.OutputContex.CloseOutputAndRelease()
}
