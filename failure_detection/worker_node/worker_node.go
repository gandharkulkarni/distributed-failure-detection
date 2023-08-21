package main

import (
	"failure_detection/msg_handler"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		panic("Insufficient arguments given")
	}
	var host string
	var worker string
	if strings.Contains(os.Args[1], ":") {
		host = os.Args[1]
	} else {
		panic("Invalid hostname. Please specify port [hostname:port]")
	}

	if strings.Contains(os.Args[2], ":") {
		worker = os.Args[2]
	} else {
		panic("Invalid hostname for worker node. Please specify port [hostname:port]")
	}
	registerNode(host, worker)
	listener, err := net.Listen("tcp", ":"+strings.Split(worker, ":")[1])
	checkErr(err)
	for {

		if conn, err := listener.Accept(); err == nil {
			msgHandler := msg_handler.NewMsgHandler(conn)
			msg, err := msgHandler.Receive()
			checkErr(err)
			if !msg.GetRegistration() {
				fmt.Println("Reregistering...")
				registerNode(host, worker)
			}
			if !msg.GetHeartbeat() {
				announceLiveness(msgHandler, worker)
			}
			msgHandler.Close()
			conn.Close()
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
func announceLiveness(msgHandler *msg_handler.MsgHandler, workerNodeName string) {
	msgHandler.Send(&msg_handler.NodeMsg{Hostname: workerNodeName, Heartbeat: true})
	fmt.Println("Liveness announced")
	msgHandler.Close()
	fmt.Println()
}
func registerNode(host string, worker string) {
	fmt.Println("Connecting to ", host)
	conn, err := net.Dial("tcp", host)
	checkErr(err)
	fmt.Println("Connection successful")
	msgHandler := msg_handler.NewMsgHandler(conn)
	registerNodeWithController(msgHandler, worker)
	msgHandler.Close()
	conn.Close()
	fmt.Println("Registration successful")
}
func registerNodeWithController(msgHandler *msg_handler.MsgHandler, workerNodeName string) {
	msgHandler.Send(&msg_handler.NodeMsg{Hostname: workerNodeName, Heartbeat: true})
	fmt.Println("Hostname registered in topology")
}
