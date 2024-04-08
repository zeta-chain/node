package common

type ClientQueue struct {
	clientList []any
}

func NewClientQueue() *ClientQueue {
	return &ClientQueue{
		clientList: []any{},
	}
}

func (c *ClientQueue) Append(item any) {
	c.clientList = append(c.clientList, item)
}

func (c *ClientQueue) Length() int {
	return len(c.clientList)
}

func (c *ClientQueue) Next() {
	if len(c.clientList) < 1 {
		return
	}
	// Rotate client order by 1
	c.clientList = append(c.clientList[1:], c.clientList[:1]...)
}

func (c *ClientQueue) First() any {
	if len(c.clientList) < 1 {
		return nil
	}
	return c.clientList[0]
}
