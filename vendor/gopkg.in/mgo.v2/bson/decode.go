// BSON library for Go
//
// Copyright (c) 2010-2012 - Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// gobson - BSON library for Go.

package bson

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type decoder struct {
	in      []byte
	i       int
	docType reflect.Type
}

var typeM = reflect.TypeOf(M{})

func newDecoder(in []byte) *decoder {
	return &decoder{in, 0, typeM}
}

// --------------------------------------------------------------------------
// Some helper functions.

func corrupted() {
	panic("Document is corrupted")
}

func settableValueOf(i interface{}) reflect.Value {
	v := reflect.ValueOf(i)
	sv := reflect.New(v.Type()).Elem()
	sv.Set(v)
	return sv
}

// --------------------------------------------------------------------------
// Unmarshaling of documents.

const (
	setterUnknown = iota
	setterNone
	setterType
	setterAddr
)

var setterStyles map[reflect.Type]int
var setterIface reflect.Type
var setterMutex sync.RWMutex

func init() {
	var iface Setter
	setterIface = reflect.TypeOf(&iface).Elem()
	setterStyles = make(map[reflect.Type]int)
}

func setterStyle(outt reflect.Type) int {
	setterMutex.RLock()
	style := setterStyles[outt]
	setterMutex.RUnlock()
	if style == setterUnknown {
		setterMutex.Lock()
		defer setterMutex.Unlock()
		if outt.Implements(setterIface) {
			setterStyles[outt] = setterType
		} else if reflect.PtrTo(outt).Implements(setterIface) {
			setterStyles[outt] = setterAddr
		} else {
			setterStyles[outt] = setterNone
		}
		style = setterStyles[outt]
	}
	return style
}

func getSetter(outt reflect.Type, out reflect.Value) Setter {
	style := setterStyle(outt)
	if style == setterNone {
		return nil
	}
	if style == setterAddr {
		if !out.CanAddr() {
			return nil
		}
		out = out.Addr()
	} else if outt.Kind() == reflect.Ptr && out.IsNil() {
		out.Set(reflect.New(outt.Elem()))
	}
	return out.Interface().(Setter)
}

func clearMap(m reflect.Value) {
	var none reflect.Value
	for _, k := range m.MapKeys() {
		m.SetMapIndex(k, none)
	}
}

func (d *decoder) readDocTo(out reflect.Value) {
	var elemType reflect.Type
	outt := out.Type()
	outk := outt.Kind()

	for {
		if outk == reflect.Ptr && out.IsNil() {
			out.Set(reflect.New(outt.Elem()))
		}
		if setter := getSetter(outt, out); setter != nil {
			var raw Raw
			d.readDocTo(reflect.ValueOf(&raw))
			err := setter.SetBSON(raw)
			if _, ok := err.(*TypeError); err != nil && !ok {
				panic(err)
			}
			return
		}
		if outk == reflect.Ptr {
			out = out.Elem()
			outt = out.Type()
			outk = out.Kind()
			continue
		}
		break
	}

	var fieldsMap map[string]fieldInfo
	var inlineMap reflect.Value
	start := d.i

	origout := out
	if outk == reflect.Interface {
		if d.docType.Kind() == reflect.Map {
			mv := reflect.MakeMap(d.docType)
			out.Set(mv)
			out = mv
		} else {
			dv := reflect.New(d.docType).Elem()
			out.Set(dv)
			out = dv
		}
		outt = out.Type()
		outk = outt.Kind()
	}

	docType := d.docType
	keyType := typeString
	convertKey := false
	switch outk {
	case reflect.Map:
		keyType = outt.Key()
		if keyType.Kind() != reflect.String {
			panic("BSON map must have string keys. Got: " + outt.String())
		}
		if keyType != typeString {
			convertKey = true
		}
		elemType = outt.Elem()
		if elemType == typeIface {
			d.docType = outt
		}
		if out.IsNil() {
			out.Set(reflect.MakeMap(out.Type()))
		} else if out.Len() > 0 {
			clearMap(out)
		}
	case reflect.Struct:
		if outt != typeRaw {
			sinfo, err := getStructInfo(out.Type())
			if err != nil {
				panic(err)
			}
			fieldsMap = sinfo.FieldsMap
			out.Set(sinfo.Zero)
			if sinfo.InlineMap != -1 {
				inlineMap = out.Field(sinfo.InlineMap)
				if !inlineMap.IsNil() && inlineMap.Len() > 0 {
					clearMap(inlineMap)
				}
				elemType = inlineMap.Type().Elem()
				if elemType == typeIface {
					d.docType = inlineMap.Type()
				}
			}
		}
	case reflect.Slice:
		switch outt.Elem() {
		case typeDocElem:
			origout.Set(d.readDocElems(outt))
			return
		case typeRawDocElem:
			origout.Set(d.readRawDocElems(outt))
			return
		}
		fallthrough
	default:
		panic("Unsupported document type for unmarshalling: " + out.Type().String())
	}

	end := int(d.readInt32())
	end += d.i - 4
	if end <= d.i || end > len(d.in) || d.in[end-1] != '\x00' {
		corrupted()
	}
	for d.in[d.i] != '\x00' {
		kind := d.readByte()
		name := d.readCStr()
		if d.i >= end {
			corrupted()
		}

		switch outk {
		case reflect.Map:
			e := reflect.New(elemType).Elem()
			if d.readElemTo(e, kind) {
				k := reflect.ValueOf(name)
				if convertKey {
					k = k.Convert(keyType)
				}
				out.SetMapIndex(k, e)
			}
		case reflect.Struct:
			if outt == typeRaw {
				d.dropElem(kind)
			} else {
				if info, ok := fieldsMap[name]; ok {
					if info.Inline == nil {
						d.readElemTo(out.Field(info.Num), kind)
					} else {
						d.readElemTo(out.FieldByIndex(info.Inline), kind)
					}
				} else if inlineMap.IsValid() {
					if inlineMap.IsNil() {
						inlineMap.Set(reflect.MakeMap(inlineMap.Type()))
					}
					e := reflect.New(elemType).Elem()
					if d.readElemTo(e, kind) {
						inlineMap.SetMapIndex(reflect.ValueOf(name), e)
					}
				} else {
					d.dropElem(kind)
				}
			}
		case reflect.Slice:
		}

		if d.i >= end {
			corrupted()
		}
	}
	d.i++ // '\x00'
	if d.i != end {
		corrupted()
	}
	d.docType = docType

	if outt == typeRaw {
		out.Set(reflect.ValueOf(Raw{0x03, d.in[start:d.i]}))
	}
}

func (d *decoder) readArrayDocTo(out reflect.Value) {
	end := int(d.readInt32())
	end += d.i - 4
	if end <= d.i || end > len(d.in) || d.in[end-1] != '\x00' {
		corrupted()
	}
	i := 0
	l := out.Len()
	for d.in[d.i] != '\x00' {
		if i >= l {
			panic("Length mismatch on array field")
		}
		kind := d.readByte()
		for d.i < end && d.in[d.i] != '\x00' {
			d.i++
		}
		if d.i >= end {
			corrupted()
		}
		d.i++
		d.readElemTo(out.Index(i), kind)
		if d.i >= end {
			corrupted()
		}
		i++
	}
	if i != l {
		panic("Length mismatch on array field")
	}
	d.i++ // '\x00'
	if d.i != end {
		corrupted()
	}
}

func (d *decoder) readSliceDoc(t reflect.Type) interface{} {
	tmp := make([]reflect.Value, 0, 8)
	elemType := t.Elem()
	if elemType == typeRawDocElem {
		d.dropElem(0x04)
		return reflect.Zero(t).Interface()
	}

	end := int(d.readInt32())
	end += d.i - 4
	if end <= d.i || end > len(d.in) || d.in[end-1] != '\x00' {
		corrupted()
	}
	for d.in[d.i] != '\x00' {
		kind := d.readByte()
		for d.i < end && d.in[d.i] != '\x00' {
			d.i++
		}
		if d.i >= end {
			corrupted()
		}
		d.i++
		e := reflect.New(elemType).Elem()
		if d.readElemTo(e, kind) {
			tmp = append(tmp, e)
		}
		if d.i >= end {
			corrupted()
		}
	}
	d.i++ // '\x00'
	if d.i != end {
		corrupted()
	}

	n := len(tmp)
	slice := reflect.MakeSlice(t, n, n)
	for i := 0; i != n; i++ {
		slice.Index(i).Set(tmp[i])
	}
	return slice.Interface()
}

var typeSlice = reflect.TypeOf([]interface{}{})
var typeIface = typeSlice.Elem()

func (d *decoder) readDocElems(typ reflect.Type) reflect.Value {
	docType := d.docType
	d.docType = typ
	slice := make([]DocElem, 0, 8)
	d.readDocWith(func(kind byte, name string) {
		e := DocElem{Name: name}
		v := reflect.ValueOf(&e.Value)
		if d.readElemTo(v.Elem(), kind) {
			slice = append(slice, e)
		}
	})
	slicev := reflect.New(typ).Elem()
	slicev.Set(reflect.ValueOf(slice))
	d.docType = docType
	return slicev
}

func (d *decoder) readRawDocElems(typ reflect.Type) reflect.Value {
	docType := d.docType
	d.docType = typ
	slice := make([]RawDocElem, 0, 8)
	d.readDocWith(func(kind byte, name string) {
		e := RawDocElem{Name: name}
		v := reflect.ValueOf(&e.Value)
		if d.readElemTo(v.Elem(), kind) {
			slice = append(slice, e)
		}
	})
	slicev := reflect.New(typ).Elem()
	slicev.Set(reflect.ValueOf(slice))
	d.docType = docType
	return slicev
}

func (d *decoder) readDocWith(f func(kind byte, name string)) {
	end := int(d.readInt32())
	end += d.i - 4
	if end <= d.i || end > len(d.in) || d.in[end-1] != '\x00' {
		corrupted()
	}
	for d.in[d.i] != '\x00' {
		kind := d.readByte()
		name := d.readCStr()
		if d.i >= end {
			corrupted()
		}
		f(kind, name)
		if d.i >= end {
			corrupted()
		}
	}
	d.i++ // '\x00'
	if d.i != end {
		corrupted()
	}
}

// --------------------------------------------------------------------------
// Unmarshaling of individual elements within a document.

var blackHole = settableValueOf(struct{}{})

func (d *decoder) dropElem(kind byte) {
	d.readElemTo(blackHole, kind)
}

// Attempt to decode an element from the document and put it into out.
// If the types are not compatible, the returned ok value will be
// false and out will be unchanged.
func (d *decoder) readElemTo(out reflect.Value, kind byte) (good bool) {

	start := d.i

	if kind == 0x03 {
		// Delegate unmarshaling of documents.
		outt := out.Type()
		outk := out.Kind()
		switch outk {
		case reflect.Interface, reflect.Ptr, reflect.Struct, reflect.Map:
			d.readDocTo(out)
			return true
		}
		if setterStyle(outt) != setterNone {
			d.readDocTo(out)
			return true
		}
		if outk == reflect.Slice {
			switch outt.Elem() {
			case typeDocElem:
				out.Set(d.readDocElems(outt))
			case typeRawDocElem:
				out.Set(d.readRawDocElems(outt))
			default:
				d.readDocTo(blackHole)
			}
			return true
		}
		d.readDocTo(blackHole)
		return true
	}

	var in interface{}

	switch kind {
	case 0x01: // Float64
		in = d.readFloat64()
	case 0x02: // UTF-8 string
		in = d.readStr()
	case 0x03: // Document
		panic("Can't happen. Handled above.")
	case 0x04: // Array
		outt := out.Type()
		if setterStyle(outt) != setterNone {
			// Skip the value so its data is handed to the setter below.
			d.dropElem(kind)
			break
		}
		for outt.Kind() == reflect.Ptr {
			outt = outt.Elem()
		}
		switch outt.Kind() {
		case reflect.Array:
			d.readArrayDocTo(out)
			return true
		case reflect.Slice:
			in = d.readSliceDoc(outt)
		default:
			in = d.readSliceDoc(typeSlice)
		}
	case 0x05: // Binary
		b := d.readBinary()
		if b.Kind == 0x00 || b.Kind == 0x02 {
			in = b.Data
		} else {
			in = b
		}
	case 0x06: // Undefined (obsolete, but still seen in the wild)
		in = Undefined
	case 0x07: // ObjectId
		in = ObjectId(d.readBytes(12))
	case 0x08: // Bool
		in = d.readBool()
	case 0x09: // Timestamp
		// MongoDB handles timestamps as milliseconds.
		i := d.readInt64()
		if i == -62135596800000 {
			in = time.Time{} // In UTC for convenience.
		} else {
			in = time.Unix(i/1e3, i%1e3*1e6)
		}
	case 0x0A: // Nil
		in = nil
	case 0x0B: // RegEx
		in = d.readRegEx()
	case 0x0C:
		in = DBPointer{Namespace: d.readStr(), Id: ObjectId(d.readBytes(12))}
	case 0x0D: // JavaScript without scope
		in = JavaScript{Code: d.readStr()}
	case 0x0E: // Symbol
		in = Symbol(d.readStr())
	case 0x0F: // JavaScript with scope
		d.i += 4 // Skip length
		js := JavaScript{d.readStr(), make(M)}
		d.