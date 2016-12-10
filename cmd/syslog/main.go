package main

import (
	"log"

	"github.com/codeship/clio/lib/storage"
	"github.com/codeship/clio/lib/storage/cassandra"

	syslog "gopkg.in/mcuadros/go-syslog.v2"
)

var (
	host     = "127.0.0.1"
	keyspace = "clio"
	count    = 10
)

func hosts(s ...string) []string { return s }

func main() {
	store, err := cassandra.New(hosts(host), cassandra.Credentials("cassandra", "cassandra"))
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	log.Println("connected to cassandra")

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)
	server.ListenTCP("0.0.0.0:34567")
	server.Boot()

	log.Println("listening on address 0.0.0.0:34567")

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			var line storage.SysLogLineV0
			if err := line.UnmarshalLogParts(logParts); err != nil {
				panic(err)
			}

			if err := store.Write(&line); err != nil {
				log.Println(err)
			}
		}
	}(channel)

	server.Wait()
}
