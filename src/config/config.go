package config

type Logic struct {
	Addr          string
	Name          string
	KafkaProducer KafkaProducer
	Etcd          Etcd
}

type Comet struct {
	Addr string
	Name string
	Etcd Etcd

	BucketSize    int
	BucketRoutine int
}

type Job struct {
	Addr          string
	Name          string
	KafkaComsumer KafkaComsumer
	Etcd          Etcd

	BucketSize       int
	CometRoutineSize int
}

type Etcd struct {
	Addr []string
}

type KafkaProducer struct {
	Addr []string
}

type KafkaComsumer struct {
	Addr    []string
	Topics  []string
	GroupId string
}
