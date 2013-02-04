package gsnmpgo

// gsnmpgo is a go/cgo wrapper around gsnmp.
//
// Copyright (C) 2012-2013 Sonia Hamilton sonia@snowfrog.net.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

type QueryResult struct {
	Oid   string
	Value Varbinder
}

type QueryResults []QueryResult

// Query takes a URI in RFC 4088 format, does an SNMP query and returns the results.
func Query(uri string, version SnmpVersion) (results QueryResults, err error) {
	parsed_uri, err := parseURI(uri)
	if err != nil {
		return nil, err
	}

	vbl, uritype, err := parsePath(uri, parsed_uri)
	defer uriDelete(parsed_uri)
	if err != nil {
		return nil, err
	}

	session, err := newUri(uri, version, parsed_uri)
	if err != nil {
		return nil, err
	}

	// TODO must do a free (g_list_foreach(gnet_snmp_varbind_delete), g_list_free) on vbl_results
	vbl_results, err := querySync(session, vbl, uritype)
	if err != nil {
		return nil, err
	}
	return convertResults(vbl_results), nil // TODO no err from decode?
}

// Dump is a convenience function for printing the results of a Query.
func Dump(results QueryResults) {
	fmt.Println("Dump:")
	for _, result := range results {
		fmt.Printf("%T:%s STRING:%s INTEGER:%d\n", result.Value, result.Oid, result.Value, result.Value.Integer())
	}
}

// parseURI parses an SNMP URI into fields.
//
// The generic URI parsing is done by gnet_uri_new(), and the SNMP specific
// portions by gnet_snmp_parse_uri(). Only basic URI validation is done here,
// more is done by parsePath()
//
// Example:
//
//    uri := `snmp://public@192.168.1.10//(1.3.6.1.2.1.1.1.0,1.3.6.1.2.1.1.2.0)`
//    parsed_uri, err := gsnmpgo.parseURI(uri)
//    if err != nil {
//    	fmt.Println(err)
//    	os.Exit(1)
//    }
//    fmt.Println("parseURI():", parsed_uri)
func parseURI(uri string) (parsed_uri *_Ctype_GURI, err error) {
	curi := (*C.gchar)(C.CString(uri))
	defer C.free(unsafe.Pointer(curi))

	var gerror *C.GError
	parsed_uri = C.gnet_snmp_parse_uri(curi, &gerror)
	if parsed_uri == nil {
		return nil, fmt.Errorf("%s: invalid snmp uri: %s", libname(), uri)
	}
	return parsed_uri, nil
}

// parsePath parses an SNMP URI.
//
// The uritype will default to GNET_SNMP_URI_GET. If the uri ends in:
//
// '*' the uritype will be GNET_SNMP_URI_WALK
//
// '+' the uritype will be GNET_SNMP_URI_NEXT
//
// See RFC 4088 "Uniform Resource Identifier (URI) Scheme for the Simple
// Network Management Protocol (SNMP)" for further documentation.
func parsePath(uri string, parsed_uri *_Ctype_GURI) (vbl *_Ctype_GList, uritype _Ctype_GNetSnmpUriType, err error) {
	var gerror *C.GError
	rv := C.gnet_snmp_parse_path(parsed_uri.path, &vbl, &uritype, &gerror)
	if rv == 0 {
		err_string := C.GoString((*_Ctype_char)(gerror.message))
		return vbl, uritype, fmt.Errorf("%s: %s: <%s>", libname(), err_string, uri)
	}
	return vbl, uritype, nil
}

// uriDelete frees the memory used by a parsed_uri.
//
// A defered call to uriDelete should be made after parsePath().
func uriDelete(parsed_uri *_Ctype_GURI) {
	C.gnet_uri_delete(parsed_uri)
}

// newUri creates a session from a parsed uri.
func newUri(uri string, version SnmpVersion, parsed_uri *_Ctype_GURI) (session *_Ctype_GNetSnmp, err error) {
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

// Do an gsnmp library sync_* query
//
// Results are returned in C form, use convertResults() to convert to a Go struct.
func querySync(session *_Ctype_GNetSnmp, vbl *_Ctype_GList,
	uritype _Ctype_GNetSnmpUriType) (*_Ctype_GList, error) {
	var gerror *C.GError
	var out *_Ctype_GList

	switch UriType(uritype) {
	case GNET_SNMP_URI_GET:
		out = C.gnet_snmp_sync_get(session, vbl, &gerror)
	case GNET_SNMP_URI_NEXT:
		out = C.gnet_snmp_sync_getnext(session, vbl, &gerror)
	case GNET_SNMP_URI_WALK:
		out = C.gnet_snmp_sync_walk(session, vbl, &gerror)
	default:
		panic(fmt.Sprintf("%s: QueryC(): unknown uritype", libname()))
	}

	// error handling
	if gerror != nil {
		err_string := C.GoString((*_Ctype_char)(gerror.message))
		C.g_clear_error(&gerror)
		return out, fmt.Errorf("%s: %s", libname(), err_string)
	}
	err_status := PduError(session.error_status)
	switch UriType(uritype) {
	case GNET_SNMP_URI_WALK:
		if err_status != GNET_SNMP_PDU_ERR_NOERROR && err_status != GNET_SNMP_PDU_ERR_NOSUCHNAME {
			es := C.get_err_label(session.error_status)
			err_string := C.GoString((*_Ctype_char)(es))
			return out, fmt.Errorf("%s: %s for uri %d", libname(), err_string, session.error_index)
		}
	default:
		if err_status != GNET_SNMP_PDU_ERR_NOERROR {
			es := C.get_err_label(session.error_status)
			err_string := C.GoString((*_Ctype_char)(es))
			return out, fmt.Errorf("%s: %s for uri %d", libname(), err_string, session.error_index)
		}
	}

	// results
	return out, nil
}

// convertResults converts C results to a Go struct.
func convertResults(out *_Ctype_GList) (results QueryResults) {
	for {
		if out == nil {
			// finished
			return results
		}

		// another result: initialise
		data := (*C.GNetSnmpVarBind)(out.data)
		oid := gIntArrayOidString(data.oid, data.oid_len)
		result := QueryResult{Oid: oid}
		var value Varbinder

		// convert C values to Go values
		vbt := VarBindType(data._type)
		switch vbt {

		case GNET_SNMP_VARBIND_TYPE_NULL:
			value = new(VBT_Null)

		case GNET_SNMP_VARBIND_TYPE_OCTETSTRING:
			value = VBT_OctetString(union_ui8v_string(data.value, data.value_len))

		case GNET_SNMP_VARBIND_TYPE_OBJECTID:
			guint32_ptr := union_ui32v(data.value)
			value = VBT_ObjectID(gIntArrayOidString(guint32_ptr, data.value_len))

		case GNET_SNMP_VARBIND_TYPE_IPADDRESS:
			value = VBT_IPAddress(union_ui8v_ipaddress(data.value, data.value_len))

		case GNET_SNMP_VARBIND_TYPE_INTEGER32:
			value = VBT_Integer32(union_i32(data.value))

		case GNET_SNMP_VARBIND_TYPE_UNSIGNED32:
			value = VBT_Unsigned32(union_ui32(data.value))

		case GNET_SNMP_VARBIND_TYPE_COUNTER32:
			value = VBT_Counter32(union_ui32(data.value))

		case GNET_SNMP_VARBIND_TYPE_TIMETICKS:
			value = VBT_Timeticks(union_ui32(data.value))

		case GNET_SNMP_VARBIND_TYPE_OPAQUE:
			value = VBT_Opaque(union_ui8v_hexstring(data.value, data.value_len))

		case GNET_SNMP_VARBIND_TYPE_COUNTER64:
			value = VBT_Counter64(union_ui64(data.value))

		case GNET_SNMP_VARBIND_TYPE_NOSUCHOBJECT:
			value = new(VBT_NoSuchObject)

		case GNET_SNMP_VARBIND_TYPE_NOSUCHINSTANCE:
			value = new(VBT_NoSuchInstance)

		case GNET_SNMP_VARBIND_TYPE_ENDOFMIBVIEW:
			value = new(VBT_EndOfMibView)
		}
		result.Value = value
		results = append(results, result)

		// move on to next element in list
		out = out.next
	}
	panic(fmt.Sprintf("%s: convertResults(): fell out of for loop", libname()))
}

// libname returns the name of this library, for generating error messages.
func libname() string {
	return "gsnmpgo"
}
