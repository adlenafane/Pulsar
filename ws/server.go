package ws

import (
	// System
	"encoding/json"
	"log"

	// Remote
	"golang.org/x/net/websocket"

	// Local
	"github.com/samuelramond/Pulsar/model"
)

type DashboardClient struct {
	ws     *websocket.Conn
	Galaxy *model.Galaxy
}

type WsServer struct {
	Gc      *model.GalaxyCluster
	Clients map[*DashboardClient]int
}

type ActionMessage struct {
	Action string  `json:"action"`
	Data   string  `json:"data"`
	Cx     float64 `json:"cx"`
	Cy     float64 `json:"cy"`
}

func (this *WsServer) Server(ws *websocket.Conn) {
	log.Println("WsServer: New client connected:", ws.Request().RemoteAddr)

	var err error
	var rcvBuffer []byte

	defer func() {
		if err = ws.Close(); err != nil {
			log.Println("WsServer: Websocket could not be closed", err.Error())
		}
	}()

	client := &DashboardClient{ws, nil}
	if this.Clients == nil {
		this.Clients = make(map[*DashboardClient]int)
	}
	this.Clients[client] = 0

	for {
		if err = websocket.Message.Receive(client.ws, &rcvBuffer); err != nil {
			log.Println("WsServer: Client Disconnected....", err.Error())
			delete(this.Clients, client)
			return
		}
		var am ActionMessage
		json.Unmarshal(rcvBuffer, &am)
		log.Println("WsServer: Message received:", am)
		this.ProcessAction(&am, client)
	}
}

func (this *WsServer) ProcessAction(am *ActionMessage, client *DashboardClient) {
	if am.Action == "joingalaxy" {
		client.Galaxy = nil
		if client.Galaxy = this.Gc.Find(am.Data); client.Galaxy != nil {
			log.Println("WsServer: Client joined Galaxy")
		}
	}
}

func (this *WsServer) BroadcastMessage(cgalaxy *model.Galaxy, cpulsar *model.Pulsar) {
	//@todo(sam): wrap with a worker pool (scale issue)
	am := ActionMessage{"pulse", cpulsar.ClientId, cpulsar.Cx, cpulsar.Cy}
	js, err := json.Marshal(am)
	if err != nil {
		log.Println("WsServer: Unable to marshal message", err.Error())
		return
	}
	for cl, _ := range this.Clients {
		if cgalaxy == cl.Galaxy {
			if err := websocket.JSON.Send(cl.ws, js); err != nil {
				log.Println("Could not send message", err.Error())
			}
		}
	}
}
