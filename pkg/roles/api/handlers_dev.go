package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type getBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (ro *APIRole) handleAPIGet(rw http.ResponseWriter, r *http.Request) {
	k, err := ro.i.GetKV().KV.Get(r.Context(), ro.i.GetKV().Key(), clientv3.WithPrefix())
	if err != nil {
		ro.log.WithError(err).Warning("failed to get keys")
		return
	}
	b := make([]getBody, len(k.Kvs))
	for idx, kvs := range k.Kvs {
		b[idx] = getBody{
			Key:   string(kvs.Key),
			Value: string(kvs.Value),
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(b)
}

func (ro *APIRole) handleAPIPost(rw http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ro.log.WithError(err).Warning("failed to read body")
		return
	}
	_, err = ro.i.GetKV().Put(r.Context(), ro.i.GetKV().Key(r.URL.Query().Get("key")), string(b))
	if err != nil {
		ro.log.WithError(err).Warning("failed to put")
	}
}

func (ro *APIRole) handleAPIDel(rw http.ResponseWriter, r *http.Request) {
	_, err := ro.i.GetKV().Delete(r.Context(), ro.i.GetKV().Key(r.URL.Query().Get("key")))
	if err != nil {
		ro.log.WithError(err).Warning("failed to delete")
	}
}
