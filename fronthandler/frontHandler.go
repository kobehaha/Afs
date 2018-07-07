package fronthandler

import (
    "net/http"
    "github.com/kobehaha/Afs/objects"
)



func Handler(w http.ResponseWriter, r *http.Request) {

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
