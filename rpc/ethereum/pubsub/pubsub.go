package pubsub

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
)

type Option func(bus *memEventBus)

func WithMaxSubscribers(n int) Option {
	return func(bus *memEventBus) {
		bus.maxTotalSubscribers = n
	}
}

const (
	DefaultMaxSubscribers = 500_000
)

var (
	ErrTooManySubscribers = errors.New("too many subscribers")
	ErrTopicNotFound      = errors.New("topic not found")
)

type UnsubscribeFunc func()

type EventBus interface {
	AddTopic(name string, src <-chan coretypes.ResultEvent) error
	RemoveTopic(name string)
	Subscribe(name string) (<-chan coretypes.ResultEvent, UnsubscribeFunc, error)
	Topics() []string
}

type memEventBus struct {
	topics          map[string]<-chan coretypes.ResultEvent
	topicsMux       *sync.RWMutex
	subscribers     map[string]map[uint64]chan<- coretypes.ResultEvent
	subscribersMux  *sync.RWMutex
	currentUniqueID uint64

	maxTotalSubscribers int
	totalSubscribers    atomic.Int64
}

func NewEventBus(opts ...Option) EventBus {
	bus := &memEventBus{
		topics:              make(map[string]<-chan coretypes.ResultEvent),
		topicsMux:           new(sync.RWMutex),
		subscribers:         make(map[string]map[uint64]chan<- coretypes.ResultEvent),
		subscribersMux:      new(sync.RWMutex),
		maxTotalSubscribers: DefaultMaxSubscribers,
	}
	for _, opt := range opts {
		opt(bus)
	}
	return bus
}

func (m *memEventBus) GenUniqueID() uint64 {
	return atomic.AddUint64(&m.currentUniqueID, 1)
}

func (m *memEventBus) Topics() (topics []string) {
	m.topicsMux.RLock()
	defer m.topicsMux.RUnlock()

	topics = make([]string, 0, len(m.topics))
	for topicName := range m.topics {
		topics = append(topics, topicName)
	}

	return topics
}

func (m *memEventBus) AddTopic(name string, src <-chan coretypes.ResultEvent) error {
	m.topicsMux.RLock()
	_, ok := m.topics[name]
	m.topicsMux.RUnlock()

	if ok {
		return errors.New("topic already registered")
	}

	m.topicsMux.Lock()
	m.topics[name] = src
	m.topicsMux.Unlock()

	go m.publishTopic(name, src)

	return nil
}

func (m *memEventBus) RemoveTopic(name string) {
	m.topicsMux.Lock()
	delete(m.topics, name)
	m.topicsMux.Unlock()
}

func (m *memEventBus) Subscribe(name string) (<-chan coretypes.ResultEvent, UnsubscribeFunc, error) {
	m.topicsMux.RLock()
	_, ok := m.topics[name]
	m.topicsMux.RUnlock()

	if !ok {
		return nil, nil, errors.Wrapf(ErrTopicNotFound, name)
	}

	ch := make(chan coretypes.ResultEvent)
	m.subscribersMux.Lock()
	defer m.subscribersMux.Unlock()

	if m.maxTotalSubscribers > 0 && m.totalSubscribers.Load() >= int64(m.maxTotalSubscribers) {
		return nil, nil, errors.Wrap(ErrTooManySubscribers, fmt.Sprintf("%d", m.maxTotalSubscribers))
	}

	id := m.GenUniqueID()
	if _, ok := m.subscribers[name]; !ok {
		m.subscribers[name] = make(map[uint64]chan<- coretypes.ResultEvent)
	}
	m.subscribers[name][id] = ch
	m.totalSubscribers.Add(1)

	unsubscribe := func() {
		m.subscribersMux.Lock()
		defer m.subscribersMux.Unlock()
		if _, ok := m.subscribers[name][id]; ok {
			close(m.subscribers[name][id])
			delete(m.subscribers[name], id)
			m.totalSubscribers.Add(-1)
		}
	}

	return ch, unsubscribe, nil
}

func (m *memEventBus) publishTopic(name string, src <-chan coretypes.ResultEvent) {
	for {
		msg, ok := <-src
		if !ok {
			m.closeAllSubscribers(name)
			m.topicsMux.Lock()
			delete(m.topics, name)
			m.topicsMux.Unlock()
			return
		}
		m.publishAllSubscribers(name, msg)
	}
}

func (m *memEventBus) closeAllSubscribers(name string) {
	m.subscribersMux.Lock()
	defer m.subscribersMux.Unlock()

	subscribers := m.subscribers[name]
	delete(m.subscribers, name)
	m.totalSubscribers.Add(int64(-len(subscribers)))
	// #nosec G705
	for _, sub := range subscribers {
		close(sub)
	}
}

func (m *memEventBus) publishAllSubscribers(name string, msg coretypes.ResultEvent) {
	m.subscribersMux.RLock()
	defer m.subscribersMux.RUnlock()
	subscribers := m.subscribers[name]
	// #nosec G705
	for _, sub := range subscribers {
		select {
		case sub <- msg:
		default:
		}
	}
}
