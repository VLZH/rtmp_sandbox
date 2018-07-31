package main

import (
	"log"

	"github.com/imkira/go-libav/avcodec"
	"github.com/imkira/go-libav/avformat"
)

func CreateWriter(ch chan *avcodec.Packet, dst string) (Writer, error) {
	return Writer{
		Ch:          ch,
		Destination: dst,
	}, nil
}

// Writer is struct of writer to a muxer
type Writer struct {
	Ch          chan *avcodec.Packet
	Destination string
	//
	encFmt    *avformat.Context
	encStream *avformat.Stream
	encIO     *avformat.IOContext
}

func (wr *Writer) Prepare() {
	wr.PrepareFormatContext()
	wr.PrepareStream()
	wr.PrepareIOContext()
	// wr.WriteHeader()
}

func (wr *Writer) PrepareFormatContext() {
	var err error
	output := avformat.GuessOutputFromFileName(wr.Destination)
	// output.
	if output == nil {
		log.Fatalf("Failed to guess output for output file: %s\n", wr.Destination)
	}
	if wr.encFmt, err = avformat.NewContextForOutput(output); err != nil {
		log.Fatalf("Failed to open output context: %v\n", err)
	}
	wr.encFmt.SetFileName(wr.Destination)
}

func (wr *Writer) PrepareStream() {
	var err error
	wr.encStream, err = wr.encFmt.NewStreamWithCodec(nil)
	if err != nil {
		log.Fatalf("Failed to open output video stream: %v\n", err)
	}
}

func (wr *Writer) PrepareIOContext() {
	var err error
	flags := avformat.IOFlagWrite
	if wr.encIO, err = avformat.OpenIOContext(wr.Destination, flags, nil, nil); err != nil {
		log.Fatalf("Failed to open I/O context: %v\n", err)
	}
	wr.encFmt.SetIOContext(wr.encIO)
}

// StartLoop is function for starting listen chan of packets and write they to muxer
func (wr *Writer) StartLoop() {
	var err error
	for {
		pkt := <-wr.Ch
		if pkt != nil {
			pkt.SetPosition(-1)
			pkt.SetStreamIndex(wr.encStream.Index())
			err = wr.encFmt.WriteFrame(pkt)
			if err != nil {
				log.Fatalf("Error on write packet; %v", err)
			}
			log.Println("Write packet")
		} else {
			// wr.WriteTrailer()
		}
	}
}

func (wr *Writer) WriteHeader() {
	log.Println("Write Header")
	if err := wr.encFmt.WriteHeader(nil); err != nil {
		log.Fatalf("Failed to write header: %v\n", err)
	}
}

func (wr *Writer) WriteTrailer() {
	log.Println("Write Trailer")
	if err := wr.encFmt.WriteTrailer(); err != nil {
		log.Fatalf("Failed to write trailer: %v\n", err)
	}
}
