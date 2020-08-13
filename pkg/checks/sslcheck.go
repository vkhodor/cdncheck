package checks

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type SSLCheck struct {
	Logger      *logrus.Logger
	CertDomains []string
	Host        string
	Port        int
}

func (h *SSLCheck) Check() (bool, error) {
	//	now := time.Now()
	crt, err := h.getCert()
	if err != nil {
		return false, err
	}

	if !h.dnsNameCheck(crt) {
		return false, nil
	}
	if !h.expirationCheck(crt, time.Now()) {
		return false, nil
	}
	return true, nil
}

func (h *SSLCheck) dnsNameCheck(crt *x509.Certificate) bool {
	for _, domain := range h.CertDomains {
		h.Logger.Debug("SSLCheck for domain ", domain)
		for _, crtDomain := range crt.DNSNames {
			if domain == crtDomain {
				h.Logger.Debug(domain, " == ", crtDomain)
				return true
			}
			h.Logger.Debug(domain, " != ", crtDomain)
		}
	}
	return false
}

func (h *SSLCheck) expirationCheck(crt *x509.Certificate, now time.Time) bool {
	expirationIn := crt.NotAfter.Sub(now).Hours()
	h.Logger.Debug("expirationCheck: expirationIn = ", expirationIn)
	if expirationIn <= 1 {
		return false
	}
	return true
}

func (h *SSLCheck) getCert() (*x509.Certificate, error) {

	strPort := strconv.Itoa(h.Port)
	h.Logger.Debug("Get SSL Cert from ", h.Host+":"+strPort)

	conn, err := tls.Dial("tcp", h.Host+":"+strPort, &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		h.Logger.Debug(err)
		return nil, err
	}

	if len(conn.ConnectionState().PeerCertificates) < 1 {
		return nil, errors.New("no SSL certs found")
	}
	return conn.ConnectionState().PeerCertificates[0], nil
}
