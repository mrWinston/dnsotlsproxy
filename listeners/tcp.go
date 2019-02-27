package listeners

import (
	"fmt"
	"net"

	"github.com/mrWinston/dnsotlsproxy/resolver"
	"github.com/sirupsen/logrus"
)

type TcpListener struct {
	listener    net.Listener
	stop        chan bool
	stopped     bool
	dnsUpstream *resolver.Resolver
}

func NewTcpListener(ipString string, port string, resolver *resolver.Resolver) (*TcpListener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ipString, port))
	if err != nil {
		logError(err, "error Initializing TCP Listener")
		return nil, err
	}
	tcpListener := &TcpListener{
		listener:    listener,
		stop:        make(chan bool),
		stopped:     false,
		dnsUpstream: resolver,
	}
	go tcpListener.listenForMessages()
	return tcpListener, nil
}

func (tcpListener *TcpListener) listenForMessages() {
	for {
		conn, err := tcpListener.listener.Accept()

		if err == nil {
			go tcpListener.handleRequest(conn)

		} else {

			if tcpListener.stopped {
				logrus.Info("Ending Receive Loop")
				return
			} else {
				logError(err, "error Getting connection")
			}

		}
	}
}

func (tcpListener *TcpListener) handleRequest(conn net.Conn) {
	defer conn.Close()
	logrus.Info("Got a Tcp Request")

	buf := make([]byte, 65535)

	reqLen, err := conn.Read(buf)
	if err != nil {
		logError(err, "Error reading DNS Request")
	}
	n, tlsbuf, err := tcpListener.dnsUpstream.ForwardDns(buf[:reqLen])

	if err != nil {
		logError(err, "Error getting DNS over Tls")
	}

	conn.Write(tlsbuf[:n])

}

func (tcpListener *TcpListener) Shutdown() {
	logrus.Info("Shutting down Listener")
	tcpListener.stopped = true
	tcpListener.listener.Close()
	logrus.Info("Closed Listener")
	logrus.Info("Sent Shutdown Signal")
	close(tcpListener.stop)
}
