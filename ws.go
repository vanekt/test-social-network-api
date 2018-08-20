package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"vanekt/test-social-network-api/entity"
	"vanekt/test-social-network-api/model"
	"vanekt/test-social-network-api/util"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // TODO
}

type WS struct {
	logger             *logging.Logger
	httpServeMux       *http.ServeMux
	wsConnMap          map[string]*wsConnection
	wsConnMapMu        *sync.RWMutex
	userConnMap        map[uint32]map[string]struct{}
	userConnMapMu      *sync.RWMutex
	messagesModel      *model.MessagesModel
	sendMessageChannel chan *entity.WSMessage
}

func NewWebsocket(logger *logging.Logger, messagesModel *model.MessagesModel) *WS {
	return &WS{
		logger:             logger,
		httpServeMux:       http.NewServeMux(),
		wsConnMap:          make(map[string]*wsConnection),
		wsConnMapMu:        new(sync.RWMutex),
		userConnMap:        make(map[uint32]map[string]struct{}),
		userConnMapMu:      new(sync.RWMutex),
		messagesModel:      messagesModel,
		sendMessageChannel: make(chan *entity.WSMessage),
	}
}

func (ws *WS) Run(port string) {
	go (func() {
		for {
			wsMessage := <-ws.sendMessageChannel
			var targetConnections map[string]struct{}
			switch wsMessage.Type {
			case entity.WSMessageTypeCreateMessageSuccess:
				newMessage := wsMessage.Payload.(*entity.Message)
				var ok bool
				targetConnections, ok = ws.getUserConn(newMessage.AuthorId)
				if !ok {
					ws.logger.Warningf("Cannot find ws user connections for %d", newMessage.AuthorId)
					continue
				}
			case entity.WSMessageTypeNewMessage:
				newMessage := wsMessage.Payload.(*entity.Message)
				if newMessage.PeerId == newMessage.AuthorId {
					ws.logger.Warning("WSMessageTypeNewMessage skip: newMessage.PeerId == newMessage.AuthorId")
					continue
				}
				var ok bool
				targetConnections, ok = ws.getUserConn(newMessage.PeerId)
				if !ok {
					ws.logger.Warningf("Cannot find ws user connections for %d", newMessage.PeerId)
					continue
				}
			default:
				ws.logger.Warning("Unknown WS message type")
				continue
			}

			for connectionId := range targetConnections {
				connection, ok := ws.getWsConn(connectionId)
				if !ok {
					ws.logger.Warningf("Cannot find ws user connection: %s", connectionId)
					continue
				}
				payload, err := json.Marshal(wsMessage)
				if err != nil {
					ws.logger.Error("Cannot send ws message: %s", err.Error())
				}
				connection.SendTextMessage(payload)
			}
		}
	})()

	ws.httpServeMux.HandleFunc("/ws", ws.wsHandler)
	ws.logger.Fatal(http.ListenAndServe(port, ws.httpServeMux))
}

func (ws *WS) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws.logger.Debugf("Request for WebSocket connection by \"%s%s\"", r.Host, r.RequestURI)

	var tokenString string
	if authCookie, err := r.Cookie(os.Getenv("TOKEN_COOKIE_NAME")); err == nil {
		tokenString = authCookie.Value
	}

	authUserId, err := util.GetUserIdFromToken(tokenString)
	if err != nil {
		ws.logger.Warningf("[wsHandler] Cannot fetch authUserId from token. Trace %s", err.Error())
	}

	httpHeader := http.Header{}
	conn, err := wsUpgrader.Upgrade(w, r, httpHeader)
	if err != nil {
		ws.sendHTTP500(w, "Failed to upgrade the connection", err)
		return
	}
	ws.logger.Debugf("Connection opened for %d", authUserId)

	currentWsConnection := &wsConnection{
		wsConn: conn,
		mu:     &sync.Mutex{},
		logger: ws.logger,
	}

	connectionId := strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(int(authUserId))
	ws.addWsConn(connectionId, currentWsConnection)
	ws.addUserConn(authUserId, connectionId)

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			ws.logger.Error(err)
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

		//ws.logger.Warning("mt", mt)
		//ws.logger.Warning("message", string(message))
		//ws.logger.Warning("msg", msg)

		newMessage, err := ws.messagesModel.CreateMessage(&msg)
		ws.logger.Notice("newMessage", newMessage)

		ws.sendMessageChannel <- &entity.WSMessage{
			Type:    entity.WSMessageTypeCreateMessageSuccess,
			Payload: newMessage,
		}

		ws.sendMessageChannel <- &entity.WSMessage{
			Type:    entity.WSMessageTypeNewMessage,
			Payload: newMessage,
		}
	}

	currentWsConnection.Close()
	ws.logger.Debugf("Connection closed for %d", authUserId)
	ws.delWsConn(connectionId)
	ws.delUserConn(authUserId, connectionId)
}

func (ws *WS) sendHTTP500(w http.ResponseWriter, message string, err error) {
	ws.logger.Errorf("%s: %v", message, err)
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, message)
}

func (ws *WS) addWsConn(sID string, wsConn *wsConnection) {
	ws.wsConnMapMu.Lock()
	ws.wsConnMap[sID] = wsConn
	ws.logger.Debug(ws.wsConnMap)
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

func (ws *WS) addUserConn(uID uint32, connectionId string) {
	ws.userConnMapMu.Lock()
	if connectionsList, ok := ws.userConnMap[uID]; ok {
		connectionsList[connectionId] = struct{}{}
	} else {
		ws.userConnMap[uID] = map[string]struct{}{
			connectionId: struct{}{},
		}
	}
	ws.logger.Debug(ws.userConnMap)
	ws.userConnMapMu.Unlock()
}

func (ws *WS) delUserConn(uID uint32, connectionId string) {
	ws.userConnMapMu.Lock()
	if connectionsList, ok := ws.userConnMap[uID]; ok {
		connectionsList[connectionId] = struct{}{}
		delete(connectionsList, connectionId)
	}
	ws.userConnMapMu.Unlock()
}

func (ws *WS) getUserConn(uID uint32) (userConn map[string]struct{}, ok bool) {
	ws.userConnMapMu.RLock()
	userConn, ok = ws.userConnMap[uID]
	ws.userConnMapMu.RUnlock()
	return
}
