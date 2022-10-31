package mapwatcher

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestWatcher(t *testing.T) {

	m := NewResponsiveMap()
	l := NewMapListener("long.addr", func(me MapEvent) {
		fmt.Printf("listen1 : map changed. type : %s, key : %s, old : %+v, new : %+v\n", me.Option, me.Key, me.OldVal, me.NewVal)
	})

	m.Watch(l)

	m.Set("long.name", "longalong")
	m.Set("long.age", 18)
	m.Set("long.name", "longsang")

	v, err := m.Get("long.age")
	assert.Nil(t, err)
	assert.Equal(t, 18, v)

	v, err = m.Get("long.name")
	assert.Nil(t, err)
	assert.Equal(t, "longsang", v)

	v, err = m.Get("long.addr")
	assert.Nil(t, err)
	assert.Equal(t, nil, v)

	err = m.Set("long.addr", "sichuan")
	assert.Nil(t, err)

	v, err = m.Get("long.addr")
	assert.Nil(t, err)
	assert.Equal(t, "sichuan", v)

	v, err = m.Get("long")
	assert.Nil(t, err)
	spew.Dump(v)

	l2 := NewMapListener("long", func(me MapEvent) {
		fmt.Printf("listen2 : map changed. type : %s, key : %s, old : %+v, new : %+v\n", me.Option, me.Key, me.OldVal, me.NewVal)
	})
	m.Watch(l2)

	err = m.Set("long.wife", "shanzhu")
	assert.Nil(t, err)

	v, err = m.Get("long.addr")
	assert.Nil(t, err)
	assert.Equal(t, "sichuan", v)

	m.UnWatch(l)

	_, err = m.Get("long.addr")
	assert.Nil(t, err)

	err = m.Set("long.home.addr", "wuhan")
	assert.Nil(t, err)

	v, err = m.Get("long.home")
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"addr": "wuhan"}, v)

	err = m.Set("long.addr.province", "sichuan")
	assert.NotNil(t, err)
	spew.Dump(err)
}
