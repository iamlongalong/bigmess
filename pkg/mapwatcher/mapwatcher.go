package mapwatcher

import "sync"

func NewMapListener(key string, cb func(MapEvent)) IListener {
	return &MapListener{key: key, cb: cb}
}

type MapListener struct {
	key string
	cb  func(MapEvent)
}

func (m *MapListener) Name() string {
	return m.key
}

func (m *MapListener) Notify(ie IEvent) {
	e := ie.Msg().(*MapEvent)

	m.cb(*e)
}

func NewResponsiveMap() *ResponsiveMap {
	return &ResponsiveMap{
		val:      NewMapMap(),
		eventhub: NewEventHub(),
	}

}

type ResponsiveMap struct {
	mu sync.Mutex

	val      *MapMap
	eventhub *EventHub
}

func (m *ResponsiveMap) UnWatch(i IListener) {
	m.eventhub.Off(i)
}

func (m *ResponsiveMap) Watch(i IListener) {
	m.eventhub.On(i)
}

func (m *ResponsiveMap) Del(k string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	old, err := m.val.Del(k)
	if err != nil {
		return err
	}

	e := &MapEvent{Option: MapKeyOption(MapKeyDel), Key: k, OldVal: old, NewVal: nil}
	m.eventhub.Process(e)

	return nil
}

func (m *ResponsiveMap) Get(k string) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	val, err := m.val.Get(k)
	if err != nil {
		return nil, err
	}

	e := &MapEvent{Option: MapKeyOption(MapKeyGet), Key: k, OldVal: val, NewVal: nil}
	m.eventhub.Process(e)

	return val, nil
}

func (m *ResponsiveMap) Set(k string, v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	old, err := m.val.Set(k, v)
	if err != nil {
		return err
	}

	e := &MapEvent{Option: MapKeyOption(MapKeySet), Key: k, OldVal: old, NewVal: v}
	m.eventhub.Process(e)

	return nil
}

type MapKeyOption string

var (
	MapKeyGet = "get"
	MapKeySet = "set"
	MapKeyDel = "del"
)

type MapEvent struct {
	Option MapKeyOption

	Key    string
	OldVal interface{}
	NewVal interface{}
}

func (e *MapEvent) Name() string {
	return string(e.Key)
}

func (e *MapEvent) Msg() interface{} {
	return e
}
