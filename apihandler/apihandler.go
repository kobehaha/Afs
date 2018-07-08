package apihandler

import (
	"encoding/json"
	"github.com/kobehaha/Afs/locate"
	"github.com/kobehaha/Afs/objects"
	"net/http"
	"strings"
	"github.com/kobehaha/Afs/utils"
	"github.com/kobehaha/Afs/log"
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
	if m == http.MethodDelete {
		objects.GetObjectHandler().Del(w,r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)

}

func VersionHandler(w http.ResponseWriter, r *http.Request) {

	m := r.Method

	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	from := 0
	size := 1000

	name := strings.Split(r.URL.EscapedPath(), "/")[2]

	for {

		metas, e := utils.NewEs().SerachAllVersions(name, from, size)
		if e != nil {
		    log.GetLogger().Error(e)
		    w.WriteHeader(http.StatusInternalServerError)
		    return
		}

		for i := range metas {
			b,_ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}

		if len(metas) != size {
			return
		}

		from += size
	}

}
