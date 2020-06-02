package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

func main() {
	fi, _ := ioutil.ReadFile("../bbb_sunflower_1080p_30fps_normal.mp4")

	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-listen", "1", "-f", "matroska", "-c:v", "libx264", "-b", "300k", "-preset", "fast", "-tune", "zerolatency", "pipe:1")
	stdin, err := cmd.StdinPipe()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err != nil {
		log.Panic(err)
	}

	go func() {
		defer stdin.Close()
		stdin.Write(fi)
	}()

	go func() {
		for {
			log.Println(out.Len())
			time.Sleep(1 * time.Second)
		}
	}()

	// io.Copy(stdin, bytes.NewReader(fi))
	// err = cmd.Start()
	// if err != nil {
	// 	log.Panic(err)
	// }

	err = cmd.Start()
	if err != nil {
		log.Println("Error starting")
	}

	cmd.Wait()

}
