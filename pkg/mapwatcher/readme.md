## 说明

这是一个可监听 map，当 map 中的 key 发生变化时，可以触发回调通知。

这个功能在 js 中非常常见，前端的几大响应式框架的核心能力便是基于此。

js 中，由于可以为 Object 设置 getter/setter，或者为 Object 增加 proxy 机制，就能够非常便捷地实现该能力，语言层面直接提供最基础的支持。

而在 golang 中，该能力只能自行封装 map 来实现。

本库，则是对该能力的一个实现。 本库暂未实现对 slice 的响应式封装，后续将会考虑。

另外，EventHub 是一个事件通知的模块，实现上接受的参数均为接口，意味着只要实现了 IEvent 和 INotify ，都可以使用该模块作为通知模块。

使用起来是比较简单的：

```go
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

```
