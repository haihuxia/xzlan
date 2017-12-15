package es

import (
	"testing"
	"gopkg.in/olivere/elastic.v5"
	"context"
	"fmt"
	"encoding/json"
)

func TestNewEs(t *testing.T) {
	client, err := elastic.NewClient(elastic.SetURL("http://10.0.21.44:9200"))
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
	query = query.Filter(elastic.NewRangeQuery("elapsed").Gte(5))
	query = query.Filter(elastic.NewRangeQuery("@timestamp").Gte("now-4m").Lt("now"))

	src, err := query.Source()
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	fmt.Printf("data: %s \n", data)

	result, err := client.Search("logstash-2017.12.14").Query(query).Do(context.Background())
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	fmt.Printf("result: %d \n", result.Hits.TotalHits)
}
