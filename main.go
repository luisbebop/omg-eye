package main

import (
	"net/http"
	"log"
	"code.google.com/p/go.net/websocket"
	"os/exec"
	"bufio"
)

func main () {
	log.Printf("omg-eye listening websocket and http on port 8080")
	log.Printf("omg-eye S2 U")
	http.Handle("/see", websocket.Handler(See))
	http.Handle("/", http.FileServer(http.Dir("./html")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("omg-eye listen error on http port 8080: %v", err.Error())
	}
}

func See(ws *websocket.Conn) {
	log.Printf("see %v\n", ws.Config())

	// filtering the log by company name ..
	// tail -f log.txt | grep '\[   \]'
	
	//streaming all log
	cmd := exec.Command("tail", "-f", "/var/log/omg-switch.log")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("err:%v\n", err)
		return
	}
	
	if err := cmd.Start(); err != nil {
		log.Printf("err:%v\n", err)
		return
	}
	
	reader := bufio.NewReader(stdout)
	
	for {
		s, err := reader.ReadString('\n')
		err = websocket.Message.Send(ws, s)
		if err != nil {
			log.Printf("%s", err)
			break
		}		
		log.Printf("send:%s", s)
	}
}