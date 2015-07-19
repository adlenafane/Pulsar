package model

import (
	"log"
)

type Pulsar struct {
	ClientId string `json:"data"`
	Hit      int
	Cx       float64 `json:"cx"`
	Cy       float64 `json:"cy"`
}

type Nebula struct {
	GroupId string
	Hit     int
	Pulsars map[string]*Pulsar
}

func (this *Nebula) Register(pulsar *Pulsar) {
	cpulsar := this.Pulsars[pulsar.ClientId]
	if cpulsar == nil {
		log.Println("Nebula: Pulsar", pulsar.ClientId, "joining Nebula", this.GroupId)
		this.Pulsars[pulsar.ClientId] = pulsar
	}
}
