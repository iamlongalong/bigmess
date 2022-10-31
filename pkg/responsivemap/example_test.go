package rmap

import "fmt"

func ExampleResponsiveMap() {
	// init a new responsive map
	rmap := NewResponsiveMap()

	// create a key listener
	cityKeyListener := NewMapListener("city", func(me MapEvent) {
		fmt.Printf("city changed, key: %s, option type: %s, oldval: %+v, newval: %+v\n", me.Key, me.Option, me.OldVal, me.NewVal)
	})

	// watch key
	rmap.Watch(cityKeyListener)

	// set a city, this will trigger the listener above
	rmap.Set("city.name", "beijing")
	rmap.Set("city.size", "huge")

	// try get, this will trigger too
	v, err := rmap.Get("city.name")
	if err != nil {
		panic(err)
	}
	fmt.Printf("got city name : %s\n", v)

	// try del, this will trigger as well
	err = rmap.Del("city")
	if err != nil {
		panic(err)
	}

}
