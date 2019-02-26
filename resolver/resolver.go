package resolver

import (
	"crypto/tls"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Resolver struct {
	RemoteIp   string
	RemotePort int
}

// ForwardDns takes a dns request in TCP format, forwards it to a dns over TLS
// service and returns the number of bytes in the answer and the answer or an
// error if something goes wrong
func (resolver *Resolver) ForwardDns(request []byte) (int, []byte, error) {
	logrus.Debug("Initiating TLS Connection")

	conn, err := tls.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", resolver.RemoteIp, resolver.RemotePort),
		&tls.Config{},
	)

	if err != nil {
		logError(err, "Error initiaiting TLS connection")
		return 0, nil, err
	}
	defer conn.Close()

	buf := make([]byte, 65535)

	_, err = conn.Write(request)

	if err != nil {
		logError(err, "Error sending buffer to dns service")
		return 0, nil, err
	}

	n, err := conn.Read(buf)

	if err != nil {
		logError(err, "Error receiving Dns Response")
		return 0, nil, err
	}
	logrus.WithFields(logrus.Fields{
		"buf": buf[:n],
	}).Debug("Got a buffer back")

	return n, buf[:n], nil
}

func logError(err error, msg string) {
	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Error(msg)
}
