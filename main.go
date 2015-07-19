package main

import (
	// System
	"encoding/json"
	"log"
	"net/http"

	// Remote
	"golang.org/x/net/websocket"

	// Local
	"github.com/samuelramond/Pulsar/model"
	"github.com/samuelramond/Pulsar/stats"
	"github.com/samuelramond/Pulsar/ws"
)

const (
	listenAddr = "localhost:4000" // server address
)

func PulseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}

	token := r.URL.Query().Get("token")
	client_id := r.URL.Query().Get("client_id")
	group_id := r.URL.Query().Get("group_id")
	log.Println("HTTP: +------------------------------+")
	log.Println("HTTP: |	Token :", token)
	log.Println("HTTP: |	Client:", client_id)
	log.Println("HTTP: |	Group :", group_id)
	log.Println("HTTP: +------------------------------+")
	cur_galaxy := ss.Gc.Find(token)
	if cur_galaxy == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	pulsar, nebula := cur_galaxy.Find(client_id, group_id)
	pulsar.Hit += 1
	nebula.Hit += 1
	st.Add(cur_galaxy, pulsar)
	ss.BroadcastMessage(cur_galaxy, pulsar)
	//@Todo(sam): Broadcast event to cluster (RPC)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}

	token := r.URL.Query().Get("token")
	group_by := r.URL.Query().Get("groupBy")
	past := r.URL.Query().Get("past")
	if token == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cur_galaxy := ss.Gc.Find(token)
	if cur_galaxy == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Println("HTTP: +------------------------------+")
	log.Println("HTTP: |	Token :", token)
	log.Println("HTTP: |	GrouBy:", group_by)
	log.Println("HTTP: |	Past  :", past)
	log.Println("HTTP: +------------------------------+")
	as := st.Get(cur_galaxy, group_by, past)

	js, err := json.Marshal(as)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func InitStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cur_galaxy := ss.Gc.Find(token)
	if cur_galaxy == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Println("HTTP: +------------------------------+")
	log.Println("HTTP: |	Token :", token)
	log.Println("HTTP: +------------------------------+")
	st.Load(cur_galaxy)

	js, err := json.Marshal(struct {
		Name    string                   `json:"name"`
		Pulsars map[string]*model.Pulsar `json:"pulsars"`
	}{cur_galaxy.Name, cur_galaxy.Pulsars})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

var ss *ws.WsServer
var st *stats.GalacticStats

func main() {
	ss = &ws.WsServer{}
	st = &stats.GalacticStats{}
	st.InitDB()

	gc := &model.GalaxyCluster{}

	// @todo(sam): Add a configuration backend
	gc.Add(&model.Galaxy{
		Name:        "ProductStream",
		Description: "Alkemics platform https://stream.alkemics.com",
		Token:       "7c851e6489f0900a9c9cf80d55079e8beb3ef63a",
	})
	gc.Add(&model.Galaxy{
		Name:        "GDSN",
		Description: "Alkemics platform GDSN feed",
		Token:       "b68e070d2e229cd66aca21dd2c4eb488e00d6346",
	})
	gc.Add(&model.Galaxy{
		Name:        "Corsair",
		Description: "Corsair.space Pulsar",
		Token:       "1b8463c5bcf1afef55874524157a91dc22576b6e",
	})
	ss.Gc = gc

	// Http Handler
	// @todo(sam): Add proper routing system
	http.HandleFunc("/pulse", PulseHandler)
	http.HandleFunc("/stats", StatsHandler)
	http.HandleFunc("/stats/load", InitStatsHandler)

	// WebSocket Handler
	http.Handle("/sock", websocket.Handler(ss.Server))

	log.Println("Pulsar: Starting server...")
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
