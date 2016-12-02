Clio - Codeships Syslog Sink
============================

Parses syslog lines delivered using a predefined tag format, based on codeship builds and docker container IDs in to a backend database.

Current draft of the syslog tag format is:

`docker/v0/<build-id>/<container-name>/<container-id>`

The idea being that a container run with the correct log options, forwarded to Clio, should have its log lines handled appropriately by Clio.

e.g.

```bash
docker run -it --rm \
    --log-driver syslog \
    --log-opt tag="v0/{{ (.ExtraAttributes nil).CODESHIP_BUILD_ID }}/{{ .Name }}/{{ .ID }}" \
    --log-opt env=CODESHIP_BUILD_ID \
    -e CODESHIP_BUILD_ID=aaaaaa-bbbb-cccc-dddd-eeeeeeee \ 
    --name "my-container"
```

should result in syslog output similar to the following:

```
Apr  1 15:22:17 ip-10-27-39-73 docker/v0/aaaaaa-bbbb-cccc-dddd-eeeeeeee/some-container-id/my-container[12345]: some application log line 
Apr  1 15:22:17 ip-10-27-39-73 docker/v0/aaaaaa-bbbb-cccc-dddd-eeeeeeee/some-container-id/my-container[12345]: 2016-04-01 15:22:17.075416751 +0000 UTC stderr msg: 1
```

## Backends for Clio


This could be configurable and defined as a software interface. Initially a Cassandra implementation will be used. However, this should be pluggable.

> We could use managed alternatives like GCP's BigTable

### important types 

```go
type Storage interface {
  PutLine(build uuid.UUID, when time.Time, name, id, version string, []byte payload) error
}
```

### Cassandra Backend

Initially we will have an implementaion of Storage using cassandra as the backend. 
