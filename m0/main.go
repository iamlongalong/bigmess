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

	wsEngine := NewEngine()
	wsEngine.Register(MessageCode("/fileroom/publish"), NewHandler(HandleFileRoomPub))

	engine.GET("/ws", func(ctx *gin.Context) {
		conn, err := dfupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("upgrate fail : %s", err)
			return
		}

		cli := NewClient(conn, wsEngine)

		cli.Start()
	})

	// http.ListenAndServeTLS(":8081", "", "", engine)
	if err := http.ListenAndServe(":8080", engine); err != nil {
		log.Fatalf("listen fail : %s", err)
	}
}
