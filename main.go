package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	//parse command line arguments
	var host string
	var port string
	flag.StringVar(&host, "host", "", "Host name to download certificate from")
	flag.StringVar(&port, "port", "443", "Port number to connect to")
	flag.Parse()
	if host == "" {
		fmt.Println("Host name is required")
		fmt.Println("Usage: sslcert-capture.exe -host <host name> -port <port number>")
		flag.PrintDefaults()
		return
	}

	hostname := host + ":" + port

	conn, err := tls.Dial("tcp", hostname, nil)
	if err != nil {
		fmt.Printf("Server doesn't support SSL certificate err: " + err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Saving certificates...")

	for _, cert := range conn.ConnectionState().PeerCertificates {
		fmt.Printf("Issuer: %s\nExpiry: %v\n", cert.Issuer, cert.NotAfter.Format(time.RFC850))
		saveCertificate(cert)
	}
}

func saveCertificate(cert *x509.Certificate) {
	//save certificate to file
	time := fmt.Sprintf("%v", time.Now().UnixNano())
	fmt.Println("Saving certificate to file: cert" + time + ".cer")
	certFile, err := os.Create("cert" + time + ".cer")
	if err != nil {
		panic("Couldn't create file: " + err.Error())
	}
	defer certFile.Close()
	certFile.Write(cert.Raw)

}
