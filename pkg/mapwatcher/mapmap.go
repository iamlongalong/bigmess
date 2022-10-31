package mapwatcher

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrTypeNotMap = errors.New("val of key is not map")
)

type M = map[string]interface{}

func NewMapMap() *MapMap {
	return &MapMap{val: make(map[string]interface{})}
}

type MapMap struct {
	mu sync.Mutex // 姑且先用一把大锁

	val M
}

func (mm *MapMap) Set(k string, val interface{}) (oldVal interface{}, err error) {
	subs, key := mm.subkeys(k)

	mm.mu.Lock()
	defer mm.mu.Unlock()

	m, err := mm.deepGetMapOrInit(subs, true)
	if err != nil {
		return nil, err
	}

	oldVal = m[key]
	m[key] = val

	return oldVal, nil
}

func (mm *MapMap) Get(k string) (val interface{}, err error) {
	subs, key := mm.subkeys(k)

	mm.mu.Lock()
	defer mm.mu.Unlock()

	m, err := mm.deepGetMapOrInit(subs, false)
	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, nil
	}

	return m[key], nil
}

func (mm *MapMap) Del(k string) (oldVal interface{}, err error) {
	subs, key := mm.subkeys(k)

	mm.mu.Lock()
	defer mm.mu.Unlock()

	m, err := mm.deepGetMapOrInit(subs, false)
	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, nil
	}

	oldVal = m[key]
	delete(m, key)

	return oldVal, nil
}

func (mm *MapMap) deepGetMapOrInit(sks []string, initIfNotFount bool) (map[string]interface{}, error) {
	preval := mm.val

	for i := 0; i < len(sks); i++ {
		k := sks[i]
		iv, ok := preval[k]
		if !ok {
			if initIfNotFount {
				iv = M{}
				preval[k] = iv
			}
		}

		v, ok := iv.(M)
		if !ok {
			return nil, ErrTypeNotMap
		}

		preval = v
	}

	return preval, nil
}

// subkeys 获取前缀 以及 当前key
func (mm *MapMap) subkeys(k string) ([]string, string) {
	sts := strings.Split(k, ".")

	if len(sts) == 1 {
		return make([]string, 0), k
	}

	return sts[0 : len(sts)-1], sts[len(sts)-1]
}
