package utils

import (
	"encoding/json"
	"fmt"
	"github.com/kobehaha/Afs/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Metadata struct {
	Name string

	Version int

	Size int64

	Hash string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int

		Hits []hit
	}
}

type Es struct {
}

func NewEs() *Es{

	return &Es{}
}

func (es *Es) AddVersion(name, hash string, size int64) error {

	version, e := es.SearchLatestVersion(name)

	if e != nil {
		return e
	}

	return es.PutMetadata(name, version.Version+1, size, hash)

}

func (es *Es) GetMetadata(name string, version int) (Metadata, error) {

	if version == 0 {
		return es.SearchLatestVersion(name)
	}
	return getMetadata(name, version)

}

func (es *Es) SerachAllVersions(name string, from, size int) ([]Metadata, error) {

	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size%d", os.Getenv("ES_SERVER"), from, size)

	if name != "" {
		url += "&q=name:" + name
	}

	r, e := http.Get(url)

	if e != nil {
		return nil, e
	}

	metas := make([]Metadata, 0)

	result, _ := ioutil.ReadAll(r.Body)

	var sr searchResult

	json.Unmarshal(result, &sr)

	for i := range sr.Hits.Hits {

		metas = append(metas, sr.Hits.Hits[i].Source)
	}

	return metas, nil
}

func (es *Es) PutMetadata(name string, version int, size int64, hash string) error {

	doc := fmt.Sprintf(`{"name":"%s", "version":%d, "size":%d, "hash":"%s"}`, name, version, size, hash)


	client := http.Client{}

	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create", os.Getenv("ES_SERVER"), name, version)


	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))

	r, e := client.Do(request)

	if e != nil {
		return e
	}

	if r.StatusCode == http.StatusConflict {
		return es.PutMetadata(name, version+1, size, hash)
	}

	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("failed to put metadata: %d %s", r.StatusCode, string(result))
	}

	return nil

}

func (es *Es) SearchLatestVersion(name string) (meta Metadata, e error) {

	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc", os.Getenv("ES_SERVER"), url.PathEscape(name))


	r, e := http.Get(url)


	if e != nil {
		log.GetLogger().Error(e)
		return
	}

	if r.StatusCode != http.StatusOK {
		log.GetLogger().Error("failed to search latest metadata: %d", r.StatusCode)
		return
	}

	result, _ := ioutil.ReadAll(r.Body)

	var sr searchResult
	json.Unmarshal(result, &sr)

	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}

	return meta, nil

}

func getMetadata(name string, versionId int) (meta Metadata, e error) {

	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source", os.Getenv("ES_SERVER"), name, versionId)

	r, e := http.Get(url)

	if e != nil {
		log.GetLogger().Error(e)
		return
	}

	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("failed to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}

	result, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(result, &meta)

	return

}
