package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)
	go handleOsSignal()

	settings := getSettingsFromEnvVars()

	logrus.SetLevel(settings.LogLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.Info(fmt.Sprintf("Starting DnsProxy with these settings: %v", settings))

	dnsUpstream := resolver.Resolver{
		RemoteIp:   settings.DnsUpstreamAddress,
		RemotePort: settings.DnsUpstreamPort,
	}

	udpListener, err := listeners.NewUdpListener(
		settings.UdpAddress,
		settings.UdpPort,
		&dnsUpstream,
	)
	if err != nil {
		logError(err, "Error creating UDP Listener")
		return
	}

	logrus.Info("Started UDP Listener")

	tcpListener, err := listeners.NewTcpListener(
		settings.TcpAddress,
		settings.TcpPort,
		&dnsUpstream,
	)
	if err != nil {
		logError(err, "Error creating TCP Listener")
		udpListener.Shutdown()
		return
	}
	logrus.Info("Started TCP Listener")

	defer tcpListener.Shutdown()
	defer udpListener.Shutdown()
	<-shutdown

}

func logError(err error, msg string) {
	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Error(msg)
}
