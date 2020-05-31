package main

import (
	"github.com/nareix/joy5/av"
	"github.com/nareix/joy5/av/avutil"
	"github.com/nareix/joy5/av/transcode"
	"github.com/nareix/joy5/cgo/ffmpeg"
	"github.com/nareix/joy5/format"
)

// need ffmpeg with libfdkaac installed

func init() {
	format.RegisterAll()
}

func main() {
	infile, _ := avutil.Open("../bbb_sunflower_1080p_30fps_normal")

	findcodec := func(stream av.AudioCodecData, i int) (need bool, dec av.AudioDecoder, enc av.AudioEncoder, err error) {
		need = true
		dec, _ = ffmpeg.NewAudioDecoder(stream)
		enc.SetSampleRate(stream.SampleRate())
		enc.SetChannelLayout(av.CH_STEREO)
		enc.SetBitrate(300)
		return
	}

	trans := &transcode.Demuxer{
		Options: transcode.Options{
			FindAudioDecoderEncoder: findcodec,
		},
		Demuxer: infile,
	}

	outfile, _ := avutil.Create("out.mp4")
	avutil.CopyFile(outfile, trans)

	outfile.Close()
	infile.Close()
	trans.Close()
}
