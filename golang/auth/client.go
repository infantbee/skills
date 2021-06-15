package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func signature(data []byte) ([]byte, error) {
	contents, err := ioutil.ReadFile("certs/rsa_private_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(contents)
	if block == nil {
		return nil, fmt.Errorf("failed to pem decode private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var h = sha256.New()
	h.Write(data)
	var d = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, d)
}

func clientHandshake(w io.Writer) error {
	//step := 1 // data
	data := make([]byte, 32)
	n, err := rand.Reader.Read(data)
	if err != nil {
		log.Println("get rand data error:", err)
		return err
	}

	//step = 2 // summery
	h := md5.New()
	h.Write(data)
	summery := h.Sum(nil)

	//step = 3 // sign for summery
	sign, err := signature(summery)
	if err != nil {
		log.Println("sign msg error:", err)
		return err
	}

	//step = 4 // data.sign
	data = append(data, sign...)
	n, err = w.Write(data)
	if err != nil {
		log.Println(n, err)
		return err
	}

	return nil
}

func main() {
	conf := &tls.Config{
		//RootCAs:            caCertPool,
		//Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	err = clientHandshake(conn)
	if err != nil {
		log.Println("client handshake request error: ", err.Error())
		return
	}

	msg := "my name is leizi\n"
	n, err := conn.Write([]byte(msg))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}
	println("read from server msg:", string(buf[:n]))
}
