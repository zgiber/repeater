package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"sync"
)

type payload struct {
	NumParallelConsumers int       `json:"number_of_parallel_consumers"`
	ConsumerBufferSize   int       `json:"parallel_buffer_size"`
	ConsumerBufferUsage  []float64 `json:"parallel_buffers_usage_perc"`
}

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

func (srv *service) handleGetGUI(w http.ResponseWriter, r *http.Request) {

	data := &payload{}
	err := json.Unmarshal(srv.store.data.Bytes(), data) // if unmarshal fails we just display empty page
	if err != nil {
		log.Print(err)
	}

	viewTemplate := `<!DOCTYPE html><html>
	<head><meta http-equiv="refresh" content="2"></head>
	<body>
	Load on workers <br>
	{{range .ConsumerBufferUsage}}<div>{{ . }}%</div>{{end}}</body></html>`

	t, err := template.New("view").Parse(viewTemplate)
	if err != nil {
		log.Print("1", err)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Print("2", err)
	}
}

func main() {
	srv := newService()
	http.HandleFunc("/", srv.handleRequests)
	http.HandleFunc("/gui", srv.handleGetGUI)
	http.ListenAndServe(":8008", nil)
}

func newService() *service {
	return &service{
		store: store{data: bytes.NewBuffer([]byte{})},
	}
}
