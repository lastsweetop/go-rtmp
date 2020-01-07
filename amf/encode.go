package amf

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"reflect"
	"sync"
)

type encodeState struct {
	bytes.Buffer
}

var encodeStatePool sync.Pool

func newEncodeState() *encodeState {
	if v := encodeStatePool.Get(); v != nil {
		e := v.(*encodeState)
		e.Reset()
		return e
	}
	return new(encodeState)
}

func (e *encodeState) marshal(v interface{}) error {
	t := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		t = t.Elem()
	}
	log.Println(t.String())
	log.Println(value)
	for i := 0; i < value.NumField(); i++ {
		ft := t.Field(i)
		fv := value.Field(i)
		switch ft.Type.String() {
		case "string":
			e.writeAMFString(fv.String())
			break
		case "float64":
			e.writeAMFNumber(fv.Float())
			break
		case "bool":
			e.writeAMFBool(fv.Bool())
		default:
			e.writeAMFObject(fv.Interface())
			break
		}
	}
	return nil
}

func (e *encodeState) writeAMFString(s string) {
	e.WriteByte(0x02)
	strlen := len(s)
	e.WriteByte(byte(strlen >> 8))
	e.WriteByte(byte(strlen))
	e.WriteString(s)
}

func (e *encodeState) writeAMFBool(b bool) {
	e.WriteByte(0x01)
	if b {
		e.WriteByte(0x01)
	} else {
		e.WriteByte(0x00)
	}
}

func (e *encodeState) writeAMFObjectKey(s string) {
	strlen := len(s)
	e.WriteByte(byte(strlen >> 8))
	e.WriteByte(byte(strlen))
	e.WriteString(s)
}

func (e *encodeState) writeAMFNumber(i float64) {
	e.WriteByte(0x00)
	temp := make([]byte, 8)
	binary.BigEndian.PutUint64(temp, math.Float64bits(i))
	e.Write(temp)
}

func (e *encodeState) writeAMFObject(i interface{}) {
	if i == nil {
		e.WriteByte(0x05)
		return
	}
	e.WriteByte(0x03)

	fv := reflect.ValueOf(i)
	ft := reflect.TypeOf(i)
	if fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
		ft = ft.Elem()
	}

	for j := 0; j < fv.NumField(); j++ {
		e.writeAMFObjectKey(ft.Field(j).Tag.Get("amf"))

		switch ft.Field(j).Type.String() {
		case "string":
			e.writeAMFString(fv.Field(j).String())
			break
		case "float64":
			e.writeAMFNumber(fv.Field(j).Float())
			break
		case "boolean":
			e.writeAMFBool(fv.Field(j).Bool())
			break
		default:
			e.writeAMFObject(fv.Field(j).Interface())
			break
		}
	}
	e.Write([]byte{0x00, 0x00, 0x09})
}

func Marshal(v interface{}) ([]byte, error) {
	e := newEncodeState()
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)
	e.Reset()
	encodeStatePool.Put(e)
	return buf, nil
}

//func Unmarshal(data []byte, v interface{}) error {
//
//}
