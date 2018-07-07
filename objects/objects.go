package objects

import (
	"fmt"
	"github.com/kobehaha/Afs/heartbeat"
	"github.com/kobehaha/Afs/log"
	"github.com/kobehaha/Afs/objectstreaming"
	"io"
	"net/http"
	"strings"
	"github.com/kobehaha/Afs/locate"
)

var objectHandler *ObjectHandler

type ObjectHandler struct {

	heartbeat *heartbeat.Heartbeat
	locate *locate.Locate
}

func NewObjectHandler() *ObjectHandler{

    heartbeat := heartbeat.NewHeartbeat()
    locate := locate.NewLocate()

    go heartbeat.StartHeartbeat()
    go heartbeat.ListenHeartbeat()

    return &ObjectHandler{heartbeat, locate}


}


func GetObjectHandler() *ObjectHandler{

    if objectHandler == nil {

        objectHandler = NewObjectHandler()

        return objectHandler
    }
    return objectHandler
}

func (o *ObjectHandler) Get(w http.ResponseWriter, r *http.Request) {

	object := strings.Split(r.URL.EscapedPath(), "/")[2]

	stream , e := o.getStreaming(object)

	if e != nil {

		log.GetLogger().Error("Get Object error And message %s", e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w,stream)

}

func (o *ObjectHandler) Put(w http.ResponseWriter, r *http.Request) {

	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, e := o.storeObject(r.Body, object)
	if e != nil {
		log.GetLogger().Error("Store Object error And message is %s", e)
	}

	w.WriteHeader(c)

}

func (o *ObjectHandler) storeObject(r io.Reader, object string) (int, error) {

	stream, e := o.putStreaming(object)

	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	io.Copy(stream, r)

	e = stream.Close()

	if e != nil {
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

func (o *ObjectHandler) getStreaming(object string) (io.Reader, error){


	server := o.locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("can not locate this object %s", object)
	}

	return objectstreaming.NewGetStream(server, object)

}


