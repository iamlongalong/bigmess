## 说明

我们回到最最最初的目标，实现一个最最最简单的消息分发系统。

## 需求

客户端有一些消息，想要发送给其他客户端。


## 需求分析

- 服务端有提供消息上传和分发的能力
- 客户端上传消息可以通过 http 也可以通过 websocket
- 考略到消息的实时性要求比较高，通过 websocket 分发

## 需求的技术分析以及调研

- 有竞争力一点的是 nodejs 生态 和 golang 生态。 nodejs 是因为和客户端语言生态一致，或许将来可以直接提供一致的SDK，nodejs 的异步机制能提供高并发。服务端目前都是以 golang 为基本语言，要提供的这个能力在 golang 上没有什么硬伤，golang 也天然支持高并发。 考虑到项目的主维护团队是 服务端团队，且 nodejs 在部署侧、后续维护性等方面弱于 goalng，因此选择 golang 为基础生态。

- websocket 需要先有 http 服务，可以使用原生的 golang http 包，也可以使用第三方的包。这不是个大问题的抉择，但考虑到其他服务主要使用的 [gin](https://github.com/gin-gonic/gin) 这个 http 框架，为了少踩坑，也为了降低将来团队内维护的阻力，我们也采用 gin 框架作为 http 部分。

- websocket 能力在 golang 下主要有 4 个实现，分别是 ① golang 官方库  ② [gobws/ws](https://github.com/gobwas/ws) ③ [gorilla/websocket](https://github.com/gorilla/websocket)④ [socket.io](https://github.com/googollee/go-socket.io) 。从协议实现完善度上，官方库稍次，gobws 和 gorilla 和 socket.io 差不多。 从性能上，gobws 是最好的。 从易用性上，socket.io 是最好的，这是一个封装了很多东西的库，从 nodejs 版本复刻过来的，底层实现是基于 gorilla 的。gorilla 是一个没有封装太多东西的纯 websocket 实现，可定制空间更大。 考虑到社区中使用最多的是 gorilla 的库，并且我们当前对性能的需求没那么大，更在意易用性，因此选择 gorilla/websocket 的方案了。 
> 没有选 socket.io 是由于整体来看它没有那么直接，有些东西算是黑盒，例如 http 降级、session封装、二进制消息封装等等，对我们来说不够透明，现在也没有那么多的时间和精力去细致地调研

- 业务逻辑上，存在 ”房间“ 的概念，这个房间内的消息是可以广播出去的。
- 在接收到消息后，有一些消息不需要做业务逻辑处理，直接分发到对应的通道中 (不那么重要的消息)。也有一些类型的消息需要做一些业务逻辑处理，然后由业务逻辑决定是否分发、怎么分发 (有特殊含义的消息)。
- 客户端可以做消息过滤，有些消息可以选择不接收，有些消息要接收。但不是每个消息都能由客户端决定是否接收。
- 客户端加入房间需要做权限校验，未通过鉴权的不给加入房间。
- 客户端发送的消息有限制，不是每个客户端都能发送每种消息。

- 当前业务的量非常小，仅仅是 MVP 实验阶段，因此直接用单机部署即可。

## 技术设计

- 实现 room 的概念，room hub 的概念，room 提供 `pub`/`join`/`leave` 三个接口。
- 在 ws 的协议之上，设计 `消息码` 的概念，对应到 http 请求中的 `请求行`，用来标识不同的请求类型，以此分发到不同的 hanlder 上。 handler 是实际提供业务逻辑的方法。
- 根据业务场景，auth 受到 role 的约束，例如 `管理员` 拥有发送一些特定请求的权力，而 `成员` 仅能发送通用请求的能力，`观察者` 甚至无法发送消息，仅能被动接受消息。 当前角色的权限模型不清晰，在程序中先体现为简单的 `if/else` 判断。
- 消息的结构需要有统一的约定，考虑到请求可能会带一些 `元信息`，因此和 http 请求类似，存在 `header` 部分，除了 header，还有 `消息码` 和 `body`。 为了开发简单，序列化直接用 json 进行传输即可。

## 前期总结

以上就是我们第一版的需求情况，简单来说，就是一个简单的、常规的 websocket 服务。

## 其他层面的思考

1. 框架中是否应该直接包含 room 的机制？ 还是交给业务层自己做？
2. http 是 ping/pong 机制，ws 是什么机制？应该是什么机制？有什么利弊？
3. 框架是怎么产生的？上来就设计？还是写业务代码，然后抽象、封装？
4. 在消息驱动的系统中，接口名设计有些什么内在需求？

## 技术设计的细节问题

### 消息码的设计

消息分发时，需要有标识来将请求分发到不同的流程。在 http 中, 这个标识是 `请求行` 也就是 `method + path` 的方式， 在 grpc 中，这是一个字符串，由 `包名 + 服务名 + 方法名` 构成，例如 `/xxxService.xxService/Doxxx` 。 在 ws 的设计中，我们可以用同样的方式来设计，结构可能如下：
```json
{
    "messageCode": "publish",
    "header": {},
    "body": null
}
```
但这样的区分度就不高了，没有 `分组` 的概念在当中，为了增加将来的可扩展性，增加分组，变成和 path 类似的格式：
```json
{
    "messageCode": "/fileroom/publish",
    "header": {
        "customMeta": {"roomID": "123456"}
    },
    "body": "A79a11209b123b90=="
}
```

虽然长得像是个阉割版的 http，但能解决我们的问题了，不错。
> 为什么叫 messageCode ? 没为啥，一个称呼而已，就跟 http 叫 path、 rpc 中叫 method 一样，ws 中可以认为万物皆 "消息"，不妨就叫 消息码。

在实现消息分发时，先不做分组的实现了，相当于虽然不在程序中处理分组，一个 messageCode 就是一个路由。但这种方式存在了分组的潜力，将来可以加上分组的概念，比如增加 分组 Middleware 之类的能力。

### 房间的设计

在需求中，主要的场景是： 把同一个房间内的消息广播出去。 因此，我们要在程序中实现 `房间`。 

从行为上看，一个房间最起码需要具备这些能力： 
- 把一个 client 加入到房间中
- 把一个 client 移出某个房间
- 把一个消息在房间内广播

client 最好有一个 id，用来全局唯一标识一个 client。 我们给 client 增加一个 `func ID() string` 的方法。

那么 room 的结构姑且为这样：
```go
type IRoom interface {
     JoinRoom(c *Client) error 
     Leave(ID string) error
     BoardCast(msg *Message) error 
}
```

这是一个最简的 room 了，有什么更多的需求之后再说。

除了 `房间` 外，我们还需要一个 `房间管理器`，这个东西从行为上，需要有这些能力：
- 在全局创建房间
- 获取房间的实例
- 销毁一个房间
- 给一些/所有房间广播消息

那么这个管理器姑且为这样：
```go
type IRoomHub interface {
    GetRoom(ID string) (IRoom, bool)
    CreateRoom(ID string) (IRoom, error)
    GetOrCreateRoom(ID string) (IRoom,error)
    DestoryRoom(ID string) error
    BoardCast(msg *Message, roomIds ...string) error
}
```
销毁房间的方法中，我们姑且只是把这个房间从管理器中移除，暂时不错更复杂的交互处理。

roomhub 姑且采用单例模式，毕竟现在来看，我们的 hub 只需要有一个就行了。

### message 的设计

在消息的处理过程中 ，我们需要持续将 message 往下传，那么 message 的格式应该是什么样的呢？

在 http 体系下，最常用的方式，就是不同接口定义不同的 model 类型，在接口中去做 Unmarshal 操作，然后将自定义的 model 往下传。 这样的好处在于，流程清晰，实现简单。 坏处在于，基于 model 的传递，意味着所有处理函数都是 `硬编码` 的，只能接收特定的参数，增加了冗余代码的工作量。

我们可以选择用 接口 的方式，message 分成的几个部分，分别用接口获取，接口类似于如下：
```go
type IMessage interface {
    MessageCode() MessageCode
    Header() MessageHeader
    Body() ???
}
```

这里就会发现，前两部分都是确定的类型，但是 body 的类型是不定的，是根据接口自行定义的，如何解决 Imessage 往下传递时，保持同样的接口，而 body 又能实现多态呢？

有这么几种思路：
1. 全局统一实现 json decoder、pb decoder 等等，直接在解 IMessage 时就转成结构体 (需要在初始化时注册结构体)，然后用 `interface{}` 往下传递，使用的时候直接类型断言。
2. IMessage 不处理 body，而是直接以 `[]byte` 的结构传给下游处理，handler 增加一个 Decode() 的方法。
3. 使用接口定义文档，定义好接口后，自动生成对应的 decode 方法，并生成对应的方法函数，以此降低手工维护冗余代码的工作量。

这几种思路各有优劣，第一种，优点是，能够保证接出来的结构体能在多个 middleware 中使用，缺点是，用 interface{} 做类型断言比较难看。 
第二种，优点是，流程清晰易理解，是现有的 http 体系下常用的方案，缺点是，middleware 中如果想使用 body，就得重复解body。
第三种，优点是，生成后的代码很好理解，开发时也很友好，缺点是，实现的成本较高，需要开发一系列的工具。

考虑到 第二种整体容易理解，开发成本也比较低，多个 middleware 要使用 body 的假设也不知道是否真的成立，因此先采用这种方案，后面遇到实际的问题再改吧。

另外，第三种方案之后可以考虑做一下，思路比较有意思。实际上，像 grpc 和 go-zero 都是采用的这种方案，确实十分好用。

### 连接的建立过程

在通过了 ws 的格式校验后 (upgrade)，算是把连接建立好了，按理下一步就是业务逻辑，但是，很多业务逻辑都有相同的前置条件，例如 `鉴权`。 

而 ws 和 http 很大的不同就是，http 的鉴权往往是直接通过 header 中的 token ，或者通过 Cookie 中的 sessionID，然后从 `内存缓存` 中、`redis` 中、`数据库` 中、`外部系统` 等方式进行 校验 以及 获取 user 信息 (用户基本信息，看场景有时候也可以包含 role 等权限信息)。

在 http 中，我们鉴权是直接通过 MiddleWare 对所有接口 或者 一些接口组 进行鉴权，这是增量法，有时候根据情况会使用接口白名单的方式反向剔除一些，这是排除法。

ws 是一个 `有状态` 的协议。 相比于 http 需要每次传递 token， ws 可以直接把各自认可的东西放在内存里，我们可以称之为 `已认可状态`，这样就可以不用每次都传递相同的东西了，例如 token。 比如鉴权过程，不同接口的鉴权，我们直接在 client 中添加一个 key 为 "user", 当这个值存在内容时，我们就认为已经做过登录验证了，其他类似。

实际上，ws 的这种 `状态` ，在 http 的生态下也是存在的，session 和它实现的功能是一样的，不同点在于，http 是不保证 `连接持续存在` 的，因此 http 不能使用 `连接` 来关联 session，所以呈现出了用一个唯一 id 来关联的现状。 而这种状态的关联，在 ws 下则直接以 `tcp连接` 来保证存续和映射。

基于以上分析，我们把 状态 交给 client 实例去负责，各个接口通过各自的 middleware 去做各自关心的逻辑，middleware 直接从 client 实例中去获取自己关心的状态值。

### 代码实现

[bigmess 零号机](https://github.com/iamlongalong/bigmess/tree/main/m0)


## 稍微总结一下

这个项目写的比较简单，有很多地方都写得很丑，也没有可扩展性，但时间紧急，这版能基本能满足我们一期的需求了，之后我们把这些地方找出来，进行一些改进。

二期的需求正在策划中，详情可以从 [一个消息系统的进化史](https://blog.longalong.cn/posts/22_10_28_15_04_envolution_of_a_message_system.html) 查看后续进度。

