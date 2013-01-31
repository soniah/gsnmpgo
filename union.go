package gsnmpgo

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

// union.go contains functions for extracting values from the union
// in the following struct. Go will read the struct as a byte sequence
// as long as the longest field in the struct ie [8]byte. Hence these
// functions take "cbytes [8]byte" as a parameter.
//
// struct _GNetSnmpVarBind {
//     guint32		*oid;		/* name of the variable */
//     gsize		oid_len;	/* length of the name */
//     GNetSnmpVarBindType	type;		/* variable type / exception */
//     union {
//         gint32   i32;			/* 32 bit signed   */
//         guint32  ui32;			/* 32 bit unsigned */
//         gint64   i64;			/* 64 bit signed   */
//         guint64  ui64;			/* 64 bit unsigned */
//         guint8  *ui8v;			/*  8 bit unsigned vector */
//         guint32 *ui32v;			/* 32 bit unsigned vector */
//     }			value;		/* value of the variable */
//     gsize		value_len;	/* length of a vector in bytes */
// };

/*
// no C code
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"unsafe"
)

// return i32 field
func union_i32(cbytes [8]byte) (result int32) {
	buf := bytes.NewBuffer(cbytes[:])
	if err := binary.Read(buf, binary.LittleEndian, &result); err == nil { // read bytes as int32
		return result
	}
	return 0
}

// return ui32 field
func union_ui32(cbytes [8]byte) (result uint32) {
	buf := bytes.NewBuffer(cbytes[:])
	if err := binary.Read(buf, binary.LittleEndian, &result); err == nil { // read bytes as uint32
		return result
	}
	return 0
}

// i64 field isn't used (??)

// return ui64 field
func union_ui64(cbytes [8]byte) (result uint64) {
	buf := bytes.NewBuffer(cbytes[:])
	if err := binary.Read(buf, binary.LittleEndian, &result); err == nil { // read bytes as uint64
		return result
	}
	return 0
}

// return ui8v field as an ip address
func union_ui8v_ipaddress(cbytes [8]byte, value_len _Ctype_gsize) (result string) {
	if int(value_len) != 4 { // an ip4 address must have 4 bytes
		return
	}
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil { // read bytes as uint64
		up := (unsafe.Pointer(uintptr(ptr))) // convert the uint64 into a pointer
		gobytes := C.GoBytes(up, 4)
		for i := 0; i < 4; i++ {
			result += fmt.Sprintf(".%d", gobytes[i])
		}
		return result[1:] // strip leading dot
	}
	return
}

// return ui8v field as a string
func union_ui8v_string(cbytes [8]byte, value_len _Ctype_gsize) (result string) {
	var ptr uint64
	var err error
	buf := bytes.NewBuffer(cbytes[:])
	if err = binary.Read(buf, binary.LittleEndian, &ptr); err != nil { // read bytes as uint64
		return
	}
	up := (unsafe.Pointer(uintptr(ptr))) // convert the uint64 into a pointer

	length := (_Ctype_int)(value_len)
	gobytes := C.GoBytes(up, length)
	for i := 0; i < int(length); i++ {
		if !strconv.IsPrint(rune(gobytes[i])) {
			// can't pass gobytes & length to union_ui8v_hexstring() -
			// it's also used for TYPE_OPAQUE
			return union_ui8v_hexstring(cbytes, value_len)
		}
	}

	char_ptr := (*_Ctype_char)(up)
	return C.GoString(char_ptr)
}

// the ui8v field contains unprintable characters - convert to "hex string"
//
// eg 00 25 89 27 56 1B
func union_ui8v_hexstring(cbytes [8]byte, value_len _Ctype_gsize) (result string) {
	var ptr uint64
	var err error
	buf := bytes.NewBuffer(cbytes[:])
	if err = binary.Read(buf, binary.LittleEndian, &ptr); err != nil { // read bytes as uint64
		return
	}
	up := (unsafe.Pointer(uintptr(ptr))) // convert the uint64 into a pointer

	length := (_Ctype_int)(value_len)
	gobytes := C.GoBytes(up, length)
	for i := 0; i < int(length); i++ {
		result += fmt.Sprintf(" %02X", gobytes[i])
	}
	return result[1:] // strip leading space
}

// return ui32v field
func union_ui32v(cbytes [8]byte) (result *_Ctype_guint32) {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil { // read bytes as uint64
		return (*_Ctype_guint32)(unsafe.Pointer(uintptr(ptr))) // convert the uint64 into a pointer
	}
	return nil
}
