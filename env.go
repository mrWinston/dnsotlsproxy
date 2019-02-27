package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Settings struct {
	UdpPort            string
	UdpAddress         string
	TcpPort            string
	TcpAddress         string
	DnsUpstreamPort    string
	DnsUpstreamAddress string
	LogLevel           logrus.Level
}

func (s Settings) String() string {
	return fmt.Sprintf("UdpPort: %s, UdpAddress: %s, TcpPort: %s, TcpAddress: %s, DnsUpstreamPort: %s, DnsUpstreamAddress: %s, LogLevel: %s",
		s.UdpPort,
		s.UdpAddress,
		s.TcpPort,
		s.TcpAddress,
		s.DnsUpstreamPort,
		s.DnsUpstreamAddress,
		s.LogLevel,
	)
}

func getSettingsFromEnvVars() *Settings {
	var udpport, udpaddress, tcpport, tcpaddress, dnsupstreamport, dnsupstreamaddress, loglevelstring string

	if udpport = os.Getenv("UDP_PORT"); udpport == "" {
		udpport = "53"
	}
	if udpaddress = os.Getenv("UDP_ADDRESS"); udpaddress == "" {
		udpaddress = "0.0.0.0"
	}

	if tcpport = os.Getenv("TCP_PORT"); tcpport == "" {
		tcpport = "53"
	}
	if tcpaddress = os.Getenv("TCP_ADDRESS"); tcpaddress == "" {
		tcpaddress = "0.0.0.0"
	}

	if dnsupstreamport = os.Getenv("DNS_UPSTREAM_PORT"); dnsupstreamport == "" {
		dnsupstreamport = "853"
	}
	if dnsupstreamaddress = os.Getenv("DNS_UPSTREAM_ADDRESS"); dnsupstreamaddress == "" {
		dnsupstreamaddress = "1.1.1.1"
	}

	if loglevelstring = os.Getenv("LOG_LEVEL"); loglevelstring == "" {
		loglevelstring = "error"
	}

	loglevel, err := logrus.ParseLevel(loglevelstring)
	if err != nil {
		logrus.Error(fmt.Sprintf("%s is not a valid LogLevel, Using default", loglevelstring))
		loglevel = logrus.ErrorLevel
	}

	return &Settings{
		UdpPort:            udpport,
		UdpAddress:         udpaddress,
		TcpPort:            tcpport,
		TcpAddress:         tcpaddress,
		DnsUpstreamPort:    dnsupstreamport,
		DnsUpstreamAddress: dnsupstreamaddress,
		LogLevel:           loglevel,
	}

}
