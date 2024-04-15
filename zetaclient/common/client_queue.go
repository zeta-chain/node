package common

// ClientQueue - Generic data structure to keep track of fallback clients used to connect to EVM and Bitcoin nodes
type ClientQueue struct {
	clientList []any
}

// NewClientQueue - constructor
func NewClientQueue() *ClientQueue {
	return &ClientQueue{
		clientList: []any{},
	}
}

// Append - adds another client to the end of the queue
func (c *ClientQueue) Append(item any) {
	c.clientList = append(c.clientList, item)
}

// Length - returns total length of queue or number of clients
func (c *ClientQueue) Length() int {
	return len(c.clientList)
}

// Next - rotates list of clients by moving the first to the last
func (c *ClientQueue) Next() {
	if len(c.clientList) < 1 {
		return
	}
	// Rotate client order by 1
	c.clientList = append(c.clientList[1:], c.clientList[:1]...)
}

// First - returns the first client in the list or highest priority
func (c *ClientQueue) First() any {
	if len(c.clientList) < 1 {
		return nil
	}
	return c.clientList[0]
}
