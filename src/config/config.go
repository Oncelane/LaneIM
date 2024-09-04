package config

import (
	"laneIM/src/pkg/util"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Logic struct {
	Addr          string
	Name          string
	Type          string
	KafkaProducer KafkaProducer
	Etcd          Etcd
}

func (c *Logic) Default() {
	ip, err := util.GetOutBoundIP()
	if err != nil {
		log.Panicln("faild to get outbound ip:", err)
	}
	*c = Logic{
		Addr: ip + ":50060",
		Name: "0",
		Type: "logic",
		KafkaProducer: KafkaProducer{
			Addr: []string{ip + ":9092"},
		},
		Etcd: Etcd{
			Addr: []string{ip + ":2379"},
		},
	}
}

func (c *Logic) GetName() string {
	return c.Name
}

func (c *Logic) GetType() string {
	return c.Type
}

func (c *Logic) ReadLocal() error {
	log.Println("read from ", c.Type+".yml")
	data, err := os.ReadFile(c.Type + ".yml")
	if err != nil {
		log.Panicln("config.yaml does not exist")
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		log.Panicln("can't not read config.yml")
	}
	return err
}

type Comet struct {
	Addr string
	Name string
	Type string
	Etcd Etcd

	BucketSize int
}

func (c *Comet) GetName() string {
	return c.Name
}
func (c *Comet) GetType() string {
	return c.Type
}
func (c *Comet) Default() {
	ip, err := util.GetOutBoundIP()
	if err != nil {
		log.Panicln("faild to get outbound ip:", err)
	}
	*c = Comet{
		Addr: ip + ":50050",
		Name: "0",
		Type: "comet",
		Etcd: Etcd{
			Addr: []string{ip + ":2379"},
		},
		BucketSize: 32,
	}
}
func (c *Comet) ReadLocal() error {
	log.Println("read from ", c.Type+".yml")
	data, err := os.ReadFile(c.Type + ".yml")
	if err != nil {
		log.Panicln("config.yaml does not exist")
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		log.Panicln("can't not read config.yml")
	}
	return err
}

type Job struct {
	Addr          string
	Name          string
	Type          string
	KafkaComsumer KafkaComsumer
	Etcd          Etcd

	BucketSize       int
	CometRoutineSize int
}

func (c *Job) GetName() string {
	return c.Name
}
func (c *Job) GetType() string {
	return c.Type
}
func (c *Job) Default() {
	ip, err := util.GetOutBoundIP()
	if err != nil {
		log.Panicln("faild to get outbound ip:", err)
	}
	*c = Job{
		Addr: ip + ":50070",
		Name: "0",
		Type: "job",
		KafkaComsumer: KafkaComsumer{
			Addr:    []string{ip + ":9092"},
			Topics:  []string{"laneIM"},
			GroupId: "job",
		},
		Etcd: Etcd{
			Addr: []string{ip + ":2379"},
		},
		BucketSize:       32,
		CometRoutineSize: 32,
	}
}
func (c *Job) ReadLocal() error {
	log.Println("read from ", c.Type+".yml")
	data, err := os.ReadFile(c.Type + ".yml")
	if err != nil {
		log.Panicln("config.yaml does not exist")
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		log.Panicln("can't not read config.yml")
	}
	return err
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

type LaneConfig interface {
	Default()
	GetName() string
	GetType() string
	ReadLocal() error
}

func WriteLocal(conf LaneConfig) error {
	out, err := yaml.Marshal(conf)
	if err != nil {
		log.Panicln("failed to marshal config", conf.GetType(), ":", err)
		return err
	}

	err = os.WriteFile(conf.GetType()+".yml", out, 0644)
	if err != nil {
		log.Panicln("failed to write ", conf.GetType(), err)
		return err
	}
	return nil
}

func WriteRemote(conf LaneConfig) error {

	return nil
}

// TODO etcd dynamic config
func ReadRemote(conf LaneConfig) {

}

func Init(conf LaneConfig) {
	_, err := os.Stat(conf.GetType() + ".yml")
	if err != nil {
		if os.IsNotExist(err) {
			conf.Default()
			WriteLocal(conf)
			log.Println("please check for the config.yml if needed to be modified, then run again")
			os.Exit(1)
		} else {
			log.Panicln("wrong err:", err)
		}
	}
	conf.ReadLocal()
}
