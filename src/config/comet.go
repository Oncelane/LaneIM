package config

type Comet struct {
	Id            int
	UbuntuIP      string
	WindowIP      string
	GrpcPort      string
	WebsocketPort string
	Name          string
	Etcd          Etcd
	BucketSize    int
}

func (c *Comet) Default() {
	// ip, err := util.GetOutBoundIP()
	// if err != nil {
	// 	log.Panicln("faild to get outbound ip:", err)
	// }
	*c = Comet{
		Id:            100,
		UbuntuIP:      "172.xx",
		WindowIP:      "192.xx",
		GrpcPort:      ":50050",
		WebsocketPort: ":40050",
		Name:          "0",
		Etcd:          DefaultEtcd(),
		BucketSize:    32,
	}
}
