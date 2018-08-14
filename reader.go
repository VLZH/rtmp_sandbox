package main

import (
	"fmt"
	"log"
	"time"

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

// StartLoop is
func (r *Reader) StartLoop() {
	var cFile *VFile
	startTime := time.Now().UnixNano()
	prevFilePts := int64(0)
	prevPacketPts := int64(0)
	timeDiff := int64(0)
	for {
		cFile = r.GetNextFile()
		log.Printf("INFO: File: %v \n", cFile.Path)
		cFile.Prepare()
		for {
			pkt := cFile.ReadPacket()
			if pkt == nil {
				prevFilePts = prevPacketPts
				break
			}
			pkt.SetPts(pkt.Pts() + prevFilePts)
			pkt.SetDts(pkt.Dts() + prevFilePts)
			prevPacketPts = pkt.Pts()
			timeDiff = (time.Now().UnixNano() - startTime) / int64(1000000)
			fmt.Println("Sleep: ", prevPacketPts-timeDiff, " Millisecond", timeDiff)
			time.Sleep(time.Millisecond * time.Duration(prevPacketPts-timeDiff))
			r.Ch <- pkt
		}
		cFile.free()
		fmt.Println("End of file", cFile)
	}
}

// GetNextFile is
func (r *Reader) GetNextFile() *VFile {
	nextIndex := r.Idx + 1
	if nextIndex >= len(r.Files) {
		nextIndex = 0
	}
	r.Idx = nextIndex
	return r.Files[r.Idx]
}
