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
	stopped     bool
	dnsUpstream *resolver.Resolver
}

func NewUdpListener(ipString string, port string, resolver *resolver.Resolver) (*UdpListener, error) {
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", ipString, port))
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
		stopped:     false,
		dnsUpstream: resolver,
	}
	go listener.listenForMessages()
	return listener, nil
}

func (udpListener *UdpListener) listenForMessages() {
	for {
		buf := make([]byte, 512)
		n, addr, err := udpListener.udpConn.ReadFromUDP(buf)

		if err != nil {
			if udpListener.stopped {
				logrus.Info("Shutting down Udp Listen loop")
				return
			} else {
				logError(err, "Error during ReadFromUDP")
				break
			}
		}
		logrus.Info("Got a Udp Request")
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
		//		strippedRes := stripTrailingNull(dnsRes[2:], 0) //stip trailing zeroes from end and remove lenght from beginning
		_, _, err = udpListener.udpConn.WriteMsgUDP(dnsRes[2:], nil, addr)
		if err != nil {
			logError(err, "Error sending UDP Datagram back to client")
			break
		}

	}
}

func (udpListener *UdpListener) Shutdown() {
	logrus.Info("Shutting down Listener")
	udpListener.stopped = true
	udpListener.udpConn.Close()
}
