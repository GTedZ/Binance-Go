package Binance

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Serializes any number into binary format
// Accepts any type (int, uint, float) and any size (8, 16, 32, 64)
func SerializeNumber(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, fmt.Errorf("serialization error: %v", err)
	}
	return buf.Bytes(), nil
}

// Deserializes any number into binary format
// Accepts any type (int, uint, float) and any size (8, 16, 32, 64)
func DeserializeNumber(data []byte, value interface{}) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, value)
	if err != nil {
		return fmt.Errorf("deserialization error: %v", err)
	}
	return nil
}

// Serialize a string from string to binary format
func SerializeString(value string) ([]byte, error) {
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
func DeserializeString(data []byte) (string, error) {
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

// Checks if a value is different from its default value.
func IsDifferentFromDefault(value any) bool {
	// Get the reflect.Value of the input
	val := reflect.ValueOf(value)

	// Get the default value of the type
	defaultValue := reflect.Zero(val.Type()).Interface()

	// Compare the input value with the default value
	return !reflect.DeepEqual(value, defaultValue)
}
