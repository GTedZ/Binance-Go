package Binance

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BinaryUtils struct{}

// Serializes any number into binary format
// Accepts any type (int, uint, float) and any size (8, 16, 32, 64)
func (utils *BinaryUtils) SerializeNumber(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, fmt.Errorf("serialization error: %v", err)
	}
	return buf.Bytes(), nil
}

// Deserializes any number into binary format
// Accepts any type (int, uint, float) and any size (8, 16, 32, 64)
func (utils *BinaryUtils) DeserializeNumber(data []byte, value interface{}) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, value)
	if err != nil {
		return fmt.Errorf("deserialization error: %v", err)
	}
	return nil
}

// Serialize a string from string to binary format
func (utils *BinaryUtils) SerializeString(value string) ([]byte, error) {
	length := int32(len(value))
	buf := new(bytes.Buffer)
	// Write the length of the string first
	err := binary.Write(buf, binary.LittleEndian, length)
	if err != nil {
		return nil, fmt.Errorf("serialization error: %v", err)
	}
	// Write the string bytes
	_, err = buf.Write([]byte(value))
	if err != nil {
		return nil, fmt.Errorf("serialization error: %v", err)
	}
	return buf.Bytes(), nil
}

// Deserialize a string from binary to string format
func (utils *BinaryUtils) DeserializeString(data []byte) (string, error) {
	buf := bytes.NewReader(data)
	var length int32
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return "", fmt.Errorf("deserialization error: %v", err)
	}
	strBytes := make([]byte, length)
	_, err = buf.Read(strBytes)
	if err != nil {
		return "", fmt.Errorf("deserialization error: %v", err)
	}
	return string(strBytes), nil
}
