package dns

import (
	"net/http"
	"time"

	"beryju.io/gravity/internal/resources"
	"github.com/gorilla/mux"
)

var blockyListListening = false

const blockyListAddr = "127.0.0.1:8100"
const blockyListBase = "http://" + blockyListAddr + "/"

func startBlockyListServer() {
	if blockyListListening {
		return
	}
	s := mux.NewRouter()
	s.Methods("GET").Handler(http.FileServer(http.FS(resources.BlockyLists)))
	blockyListListening = true
	go func() {
		err := http.ListenAndServe(blockyListAddr, s)
		if err != nil {
			time.Sleep(5 * time.Millisecond)
			startBlockyListServer()
		}
	}()
}
