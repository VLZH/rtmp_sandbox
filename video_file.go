package main

import (
	"fmt"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pktque"
	"github.com/nareix/joy4/av/transcode"
	"github.com/nareix/joy4/cgo/ffmpeg"
)

type VFile struct {
	Name string
}

// GetDemuxer is function for getting demuxer from file
func (vf *VFile) GetDemuxer() (av.DemuxCloser, error) {
	file, err := avutil.Open(vf.Name)
	if err != nil {
		fmt.Println("Error on open mp4 file;", err.Error())
		return nil, err
	}
	demuxer := &pktque.FilterDemuxer{Demuxer: file, Filter: &pktque.Walltime{}}
	findcodec := func(stream av.AudioCodecData, i int) (need bool, dec av.AudioDecoder, enc av.AudioEncoder, err error) {
		need = true
		dec, _ = ffmpeg.NewAudioDecoder(stream)
		enc, _ = ffmpeg.NewAudioEncoderByName("aac_at")
		err = enc.SetSampleRate(stream.SampleRate())
		if err != nil {
			fmt.Println("Error on SetSampleRate")
		}
		fmt.Println(stream.SampleRate(), "rate")
		enc.SetChannelLayout(av.CH_STEREO)
		enc.SetBitrate(48000)
		enc.SetOption("profile", "HE-AACv2")
		return
	}

	trans := &transcode.Demuxer{
		Options: transcode.Options{
			FindAudioDecoderEncoder: findcodec,
		},
		Demuxer: demuxer,
	}
	return trans, nil
}
