package main

import (
	"fmt"

	"github.com/nareix/joy4/av"
)

func CreateWriter(ch chan av.Packet, headCh chan []av.CodecData, dst av.MuxCloser) Writer {
	return Writer{
		Ch:          ch,
		HeadCh:      headCh,
		Destination: dst,
	}
}

// Writer is struct of writer to a muxer
type Writer struct {
	Ch          chan av.Packet
	HeadCh      chan []av.CodecData
	Destination av.MuxCloser
}

// StartLoop is function for starting listen chan of packets and write they to muxer
func (wr *Writer) StartLoop() {
	for {
		select {
		case hDat := <-wr.HeadCh:
			fmt.Println("WriteHeader")
			// wr.Destination.WriteTrailer()
			wr.Destination.WriteHeader(hDat)

		case pkg := <-wr.Ch:
			err := wr.Destination.WritePacket(pkg)
			if err != nil {
				fmt.Println("Error on write packet to Destination", err)
			}
		default:

		}
	}
}
