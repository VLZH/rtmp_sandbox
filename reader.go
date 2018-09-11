package main

import (
	"fmt"
	"log"
)

// CreateReader return new Reader with chanel
func CreateReader(ch chan *SFrame, chclose chan bool, files []*VFile) (Reader, chan bool) {
	r := Reader{
		Ch:      ch,
		CloseCh: chclose,
		Files:   files,
	}
	return r, r.CloseCh
}

// Reader is
type Reader struct {
	Ch      chan *SFrame
	CloseCh chan bool
	Files   []*VFile
	Idx     int
}

type FilePts map[int]int64

// StartLoop is
func (r *Reader) StartLoop() {
	var cFile *VFile
	prevFilePts := make(FilePts)
	prevPacketPts := make(FilePts)
	for {
		cFile = r.GetNextFile()
		log.Printf("INFO: File: %v \n", cFile.Path)
		cFile.Prepare()
		cFile.LogStreams()
		for {
			cf := cFile.ReadFrames()
			if cf == nil {
				// end of file
				prevFilePts[0] = prevPacketPts[0]
				prevFilePts[1] = prevPacketPts[1]
				break
			}
			// change pts
			for _, frame := range cf.Frames {
				pts := frame.Pts() + prevFilePts[cf.StreamIndex]
				prevPacketPts[cf.StreamIndex] = pts
				frame.SetPts(pts)
				frame.SetPktPts(pts)
				frame.SetPktDts(int(pts))
				log.Printf("Pts: %v; SetPktPts: %v; SetPktDts: %v; Time Base: %+v", frame.Pts(), frame.PktPts(), frame.PktDts(), cf.TimeBase)

			}
			r.Ch <- cf
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
