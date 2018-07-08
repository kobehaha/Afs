package main

import (
	"github.com/kobehaha/Afs/datahandler"
	"github.com/kobehaha/Afs/heartbeat"
	"github.com/kobehaha/Afs/locate"
	"github.com/kobehaha/Afs/log"
	"net/http"
	"os"
)

func main() {

	log.Init()

	go heartbeat.NewHeartbeat().StartHeartbeat()
	go locate.NewLocate().StartLocate()

	http.HandleFunc("/objects/", datahandler.ObjectHandler)

	err := http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil)
	if err != nil {
		log.GetLogger().Error("Listen_Addrss: "+os.Getenv("LISTEN_ADDRESS")+"Error And Message is %s", err)
		os.Exit(1)
	}
	log.GetLogger().Info("Listen_Address: " + os.Getenv("LISTEN_ADDRESS"))
}
