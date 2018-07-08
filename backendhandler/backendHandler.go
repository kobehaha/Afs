package backendhandler

import (
    "net/http"
    "strings"
    "encoding/json"
    "github.com/kobehaha/Afs/locate"
)

func Handler(w http.ResponseWriter, r *http.Request) {

    httpMethod := r.Method

    if httpMethod != http.MethodGet {
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