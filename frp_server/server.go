package frp_server

import (
	"frp-demo/libs"
	"fmt"
	"net"
)

type ProxyType string

type FrpListener struct {
	ProxyType ProxyType
	Listener  net.Listener
}

type server struct {
	BindAddr string
	BindPort int

	Protocol string
	FrpListener FrpListener
	conn net.Conn
	readChan chan libs.MessageInterface
	sendChan chan libs.MessageInterface
}

func NewServer(config []byte) (ser *server, err error) {
	bindAddr := "0.0.0.0"
	bindPort := 7000
	protocol := "tcp"

	ln, err := net.Listen(protocol, fmt.Sprintf("%s:%d", bindAddr, bindPort))
	if err != nil {
		return nil ,nil
	}

	return &server{
		BindAddr: bindAddr,
		BindPort: bindPort,
		FrpListener: FrpListener{ProxyType: "tcp", Listener: ln},
	}, nil
}

func (ser *server) Run() (err error) {
	for {
		conn, err := ser.FrpListener.Listener.Accept()
		if err != nil {
			break
		}
		ser.conn = conn
		go ser.receive()
	}
	return err
}

func (ser *server) send() {
	//send message to server
	for {
		select {
		case msg, ok := <-ser.sendChan:
			if !ok {
				return
			}
			_, err := ser.conn.Write(msg.Dumps())
			if err != nil {
				return
			}
		}
	}
}

func (ser *server) receive() {
	//receive message from server
	for {
		select {
		case msg, ok := <-ser.readChan:
			if !ok {
				return
			}
			switch msg.Type() {
			case libs.Ping:
				fmt.Println("receive ping from client")
				ser.sendChan <- libs.NewMessage(libs.Pong, nil)
			case libs.Login:
				fmt.Println("receive login from client")
				ser.sendChan <- libs.NewMessage(libs.LoginResp, nil)
			case libs.NewProxy:
				fmt.Println("receive newProxy from client")
				ser.sendChan <- libs.NewMessage(libs.NewProxyResp, nil)
			case libs.NewWorkConn:
				fmt.Println("receive newWorkConn from client")
			default:
				fmt.Println(msg)
				return
			}
		}
	}
}
