package stats

import (
	// Sys
	"log"
	"os"
	"net/url"
	"fmt"
	"time"
	"encoding/json"

	// Misc
	"github.com/influxdb/influxdb/client"

	// Local
	"github.com/samuelramond/Pulsar/model"
)

type AggregatedHit struct {
	Hit 	int64 	`json:"hit"`
	Time	string	`json:"time"`
}

type AggregateStats []*AggregatedHit

type GalacticStats struct {
	dbName		string
	con 		*client.Client
}

func (this *GalacticStats) Add(cur_galaxy *model.Galaxy, cur_pulsar *model.Pulsar) {
	log.Println("GalacticStats: Adding new hit to", cur_galaxy.Token, "by", cur_pulsar.ClientId)
	var pts = make([]client.Point, 1)
	pts[0] = client.Point{
        Measurement: cur_galaxy.Name,
        Tags: map[string]string{
            "pulsar": cur_pulsar.ClientId,
        },
        Fields: map[string]interface{}{
            "value": 1,
        },
        Time: time.Now(),
        Precision: "s",
    }
    bps := client.BatchPoints{
        Points:          pts,
        Database:        this.dbName,
        RetentionPolicy: "default",
    }
    _, err := this.con.Write(bps)
    if err != nil {
        log.Fatal(err)
    }
}

func (this *GalacticStats) Get(cur_galaxy *model.Galaxy, group_by string, past string) *AggregateStats {
	log.Println("GalacticStats: Getting hit for", cur_galaxy.Token)
	if group_by == "" {
		group_by = "10m"
	}
	if past == "" {
		past = "8h"
	}
	query := fmt.Sprintf("SELECT COUNT(value) FROM %s WHERE time > now() - %s GROUP BY time(%s)", cur_galaxy.Name, past, group_by)
	res, _ := this.queryDB(query)
	log.Println("GalacticStats: query", query)
	if res == nil {
		log.Println("TODO: Handle error")
		return nil		
	}
	var payload AggregateStats

    for _, row := range res[0].Series[0].Values {
        t, err := time.Parse(time.RFC3339, row[0].(string))
        if err != nil {
            log.Fatal(err)
	    }
	    var val int64
	    val = 0
	    if row[1] != nil {
	        val, err = row[1].(json.Number).Int64()
	    }
	    s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()) 
        payload = append(payload, &AggregatedHit{val, s})
    } 
	return &payload
}

func (this *GalacticStats) Load(cur_galaxy *model.Galaxy) {
	log.Println("GalacticStats: Loading tags", cur_galaxy.Token)
	query := fmt.Sprintf("SHOW TAG VALUES FROM %s WITH KEY = pulsar", cur_galaxy.Name)
	res, _ := this.queryDB(query)
	log.Println("GalacticStats: query", query)
	for _, row := range res[0].Series[0].Values {
		if str, ok := row[0].(string); ok {
			cur_galaxy.Find(str, "")
		}
		
	}
}

var con *client.Client

func (this *GalacticStats) InitDB() {
	if this.con != nil {
		return
	}

	log.Println("GalacticStats: Connecting to influxdb Cluster")

	this.dbName = "test" //Todo	: replace with proper config

	//-- DB
    u, err := url.Parse(fmt.Sprintf("http://127.0.0.1:8086"))
    if err != nil {
        log.Fatal(err)
    }

    conf := client.Config{
        URL:      *u,
        Username: os.Getenv("INFLUX_USER"),
        Password: os.Getenv("INFLUX_PWD"),
    }

    this.con, err = client.NewClient(conf)
    if err != nil {
        log.Fatal(err)
    }

    dur, ver, err := this.con.Ping()
    if err != nil {
        log.Fatal("GalacticStats:", err)
    }	
    log.Printf("Ping Influx OK %v, %s", dur, ver)	
}

func (this *GalacticStats) queryDB(cmd string) (res []client.Result, err error) {
    q := client.Query{
        Command:  cmd,
        Database: this.dbName,
    }
    if response, err := this.con.Query(q); err == nil {
        if response.Error() != nil {
            return res, response.Error()
        }
        res = response.Results
    }
    return
}
