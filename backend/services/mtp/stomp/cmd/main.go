package main

import (
	"log"
	"net"
	"os"

	"github.com/go-stomp/stomp/v3/server"
	"github.com/joho/godotenv"
)

type Credentials struct {
	Login  string
	Passwd string
}

func (c *Credentials) Authenticate(login, passwd string) bool {

	if c.Login == "" && c.Passwd == "" {
		log.Println("NEW CLIENT AUTH: Server is running with no auth")
		return true
	}

	if login != c.Login || passwd != c.Passwd {
		log.Println("NEW CLIENT AUTH: Invalid Credentials")
		return false
	}
	log.Println("NEW CLIENT AUTH: OK")
	return true
}

func main() {

	err := godotenv.Load()
	localEnv := ".env.local"
	if _, err := os.Stat(localEnv); err == nil {
		_ = godotenv.Overload(localEnv)
		log.Println("Loaded variables from '.env.local'")
	} else {
		log.Println("Loaded variables from '.env'")
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	creds := Credentials{
		Login:  os.Getenv("STOMP_USERNAME"),
		Passwd: os.Getenv("STOMP_PASSWORD"),
	}

	l, err := net.Listen("tcp", server.DefaultAddr)
	if err != nil {
		log.Println("Error to open tcp port: ", err)
	}

	s := server.Server{
		Addr:          server.DefaultAddr,
		HeartBeat:     server.DefaultHeartBeat,
		Authenticator: &creds,
	}

	log.Println("Started STOMP server at port", s.Addr)
	err = s.Serve(l)
	if err != nil {
		log.Println("Error to start stomp server: ", err)
	}
}
