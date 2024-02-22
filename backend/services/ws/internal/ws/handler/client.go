package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/OktopUSP/oktopus/ws/internal/usp_record"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	//maxMessageSize = 512

	// Websockets version of the protocol
	wsVersion = "13"

	// USP specification version
	uspVersion = "v1.usp"
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	//Websockets client endpoint id, eid follows usp specification
	eid string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan message
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
// cEID = controller endpoint id
func (c *Client) readPump(cEID string) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	//c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message := constructMsg(cEID, c.eid, data)
		c.hub.broadcast <- message
	}
}

func constructMsg(eid string, from string, data []byte) message {
	if eid == "" {
		var record usp_record.Record
		err := proto.Unmarshal(data, &record)
		if err != nil {
			log.Println(err)
		}
		eid = record.ToId
	}
	return message{
		eid:     eid,
		from:    from,
		data:    data,
		msgType: websocket.BinaryMessage,
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Println("The hub closed the channel of", c.eid)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(message.msgType)
			if err != nil {
				return
			}
			w.Write(message.data)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				send := <-c.send
				w.Write(send.data)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Handle USP Controller events
func ServeController(w http.ResponseWriter, r *http.Request, token, cEID string, authEnable bool) {
	if authEnable {
		recv_token := r.URL.Query().Get("token")
		if recv_token != token {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, eid: cEID, conn: conn, send: make(chan message)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump("")
}

// Handle USP Agent events, cEID = controller endpoint id
func ServeAgent(w http.ResponseWriter, r *http.Request, cEID string) {

	//TODO: find out a way to authenticate agents

	header := http.Header{
		"Sec-Websocket-Protocol": {uspVersion},
		"Sec-Websocket-Version":  {wsVersion},
	}

	deviceid := extractDeviceId(r.Header)
	if deviceid == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Device id not found")
		w.Write([]byte("Device id not found"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, header)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, eid: deviceid, conn: conn, send: make(chan message)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	//TODO: get cEID from device message toId record field (must refact nice part of the code for this to be dynamic)
	go client.readPump(cEID)
}

// gets device id from websockets header
func extractDeviceId(header http.Header) string {

	// Header must be like that: bbf-usp-protocol; eid="<endpoint-id>" <endpoint-id> is the same ar the record.FromId/record.ToId
	// log.Println("Header sec-websocket-extensions:", header.Get("sec-websocket-extensions"))
	wsHeaderExtension := header.Get("sec-websocket-extensions")

	// Split the input string by double quotes
	deviceid := strings.Split(wsHeaderExtension, "\"")
	if len(deviceid) < 2 {
		return ""
	}

	return deviceid[1]
}
