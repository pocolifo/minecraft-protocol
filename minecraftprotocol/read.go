package minecraftprotocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// translated to Go from https://wiki.vg/Protocol#VarInt_and_VarLong
const (
	segmentByte  = 0x7F
	continueByte = 0x80
)

// translated to Go from https://wiki.vg/Protocol#VarInt_and_VarLong
func ReadNextVarInt(buffer *bytes.Buffer) (int32, int, error) {
	pos := 0
	size := 0
	var value int32 = 0

	for {
		// read next byte
		b, err := buffer.ReadByte()

		// return error if there was one
		if err != nil {
			return 0, 0, err
		}

		// convert it
		value |= int32(b&segmentByte) << pos
		size++

		// check if that byte is signalling the end, if so break
		if (b & continueByte) == 0 {
			break
		}

		// move up to the next position
		pos += 7

		// make sure that it isn't too big
		if pos > 32 {
			return 0, 0, errors.New("VarInt cannot be larger than 32 bytes")
		}
	}

	return value, size, nil
}

// translated to Go from https://wiki.vg/Protocol#VarInt_and_VarLong
func ReadNextVarLong(buffer []byte) (int64, int, error) {
	pos := 0
	size := 0
	var value int64 = 0

	for _, b := range buffer {
		// convert it
		value |= int64(b&segmentByte) << pos
		size++

		// check if that byte is signalling the end, if so break
		if (b & continueByte) == 0 {
			break
		}

		// move up to the next position
		pos += 7

		// make sure that it isn't too big
		if pos > 64 {
			return 0, 0, errors.New("VarLong cannot be larger than 64 bytes")
		}
	}

	return value, size, nil
}

func ReadNextByteArray(buffer *bytes.Buffer, length int32) ([]byte, error) {
	buf := make([]byte, length)
	_, err := buffer.Read(buf)
	return buf, err
}

func ReadNextString(buffer *bytes.Buffer) (string, int, error) {
	// strings are prefixed with the string size
	stringSize, bytesSize, err := ReadNextVarInt(buffer)

	if err != nil {
		return "", 0, err
	}

	buf := make([]byte, stringSize)
	_, err = buffer.Read(buf)

	if err != nil {
		return "", 0, err
	}

	return string(buf), len(buf) + bytesSize, nil
}

func ReadNextUnsignedShort(buffer *bytes.Buffer) (uint16, error) {
	buf := make([]byte, 2)
	_, err := buffer.Read(buf)
	return binary.BigEndian.Uint16(buf), err
}
