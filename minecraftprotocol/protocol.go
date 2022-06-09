package minecraftprotocol

import (
	"bytes"
)

// From https://wiki.vg/Protocol#Packet_format
type MinecraftPacket struct {
	Length   int32
	PacketID int32
	Data     []byte
}

func ReadNextPacket(buffer *bytes.Buffer) (*MinecraftPacket, error) {
	packet := &MinecraftPacket{}

	// read packet length
	length, _, err := ReadNextVarInt(buffer)

	if err != nil {
		return nil, err
	}

	packet.Length = length

	// read packet id
	packetID, packetIDSize, err := ReadNextVarInt(buffer)

	if err != nil {
		return nil, err
	}

	packet.PacketID = packetID

	// read packet data
	data, err := ReadNextByteArray(buffer, length-int32(packetIDSize))
	packet.Data = data

	if err != nil {
		return nil, err
	}

	// return
	return packet, nil
}
