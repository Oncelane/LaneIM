package pkg

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

type UuidGenerator struct {
	node *snowflake.Node
}

func NewUuidGenerator(name int64) *UuidGenerator {
	node, err := snowflake.NewNode(name)
	if err != nil {
		log.Panicln("faild to create snowflake node")
	}
	return &UuidGenerator{
		node: node,
	}
}

func (u *UuidGenerator) Generator() (rt int64) {
	return u.node.Generate().Int64()
}
