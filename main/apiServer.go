package main

import (
	"github.com/kobehaha/Afs/apihandler"
	"github.com/kobehaha/Afs/heartbeat"
	"github.com/kobehaha/Afs/log"
	"net/http"
	"os"
)

func main() {

	log.Init()
	go heartbeat.NewHeartbeat().ListenHeartbeat()

	http.HandleFunc("/objects/", apihandler.ObjectHandler)
	http.HandleFunc("/locate/", apihandler.LocateHandler)
	http.HandleFunc("/versions/", apihandler.VersionHandler)

	err := http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil)
	if err != nil {
		log.GetLogger().Error("Listen_Addrss: "+os.Getenv("LISTEN_ADDRESS")+"Error And Message is %s", err)
		os.Exit(1)
	}
	log.GetLogger().Info("Listen_Address: " + os.Getenv("LISTEN_ADDRESS"))

}
