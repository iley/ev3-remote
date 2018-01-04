package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type JsonCommand struct {
	Control string  `json:"control"`
	Id      int     `json:"id"`
	Value   float64 `json:"value"`
}

func Verbose() bool {
	return false
}

func MakeWebsocketServer() func(*websocket.Conn) {
	robot, err := MakeTractor()
	if err != nil {
		panic(err)
	}
	return func(ws *websocket.Conn) {
		defer robot.Shutdown()
		for {
			var message []byte
			if err := websocket.Message.Receive(ws, &message); err != nil {
				log.Printf("Error receiving: %s", err.Error())
			}
			var command JsonCommand
			if err := json.Unmarshal(message, &command); err != nil {
				log.Printf("Error unmarshalling: %s", err.Error())
			}
			if Verbose() {
				log.Printf("Control=%s ID=%d value=%f", command.Control,
					command.Id, command.Value)
			}
			if command.Control == "button" {
				pressed := true
				if command.Value == 0 {
					pressed = false
				}
				robot.HandleButton(command.Id, pressed)
			} else if command.Control == "axis" {
				robot.HandleAxisChange(command.Id, command.Value)
			}
		}
	}
}

func PingLoop(port int, agentUrl, secret string) {
	for {
		err := NotifyAgent(port, agentUrl, secret)
		if err != nil {
			log.Printf("Couldn't notify agent %s: %s", agentUrl, err.Error())
		} else {
			log.Printf("Notified agent at %s", agentUrl)
		}
		time.Sleep(time.Minute)
	}
}

func main() {
	port := flag.Int("port", 80, "port to listen")
	agentUrl := flag.String("agent", "", "URL of redirect agent")
	secret := flag.String("secret", "", "secret for redirect agent")
	flag.Parse()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	log.Printf("Main directory: %s", dir)
	if *agentUrl != "" {
		if *secret == "" {
			fmt.Fprintln(os.Stderr, "you must provide --secret to use --agent")
			os.Exit(1)
		}
		err := NotifyAgent(*port, *agentUrl, *secret)
		if err != nil {
			log.Printf("Couldn't notify agent %s: %s", *agentUrl, err.Error())
		} else {
			log.Printf("Notified agent at %s", *agentUrl)
		}
		go PingLoop(*port, *agentUrl, *secret)
	}
	fs := http.FileServer(http.Dir(dir + "/static"))
	http.Handle("/", http.RedirectHandler("/static/index.html", 302))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/websocket", websocket.Handler(MakeWebsocketServer()))
	addr := fmt.Sprintf(":%d", *port)
	log.Println(fmt.Sprintf("Listening at %s...", addr))
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
