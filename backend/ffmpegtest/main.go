package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
)

func main() {
	fi, _ := ioutil.ReadFile("bbb_sunflower_1080p_30fps_normal.mp4")

	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-ab", "300k", "-f", "mp4", "output.mp4")
	stdin, err := cmd.StdinPipe()
	defer stdin.Close()
	if err != nil {
		log.Panic(err)
	}
	io.Copy(stdin, bytes.NewReader(fi))
	err = cmd.Start()
	if err != nil {
		log.Panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Panic(err)
	}
}
