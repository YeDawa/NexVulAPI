package tasks

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"

	"net"
	"strings"
)

func GetCertificateInfo(domain string) ([]CertInfo, error) {
	if !strings.Contains(domain, ":") {
		domain = net.JoinHostPort(domain, "443")
	}

	conn, err := tls.Dial("tcp", domain, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	var infos []CertInfo

	for _, cert := range certs {
		info := CertInfo{
			Subject:            cert.Subject.String(),
			Issuer:             cert.Issuer.String(),
			NotBefore:          cert.NotBefore,
			NotAfter:           cert.NotAfter,
			DNSNames:           cert.DNSNames,
			SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		}

		switch pub := cert.PublicKey.(type) {
		case *rsa.PublicKey:
			info.PublicKey = "RSA " + string(rune(pub.N.BitLen())) + " bits"
		case *ecdsa.PublicKey:
			info.PublicKey = "ECDSA " + pub.Curve.Params().Name
		default:
			info.PublicKey = "Unknown"
		}

		infos = append(infos, info)
	}

	return infos, nil
}
