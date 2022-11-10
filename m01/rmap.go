package main

import "fmt"

func init() {
	rmap := rmap.NewResponsiveMap()

	// create a key listener
	cityKeyListener := NewMapListener("city", func(me rmap.MapEvent) {
		fmt.Printf("city changed, key: %s, option type: %s, oldval: %+v, newval: %+v\n", me.Key, me.Option, me.OldVal, me.NewVal)
	})

	// watch key
	rmap.Watch(cityKeyListener)
}
