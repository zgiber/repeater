package main

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

type store struct {
	sync.RWMutex
	data *bytes.Buffer
}

type service struct {
	store store
}

func (srv *service) handleRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		srv.store.Lock()
		srv.store.data.Reset()
		io.Copy(srv.store.data, r.Body)
		srv.store.Unlock()
	case "GET":
		srv.store.Lock()
		w.Write(srv.store.data.Bytes())
		srv.store.Unlock()
	}
}

func main() {
	srv := newService()
	http.HandleFunc("/", srv.handleRequests)
	http.ListenAndServe(":8000", nil)
}

func newService() *service {
	return &service{
		store: store{data: bytes.NewBuffer([]byte{})},
	}
}
