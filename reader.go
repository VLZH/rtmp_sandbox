package main

import (
	"fmt"
	"time"

	"github.com/nareix/joy4/av"
)

// CreateReader return new Reader with chanel
func CreateReader(ch chan av.Packet, headCh chan []av.CodecData, files []*VFile) (Reader, chan bool) {
	r := Reader{
		Ch:      ch,
		HeadCh:  headCh,
		CloseCh: make(chan bool),
		Files:   files,
	}
	return r, r.CloseCh
}

// Reader is
type Reader struct {
	Ch      chan av.Packet
	HeadCh  chan []av.CodecData
	CloseCh chan bool
	Files   []*VFile
	Idx     int
}

func (r *Reader) StartLoop() {
	var cFile *VFile
	times := make(map[int8]time.Duration)
	for {
		cFile = r.GetNextFile()
		fmt.Printf("Next file; %s \n", cFile)
		demuxer, err := cFile.GetDemuxer()
		if err != nil {
			fmt.Printf("Error on getting demuxer from VFile %s \n", cFile)
			fmt.Println(err)
			r.CloseCh <- true
		}
		// Write headers
		codecDat, err := demuxer.Streams()
		if err != nil {
			fmt.Printf("Error on getting streams from Demuxer %s \n", cFile)
			fmt.Println(err)
			r.CloseCh <- true
		}
		r.HeadCh <- codecDat
		// Write packets
		for {
			pkg, err := demuxer.ReadPacket()
			// update time
			t, ok := times[pkg.Idx]
			if !ok {
				times[pkg.Idx] = 0
			}
			if pkg.Time < t {
				times[pkg.Idx] = t + pkg.Time
			} else {
				times[pkg.Idx] = pkg.Time
			}
			pkg.Time = times[pkg.Idx]
			if err != nil {
				fmt.Println("Error on getting packet;", err)
				break
			}
			r.Ch <- pkg
		}
		demuxer.Close()
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
