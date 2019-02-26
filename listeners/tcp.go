package listeners

import (
	"fmt"
	"net"

	"github.com/mrWinston/dnsotlsproxy/resolver"
	"github.com/sirupsen/logrus"
)

type TcpListener struct {
	listener    *net.Listener
	stop        chan bool
	dnsUpstream *resolver.Resolver
}

func NewTcpListener(ipString string, port int, resolver *resolver.Resolver) (*TcpListener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ipString, port))
	if err != nil {
		logError(err, "error Initializing TCP Listener")
		return nil, err
	}
	tcpListener := &TcpListener{
		listener:    &listener,
		stop:        make(chan bool),
		dnsUpstream: resolver,
	}
	go tcpListener.listenForMessages()
	return tcpListener, nil
}

func (tcpListener *TcpListener) listenForMessages() {
	for {
		select {
		case <-tcpListener.stop:
			logrus.Info("Ending Receive Loop")
			return
		default:
			_, err := tcpListener.listener.Accept()
			if err != nil {
				logError(err, "error Getting connection")

				break
			}

		}

	}
}

func (tcpListener *TcpListener) Shutdown() {
	logrus.Info("Shutting down Listener")
	tcpListener.stop <- true
	tcpListener.udpConn.Close()
	close(tcpListener.stop)
}

func logError(err error, msg string) {
	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Error(msg)
}

func stripTrailingNull(buf []byte) []byte {
	for i := len(buf) - 1; i >= 0; i-- {
		if buf[i] != 0x00 {
			return buf[:i+1]
		}
	}
	return make([]byte, 0)
}
