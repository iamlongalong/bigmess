package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var dfupgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 5,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	// WriteBufferPool: ,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	wsEngine := NewEngine(EngineOpt{})
	wsEngine.Register(MessageCode("/publish"), &MiddleWareLog{}, NewHandler(HandlePubToRoom))
	wsEngine.Register(MessageCode("/fileroom/join"), &MiddleWareLog{}, NewHandler(HandleJoinRoom))
	wsEngine.Register(MessageCode("/fileroom/pub"), &MiddleWareLog{}, NewHandler(HandlePubMsg))

	engine.GET("/ws", func(ctx *gin.Context) {
		conn, err := dfupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("upgrate fail : %s", err)
			return
		}

		log.Printf("get new client : %s", ctx.Request.Host)

		cli := NewClient(conn, wsEngine)

		cli.Start()
	})

	// http.ListenAndServeTLS(":8081", "", "", engine)
	addr := ":8080"
	log.Printf("ws server is running at %s", addr)
	if err := http.ListenAndServe(addr, engine); err != nil {
		log.Fatalf("listen fail : %s", err)
	}
}
