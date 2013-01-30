package gsnmp

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

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
