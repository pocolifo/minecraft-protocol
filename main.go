package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"

	"pocolifo.com/minecraft-protocol/v2/minecraftprotocol"
)

func main() {
	socket, _ := net.Listen("tcp", "0.0.0.0:25565")

	defer socket.Close()

	fmt.Println("Listening!")

	buffer := make([]byte, 1024)

	for {
		client, _ := socket.Accept()
		length, _ := client.Read(buffer)
		bs := buffer[:length]

		buf := bytes.NewBuffer(bs)
		packet, _ := minecraftprotocol.ReadNextPacket(buf)

		fmt.Println("Packet receieved")
		fmt.Printf("Length : %d\n", packet.Length)
		fmt.Printf("ID     : %d\n", packet.PacketID)

		if packet.PacketID == 0 { // server list ping
			db := bytes.NewBuffer(packet.Data)

			pv, _, _ := minecraftprotocol.ReadNextVarInt(db)
			fmt.Printf("Protocol Version: %d\n", pv)

			sa, _, _ := minecraftprotocol.ReadNextString(db)
			fmt.Printf("Server Address: %s\n", sa)

			port, _ := minecraftprotocol.ReadNextUnsignedShort(db)
			fmt.Printf("Port: %d\n", port)

			nextState, _, _ := minecraftprotocol.ReadNextVarInt(db)
			fmt.Printf("Next State: %d\n", nextState)
		}

		client.Close()
	}
}

func PrintHex(bs []byte) {
	hexString := hex.EncodeToString(bs)
	outputString := ""

	for i, e := range hexString {
		outputString += string(e)

		if (i-1)%2 == 0 {
			outputString += " "
		}
	}

	fmt.Println(outputString)
}
