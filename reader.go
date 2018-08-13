package main

import (
	"fmt"
	"log"

	"github.com/3d0c/gmf"
)

// CreateReader return new Reader with chanel
func CreateReader(ch chan *gmf.Packet, chclose chan bool, files []*VFile) (Reader, chan bool) {
	r := Reader{
		Ch:      ch,
		CloseCh: chclose,
		Files:   files,
	}
	return r, r.CloseCh
}

// Reader is
type Reader struct {
	Ch      chan *gmf.Packet
	CloseCh chan bool
	Files   []*VFile
	Idx     int
}

func (r *Reader) StartLoop() {
	var cFile *VFile
	for {
		cFile = r.GetNextFile()
		log.Printf("INFO: File: %v \n", cFile.Name)
		cFile.prepare()
		for {
			pkt, err := cFile.readPacket()
			if err != nil {
				log.Printf("Error on read packet: %v", err.Error())
				close(r.Ch)
				break
			}
			r.Ch <- pkt
		}
		cFile.free()
		fmt.Println("End of file", cFile)
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
