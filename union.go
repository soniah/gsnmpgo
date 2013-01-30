package gsnmp

// gsnmp is a Go wrapper around the C gsnmp library.
//
// Copyright (C) 2013 Sonia Hamilton sonia@snowfrog.get.
//
// gsnmp is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// gsnmp is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser Public License for more details.
//
// You should have received a copy of the GNU Lesser Public License
// along with gsnmp.  If not, see <http://www.gnu.org/licenses/>.

// union.go contains functions for extracting values from the union
// in this struct:
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
#cgo pkg-config: glib-2.0 gsnmp

#include <gsnmp/ber.h>
#include <gsnmp/pdu.h>
#include <gsnmp/dispatch.h>
#include <gsnmp/message.h>
#include <gsnmp/security.h>
#include <gsnmp/session.h>
#include <gsnmp/transport.h>
#include <gsnmp/utils.h>
#include <gsnmp/gsnmp.h>

#include <stdlib.h>
#include <stdio.h>
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

// return ui32v field as a guint32 ptr
func union_to_guint32_ptr(cbytes [8]byte) (result *_Ctype_guint32) {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil {
		uptr := uintptr(ptr)
		return (*_Ctype_guint32)(unsafe.Pointer(uptr))
	}
	return nil
}

// return ui8v as a string
func union_to_string(cbytes [8]byte) string {
	buf := bytes.NewBuffer(cbytes[:])
	var ptr uint64
	if err := binary.Read(buf, binary.LittleEndian, &ptr); err == nil {
		uptr := uintptr(ptr)
		char_ptr := (*_Ctype_char)(unsafe.Pointer(uptr))
		return C.GoString(char_ptr)
	}
	return ""
}
