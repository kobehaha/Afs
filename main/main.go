package main

import (
	"github.com/kobehaha/Afs/log"
	"net/http"
	"os"
	"github.com/kobehaha/Afs/fronthandler"
	"github.com/kobehaha/Afs/backendhandler"
)

func main() {

	log.Init()

	http.HandleFunc("/objects/", fronthandler.Handler)
	http.HandleFunc("/locate/", backendhandler.Handler)
	err := http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil)
	if err != nil {
		log.GetLogger().Error("Listen_Addrss: "+os.Getenv("LISTEN_ADDRESS")+"Error And Message is %s", err)
		os.Exit(1)
	}
	log.GetLogger().Info("Listen_Address: " + os.Getenv("LISTEN_ADDRESS"))
}
