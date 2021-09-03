package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/xunlbz/go-restful-template/pkg/collector"
	"github.com/xunlbz/go-restful-template/pkg/lib"
	"github.com/xunlbz/go-restful-template/pkg/log"
)

var mutex sync.Mutex

type Message struct {
	Type     string `json:type`
	Value    string `json:"value"`
	Interval int    `json:"interval" default:"1"`
}

var TypeCollect string = "collect"
var TypeServiceLog string = "service_log"
var TypeContainerLog string = "container_log"

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

var ec = collector.NewEdgeCollector()
var dc = lib.NewDockerClient()

func checkOrigin(r *http.Request) bool {
	return true
}

type NewWriter struct {
	Conn        *ws.Conn
	MessageType int
}

func (w NewWriter) Write(data []byte) (n int, err error) {
	defer mutex.Unlock()
	mutex.Lock()
	return len(data), w.Conn.WriteMessage(w.MessageType, data)
}

func HandleServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade error: %s", err)
		return
	}
	//defer conn.Close()
	conn.SetCloseHandler(func(code int, text string) error {
		log.Infof("client closed %d, %s", code, text)
		conn.Close()
		return nil
	})
	disconnectCh := make(chan struct{})
	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Errorf("read: %v", err)
				close(disconnectCh)
				break
			}
			msg := Message{}
			json.Unmarshal(message, &msg)
			log.Infof("recv message: %+v", msg)
			go handleMessage(msg, conn, mt, disconnectCh)
		}
	}()
}

func syncWriteMessage(conn *ws.Conn, messageType int, data []byte) error {
	defer mutex.Unlock()
	mutex.Lock()
	return conn.WriteMessage(messageType, data)
}

func handleMessage(msg Message, conn *ws.Conn, messageType int, disconnectCh chan struct{}) {
	switch msg.Type {
	case TypeCollect:
		for {
			metrics := ec.GetModuleMitrics(msg.Value)
			data, err := json.Marshal(metrics)
			if err != nil {
				log.Error(err)
				return
			}
			err = syncWriteMessage(conn, messageType, data)
			if err != nil {
				log.Debugf("write error: %v", err)
				return
			}
			if msg.Interval > 1 {
				time.Sleep(time.Duration(msg.Interval) * time.Second)
			}
		}
	case TypeServiceLog:
		writer := NewWriter{
			Conn:        conn,
			MessageType: ws.BinaryMessage,
		}

		log.Infof("read service %s log start", msg.Value)
		end := make(chan time.Time)
		go func() {
			err := lib.JournalFollow(msg.Value, writer, end)
			if err != nil {
				log.Error(err)
			}
		}()
		if _, ok := <-disconnectCh; !ok {
			log.Infof("conn disconnect read service %s log end", msg.Value)
			end <- time.Now()
		}
	case TypeContainerLog:
		wr, err := dc.ContainerLog(msg.Value)
		if err != nil {
			log.Error(err)
			conn.WriteMessage(ws.BinaryMessage, []byte(err.Error()))
			return
		}
		defer wr.Close()
		buf := make([]byte, 1024)
		for {
			i, err := wr.Read(buf)
			if err != nil {
				log.Error(err)
				conn.WriteMessage(ws.BinaryMessage, []byte(err.Error()))
				return
			}
			if i >= 8 {
				log.Debugf("container logs : %s", string(buf[8:i]))
				err = conn.WriteMessage(ws.BinaryMessage, buf[8:i])
			}
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}
