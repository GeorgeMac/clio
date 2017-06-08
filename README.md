Clio - Syslog Sink
============================

Parses syslog lines (format RFC3164) delivered using a predefined tag format, based on builds and docker container IDs in to a backend database.

Current draft of the syslog tag format is:

`v0/<build-id>/<container-name>/<container-id>`

The idea being that a container run with the correct log options, forwarded to Clio, should have its log lines handled appropriately by Clio.

e.g.

```bash
docker run -it --rm \
    --log-driver syslog \
    --log-opt tag="v0/{{ (.ExtraAttributes nil).BUILD_ID }}/{{ .Name }}/{{ .ID }}" \
    --log-opt env=BUILD_ID \
    --log-opt syslog-format=rfc3164 \
    -e BUILD_ID=aaaaaa-bbbb-cccc-dddd-eeeeeeee \
    --name "my-container"
```

should result in syslog output similar to the following:

```
Apr  1 15:22:17 ip-10-27-39-73 v0/aaaaaa-bbbb-cccc-dddd-eeeeeeee/some-container-id/my-container[12345]: some application log line 
Apr  1 15:22:17 ip-10-27-39-73 v0/aaaaaa-bbbb-cccc-dddd-eeeeeeee/some-container-id/my-container[12345]: 2016-04-01 15:22:17.075416751 +0000 UTC stderr msg: 1
```

## Backends for Clio

This is configurable and defined as a software interface. Initially a Cassandra implementation will be used. However, this should be pluggable.

Example options:

- Cassandra
- BigTable (GCP)
- Kafka
- Kinesis

### important types 

```go
package storage

type Line interface {
	Tag() Tag
	Payload() []byte
	CreatedAt() time.Time
}

// Writer is something which you can Write Line
// types in to
type Writer interface {
	Write(Line) error
}

// Reader is something which returns a read-only
// channel of lines for a given UUID.
type Reader interface {
	Read(uuid.UUID) (<-chan PersistedLine, error)
}
```

### Cassandra Backend

Initially we will have an implementaion of Storage using cassandra as the backend. 
