package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	wsEngine := NewEngine()
	wsEngine.Register(MessageCode("/fileroom/publish"), NewHandler(HandleFileRoomPub))

	// 使用简单的
	engine.GET("/ws", func(ctx *gin.Context) {
		conn, err := websocket.Upgrade(ctx.Writer, ctx.Request, nil, 1024, 1024)
		if err != nil {
			log.Printf("upgrate fail : %s", err)
			return
		}

		conn.ReadMessage()
	})

	// http.ListenAndServeTLS(":8081", "", "", engine)
	if err := http.ListenAndServe(":8080", engine); err != nil {
		log.Fatal()
	}
}
