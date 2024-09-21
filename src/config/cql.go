package config

type ScyllaDB struct {
	ClusterName string
	Addrs       []string
	Keyspace    string
	BatchWriter BatchWriter
}

func (c *ScyllaDB) Default() {
	batch := BatchWriter{}
	batch.Default()
	*c = DefaultScyllaDB()
}

func DefaultScyllaDB() ScyllaDB {
	return ScyllaDB{
		ClusterName: "laneIM",
		Addrs:       []string{"127.0.0.1"},
		BatchWriter: DefaultBatchWriter(),
	}
}
