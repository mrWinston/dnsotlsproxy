package listeners

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/mrWinston/dnsotlsproxy/resolver"
	"github.com/sirupsen/logrus"
)

type UdpListener struct {
	udpConn     *net.UDPConn
	stop        chan bool
	dnsUpstream *resolver.Resolver
}

func NewUdpListener(ipString string, port int, resolver *resolver.Resolver) (*UdpListener, error) {
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipString, port))
	if err != nil {
		logError(err, "error parsing ip and/or port")
		return nil, err
	}

	conn, err := net.ListenUDP("udp", address)

	if err != nil {
		logError(err, "Error listening for UDP Datagrams")
		return nil, err
	}

	listener := &UdpListener{
		udpConn:     conn,
		stop:        make(chan bool),
		dnsUpstream: resolver,
	}
	go listener.listenForMessages()
	return listener, nil
}

func (udpListener *UdpListener) listenForMessages() {
	for {
		select {
		case <-udpListener.stop:
			logrus.Info("Ending Receive Loop")
			return
		default:
			buf := make([]byte, 512)
			n, addr, err := udpListener.udpConn.ReadFromUDP(buf)

			sizeBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(sizeBytes, uint16(n))
			buf = buf[:n]                   //truncate the buffer
			buf = append(sizeBytes, buf...) //add size

			if err != nil {
				logError(err, "Error receiving UDP Datagram")
				break
			}

			_, dnsRes, err := udpListener.dnsUpstream.ForwardDns(buf)

			if err != nil {
				logError(err, "Error forwarding dns")
				break
			}
			strippedRes := stripTrailingNull(dnsRes[2:])
			_, _, err = udpListener.udpConn.WriteMsgUDP(strippedRes[:len(strippedRes)-4], nil, addr)
			if err != nil {
				logError(err, "Error sending UDP Datagram back to client")
				break
			}
		}

	}
}

func (udpListener *UdpListener) Shutdown() {
	logrus.Info("Shutting down Listener")
	udpListener.stop <- true
	udpListener.udpConn.Close()
	close(udpListener.stop)
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
