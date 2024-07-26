package node

import (
	"encoding/csv"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sirupsen/logrus"

	"github.com/sika365/admin-tools/pkg/file"
)

type Nodes []*Node

type Node struct {
	simutils.CommonTableFields
	simutils.PolymorphicFields
	Name          string        `json:"name,omitempty"`
	Alias         string        `json:"alias,omitempty" gorm:"index"`
	Slug          string        `json:"slug,omitempty" gorm:"index"`
	ParentID      *simutils.PID `json:"parent_id,omitempty" gorm:"default:null;index"`
	Priority      int           `json:"priority,omitempty"`
	Parent        *Node         `json:"parent,omitempty"`
	SubNodes      Nodes         `json:"sub_nodes,omitempty" gorm:"foreignkey:ParentID"`
	SubNodesCount int           `json:"sub_nodes_count" gorm:"-:all"`
}

func FromFiles(files file.MapFiles, req ScanRequest, fn func(header map[string]int, rec []string)) (Nodes, MapNodes) {
	var (
		m     = make(MapNodes)
		nodes = make(Nodes, 0, len(files))
	)

	for _, f := range files {
		reader := csv.NewReader(f.Open().Reader())
		// Read header
		header := make(map[string]int)
		if h, err := reader.Read(); err != nil {
			logrus.Errorf("failed to read header row %d: %v", 1, err)
			return nil, nil
		} else {
			for i, t := range h {
				header[t] = i
			}
		}
		// Skip the specified number of rows (offset)
		for i := 1; i < req.Offset; i++ {
			if _, err := reader.Read(); err != nil {
				logrus.Errorf("failed to skip row %d: %v", i+1, err)
				return nil, nil
			}
		}
		// Read the remaining records from the CSV file
		i := req.Offset + 1 // 1 row header
		for {
			if r, err := reader.Read(); err != nil {
				logrus.Errorf("failed to read row %d: %v", i+1, err)
				return nil, nil
			} else {
				fn(header, r)
			}
		}
	}
	return nodes, m
}

func (Node) TableName() string {
	return "nodes"
}

func (n *Node) Key() string {
	return n.ID.String()
}
