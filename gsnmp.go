// Package gsnmpgo is a go/cgo wrapper around gsnmp. It is under development,
// therefore API's may/will change, and doco/error handling/tests are minimal.
package gsnmpgo

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

// glib typedefs - http://developer.gnome.org/glib/2.35/glib-Basic-Types.html
// glib tutorial - http://www.dlhoffman.com/publiclibrary/software/gtk+-html-docs/gtk_tut-17.html
// gsnmp sourcecode browser - http://sourcecodebrowser.com/gsnmp/0.3.0/index.html

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

// convenience wrapper for gnet_snmp_enum_get_label()
gchar const *
get_err_label(gint32 const id) {
	return gnet_snmp_enum_get_label(gnet_snmp_enum_error_table, id);
}
*/
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

var _ = reflect.DeepEqual(0, 1) // dummy

// libname returns the name of this library, for generating error messages.
func libname() string {
	return "gsnmpgo"
}

// ParseURI parses an SNMP URI into fields.
//
// The generic URI parsing is done by gnet_uri_new(), and the SNMP specific
// portions by gnet_snmp_parse_uri(). Only basic URI validation is done here,
// more is done by ParsePath()
//
// Example:
//
//    uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)`
//    parsed_uri, err := gsnmpgo.ParseURI(uri)
//    if err != nil {
//    	fmt.Println(err)
//    	os.Exit(1)
//    }
//    fmt.Println("ParseURI():", parsed_uri)
func ParseURI(uri string) (parsed_uri *_Ctype_GURI, err error) {
	curi := (*C.gchar)(C.CString(uri))
	defer C.free(unsafe.Pointer(curi))

	var gerror *C.GError
	parsed_uri = C.gnet_snmp_parse_uri(curi, &gerror)
	if parsed_uri == nil {
		return nil, fmt.Errorf("%s: invalid snmp uri: %s", libname(), uri)
	}
	return parsed_uri, nil
}

// ParsePath parses an SNMP URI.
//
// The uritype will default to GNET_SNMP_URI_GET. If the uri ends in:
//
// '*' the uritype will be GNET_SNMP_URI_WALK
//
// '+' the uritype will be GNET_SNMP_URI_NEXT
//
// See RFC 4088 "Uniform Resource Identifier (URI) Scheme for the Simple
// Network Management Protocol (SNMP)" for further documentation.
func ParsePath(uri string, parsed_uri *_Ctype_GURI) (vbl *_Ctype_GList, uritype _Ctype_GNetSnmpUriType, err error) {
	var gerror *C.GError
	rv := C.gnet_snmp_parse_path(parsed_uri.path, &vbl, &uritype, &gerror)
	if rv == 0 {
		err_string := C.GoString((*_Ctype_char)(gerror.message))
		return vbl, uritype, fmt.Errorf("%s: %s: <%s>", libname(), err_string, uri)
	}
	return vbl, uritype, nil
}

// UriDelete frees the memory used by a parsed_uri.
//
// A defered call to UriDelete should be made after ParsePath().
func UriDelete(parsed_uri *_Ctype_GURI) {
	C.gnet_uri_delete(parsed_uri)
}

// NewUri creates a session from a parsed uri.
func NewUri(uri string, version SnmpVersion, parsed_uri *_Ctype_GURI) (session *_Ctype_GNetSnmp, err error) {
	var gerror *C.GError
	session = C.gnet_snmp_new_uri(parsed_uri, &gerror)

	// error handling
	if gerror != nil {
		err_string := C.GoString((*_Ctype_char)(gerror.message))
		C.g_clear_error(&gerror)
		return session, fmt.Errorf("%s: %s", libname(), err_string)
	}
	if session == nil {
		return session, fmt.Errorf("%s: unable to create session", libname())
	}
	session.version = (_Ctype_GNetSnmpVersion)(version)

	// results
	return session, nil
}

// Get does an SNMP get.
//
// It returns it results in C form, another function will convert the returned
// Glist to a Go struct.
func Get(session *_Ctype_GNetSnmp, vbl *_Ctype_GList) (*_Ctype_GList, error) {
	var gerror *C.GError
	out := C.gnet_snmp_sync_get(session, vbl, &gerror)

	// error handling
	if gerror != nil {
		err_string := C.GoString((*_Ctype_char)(gerror.message))
		C.g_clear_error(&gerror)
		return out, fmt.Errorf("%s: %s", libname(), err_string)
	}
	if PduError(session.error_status) != GNET_SNMP_PDU_ERR_NOERROR {
		es := C.get_err_label(session.error_status)
		err_string := C.GoString((*_Ctype_char)(es))
		return out, fmt.Errorf("%s: %s for uri %d", libname(), err_string, session.error_index)
	}

	// results
	return out, nil
}

func Dump(out *_Ctype_GList) {
	var result string
	for {
		if out == nil {
			fmt.Printf("result:\n%s", result)
			return
		}
		data := (*C.GNetSnmpVarBind)(out.data) // gsnmpgo._Ctype_gpointer -> *gsnmpgo._Ctype_GNetSnmpVarBind
		oid := GIntArrayOidString(data.oid, data.oid_len)
		result += oid + ":"
		result += fmt.Sprintf("%s", data._type) + ":"

		switch VarBindType(data._type) {

		case GNET_SNMP_VARBIND_TYPE_NULL:
			result += "NULL"

		case GNET_SNMP_VARBIND_TYPE_OCTETSTRING:
			result += union_ui8v_string(data.value, data.value_len)

		case GNET_SNMP_VARBIND_TYPE_OBJECTID:
			guint32_ptr := union_ui32v(data.value)
			result += GIntArrayOidString(guint32_ptr, data.value_len)

		case GNET_SNMP_VARBIND_TYPE_IPADDRESS:
			result += union_ui8v_ipaddress(data.value, data.value_len)

		case GNET_SNMP_VARBIND_TYPE_INTEGER32:
			result += fmt.Sprintf("%d", union_i32(data.value))

		case GNET_SNMP_VARBIND_TYPE_UNSIGNED32:
			result += fmt.Sprintf("%d", union_ui32(data.value))

		case GNET_SNMP_VARBIND_TYPE_COUNTER32:
			result += fmt.Sprintf("%d", union_ui32(data.value))

		case GNET_SNMP_VARBIND_TYPE_TIMETICKS:
			// TODO helper function to convert uint32 to "38 days, 15:56:15.00"
			result += fmt.Sprintf("%d", union_ui32(data.value))

		case GNET_SNMP_VARBIND_TYPE_OPAQUE:
			result += union_ui8v_hexstring(data.value, data.value_len)

		case GNET_SNMP_VARBIND_TYPE_COUNTER64:
			result += fmt.Sprintf("%d", union_ui64(data.value))

		case GNET_SNMP_VARBIND_TYPE_NOSUCHOBJECT:
			// do nothing

		case GNET_SNMP_VARBIND_TYPE_NOSUCHINSTANCE:
			// do nothing

		case GNET_SNMP_VARBIND_TYPE_ENDOFMIBVIEW:
			// do nothing

		}

		// move on to next element in list
		result += "\n"
		out = out.next
	}
	panic(fmt.Sprintf("%s: Dump(): fell out of for loop", libname()))
}
