package uidgen

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

type UidGenerator struct {
	node *snowflake.Node
}

func NewUidGenerator(nodeId int64) (*UidGenerator, error) {
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		return nil, fmt.Errorf("error generating new node: %w", err)
	}
	return &UidGenerator{node: node}, nil
}

func (ug *UidGenerator) NewId() int64 {
	return ug.node.Generate().Int64()
}
