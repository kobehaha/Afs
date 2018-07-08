package heartbeat

import (
	"github.com/kobehaha/Afs/utils"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
    "github.com/kobehaha/Afs/log"
)

type Heartbeat struct {

	heartbeat   string
	dataServers map[string]time.Time
	mutex       sync.Mutex
	rabbitmq    *utils.RabbitMq
}

func NewHeartbeat() *Heartbeat {

	rabbitmq := utils.NewRabbitMq(os.Getenv("RABBITMQ_SERVER"))
	dataServer := make(map[string]time.Time)
	metex := sync.Mutex{}

	heartbeat := &Heartbeat{
		heartbeat:   "heartbeat",
		dataServers: dataServer,
		mutex:       metex,
		rabbitmq:    rabbitmq,
	}

	return heartbeat

}

func (heartbeat *Heartbeat) StartHeartbeat() {

	q := heartbeat.rabbitmq

	defer q.Close()

	for {

		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(2 * time.Second)

	}

}

func (heartbeat *Heartbeat) ListenHeartbeat() {

	q := heartbeat.rabbitmq

	defer q.Close()

	q.Bind("apiServers")

	c := q.Consume()

	go heartbeat.removeExpiredDataServer()

	for msg := range c {

		dataServer, e := strconv.Unquote(string(msg.Body))


		if e != nil {
			panic(e)
		}

		heartbeat.mutex.Lock()

		heartbeat.dataServers[dataServer] = time.Now()

		heartbeat.mutex.Unlock()

	}

}

func (heartbeat *Heartbeat) removeExpiredDataServer() {

	for {

		time.Sleep(5 * time.Second)

		heartbeat.mutex.Lock()

		for s, t := range heartbeat.dataServers {

			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(heartbeat.dataServers, s)
			}
		}

		heartbeat.mutex.Unlock()
	}

}

func (heartbeat *Heartbeat) GetDataServers() []string {

	heartbeat.mutex.Lock()

	defer heartbeat.mutex.Unlock()

	ds := make([]string, 0)

	log.GetLogger().Info("Get data servers %s", heartbeat.dataServers)
	for s, _ := range heartbeat.dataServers {
		ds = append(ds, s)
    }

	return ds

}

func (hearbeat *Heartbeat) ChooseRandomDataServers() string {

	ds := hearbeat.GetDataServers()

	n := len(ds)

	if n == 0 {
		return ""
	}

	back := ds[rand.Intn(n)]

	return back
}
