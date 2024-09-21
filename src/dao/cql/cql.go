package cql

import (
	"fmt"
	"laneIM/src/config"
	"laneIM/src/pkg/laneLog"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v3"
	"golang.org/x/sync/singleflight"
)

type ScyllaDB struct {
	conf              config.ScyllaDB
	DB                *gocqlx.Session
	chatMessageTable  *table.Table
	SingleFlightGroup singleflight.Group
}

func NewCqlDB(conf config.ScyllaDB) *ScyllaDB {
	cluster := gocql.NewCluster(conf.Addrs...)
	// cluster.Keyspace = conf.Keyspace
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		laneLog.Logger.Fatalln("faild to conneect scyllaDB:", err)
	}
	rt := &ScyllaDB{
		conf: conf,
		DB:   &session,
	}

	rt.initKeyspace()
	return rt
}

func (s *ScyllaDB) initKeyspace() {
	err := s.DB.ExecStmt(fmt.Sprintf(
		`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`,
		s.conf.Keyspace,
	))
	if err != nil {
		laneLog.Logger.Fatalln("create keyspace:", err)
	}
}

func (s *ScyllaDB) ResetKeySpace() {
	s.DB.ExecStmt(fmt.Sprintf(`DROP KEYSPACE %s`, s.conf.Keyspace))
	err := s.DB.ExecStmt(fmt.Sprintf(
		`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`,
		s.conf.Keyspace,
	))
	if err != nil {
		laneLog.Logger.Fatalln("create keyspace:", err)
	}
}
