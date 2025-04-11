package broker

import (
	"bufio"
	"fmt"
	"github.com/FrostyKitten02/fileshare-common/util"
	"net"
	"os"
)

type Broker struct {
	Ip     net.IP
	Port   int
	conn   net.Conn
	roomId string
	reader bufio.Reader
}

func (b *Broker) ConnectToBroker() {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   b.Ip,
		Port: b.Port,
	})
	if err != nil {
		return
	}
	b.conn = conn

	roomId, roomErr := util.ReadRoomCreatedMessage(conn)
	if roomErr != nil {
		return
	}

	b.roomId = roomId
	fmt.Println("Your id:" + roomId)
	reader := bufio.NewReader(os.Stdin)
	b.reader = *reader
	b.inputHandler()
}

func (b *Broker) inputHandler() {
	fmt.Print("Select mode (receive[r], send[s]) to quite use q")
	input, readErr := b.reader.ReadString('\n')
	if readErr != nil {
		return
	}

	if input == "send" || input == "s" {
		b.sendModeHandler()
		return
	}

	if input == "receive" || input == "r" {
		b.receiveHandler()
		return
	}

	if input == "q" {
		fmt.Println("Quiting...")
		return
	}

	fmt.Println("Invalid input")
	b.inputHandler()
}

func (b *Broker) sendModeHandler() {
	connectId, connIdErr := b.reader.ReadString('\n')
	if connIdErr != nil {
		fmt.Println("Error reading connectId")
		return
	}

	err := util.WriteConnectToRoomMessage(b.conn, connectId)
	if err != nil {
		fmt.Println("error making connect request")
		return
	}

	res, connectToErr := util.ReadRoomConnectionInfo(b.conn)
	if connectToErr != nil {
		fmt.Println("error reading connect response")
		return
	}

	//magic port number? if 0 then we know our connection was refused
	if res.Port == 0 {
		fmt.Println("Connection refused")
	}

	//TODO add logic for sending file to other client

	b.inputHandler()
}

func (b *Broker) receiveHandler() {
	//TODO listen for any request from broker
	// when we get request from broker we can accept or deny it, in case we accept open up a server and send server info to broker

	fmt.Print("Waiting for anyone to connect...")

	//TODO add logic and handler for receiving files after we accept the request
	//we should just listen if any packets get to us, first should be a packet containing what other client wants to
}
