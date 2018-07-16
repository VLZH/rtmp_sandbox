package main

import (
	"fmt"

	"github.com/imkira/go-libav/avcodec"
)

// CreateReader return new Reader with chanel
func CreateReader(ch chan avcodec.Packet, files []*VFile) (Reader, chan bool) {
	r := Reader{
		Ch:      ch,
		CloseCh: make(chan bool),
		Files:   files,
	}
	return r, r.CloseCh
}

// Reader is
type Reader struct {
	Ch      chan avcodec.Packet
	CloseCh chan bool
	Files   []*VFile
	Idx     int
}

func (r *Reader) StartLoop() {
	var cFile *VFile
	for {
		cFile = r.GetNextFile()
		cFile.prepareContext()
		cFile.readPacket()
		fmt.Println("Start loop", cFile)
	}
}

func (r *Reader) GetNextFile() *VFile {
	nextIndex := r.Idx + 1
	if nextIndex >= len(r.Files) {
		nextIndex = 0
	}
	r.Idx = nextIndex
	return r.Files[r.Idx]
}
