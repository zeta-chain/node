package pubsub

import (
	"log"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
)

func TestAddTopic(t *testing.T) {
	q := NewEventBus()
	err := q.AddTopic("kek", make(<-chan coretypes.ResultEvent))
	require.NoError(t, err)

	err = q.AddTopic("lol", make(<-chan coretypes.ResultEvent))
	require.NoError(t, err)

	err = q.AddTopic("lol", make(<-chan coretypes.ResultEvent))
	require.Error(t, err)

	topics := q.Topics()
	sort.Strings(topics)
	require.EqualValues(t, []string{"kek", "lol"}, topics)
}

func TestMaxSubscribers(t *testing.T) {
	q := NewEventBus(WithMaxSubscribers(2))
	kekSrc := make(chan coretypes.ResultEvent)
	err := q.AddTopic("kek", kekSrc)
	require.NoError(t, err)
	_, _, err = q.Subscribe("kek")
	require.NoError(t, err)
	_, unsub, err := q.Subscribe("kek")
	require.NoError(t, err)

	_, _, err = q.Subscribe("kek")
	require.ErrorIs(t, err, ErrTooManySubscribers)

	unsub()
	_, _, err = q.Subscribe("kek")
	require.NoError(t, err)
}

func TestMaxSubscribersUpdatedAfterClose(t *testing.T) {
	maxSubs, topic := 5, "kek"
	q := NewEventBus(WithMaxSubscribers(maxSubs))
	kekSrc := make(chan coretypes.ResultEvent)
	err := q.AddTopic(topic, kekSrc)
	require.NoError(t, err)
	for range maxSubs {
		_, _, err = q.Subscribe(topic)
		require.NoError(t, err)
	}
	_, _, err = q.Subscribe(topic)
	require.ErrorIs(t, err, ErrTooManySubscribers)
	close(kekSrc)
	time.Sleep(1 * time.Second)

	_, _, err = q.Subscribe(topic)
	require.ErrorIs(t, err, ErrTopicNotFound)

	// should be able to subscribe back up to maximum after the topic was removed.
	kekSrc = make(chan coretypes.ResultEvent)
	err = q.AddTopic(topic, kekSrc)
	require.NoError(t, err)
	for range maxSubs {
		_, _, err = q.Subscribe(topic)
		require.NoError(t, err)
	}
}

func TestSubscribe(t *testing.T) {
	q := NewEventBus()
	kekSrc := make(chan coretypes.ResultEvent)

	err := q.AddTopic("kek", kekSrc)
	require.NoError(t, err)

	lolSrc := make(chan coretypes.ResultEvent)

	err = q.AddTopic("lol", lolSrc)
	require.NoError(t, err)

	kekSubC, _, err := q.Subscribe("kek")
	require.NoError(t, err)

	lolSubC, _, err := q.Subscribe("lol")
	require.NoError(t, err)

	lol2SubC, _, err := q.Subscribe("lol")
	require.NoError(t, err)

	wg := new(sync.WaitGroup)
	wg.Add(4)

	emptyMsg := coretypes.ResultEvent{}
	go func() {
		defer wg.Done()
		msg := <-kekSubC
		log.Println("kek:", msg)
		require.EqualValues(t, emptyMsg, msg)
	}()

	go func() {
		defer wg.Done()
		msg := <-lolSubC
		log.Println("lol:", msg)
		require.EqualValues(t, emptyMsg, msg)
	}()

	go func() {
		defer wg.Done()
		msg := <-lol2SubC
		log.Println("lol2:", msg)
		require.EqualValues(t, emptyMsg, msg)
	}()

	go func() {
		defer wg.Done()

		time.Sleep(time.Second)

		close(kekSrc)
		close(lolSrc)
	}()

	wg.Wait()
	time.Sleep(time.Second)
}

func TestConcurrentSubscribeAndPublish(t *testing.T) {
	var (
		wg        sync.WaitGroup
		eb        = NewEventBus()
		topicName = "lol"
		topicCh   = make(chan coretypes.ResultEvent)
		runsCount = 5
	)

	err := eb.AddTopic(topicName, topicCh)
	require.NoError(t, err)

	for i := 0; i < runsCount; i++ {
		subscribeAndPublish(t, eb, topicName, topicCh)
	}

	// close channel to make test end
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		close(topicCh)
	}()

	wg.Wait()
}

func subscribeAndPublish(t *testing.T, eb EventBus, topic string, topicChan chan coretypes.ResultEvent) {
	t.Helper()
	var (
		wg               sync.WaitGroup
		subscribersCount = 50
		emptyMsg         = coretypes.ResultEvent{}
	)
	for i := 0; i < subscribersCount; i++ {
		wg.Add(1)
		// concurrently subscribe to the topic
		go func() {
			defer wg.Done()
			_, _, err := eb.Subscribe(topic)
			require.NoError(t, err)
		}()

		// send events to the topic
		wg.Add(1)
		go func() {
			defer wg.Done()
			topicChan <- emptyMsg
		}()
	}
	wg.Wait()
}
