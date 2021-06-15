package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func VerifyCert(certFile string) bool {
	certBytes, err := ioutil.ReadFile(certFile)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	pemBlock, _ := pem.Decode([]byte(certBytes))
	if pemBlock == nil {
		fmt.Println("decode error")
		return false
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//fmt.Printf("Name %s\n", cert.Subject.CommonName)
	fmt.Printf("Not before %s\n", cert.NotBefore.String())
	fmt.Printf("Not after %s\n", cert.NotAfter.String())
	tmNow := time.Now()
	if tmNow.Before(cert.NotBefore) ||
		tmNow.After(cert.NotAfter) {
		return false
	}

	return true
}

func main() {
	VerifyCert("./certs/client.pem")

	cert, err := tls.LoadX509KeyPair("./certs/client.pem", "./certs/client.key")
	if err != nil {
		log.Println(err)
		return
	}

	// certBytes, err := ioutil.ReadFile("./certs/rootCA.pem")
	// if err != nil {
	// 	panic("Unable to read cert.pem")
	// }
	// caCertPool := x509.NewCertPool()
	// ok := caCertPool.AppendCertsFromPEM(certBytes)
	// if !ok {
	// 	panic("failed to parse root certificate")
	// }

	conf := &tls.Config{
		//RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	msg := "my name is leizi\n"
	n, err := conn.Write([]byte(msg))
	if err != nil {
		log.Println(n, err)
		return
	}
	fmt.Println(msg)

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}
	println("read from server msg:", string(buf[:n]))
}
