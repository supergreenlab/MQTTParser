package main

import (
	"context"
	"flag"
	"log"

	"github.com/olivere/elastic"
)

var (
	elastic_server = flag.String("elastic_server", "http://elastic:9200", "The full url of the elastic server to connect to ex: tcp://127.0.0.1:1883")
	ctx            context.Context
	client         *elastic.Client
	c              chan interface{}
)

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"log":{
			"properties":{
				"id":{
					"type":"keyword"
				},
				"channel":{
					"type":"keyword"
				},
				"topic":{
					"type":"keyword"
				},
				"payload":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"level":{
					"type":"keyword"
				},
				"timestamp":{
					"type":"date"
				},
				"tag":{
					"type":"keyword"
				},
				"module":{
					"type":"keyword"
				},
				"msg":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"kvs": {
					"type":"nested"
				},
				"kvi": {
					"type":"nested"
				}
			}
		}
	}
}`

func indexLog(l interface{}) {
	c <- l
}

func start_elastic() {
	var (
		err error
	)
	ctx = context.Background()

	for {
		client, err = elastic.NewClient(elastic.SetURL(*elastic_server))
		if err != nil {
			log.Print(err)
			continue
		}
		break
	}

	exists, err := client.IndexExists("log").Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		log.Print("Not EXIST !")
		createIndex, err := client.CreateIndex("log").BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			log.Print("Not aknowledged..")
		}
	}

	for l := range c {
		_, err := client.Index().
			Index("supergreenlab").
			Type("log").
			BodyJson(l).
			Do(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func init_elastic() {
	c = make(chan interface{}, 50)
	go start_elastic()
}
