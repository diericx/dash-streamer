package main

type crseeker struct {

}

type creader struct {
	torrent.Reader
	*os.File
}

func (r creader) struct {
	
}