package cassandra

import (
	"github.com/georgemac/clio/lib/storage"
	"github.com/gocql/gocql"
)

// Storage wraps a connection to cassandra and implements
// the storage.Storage interface. It can be used to store
// implementations of Line within Cassandra
type Storage struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

// Option is a function which takes a pointer to a storage type
// It is used to configure the cassandra.Storage type on construction
type Option func(s *Storage)

// Credentials sets a usename and password on the connection to Cassandra
func Credentials(username, password string) Option {
	return func(s *Storage) {
		s.cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}
}

// Keyspace is a storage option for overriding the default keyspace
func Keyspace(keyspace string) Option {
	return func(s *Storage) {
		s.cluster.Keyspace = keyspace
	}
}

// New constructs and configures a new cassandra Storage,
// ready for storing log lines in
func New(hosts []string, opts ...Option) (store *Storage, err error) {
	store = &Storage{
		cluster: gocql.NewCluster(hosts...),
	}

	store.cluster.Keyspace = "clio"
	store.cluster.Consistency = gocql.All

	for _, opt := range opts {
		opt(store)
	}

	store.session, err = store.cluster.CreateSession()
	return
}

func (s *Storage) Close() {
	s.session.Close()
}

// Write inserts an implementation of Line in to cassandra
func (s *Storage) Write(l storage.Line) error {
	queryStmt := `INSERT INTO clio.logs (build_id, created_at, entry_id, container_id, container_name, payload) VALUES (?,?,NOW(),?,?,?);`
	return s.session.Bind(queryStmt, func(q *gocql.QueryInfo) ([]interface{}, error) {
		tag := l.Tag()
		return []interface{}{
			gocql.UUID(tag.BuildID),
			l.CreatedAt(),
			tag.ContainerID,
			tag.ContainerName,
			l.Payload(),
		}, nil
	}).Exec()
}
