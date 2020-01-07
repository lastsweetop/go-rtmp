package amf

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	"reflect"
)

type decodeState struct {
	data []byte
	off  int // next read offset in data
}

func (d *decodeState) unmarshal(v interface{}) error {

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("nil interface{}")
	}

	t := reflect.TypeOf(v).Elem()
	value := reflect.ValueOf(v).Elem()
	for i := 0; i < value.NumField(); i++ {
		ft := t.Field(i)
		fv := value.Field(i)
		log.Println(ft, fv)
		switch ft.Type.String() {
		case "string":
			fv.SetString(d.scanString())
			break
		case "float64":
			fv.SetFloat(d.scanNumber())
			break
		case "bool":
			fv.SetBool(d.scanBool())
			break
		case "interface {}":
			d.off++
			break
		default:
			o := reflect.New(fv.Type().Elem())
			d.scanObject(o)
			fv.Set(o)
			break
		}
	}
	return nil
}

func (d *decodeState) scanObject(v reflect.Value) {
	v = v.Elem()
	t := v.Type()
	ot := d.data[d.off]
	if ot == 0x05 {
		return
	}
	d.off++
	for {
		key := d.scanProperty()
		log.Println(key)
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Tag.Get("amf") == key {
				switch d.data[d.off] {
				case 0x02:
					v.Field(i).SetString(d.scanString())
					break
				case 0x01:
					v.Field(i).SetBool(d.scanBool())
				case 0x00:
					v.Field(i).SetFloat(d.scanNumber())
					break
				default:
					o := reflect.New(v.Field(i).Type().Elem())
					d.scanObject(o)
					v.Field(i).Set(o)
					break
				}
			}
		}
		if d.data[d.off] == 0x00 && d.data[d.off+1] == 0x00 && d.data[d.off+2] == 0x09 {
			d.off += 3
			break
		}
	}

	log.Printf("scanObject %v %T\n", v, v)
}

func (d *decodeState) scanString() string {
	lens := d.data[d.off+1 : d.off+3]
	length := int(lens[0])<<8 + int(lens[1])
	d.off += 3 + length
	return string(d.data[d.off-length : d.off])
}

func (d *decodeState) scanProperty() string {
	lens := d.data[d.off : d.off+2]
	length := int(lens[0])<<8 + int(lens[1])
	d.off += 2 + length
	return string(d.data[d.off-length : d.off])
}

func (d *decodeState) scanNumber() float64 {
	d.off += 9
	return math.Float64frombits(binary.BigEndian.Uint64(d.data[d.off-8 : d.off]))
}

func (d *decodeState) scanBool() bool {
	d.off += 2
	return d.data[d.off-1] == 0x01
}

func (d *decodeState) init(data []byte) {
	d.data = data
	d.off = 0
}

func UnMarshal(data []byte, v interface{}) error {
	var d decodeState
	d.init(data)
	return d.unmarshal(v)
}
