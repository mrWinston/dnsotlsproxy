package main

import (
	"fmt"
	"net"
	"os"

	"github.com/mrWinston/dnsotlsproxy/listeners"
	"github.com/mrWinston/dnsotlsproxy/resolver"
	"github.com/sirupsen/logrus"
)

var shutdown = make(chan bool, 1)
var osSig = make(chan os.Signal, 1)

func handleOsSignal() {
	sig := <-osSig
	logrus.Info("Received ", sig)
	logrus.Info("Shutting Down")
	shutdown <- true
}

func main() {
	dnsUpstream := resolver.Resolver{
		RemoteIp:   "1.1.1.1",
		RemotePort: 853,
	}

	udpListener, err := listeners.NewUdpListener("0.0.0.0", 8053, &dnsUpstream)

	if err != nil {
		logError(err, "Error creating UDP Listener")
		return
	}
	defer udpListener.Shutdown()
	<-shutdown

}
func main2() {
	l, err := net.Listen("tcp", "0.0.0.0"+":"+"8053")
	if err != nil {
		logError(err, "Error starting TCP server")
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + "0.0.0.0" + ":" + "8053")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			logError(err, "Error listening for Messages")
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
		logError(err, "Error reading DNS Request")
	}
	dnsResolver := &resolver.Resolver{
		RemoteIp:   "1.1.1.1",
		RemotePort: 853,
	}

	n, tlsbuf, err := dnsResolver.ForwardDns(buf[:reqLen])

	if err != nil {
		logError(err, "Error getting DNS over Tls")
	}

	conn.Write(tlsbuf[:n])
}

func logError(err error, msg string) {
	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Error(msg)
}
