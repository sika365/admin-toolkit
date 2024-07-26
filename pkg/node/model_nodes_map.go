package node

type MapNodes map[string]*Node

func NewMapNodes(nodes ...*Node) MapNodes {
	mns := make(MapNodes)
	for _, n := range nodes {
		mns.Add(n)
	}
	return mns
}

func (mns MapNodes) Add(node *Node) MapNodes {
	if k := node.Key(); k == "" {
		return mns
	} else {
		mns[k] = node
		return mns
	}
}

func (mns MapNodes) Get(key string) *Node {
	return mns[key]
}

func (mns MapNodes) GetKeys() []string {
	keys := make([]string, 0, len(mns))
	for k := range mns {
		keys = append(keys, k)
	}
	return keys
}

func (mns MapNodes) GetValues() Nodes {
	ns := make(Nodes, 0, len(mns))
	for _, n := range mns {
		ns = append(ns, n)
	}
	return ns
}
