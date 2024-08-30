package node

import (
	"encoding/csv"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sirupsen/logrus"

	"github.com/sika365/admin-tools/pkg/file"
)

type LocalNodes []*LocalNode

type LocalNode struct {
	simutils.CommonTableFields
	simutils.PolymorphicFields
	Alias         simutils.Slug `json:"alias,omitempty" gorm:"index:,unique" sim:"primaryKey;"`
	Name          string        `json:"name,omitempty"`
	Slug          simutils.Slug `json:"slug,omitempty" gorm:"index"`
	ParentID      *simutils.PID `json:"parent_id,omitempty" gorm:"default:NULL"`
	ParentAlias   simutils.Slug `json:"parent_alias,omitempty" gorm:"default:NULL;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Priority      int           `json:"priority,omitempty"`
	Parent        *LocalNode    `json:"parent,omitempty" gorm:"references:ParentAlias;foreignKey:Alias;"`
	SubNodes      LocalNodes    `json:"sub_nodes,omitempty" gorm:"references:Alias;foreignkey:ParentAlias"`
	SubNodesCount int           `json:"sub_nodes_count" gorm:"-:all"`
}

func FromFiles(files file.MapFiles, req ScanRequest, fn func(header map[string]int, rec []string)) (LocalNodes, MapNodes) {
	var (
		m     = make(MapNodes)
		nodes = make(LocalNodes, 0, len(files))
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

func (LocalNode) TableName() string {
	return "local_nodes"
}

func (n *LocalNode) Key() string {
	return n.ID.String()
}
