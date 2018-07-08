package fronthandler

import (
    "net/http"
    "github.com/kobehaha/Afs/objects"
)



func Handler(w http.ResponseWriter, r *http.Request) {

    httpMethod := r.Method

    if httpMethod == http.MethodPut {
        objects.GetObjectHandler().Put(w, r)
        return
    }
    if httpMethod == http.MethodGet {
        objects.GetObjectHandler().Get(w, r)
        return
    }
    w.WriteHeader(http.StatusMethodNotAllowed)

}
