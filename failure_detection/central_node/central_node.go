package main

import (
	"failure_detection/msg_handler"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var nodeMap = map[string]int{}

func main() {
	if len(os.Args) < 2 {
		panic("Not enough arguments")
	}
	go registerNodes()
	probeLiveness()
}
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
func registerNodes() {
	listener, err := net.Listen("tcp", ":"+os.Args[1])
	checkErr(err)
	for {
		if conn, err := listener.Accept(); err == nil {
			msgHandler := msg_handler.NewMsgHandler(conn)
			handleNodeRegistration(msgHandler)
		}
	}
}
func handleNodeRegistration(msgHandler *msg_handler.MsgHandler) {
	msg, err := msgHandler.Receive()
	checkErr(err)
	_, ok := nodeMap[msg.GetHostname()]
	if !ok {
		nodeMap[msg.GetHostname()] = 0
	}
}
func updateLivenessMap(nodename string, heartbeat bool) {
	if nodeMap[nodename] == 3 || nodeMap[nodename] == -1 {
		nodeMap[nodename] = -1
		fmt.Println("Node deregistered..")
	}
	if heartbeat {
		fmt.Println(nodename, " is alive")
		nodeMap[nodename] = 0
	} else {
		if nodeMap[nodename] != -1 {
			nodeMap[nodename] += 1
			fmt.Println(nodename, "did not announce liveness : ", nodeMap[nodename], " times")
		}
	}

	return
}
func probeLiveness() {
	hostname, err := os.Hostname()
	checkErr(err)
	for {
		keys := make([]string, 0, len(nodeMap))
		for k := range nodeMap {
			keys = append(keys, k)
		}
		for _, node := range keys {

			fmt.Println("Requesting liveliness proof to ", node)
			conn, err := net.Dial("tcp", node)
			if err != nil {
				fmt.Println("No heartbeat from ", node)
				updateLivenessMap(node, false)
				continue
			}

			msgHandler := msg_handler.NewMsgHandler(conn)
			if nodeMap[node] == -1 {
				msgHandler.Send(&msg_handler.NodeMsg{Hostname: hostname, Heartbeat: false, Registration: false})
			} else {
				msgHandler.Send(&msg_handler.NodeMsg{Hostname: hostname, Heartbeat: false, Registration: true})
			}
			msg, err := msgHandler.Receive()
			checkErr(err)

			updateLivenessMap(node, msg.GetHeartbeat())

			fmt.Println()

		}
		time.Sleep(5 * time.Second)
	}

}
