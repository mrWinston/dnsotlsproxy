package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0"+":"+"8053")
	if err != nil {
		log.Error("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + "0.0.0.0" + ":" + "8053")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	log.Info("Initiating TLS Conn")
	conf := &tls.Config{}
	tlsconn, err := tls.Dial("tcp", "1.1.1.1:853", conf)
	defer tlsconn.Close()

	n, err := tlsconn.Write(buf[:reqLen])
	if err != nil {
		log.Println(n, err)
		return
	}
	tlsbuf := make([]byte, 1000)
	n, err = tlsconn.Read(tlsbuf)
	if err != nil {
		log.Println(n, err)
		return
	}

	conn.Write(tlsbuf[:n])
}
