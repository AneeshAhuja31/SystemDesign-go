package ticket

import (
	"errors"
	"sync"
)

type Slot struct{
	Start uint64
	End uint64
	Current uint64
}
type TicketServer struct {
	Mu sync.Mutex
	NextStart uint64
	BlockSize uint64
	MaxID uint64
}

func NewTicketServer(start, maxID, blockSize uint64) *TicketServer {
	return &TicketServer{
		NextStart: start, 
		BlockSize: blockSize,
		MaxID: maxID,
	}
}

func (t *TicketServer) AllocateSlot()(*Slot,error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if t.NextStart > t.MaxID{
		return nil,errors.New("no tickets left")
	}

	start := t.NextStart
	end := start + t.BlockSize - 1
	if end > t.MaxID {
		end = t.MaxID
	}

	t.NextStart = end + 1

	return &Slot{
		Start: start,
		End: end,
		Current: start,
	}, nil

}

type LocalTicketClient struct {
	mu sync.Mutex
	slot *Slot
	ts *TicketServer
}

func NewLocalTicketClient(ts *TicketServer) *LocalTicketClient {
	return &LocalTicketClient{ts: ts}
}

func (c *LocalTicketClient) NextID() (uint64, error) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.slot == nil || c.slot.Current > c.slot.End {

		slot, err := c.ts.AllocateSlot()
		if err != nil {
			return 0, err
		}

		c.slot = slot
	}

	id := c.slot.Current
	c.slot.Current++

	return id, nil
}


const base62chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func EncodeBase62(num uint64)string{
	if num == 0 {
		return "a"
	}
	result := ""

	for num > 0 {
		result = string(base62chars[num%62]) + result
		num /= 62
	}
	return result
}