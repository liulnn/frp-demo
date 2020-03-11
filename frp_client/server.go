package frp_client

import (
	"frp-demo/libs"
	"fmt"
	"net"
)

type TcpProxy struct {
	LocalIp   string
	LocalPort int

	RemotePort int
}

type UdpProxy struct {
	LocalIp   string
	LocalPort int

	RemotePort int
}

type client struct {
	ServerAddr string
	ServerPort int

	Protocol string
	conn     net.Conn
	readChan chan libs.MessageInterface
	sendChan chan libs.MessageInterface

	proxyCfgs map[string]ProxyConfig
}

func NewClient(config []byte) (cli *client, err error) {
	return &client{
		ServerAddr: "0.0.0.0",
		ServerPort: 7432,
		Protocol:   "tcp",
	}, nil
}

func (cli *client) connectServer() (err error) {
	conn, err := net.Dial(cli.Protocol, fmt.Sprintf("%s:%d", cli.ServerAddr, cli.ServerPort))
	if err != nil {
		return err
	}
	cli.conn = conn

	//send login request message
	loginMsg := libs.NewMessage(libs.Login, nil).Dumps()
	_, err = conn.Write(loginMsg)
	if err != nil {
		return err
	}
	//receive login response message
	msgHeader := make([]byte, 2)
	_, err = conn.Read(msgHeader)
	if err != nil {
		return err
	}
	return err
}

func (cli *client) Run() error {
	err := cli.connectServer()
	if err != nil {
		return err
	}

	//login success
	go cli.send()
	go cli.receive()
	return nil
}

func (cli *client) send() {
	//send message to server
	for {
		select {
		case msg, ok := <-cli.sendChan:
			if !ok {
				return
			}
			_, err := cli.conn.Write(msg.Dumps())
			if err != nil {
				return
			}
		}
	}
}

func (cli *client) receive() {
	//receive message from server
	for {
		select {
		case msg, ok := <-cli.readChan:
			if !ok {
				return
			}
			switch msg.Type() {
			case libs.Pong:
				fmt.Println("receive pong from server")
			case libs.LoginResp:
				fmt.Println("receive loginResp from server")
			case libs.NewProxyResp:
				fmt.Println("receive newProxyResp from server")
			case libs.ReqWorkConn:
				fmt.Println("receive reqWorkConn from server")
				cli.sendChan <- libs.NewMessage(libs.NewWorkConn, nil)
			default:
				fmt.Println(msg)
				return
			}
		}
	}
}

func (cli *client) reload(config []byte) {
	//update proxyCfgs
	for {
		for name, pxy := range cli.proxyCfgs {
			fmt.Println(name)
			pxy.Start()
		}

	}

}
