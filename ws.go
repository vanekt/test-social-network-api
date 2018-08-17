package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"io"
	"net/http"
	"os"
	"sync"
	"vanekt/test-social-network-api/entity"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // TODO
}

type WS struct {
	logger       *logging.Logger
	httpServeMux *http.ServeMux
	wsConnMap    map[string]*wsConnection
	wsConnMapMu  *sync.RWMutex
}

func NewWebsocket(logger *logging.Logger) *WS {
	return &WS{
		logger:       logger,
		httpServeMux: http.NewServeMux(),
		wsConnMap:    make(map[string]*wsConnection),
		wsConnMapMu:  new(sync.RWMutex),
	}
}

func (ws *WS) Run(port string) {
	ws.httpServeMux.HandleFunc("/ws", ws.wsHandler)
	ws.logger.Fatal(http.ListenAndServe(port, ws.httpServeMux))
}

func (ws *WS) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws.logger.Debugf("Request for WebSocket connection by \"%s%s\"", r.Host, r.RequestURI)

	var authString string
	if authCookie, err := r.Cookie(os.Getenv("TOKEN_COOKIE_NAME")); err == nil {
		authString = authCookie.Value
	}

	httpHeader := http.Header{}
	conn, err := wsUpgrader.Upgrade(w, r, httpHeader)
	if err != nil {
		ws.sendHTTP500(w, "Failed to upgrade the connection", err)
		return
	}
	ws.logger.Debugf("Connection opened for %q", authString)

	currentWsConnection := &wsConnection{
		wsConn: conn,
		mu:     &sync.Mutex{},
		logger: ws.logger,
	}
	ws.addWsConn(authString, currentWsConnection)

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			ws.logger.Info(err)
			break
		}

		if mt != websocket.TextMessage {
			if mt == websocket.PingMessage || mt == websocket.PongMessage {
				continue
			}
			// TODO send error to client
			ws.logger.Error("400 Received message of unsupported type. Close the connection")
			continue
		}

		var msg entity.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			// TODO send error to client
			ws.logger.Errorf("Unmarshal error: %v", err)
			continue
		}

		ws.logger.Warning("mt", mt)
		ws.logger.Warning("message", string(message))
		ws.logger.Warning("msg", msg)
	}

	currentWsConnection.Close()
	ws.logger.Debugf("Connection closed for %q", authString)
	ws.delWsConn(authString)
}

func (ws *WS) sendHTTP500(w http.ResponseWriter, message string, err error) {
	ws.logger.Errorf("%s: %v", message, err)
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, message)
}

func (ws *WS) addWsConn(sID string, wsConn *wsConnection) {
	ws.wsConnMapMu.Lock()
	ws.wsConnMap[sID] = wsConn
	ws.logger.Notice(ws.wsConnMap)
	ws.wsConnMapMu.Unlock()
}

func (ws *WS) delWsConn(sID string) {
	ws.wsConnMapMu.Lock()
	delete(ws.wsConnMap, sID)
	ws.wsConnMapMu.Unlock()
}

func (ws *WS) getWsConn(sID string) (wsConn *wsConnection, ok bool) {
	ws.wsConnMapMu.RLock()
	wsConn, ok = ws.wsConnMap[sID]
	ws.wsConnMapMu.RUnlock()
	return
}