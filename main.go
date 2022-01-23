package main

import (
	"loki-test/service"
	"strconv"
)

const (
	app     = "loki-test"
	baseUrl = "http://localhost:3100"
	count   = 1000
)

type Message struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

func main() {
	loki := service.NewLoki(baseUrl, app)

	m := Message{}

	for i := 0; i < count; i++ {
		m.Id = i
		m.Message = "go number = " + strconv.Itoa(i)

		if err := loki.Push(m); err != nil {
			panic(err)
		}
	}
}
