package locate

import (
	"github.com/kobehaha/Afs/utils"
	"os"
	"strconv"
	"time"
)


var locate *Locate

type Locate struct {
	rabbitmq *utils.RabbitMq
}


func NewLocate() *Locate {

	rabbitmq := utils.NewRabbitMq(os.Getenv("RABBITMQ_SERVER"))

	locate := &Locate{
		rabbitmq: rabbitmq,
	}

	return locate

}

func (locate *Locate) Locate(name string) string {

	q := locate.rabbitmq

	q.Publish("dataServers", name)

	c := q.Consume()

	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()

	msg := <-c

	s, _ := strconv.Unquote(string(msg.Body))

	return s

}

func (locate *Locate) Exist(name string) bool {
	return locate.Locate(name) != ""
}


func GetLocate() *Locate{

	if locate == nil {
		locate = NewLocate()
		return locate
	}
	return locate
}

