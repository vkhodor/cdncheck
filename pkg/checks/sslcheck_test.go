package checks

import (
	"crypto/x509"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestDNSNameCheck(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: false, FullTimestamp: true})
	logger.SetLevel(logrus.DebugLevel)

	sslCheck := SSLCheck{
		CertDomains: []string{
			"content.cdn.personaly.bid",
			"us-01.cdn.personaly.bid",
			"*.cdn.personaly.bid",
		},
		Logger: logger,
	}

	crt := &x509.Certificate{
		DNSNames: []string{"x.cdn.personaly.bid"},
	}
	if sslCheck.dnsNameCheck(crt) != false {
		t.Error()
	}

	crt = &x509.Certificate{
		DNSNames: []string{"content.cdn.personaly.bid"},
	}
	if sslCheck.dnsNameCheck(crt) != true {
		t.Error()
	}
}

func TestExpirationCheck(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: false, FullTimestamp: true})
	logger.SetLevel(logrus.DebugLevel)

	sslCheck := SSLCheck{
		CertDomains: []string{
			"content.cdn.personaly.bid",
			"us-01.cdn.personaly.bid",
			"*.cdn.personaly.bid",
		},
		Logger: logger,
	}

	now := time.Now()
	crt := &x509.Certificate{
		DNSNames: []string{"xxx.personaly.bid"},
		NotAfter: now,
	}
	if sslCheck.expirationCheck(crt, time.Now()) != false {
		t.Error()
	}

	now = now.Add(time.Duration(61) * time.Minute)
	crt.NotAfter = now
	if sslCheck.expirationCheck(crt, time.Now()) != true {
		t.Error()
	}
}
