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

## 总结

以上就是我们第一版的需求情况，简单来说，就是一个简单的、常规的 websocket 服务。

## 其他层面的思考

暂时没啥，先做吧。


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

### 代码实现

[bigmess 零号机](https://github.com/iamlongalong/bigmess/tree/main/m0)
