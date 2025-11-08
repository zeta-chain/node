package stream

import (
	"context"
	"sync"
)

// Stream implements a data stream, user can subscribe the stream in blocking or non-blocking way using offsets.
// it use a segmented ring buffer to store the items, with a fixed capacity, when buffer is full, the old data get pruned when new data comes.
type Stream[V any] struct {
	segments      *Queue[[]V]
	segmentSize   int
	maxSegments   int
	segmentOffset int
	cond          *Cond
	mutex         sync.RWMutex
}

func NewStream[V any](segmentSize, capacity int) *Stream[V] {
	maxSegments := (capacity + segmentSize - 1) / segmentSize
	if maxSegments < 1 {
		panic("capacity is too small")
	}

	stream := &Stream[V]{
		segments:      New[[]V](),
		segmentSize:   segmentSize,
		maxSegments:   maxSegments,
		segmentOffset: 0,
		cond:          NewCond(),
	}
	return stream
}

// Add appends items to the stream and returns the id of last one.
// item id start with 1.
func (s *Stream[V]) Add(vs ...V) int {
	if len(vs) == 0 {
		return 0
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, v := range vs {
		if s.segments.Length() == 0 || len(s.segments.Tail()) == s.segmentSize {
			var seg []V
			if s.segments.Length() > s.maxSegments {
				// reuse the free segment
				seg = s.segments.Remove()[:0]
				s.segmentOffset++
			} else {
				seg = make([]V, 0, s.segmentSize)
			}
			s.segments.Add(seg)
		}

		tail := s.segments.TailP()
		*tail = append(*tail, v)
	}

	// notify the subscribers
	s.cond.Broadcast()

	return s.lastID()
}

// Subscribe subscribes the stream in a loop, pass the chunks of items to the callback,
// it only stops if the context is canceled.
// it returns the last id of the items.
func (s *Stream[V]) Subscribe(ctx context.Context, callback func([]V, int) error) error {
	var (
		items  []V
		offset = -1
	)
	for {
		items, offset = s.ReadBlocking(ctx, offset)
		if len(items) == 0 {
			// canceled
			break
		}
		if err := callback(items, offset); err != nil {
			return err
		}
	}
	return nil
}

// ReadNonBlocking returns items with id greater than the last received id reported by user, without blocking.
// if there are no new items, it also returns the largest id of the items.
func (s *Stream[V]) ReadNonBlocking(offset int) ([]V, int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.doRead(offset)
}

// ReadAllNonBlocking returns all items in the stream, without blocking.
func (s *Stream[V]) ReadAllNonBlocking(offset int) ([]V, int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var (
		result []V
		items  []V
	)
	for {
		items, offset = s.doRead(offset)
		if len(items) == 0 {
			break
		}
		result = append(result, items...)
	}

	return result, offset
}

// ReadBlocking returns items with id greater than the last received id reported by user.
// reads at most one segment at a time.
// negative offset means read from the end.
// it also returns the largest id of the items, if there are no new items, returns the id of the last item.
func (s *Stream[V]) ReadBlocking(ctx context.Context, offset int) ([]V, int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var items []V
	for {
		items, offset = s.doRead(offset)
		if len(items) > 0 {
			return items, offset
		}

		s.mutex.RUnlock()
		r := s.cond.Wait(ctx)
		s.mutex.RLock()

		if !r {
			// context canceled
			return nil, 0
		}
	}
}

// lastID returns the id of the last item, 0 for empty stream.
func (s *Stream[V]) lastID() int {
	if s.segments.Length() == 0 {
		return 0
	}

	return (s.segmentOffset+s.segments.Length()-1)*s.segmentSize + len(s.segments.Tail())
}

// doRead is the underlying logic of Read.
func (s *Stream[V]) doRead(offset int) ([]V, int) {
	if s.segments.Length() == 0 {
		return nil, s.lastID()
	}

	if offset < 0 {
		return nil, s.lastID()
	}

	segment := offset / s.segmentSize
	if segment >= s.segmentOffset {
		segment -= s.segmentOffset
	} else {
		// the target segment is pruned, ajust to earliest segment
		segment = 0
	}

	if segment >= s.segments.Length() {
		// offset is in the future
		return nil, s.lastID()
	}

	seg := s.segments.Get(segment)
	items := seg[offset%s.segmentSize:]
	if len(items) == 0 {
		return nil, s.lastID()
	}

	// copy the slice
	clone := make([]V, len(items))
	copy(clone, items)

	return clone, (s.segmentOffset+segment)*s.segmentSize + len(seg)
}
