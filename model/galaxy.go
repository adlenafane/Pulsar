package model

import (
	"log"
	"math/rand"
)

//-------------------------Galaxy----------------------------------------------

type Galaxy struct {
	Name		string
	Description	string
	Token		string
	Pulsars		map[string]*Pulsar
	Nebulas		map[string]*Nebula
}

func (this *Galaxy) AddPulsar(client_id string) *Pulsar {
	log.Println("Galaxy: Adding new client_id", client_id)
	new_pulsar := &Pulsar{client_id, 0, rand.Float64(), rand.Float64()}
	this.Pulsars[client_id] = new_pulsar
	return new_pulsar
}

func (this *Galaxy) AddNebula(group_id string) *Nebula {
	log.Println("Galaxy: Adding new group_id", group_id)
	new_nebula := &Nebula{group_id, 0, nil}
	new_nebula.Pulsars = make(map[string]*Pulsar)
	this.Nebulas[group_id] = new_nebula
	return new_nebula
}

func (this *Galaxy) Find(client_id string, group_id string) (*Pulsar, *Nebula) {
	cclient := this.Pulsars[client_id]
	if cclient == nil {
		log.Println("Galaxy: Unable to locate client")
		cclient = this.AddPulsar(client_id)
	}
	cnebular := this.Nebulas[group_id]
	if cnebular == nil {
		cnebular = this.AddNebula(group_id)
	}
	cnebular.Register(cclient)
	return cclient, cnebular
}

//-------------------------GalaxyCluster---------------------------------------

type GalaxyCluster struct {
	galaxies 	map[string]*Galaxy
}

func (this *GalaxyCluster) Add(new_galaxy *Galaxy) {
	log.Println("GalaxyCluster: Adding new galaxy")
	if this.galaxies == nil {
		this.galaxies = make(map[string]*Galaxy)
	}
	new_galaxy.Pulsars = make(map[string]*Pulsar)
	new_galaxy.Nebulas = make(map[string]*Nebula)
	this.galaxies[new_galaxy.Token] = new_galaxy
}

func (this *GalaxyCluster) Find(token string) *Galaxy {
	cur_galaxy := this.galaxies[token]
	if cur_galaxy == nil {
		log.Println("GalaxyCluster: Unable to locate galaxy")
		return nil
	}
	return cur_galaxy
}
