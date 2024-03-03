package lib

import (
	"fmt"
	"sync"
	"time"
)

type Client struct {
	pongWait    time.Duration
	readTimeout time.Duration
	ForwardUrl  string
	wg          *sync.WaitGroup
	PingPeriod  time.Duration
}

func NewClient(port string, wg *sync.WaitGroup) (*Client, error) {
	forwardUrl := fmt.Sprintf("ws://localhost:%s", port)

	pongWait := time.Second * 3

	return &Client{

		pongWait:    pongWait,
		ForwardUrl:  forwardUrl,
		PingPeriod:  time.Second * 60,
		wg:          wg,
		readTimeout: pongWait * 2,
	}, nil
}
