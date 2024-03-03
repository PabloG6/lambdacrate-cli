package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"lambdacrate-cli/lib"
	"lambdacrate-cli/lib/webhook"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Proxy struct {
	ServerUrl string
	appID     string

	client      *lib.Client
	readTimeout time.Duration
	mu          sync.Mutex
	wg          *sync.WaitGroup
	CloseChan   chan os.Signal
	Conn        *websocket.Conn
}

func NewProxy(appID string, client *lib.Client) (proxy *Proxy) {
	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	serverUrl, err := getServerUrlFromAppID(appID)

	if err != nil {
	}

	var wg sync.WaitGroup
	return &Proxy{
		wg:          &wg,
		appID:       appID,
		ServerUrl:   serverUrl,
		client:      client,
		readTimeout: time.Minute * 10,
		CloseChan:   sigint,
	}
}

func (proxy *Proxy) Terminate() {
	proxy.CloseChan <- syscall.SIGINT
}
func (proxy *Proxy) Run() error {

	maxAttempts := 10
	attempt := 0
	msg := fmt.Sprintf("attempting to connect to remote server for app attempt: %d", attempt+1)
	log.Println(msg)
	conn, _, err := websocket.DefaultDialer.Dial(proxy.ServerUrl, nil)
	for attempt < maxAttempts && err != nil {
		msg = fmt.Sprintf("attempting to connect to remote server for app attempt: %d", attempt+1)
		log.Println(err)
		conn, _, err = websocket.DefaultDialer.Dial(proxy.ServerUrl, nil)
		attempt++
		if err != nil {
			log.Println(msg)
		} else {
			break
		}

		time.Sleep(time.Second * 3)

	}
	if err != nil {
		log.Println("Failed to connect to remote application with appID...")
		return err
	}

	conn.SetPongHandler(func(appData string) error {
		log.Println("pong handler: ", appData)
		err := conn.SetWriteDeadline(time.Now().Add(proxy.client.PingPeriod))
		return err
	})
	proxy.Conn = conn
	//keep processes alive with wait group
	proxy.wg.Add(2)

	go proxy.writePump()
	go proxy.readPump()

	proxy.wg.Wait()
	return nil
}

func (proxy *Proxy) readPump() {
	defer proxy.wg.Done()
	log.Println(proxy.readTimeout)
	for {
		err := proxy.Conn.SetReadDeadline(time.Now().Add(proxy.readTimeout))
		if err != nil {
			log.Println("failed to set read deadline for pump: ", err)

		}
		messageType, message, err := proxy.Conn.ReadMessage()
		if err != nil {
			//failed to read message.
		}

		switch messageType {
		case websocket.PongMessage:
			{
				log.Println("pong message received")
			}

		case websocket.TextMessage:
			{
				err = json.Unmarshal(message, &webhook.HttpRequest{})
				if err != nil {
					log.Println("failed to unmarshall json response: ", err)
				}
			}

		case websocket.CloseMessage:
			{
				log.Println("sent a close message")
				proxy.Terminate()
			}
		}

	}

}
func (proxy *Proxy) writePump() {

	defer proxy.wg.Done()
	ticker := time.NewTicker(proxy.client.PingPeriod / 2)
	defer ticker.Stop()

	for {

		waitTime := time.Now().Add(proxy.client.PingPeriod)
		log.Println(proxy.client.PingPeriod)
		err := proxy.Conn.SetWriteDeadline(waitTime)
		log.Println("starting up for loop again. ")
		if err != nil {
			log.Println("failed to set read deadline for your application")
		}

		select {
		case <-ticker.C:
			{
				log.Println("pinging remote server to keep connection alive...")
				pingMessage := webhook.HttpResponse{MessageType: websocket.PingMessage}
				data, err := json.Marshal(pingMessage)
				if err != nil {
					log.Println("failed to generate ping message data")
				}
				err = proxy.Conn.WriteMessage(websocket.PingMessage, data)
				if err != nil {
					log.Println("failed to ping server, terminating client connection...", err)
					proxy.CloseChan <- syscall.SIGINT
					return
				}
			}

		}

	}

}

func getServerUrlFromAppID(appID string) (string, error) {
	//todo please fix this to actually resolve to an active websocket.
	return "ws://localhost:8000/ws", nil
}
