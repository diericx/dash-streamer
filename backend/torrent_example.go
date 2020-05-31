package main

import (
	"log"
	"time"

	"github.com/anacrolix/torrent"
)

func main() {
	c, _ := torrent.NewClient(nil)
	defer c.Close()
	t, _ := c.AddMagnet("magnet:?xt=urn:btih:ZOCMZQIPFFW7OLLMIC5HUB6BPCSDEOQU")
	<-t.GotInfo()
	info := t.Info()
	log.Printf("Name", info.Name)
	log.Printf("Peers: %v", t.PeerConns())
	log.Printf("Bytes Missing: %v", t.BytesMissing())
	log.Printf("NumPieces: %v", t.NumPieces())

	t.DownloadPieces(0, 1)
	log.Printf("Bytes Missing: %v", t.BytesMissing())

	m := t.PieceBytesMissing(0)
	for m > 0 {
		log.Printf("missing: %v", m)
		m = t.PieceBytesMissing(0)
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("Bytes Missing: %v", t.BytesMissing())

	pstate := t.PieceState(0)
	log.Printf("PieceState: %+v\n", pstate)
	// t.DownloadPieces()
	// t.DownloadAll()
	// c.WaitAll()
	// log.Print("ermahgerd, torrent downloaded")
}
