package config

import (
	"laneIM/src/pkg/util"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Canal struct {
	CanalAddress     string
	CanalPort        int
	CanalName        string
	CanalPassword    string
	CanalDestination string
	SoTimeOut        int32
	IdleTimeOut      int32
	Subscribe        string
	MsgChSize        int
	KafkaProducer    KafkaProducer
	Redis            Redis
	Mysql            Mysql
}

func (c *Canal) Default() {
	mysqlC := Mysql{}
	mysqlC.Default()
	*c = Canal{
		CanalAddress:     "127.0.0.1",
		CanalPort:        11111,
		CanalName:        "canal",
		CanalPassword:    "canal",
		CanalDestination: "laneIM",
		SoTimeOut:        60000,
		IdleTimeOut:      60 * 60 * 1000,
		Subscribe:        ".*\\..*",
		MsgChSize:        128,
		KafkaProducer: KafkaProducer{
			Addr: []string{"127.0.0.1" + ":9092"},
		},
		Redis: Redis{
			Addr: []string{"127.0.0.1" + ":7001", "127.0.0.1" + ":7002", "127.0.0.1" + ":7003"},
		},
		Mysql: mysqlC,
	}
}

type Logic struct {
	Id            int
	Addr          string
	Name          string
	KafkaProducer KafkaProducer
	Etcd          Etcd
	Mysql         Mysql
}

func (c *Logic) Default() {
	ip, err := util.GetOutBoundIP()
	if err != nil {
		log.Panicln("faild to get outbound ip:", err)
	}
	mysqlC := Mysql{}
	mysqlC.Default()
	*c = Logic{
		Id:   0,
		Addr: ip + ":50060",
		Name: "0",
		KafkaProducer: KafkaProducer{
			Addr: []string{ip + ":9092"},
		},
		Etcd: Etcd{
			Addr: []string{ip + ":2379"},
		},
		Mysql: mysqlC,
	}
}

type Comet struct {
	Addr          string
	Name          string
	Etcd          Etcd
	WebsocketAddr string

	BucketSize int
}

func (c *Comet) Default() {
	ip, err := util.GetOutBoundIP()
	if err != nil {
		log.Panicln("faild to get outbound ip:", err)
	}
	*c = Comet{
		Addr: ip + ":50050",
		Name: "0",
		Etcd: Etcd{
			Addr: []string{ip + ":2379"},
		},
		BucketSize:    32,
		WebsocketAddr: ip + ":40050",
	}
}

type Job struct {
	Addr             string
	Name             string
	KafkaComsumer    KafkaComsumer
	Etcd             Etcd
	Redis            Redis
	BucketSize       int
	CometRoutineSize int
	Mysql            Mysql
}

func (c *Job) Default() {
	ip, err := util.GetOutBoundIP()
	if err != nil {
		log.Panicln("faild to get outbound ip:", err)
	}
	mysqlC := Mysql{}
	mysqlC.Default()
	*c = Job{
		Addr: ip + ":50070",
		Name: "0",
		KafkaComsumer: KafkaComsumer{
			Addr:    []string{ip + ":9092"},
			Topics:  []string{"laneIM"},
			GroupId: "job",
		},
		Etcd: Etcd{
			Addr: []string{ip + ":2379"},
		},
		Redis: Redis{
			[]string{ip + ":7001", ip + ":7002", ip + ":7003"},
		},
		BucketSize:       32,
		CometRoutineSize: 32,
		Mysql:            mysqlC,
	}
}

type Etcd struct {
	Addr []string
}

type KafkaProducer struct {
	Addr []string
}

type Redis struct {
	Addr []string
}

type KafkaComsumer struct {
	Addr    []string
	Topics  []string
	GroupId string
}

type Mysql struct {
	Type     string
	Name     string
	Username string
	Password string
	Addr     string
	DataBase string
}

func (c *Mysql) Default() {
	*c = Mysql{
		Name:     "0",
		Username: "debian-sys-maint",
		Password: "FJho5xokpFqZygL5",
		Addr:     "127.0.0.1:3306",
		DataBase: "laneIM",
	}
}

type LaneConfig interface {
	Default()
}

// TODO etcd dynamic config

func WriteRemote(conf LaneConfig) error {

	return nil
}

// TODO etcd dynamic config
func ReadRemote(conf LaneConfig) {

}

func Init(Path string, conf LaneConfig) {
	_, err := os.Stat(Path)
	if err != nil {
		if os.IsNotExist(err) {
			conf.Default()
			WriteLocal(Path, conf)
			log.Println("please check for the config.yml if needed to be modified, then run again")
			os.Exit(1)
		} else {
			log.Panicln("wrong err:", err)
		}
	}
	ReadLocal(Path, conf)
}

func WriteLocal(Path string, conf LaneConfig) error {
	out, err := yaml.Marshal(conf)
	if err != nil {
		log.Panicln("failed to marshal config", Path, ":", err)
		return err
	}

	err = os.WriteFile(Path, out, 0644)
	if err != nil {
		log.Panicln("failed to write ", Path, err)
		return err
	}
	return nil
}

func ReadLocal(Path string, conf LaneConfig) error {
	log.Println("read from ", Path)
	data, err := os.ReadFile(Path)
	if err != nil {
		log.Panicln("config.yaml does not exist")
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		log.Panicln("can't not read config.yml")
	}
	return err
}
