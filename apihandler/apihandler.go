package apihandler

import (
    "net/http"
    "github.com/kobehaha/Afs/locate"
    "strings"
    "encoding/json"
    "github.com/kobehaha/Afs/objects"
)

func LocateHandler(w http.ResponseWriter, r *http.Request) {

    m := r.Method

    if m != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    info := locate.GetLocate().Locate(strings.Split(r.URL.EscapedPath(), "/")[2])

    if len(info) == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    b, _ := json.Marshal(info)

    w.Write(b)
}


func ObjectHandler(w http.ResponseWriter, r *http.Request) {

    m := r.Method
    if m == http.MethodPut {
        objects.GetObjectHandler().Put(w, r)
        return
    }
    if m == http.MethodGet {
        objects.GetObjectHandler().Get(w, r)
        return
    }
    w.WriteHeader(http.StatusMethodNotAllowed)

}