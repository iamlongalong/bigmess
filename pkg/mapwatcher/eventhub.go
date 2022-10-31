package mapwatcher

import (
	"fmt"
	"sync"
)

type EventsMap struct {
	val      map[string]interface{}
	eventhub *EventRoom

	mu sync.Mutex
}

type MapListener struct {
	id string
}

func (m *MapListener) Name() string {
	return m.id
}

func (m *MapListener) Notify(ie IEvent) {
	e := ie.Msg().(*MapEvent)

	fmt.Printf("key : %s changed : %s\n", e.Key, e.Option)
	fmt.Printf("old : %+v, new : %+v\n", e.OldVal, e.NewVal)
}

type WatchOpt struct {
	Key          string // 事件名
	WatchSubPath bool   // 监听子路径
}

func (m *EventsMap) UnWatch(k string) error {
	return nil
}
func (m *EventsMap) Watch(k string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

func (m *EventsMap) Del(k string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

func (m *EventsMap) Get(k string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

func (m *EventsMap) Set(k string, v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

func NewEventRoom() EventRoom {

	return EventRoom{}
}

type EventRoom struct {
	Name      string
	mu        sync.Mutex
	Listeners map[string]IListener
}

type MapKeyOption string

var (
	MapKeyDel = "delete"
	MapKeyAdd = "add"
	MapKeySet = "set"
	MapKeyGet = "get"
)

type MapEvent struct {
	Option MapKeyOption

	Key    string
	OldVal interface{}
	NewVal interface{}
}

func (e *MapEvent) Name() string {
	return string(e.Option)
}

func (e *MapEvent) Msg() interface{} {
	return e
}

type IEvent interface {
	Name() string

	Msg() interface{}
}

// On 注册一个监听
func (h *EventRoom) On(l IListener) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Listeners[l.Name()] = l

	return nil
}

func (h *EventRoom) Process(msg IEvent) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, l := range h.Listeners {
		l.Notify(msg)
	}

	return nil
}

type IListener interface {
	Name() string
	Notify(IEvent)
}
