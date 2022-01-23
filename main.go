package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ayupov-ayaz/loki-test/service"

	"github.com/google/uuid"
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

type History struct {
	Id  string `json:"id"`
	Bet int    `json:"bet"`
	Win int    `json:"win"`
}

func pushMessage(wg *sync.WaitGroup, loki *service.Loki) {
	m := Message{}

	for i := 0; i < count; i++ {
		m.Id = i
		m.Message = "go number = " + strconv.Itoa(i)

		if err := loki.Push(m); err != nil {
			panic(err)
		}
	}

	wg.Done()
}

func pushHistory(wg *sync.WaitGroup, loki *service.Loki) {
	h := History{}

	for i := 0; i < count; i++ {
		h.Id = uuid.New().String()
		h.Bet = i * 5352
		h.Win = i * 343

		if err := loki.Push(h); err != nil {
			panic(err)
		}
	}

	wg.Done()
}

func readMessage(loki *service.Loki) error {
	limit := 1
	query := `|="message" `

	for i := 0; i < count; i++ {
		start := time.Now().Add(-1 * time.Second)
		end := time.Now().Add(1 * time.Second)

		list, err := loki.Read(start, end, query, limit)
		if err != nil {

			if errors.Is(err, service.ErrNotFound) {
				break
			}

			return err
		}

		for j := 0; j < len(list); j++ {
			fmt.Println("time=", list[j][0], "message = ", list[j][1])
		}
	}

	return nil
}

func main() {
	loki := service.NewLoki(baseUrl, app)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go pushMessage(wg, loki)
	go pushHistory(wg, loki)

	wg.Wait()

	if err := readMessage(loki); err != nil {
		panic(err)
	}
}
