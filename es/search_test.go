package es

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"testing"
	"xzlan/alert"
)

func TestNewEs(t *testing.T) {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	exists, err := client.IndexExists("logstash-2017.12.14").Do(context.Background())
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	fmt.Printf("exists %t \n", exists)
	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewMatchPhraseQuery("interface", "user"))
	query = query.Must(elastic.NewMatchPhraseQuery("method", "get"))
	query = query.Filter(elastic.NewRangeQuery("elapsed").Gte(8))
	query = query.Filter(elastic.NewRangeQuery("@timestamp").Gte("now-2h").Lt("now"))

	src, err := query.Source()
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	fmt.Printf("data: %s \n", data)

	result, err := client.Search("logstash-2017.12.24").Query(query).Do(context.Background())
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	fmt.Printf("result: %d \n", result.Hits.TotalHits)
	for i := 0; i < len(result.Hits.Hits); i++ {
		//fmt.Printf("result: %s \n", result.Hits.Hits[i])
		b, _ := result.Hits.Hits[i].Source.MarshalJSON()
		var m alert.Message
		json.Unmarshal(b, &m)
		fmt.Printf("%s \n", m.Message)
	}
}
