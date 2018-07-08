package objects

import (
	"fmt"
	"github.com/kobehaha/Afs/heartbeat"
	"github.com/kobehaha/Afs/locate"
	"github.com/kobehaha/Afs/log"
	"github.com/kobehaha/Afs/objectstreaming"
	"github.com/kobehaha/Afs/utils"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var objectHandler *ObjectHandler

type ObjectHandler struct {
	heartbeat *heartbeat.Heartbeat
	locate    *locate.Locate
}

func NewObjectHandler() *ObjectHandler {

	heartbeat := heartbeat.NewHeartbeat()
	locate := locate.NewLocate()

	go heartbeat.ListenHeartbeat()

	return &ObjectHandler{heartbeat, locate}

}

func GetObjectHandler() *ObjectHandler {

	if objectHandler == nil {

		objectHandler = NewObjectHandler()

		return objectHandler
	}
	return objectHandler
}

func (o *ObjectHandler) Get(w http.ResponseWriter, r *http.Request) {

	object := strings.Split(r.URL.EscapedPath(), "/")[2]

	versionId := r.URL.Query()["version"]

	version := 0

	var e error

	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.GetLogger().Error(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	meta, e := utils.NewEs().GetMetadata(object, version)

	if e != nil {
		log.GetLogger().Error(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	object_ := url.PathEscape(meta.Hash)

	stream, e := o.getStreaming(object_)

	if e != nil {

		log.GetLogger().Error("Get Object error And message %s", e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)

}

func (o *ObjectHandler) Put(w http.ResponseWriter, r *http.Request) {

	hash := utils.GetHashFromHeader(r.Header)

	if hash == "" {
		log.GetLogger().Error("msssing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, e := o.storeObject(r.Body, object)
	if e != nil {
		log.GetLogger().Error("Store Object error And message is %s", e)
	}

	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]

	size := utils.GetSizeFromHeader(r.Header)

	e = utils.NewEs().AddVersion(name, hash, size)

	if e != nil {
		log.GetLogger().Error(e)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (o *ObjectHandler) Del(w http.ResponseWriter, r *http.Request) {

	es := utils.NewEs()

	object := strings.Split(r.URL.EscapedPath(), "/")[2]

	version, e := es.SearchLatestVersion(object)

	if e != nil {
		log.GetLogger().Error(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	e = es.PutMetadata(object, version.Version+1, 0, "")

	if e != nil {
		log.GetLogger().Error(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (o *ObjectHandler) storeObject(r io.Reader, object string) (int, error) {

	stream, e := o.putStreaming(object)

	if e != nil {
		log.GetLogger().Error(e)
		return http.StatusServiceUnavailable, e
	}

	io.Copy(stream, r)

	e = stream.Close()

	if e != nil {
		log.GetLogger().Error(e)
		return http.StatusInternalServerError, e

	}
	return http.StatusOK, nil

}

func (o *ObjectHandler) putStreaming(object string) (*objectstreaming.PutStream, error) {

	servers := o.heartbeat.ChooseRandomDataServers()

	if servers == "" {
		return nil, fmt.Errorf("can not find any dataServers")
	}

	return objectstreaming.NewPutStream(servers, object), nil

}

func (o *ObjectHandler) getStreaming(object string) (io.Reader, error) {

	server := o.locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("can not locate this object %s", object)
	}

	return objectstreaming.NewGetStream(server, object)

}
