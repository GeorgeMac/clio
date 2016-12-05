package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/codeship/clio/lib/storage"
	"github.com/codeship/clio/lib/storage/cassandra"
	uuid "github.com/satori/go.uuid"
)

var (
	host     = "127.0.0.1"
	keyspace = "clio"
	count    = 10
)

func hosts(s ...string) []string { return s }

func main() {
	flag.IntVar(&count, "count", 10, "Number of lines to generate")
	flag.Parse()

	store, err := cassandra.New(hosts(host), cassandra.Credentials("cassandra", "cassandra"))
	if err != nil {
		log.Println(err)
	}

	defer store.Close()

	start := time.Now()

	defer func() {
		log.Printf("ellapsed: %v\n", time.Now().Sub(start))
	}()

	build := uuid.NewV4()
	log.Println("build: ", build)

	stamp := time.Now().Unix()
	for i := int64(0); i < int64(count); i++ {
		payload := fmt.Sprintf("skdfjhd fskdjfh skjdf skjdhf sjkdhf %d", i)
		line := storage.SysLogLineV0{
			TagV0: storage.TagV0{
				BuildID:       build,
				ContainerName: "jet-container",
				ContainerID:   "5dbd760c5b63",
			},
			Time: time.Now(),
			Data: []byte(payload),
		}

		stamp += 1000
		if err := store.Put(&line, i); err != nil {
			log.Println(err)
		}
	}
}
