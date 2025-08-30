package infrastructure

type Node struct {
	ID int `json:"id"`
}

func (n *Node) GetID() int {
	return n.ID
}
