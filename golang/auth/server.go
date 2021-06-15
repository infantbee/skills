package main

import (
	"bufio"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

var Errhandshake error = errors.New("handshake error")

func verifySign(data, sign []byte) error {
	pubkey, err := ioutil.ReadFile("certs/rsa_public_key.pem")
	if err != nil {
		return err
	}

	block, _ := pem.Decode(pubkey)
	if block == nil {
		return fmt.Errorf("failed to pem decode public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], sign)
	return err
}

func serverHandshake(r io.Reader) error {
	//step := 1 // data
	msg := make([]byte, 160)
	n, err := r.Read(msg)
	if err != nil {
		log.Println("handshake read msg from client error: ", err.Error())
		return err
	}
	if n < 32 {
		log.Println("read msg format error")
		return Errhandshake
	}

	//step = 2 // summery
	data := msg[:32]
	h := md5.New()
	h.Write(data)
	summery := h.Sum(nil)

	//step = 3 // decode sign
	sign := msg[32:]
	err = verifySign(summery, sign)
	if err != nil {
		log.Printf("verify sign data error: ", err.Error())
	}

	log.Println("server handshake success.")
	return err
}

func main() {
	cert, err := tls.LoadX509KeyPair("./certs/server.pem", "./certs/server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		//ClientCAs:    clientCertPool,
	}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		err = serverHandshake(conn)
		if err != nil {
			log.Println("server handshake error: ", err.Error())
			break
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}
		log.Println("read from client msg:", msg)

		n, err := conn.Write([]byte("hi, wellcome to wuhan\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
