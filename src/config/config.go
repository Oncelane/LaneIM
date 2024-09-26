package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Canal struct {
	Name             string
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
		Name:             "0",
		CanalAddress:     "127.0.0.1",
		CanalPort:        11111,
		CanalName:        "canal",
		CanalPassword:    "canal",
		CanalDestination: "laneIM",
		SoTimeOut:        60000,
		IdleTimeOut:      60 * 60 * 1000,
		Subscribe:        ".*\\..*",
		MsgChSize:        128,
		KafkaProducer:    DefaultKafkaProducer(),
		Redis:            DefaultRedis(),
		Mysql:            mysqlC,
	}
}

type Etcd struct {
	Addr []string
}

func DefaultEtcd() Etcd {
	return Etcd{
		Addr: []string{
			"127.0.0.1:51240",
			"127.0.0.1:51241",
			"127.0.0.1:51242"},
	}
}

type KafkaProducer struct {
	Addr []string
}

func DefaultKafkaProducer() KafkaProducer {
	return KafkaProducer{
		Addr: []string{"127.0.0.1" + ":9092"},
	}
}

type Redis struct {
	Addr []string
}

func DefaultRedis() Redis {
	return Redis{
		[]string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003"},
	}
}

type KafkaComsumer struct {
	Addr    []string
	Topics  []string
	GroupId string
}

func DefaultKafkaComsumer() KafkaComsumer {
	return KafkaComsumer{
		Addr:    []string{"127.0.0.1:9092"},
		Topics:  []string{"laneIM"},
		GroupId: "job",
	}
}

type Mysql struct {
	Name        string
	Username    string
	Password    string
	Addr        string
	DataBase    string
	BatchWriter BatchWriter
}

func (c *Mysql) Default() {
	batch := DefaultBatchWriter()
	*c = Mysql{
		Name:        "0",
		Username:    "debian-sys-maint",
		Password:    "QTLVb6BaeeaJsFMT",
		Addr:        "127.0.0.1:3306",
		DataBase:    "laneIM",
		BatchWriter: batch,
	}
}

func DefaultMysql() Mysql {
	mysqlC := Mysql{}
	mysqlC.Default()
	return mysqlC
}

type BatchWriter struct {
	MaxTime  int
	MaxCount int
}

func (c *BatchWriter) Default() {
	*c = BatchWriter{
		MaxTime:  100,
		MaxCount: 1000,
	}
}

func DefaultBatchWriter() BatchWriter {
	BatchWriter := BatchWriter{}
	BatchWriter.Default()
	return BatchWriter
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
