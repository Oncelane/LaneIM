package pkg

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v3"
)

var keyspace string = "examples"

// metadata specifies table name and columns it must be in sync with schema.
var personMetadata = table.Metadata{
	Name:    "persons",
	Columns: []string{"id", "first_name", "last_name", "email"},
	PartKey: []string{"id"},
	SortKey: []string{"last_name"},
}

// personTable allows for simple CRUD operations based on personMetadata.
var personTable = table.New(personMetadata)

// Person represents a row in person table.
// Field names are converted to snake case by default, no need to add special tags.
// A field will not be persisted by adding the `db:"-"` tag or making it unexported.
type Person struct {
	Id        gocql.UUID
	FirstName string
	LastName  string
	Email     []string
	HairColor string `db:"-"` // exported and skipped
	eyeColor  string // unexported also skipped
}

func TestDB(t *testing.T) {
	// Create gocql cluster.
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = keyspace
	// Wrap session on creation, gocqlx session embeds gocql.Session pointer.
	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		t.Fatal(err)
	}

	session.ExecStmt(`DROP KEYSPACE examples`)
	err = session.ExecStmt(fmt.Sprintf(
		`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`,
		keyspace,
	))
	if err != nil {
		t.Fatal("create keyspace:", err)
	}

	err = session.ExecStmt(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.persons (
		id uuid PRIMARY KEY,
		first_name text,
		last_name text,
		email set<text>)`, keyspace))
	if err != nil {
		t.Fatal("create table:", err)
	}

	p := Person{
		mustParseUUID("756716f7-2e54-4715-9f00-91dcbea6cf50"),
		"Michał",
		"Matczuk",
		[]string{"michal@scylladb.com"},
		"red",   // not persisted
		"hazel", // not persisted
	}

	q := session.Query(personTable.Insert()).BindStruct(p)
	if err := q.ExecRelease(); err != nil {
		t.Fatal(err)
	}

	var people []Person
	q = session.Query(personTable.Select()).BindMap(qb.M{"id": "756716f7-2e54-4715-9f00-91dcbea6cf50"})
	if err := q.SelectRelease(&people); err != nil {
		t.Fatal(err)
	}
	t.Log(people)
}

func mustParseUUID(s string) gocql.UUID {
	u, err := gocql.ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return u
}
