package dns

import (
	"net/http"
	"time"

	"beryju.io/gravity/internal/resources"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var blockyListListening = false

const blockyListAddr = "127.0.0.1:8100"
const blockyListBase = "http://" + blockyListAddr + "/blocky/"

func (bfwd *BlockyForwarder) startBlockyListServer() {
	if blockyListListening {
		return
	}
	s := mux.NewRouter()
	s.Methods("GET").Handler(http.FileServer(http.FS(resources.BlockyLists)))
	blockyListListening = true
	go func() {
		bfwd.log.Info("starting blocky list server")
		err := http.ListenAndServe(blockyListAddr, s)
		if err != nil {
			bfwd.log.Warn("failed to start blocky list server", zap.Error(err))
			time.Sleep(5 * time.Millisecond)
			bfwd.startBlockyListServer()
		}
	}()
}
