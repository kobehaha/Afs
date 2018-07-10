package objectstreaming

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type PutStream struct {
	writer *io.PipeWriter

	c chan error
}

type GetStream struct {
	reader io.Reader
}

type TempPutStream struct {
	Server string
	Uuid   string
}

func NewTempStream(server, hash string, size int64) (*TempPutStream, error) {

	request, e := http.NewRequest("POST", "http://"+server+"/temp/"+hash, nil)

	if e != nil {
		return nil, e
	}

	request.Header.Set("size", fmt.Sprintf("%d", size))

	client := http.Client{}

	response, e := client.Do(request)

	if e != nil {
		return nil, e
	}

	uuid, e := ioutil.ReadAll(response.Body)

	if e != nil {
		return nil, e
	}

	return &TempPutStream{server, string(uuid)}, nil
}

func NewPutStream(server, object string) *PutStream {

	reader, writer := io.Pipe()

	c := make(chan error)

	go func() {
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)

		client := http.Client{}
		r, e := client.Do(request)

		if e == nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("data server return http code %s", r.StatusCode)
		}

		c <- e

	}()

	return &PutStream{writer, c}
}

func NewGetStream(server, object string) (*GetStream, error) {

	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}

	url := "http://" + server + "/objects/" + object
	return newGetStream(url)

}

func newGetStream(url string) (*GetStream, error) {
	r, e := http.Get(url)

	if e != nil {
		return nil, e
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataserver return http code %s", r.StatusCode)
	}

	return &GetStream{r.Body}, nil
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (t *TempPutStream) Write(p []byte) (n int, err error) {

	request, e := http.NewRequest("PATCH", "http://"+t.Server+"/temp/"+t.Uuid, strings.NewReader(string(p)))

	if e != nil {
		return 0, e
	}

	client := http.Client{}
	r, e := client.Do(request)

	if e != nil {
		return 0, e
	}

	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("data server return http code: %d", r.StatusCode)
	}

	return len(p), nil

}

func (t *TempPutStream) Commit(good bool) {

	method := "DELETE"

	if good {
		method = "PUT"
	}

	request, _ := http.NewRequest(method, "http://"+t.Server+"/temp/"+t.Uuid, nil)

	client := http.Client{}

	client.Do(request)

}
