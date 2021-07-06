package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var connMap map[string]*websocket.Conn = make(map[string]*websocket.Conn)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//webSocket请求ping 返回pong
func WsHanderFunc(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	machID := c.Query("machid")
	ConnRegister(machID, ws)

	for {
		// 读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("websocket read message error:[%s]", err.Error())
			continue
		}

		switch mt {
		case websocket.PingMessage:
			err = ws.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				log.Printf("websocket pong response error:[%s]", err.Error())
			}
			continue

		default:
			log.Printf("websocket read message, type:[%d], data:[%s]", mt, message)
		}

		// close condition:

	}
}

func NotifyMsg(machID, msg string) {
	fd := ConnFinder(machID)

	err := fd.WriteMessage(websocket.BinaryMessage, []byte(msg))
	if err != nil {
		log.Errorf("websocket write message error:[%s]", err.Error())
	}
}

func main() {
	r := gin.Default()
	r.GET("/ws", WsHanderFunc)

	r.Run("localhost:8080")
}

//  userid -> nodes -> machineids -> connect
func ConnRegister(machID string, conn *websocket.Conn) {
	// yaqin

	connMap[machID] = conn
}

func ConnFinder(machID string) *websocket.Conn {
	return connMap[machID]
}
