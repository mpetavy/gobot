package main

import (
	"embed"
	"fmt"
	"github.com/mpetavy/common"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"
	"log"
	"time"
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "gobot", "", "", "", &resources, nil, nil, run, 0)
}

func run() error {
	mqttAdaptor := mqtt.NewAdaptor("tcp://broker.hivemq.com:1883", "pinger1")

	work := func() {
		mqttAdaptor.On("hello", func(msg mqtt.Message) {
			fmt.Printf("%+v\n", msg.Payload())
		})

		data := []byte("Marcel")
		gobot.Every(1*time.Second, func() {
			mqttAdaptor.Publish("hello", data)
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		work,
	)

	robot.AutoRun = false

	var err error

	go func() {
		err = robot.Start()
	}()

	time.Sleep(time.Second)

	if common.Error(err) {
		return err
	}

	time.Sleep(time.Second * 5)

	log.Printf("stopping server ...\n")
	err = robot.Stop()
	if common.Error(err) {
		return err
	}
	log.Printf("server stopped!\n")

	return nil
}

func main() {
	common.Run(nil)
}
