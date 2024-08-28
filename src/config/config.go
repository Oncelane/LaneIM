package config

type Logic struct {
	Addr string
	Name string
}

type Comet struct {
	Addr string
	Name string
}

type Job struct {
	Addr string
}

type KafkaProducer struct {
	Addr []string
}

type KafkaComsumer struct {
	Addr string
}
