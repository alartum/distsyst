package emulutil

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

var goodOrdering = binary.LittleEndian

// DataType provides serialization codes for each data type
type dataType uint8

const (
	tInt32 dataType = iota
	tInt64
	tString
	tVectorInt32
	tEOF
)

// Message represents a serialized data fragment
type Message struct {
	body                   bytes.Buffer
	sendTime, deliveryTime time.Time
	from, to               int32
}

// After method tells wether other message is newer than other
func (m *Message) After(other *Message) bool {
	return m.deliveryTime.After(other.deliveryTime)
}

// PutInt32 serializes int32 value to the end of message
func (m *Message) PutInt32(q int32) {
	err := binary.Write(&m.body, goodOrdering, tInt32)
	if err != nil {
		panic(err)
	}
	err = binary.Write(&m.body, goodOrdering, q)
	if err != nil {
		panic(err)
	}
}

// PutInt64 serializes int64 value to the end of message
func (m *Message) PutInt64(q int64) {
	err := binary.Write(&m.body, goodOrdering, tInt64)
	if err != nil {
		panic(err)
	}
	err = binary.Write(&m.body, goodOrdering, q)
	if err != nil {
		panic(err)
	}
}

// PutString serializes string value to the end of message
func (m *Message) PutString(s string) {
	err := binary.Write(&m.body, goodOrdering, tString)
	if err != nil {
		panic(err)
	}
	err = binary.Write(&m.body, goodOrdering, int32(len(s)))
	if err != nil {
		panic(err)
	}
	err = binary.Write(&m.body, goodOrdering, []byte(s))
	if err != nil {
		panic(err)
	}
}

// PutVectorInt32 serializes int32 vector to the end of message
func (m *Message) PutVectorInt32(v []int32) {
	err := binary.Write(&m.body, goodOrdering, tVectorInt32)
	if err != nil {
		panic(err)
	}
	err = binary.Write(&m.body, goodOrdering, int32(len(v)))
	if err != nil {
		panic(err)
	}
	err = binary.Write(&m.body, goodOrdering, v)
	if err != nil {
		panic(err)
	}
}

// GetInt32 reads int32 value from the beginning of the message
func (m *Message) GetInt32() (int32, error) {
	var res int32
	var code dataType
	err := binary.Read(&m.body, goodOrdering, &code)
	if err != nil {
		panic(err)
	}
	if code != tInt32 {
		return -1, errors.New("GetInt32: expected int32 serialization code")
	}
	err = binary.Read(&m.body, goodOrdering, &res)
	if err != nil {
		panic(err)
	}
	return res, nil
}

// GetInt64 reads int64 value from the beginning of the message
func (m *Message) GetInt64() (int64, error) {
	var res int64
	var code dataType
	err := binary.Read(&m.body, goodOrdering, &code)
	if err != nil {
		panic(err)
	}
	if code != tInt64 {
		return -1, errors.New("GetInt32: expected int64 serialization code")
	}
	err = binary.Read(&m.body, goodOrdering, &res)
	if err != nil {
		panic(err)
	}
	return res, nil
}

// GetString reads string from the beginning of the message
func (m *Message) GetString() (string, error) {
	var code dataType
	err := binary.Read(&m.body, goodOrdering, &code)
	if err != nil {
		panic(err)
	}
	if code != tString {
		return "", errors.New("GetInt32: expected string serialization code")
	}
	var sLen int32
	err = binary.Read(&m.body, goodOrdering, &sLen)
	if err != nil {
		panic(err)
	}
	s := make([]byte, sLen)
	err = binary.Read(&m.body, goodOrdering, &s)
	if err != nil {
		panic(err)
	}

	return string(s), nil
}

// GetVectorInt32 reads int32 vector from the beginning of the message
func (m *Message) GetVectorInt32() ([]int32, error) {
	var code dataType
	err := binary.Read(&m.body, goodOrdering, &code)
	if err != nil {
		panic(err)
	}
	if code != tVectorInt32 {
		return nil, errors.New("GetInt32: expected int32 vector serialization code")
	}
	var vLen int32
	err = binary.Read(&m.body, goodOrdering, &vLen)
	if err != nil {
		panic(err)
	}
	v := make([]int32, vLen)
	err = binary.Read(&m.body, goodOrdering, &v)
	if err != nil {
		panic(err)
	}

	return v, nil
}
