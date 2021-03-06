package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

var cl *torrent.Client

type CReader struct {
	torrent.Reader
}

func (c CReader) Read(p []byte) (n int, err error) {
	var p2 []byte
	n, err = c.Reader.Read(p2)

	copy(p2, p)

	// cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-listen", "1", "-f", "matroska", "-c:v", "libx264", "-b", "300k", "-preset", "fast", "-tune", "zerolatency", "pipe:1")
	// stdin, err := cmd.StdinPipe()
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// if err != nil {
	// 	log.Panic(err)
	// }

	// go func() {
	// 	defer stdin.Close()
	// 	stdin.Write(p)
	// }()

	// err = cmd.Start()
	// if err != nil {
	// 	log.Println("Error starting")
	// }

	// cmd.Wait()

	// log.Println(out.Len())

	// p = out.Bytes()

	return
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func torHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("--> %s\n\n", formatRequest(r))

	// add the magnet (in a round about way so we can log if it was already seen)
	uri := "magnet:?xt=urn:btih:88594AAACBDE40EF3E2510C47374EC0AA396C08E&dn=bbb_sunflower_1080p_30fps_normal.mp4&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80%2fannounce&tr=udp%3a%2f%2ftracker.publicbt.com%3a80%2fannounce&ws=http%3a%2f%2fdistribution.bbb3d.renderfarming.net%2fvideo%2fmp4%2fbbb_sunflower_1080p_30fps_normal.mp4"
	spec, err := torrent.TorrentSpecFromMagnetURI(uri)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "couldnt add the magnet")
		log.Printf("ERROR: %v\n", err)
		return
	}
	t, isNew, err := cl.AddTorrentSpec(spec)
	log.Println("Adding torrent, new?: %v\n", isNew)

	// wait for info
	<-t.GotInfo()
	name := t.Name()
	numPieces := t.NumPieces()
	// Cancel any pieces that were previously marked for download
	t.CancelPieces(0, numPieces)
	// mark the whole thing for download but prio the treader?
	creader := CReader{t.NewReader()}
	defer creader.Reader.Close()

	// add monitor logging if this is the first time adding the torrent
	// if isNew {
	// 	go func() {
	// 		for {
	// 			missing := t.BytesMissing()
	// 			peers := t.PeerConns()
	// 			log.Println("=~=~=~=~=~=~=~=~~")
	// 			log.Printf("Peers: %+v\n", len(peers))
	// 			log.Printf("Missing: %+v\n", missing)
	// 			log.Printf("NumPieces: %+v\n", numPieces)

	// 			// print our progresss
	// 			str := ""
	// 			for i := 0; i < numPieces; i += 2 {
	// 				pbmissing := t.PieceBytesMissing(i)
	// 				if pbmissing == 0 {
	// 					str += "="
	// 					continue
	// 				}
	// 				str += "_"
	// 			}
	// 			log.Println(str)

	// 			log.Println("=~=~=~=~=~=~=~=~~")
	// 			time.Sleep(1 * time.Second)
	// 		}
	// 	}()
	// }
	http.ServeContent(w, r, name, time.Time{}, creader)

	// Start downloading the entire file once the client leaves
	defer func() {
		log.Println("Download all...")
		// t.DownloadAll()
	}()
}

type FRSeeker struct {
	*os.File
}

func (f FRSeeker) Read(p []byte) (n int, err error) {
	log.Println("Intercepted read!", len(p))
	n, err = f.File.Read(p)
	return
}

func (f FRSeeker) Seek(offset int64, whence int) (n int64, err error) {
	log.Printf("Intercepted seek, offset: %v, whence: %v\n", offset, whence)
	n, err = f.File.Seek(offset, whence)
	return
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "./lowquality.mp4"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("%s not found\n", filePath)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<html><body style='font-size:100px'>four-oh-four</body></html>")
		return
	}
	defer file.Close()
	fmt.Printf("serve %s\n", filePath)
	_, filename := path.Split(filePath)

	frseeker := FRSeeker{file}
	http.ServeContent(w, r, filename, time.Time{}, frseeker)
}

func main() {
	// client, err := torrent.NewClient(nil)
	// if err != nil {
	// 	log.Printf("ERROR: %v\n", err)
	// }
	// defer client.Close()
	// cl = client

	http.HandleFunc("/", fileHandler)

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
