package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"time"
)

const (
	NODE_TOPIC     = "test/node"
	FLOGO_TOPIC    = "test/flogo"
	MAINFLUX_TOPIC = "channels/"
)

type Packet struct {
	AccX  json.Number `json:"accX"`
	AccY  json.Number `json:"accY"`
	AccZ  json.Number `json:"accZ"`
	GyroX json.Number `json:"gyroX"`
	GyroY json.Number `json:"gyroY"`
	GyroZ json.Number `json:"gyroZ"`
}

func (p *Packet) mock() {
	p.AccX = json.Number(fmt.Sprintf("%f", rand.Float64()))
	p.AccY = json.Number(fmt.Sprintf("%f", rand.Float64()))
	p.AccZ = json.Number(fmt.Sprintf("%f", rand.Float64()))
	p.GyroX = json.Number(fmt.Sprintf("%f", rand.Float64()))
	p.GyroY = json.Number(fmt.Sprintf("%f", rand.Float64()))
	p.GyroZ = json.Number(fmt.Sprintf("%f", rand.Float64()))
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	p := Packet{}
	for {
		p.mock()
		js, _ := json.Marshal(p)
		message := string(js)
		c.Publish("test/topic", 0, false, message)
		time.Sleep(time.Millisecond * 500)
	}
}
