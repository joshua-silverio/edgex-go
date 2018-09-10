//
// Copyright (c) 2017 Cavium
//
// SPDX-License-Identifier: Apache-2.0
//

package distro

import (
	"encoding/json"
	"fmt"

	"github.com/edgexfoundry/edgex-go/pkg/models"
	zmq "github.com/pebbe/zmq4"

)

const (
	zeroMQPort = 5563
)

func ZeroMQReceiver(eventCh chan *models.Event) {
	go initZmq(eventCh)
}

func initZmq(eventCh chan *models.Event) {
	q, _ := zmq.NewSocket(zmq.SUB)
	defer q.Close()

	logger.Info("Connecting to zmq...")
	url := fmt.Sprintf("tcp://%s:%d", configuration.DataHost, zeroMQPort)
	q.Connect(url)
	logger.Info("Connected to zmq")
	q.SetSubscribe("")

	for {
		msg, err := q.RecvMessage(0)
		if err != nil {
			id, _ := q.GetIdentity()
			logger.Error("Error getting mesage", logger.String("id", id))
		} else {
			for _, str := range msg {
				event := parseEvent(str)
				logger.Info("Event received", logger.Any("event", event))
				eventCh <- event
			}
		}
	}
}

func parseEvent(str string) *models.Event {
	event := models.Event{}

	if err := json.Unmarshal([]byte(str), &event); err != nil {
		logger.Error("Failed to parse event", logger.Error(err))
		return nil
	}
	return &event
}
