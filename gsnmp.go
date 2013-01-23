// gsnmp is a go/cgo wrapper around gsnmp. It is under development,
// therefor API's may/will change, and doco/error handling/tests are
// minimal.
//
// Copyright 2012 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// glib typedefs - http://developer.gnome.org/glib/2.35/glib-Basic-Types.html
// glib tutorial - http://www.dlhoffman.com/publiclibrary/software/gtk+-html-docs/gtk_tut-17.html
// gsnmp sourcecode browser - http://sourcecodebrowser.com/gsnmp/0.3.0/index.html
//
package gsnmp

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
	"fmt"
	"unsafe"
)

// libname returns the name of this library
//
// libname is used for generating error messages
//
func libname() string {
	return "gsnmp"
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
//    parsed_uri, err := gsnmp.ParseURI(uri)
//    if err != nil {
//    	fmt.Println(err)
//    	os.Exit(1)
//    }
//    fmt.Println("ParseURI():", parsed_uri)
//
// C:
//    GURI*
//    gnet_snmp_parse_uri(const gchar *uri_string, GError **error)
//
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

// ParsePath: gnet_snmp_parse_path
//
// uritype will default to GNET_SNMP_URI_GET. If the uri ends in:
//
// * uritype will be GNET_SNMP_URI_WALK
// + uritype will be GNET_SNMP_URI_NEXT
//
// See RFC 4088 "Uniform Resource Identifier (URI) Scheme for the Simple
// Network Management Protocol (SNMP)" for further documentation.
//
// C:
//    gsnmp-0.3.0/src/utils.c
//    gboolean
//    gnet_snmp_parse_path(const gchar *path,
//    		     GList **vbl,
//    		     GNetSnmpUriType *type,
//    		     GError **error)
//
func ParsePath(uri string, parsed_uri *_Ctype_GURI) (vbl *_Ctype_GList, uritype _Ctype_GNetSnmpUriType, err error) {
	var gerror *C.GError
	rv := C.gnet_snmp_parse_path(parsed_uri.path, &vbl, &uritype, &gerror)
	if rv == 0 {
		err_string := C.GoString((*_Ctype_char)(gerror.message))
		return vbl, uritype, fmt.Errorf("%s: %s: <%s>", libname(), err_string, uri)
	}
	return vbl, uritype, nil
}

// UriDelete frees the memory used by a parsed_uri
//
// A defered call to UriDelete should be made after ParsePath()
//
func UriDelete(parsed_uri *_Ctype_GURI) {
	C.gnet_uri_delete(parsed_uri)
}

// SnmpNewUri creates a session from a parsed uri
//
// C:
//     gsnmp-0.3.0/src/session.c
//     GNetSnmp*
//     gnet_snmp_new_uri(const GURI *uri, GError **error)
//
func SnmpNewUri(uri string, parsed_uri *_Ctype_GURI) (session *_Ctype_GNetSnmp, err error) {
	var gerror *C.GError
	session = C.gnet_snmp_new_uri(parsed_uri, &gerror)
	if gerror != nil {
		err_string := C.GoString((*_Ctype_char)(gerror.message))
		return session, fmt.Errorf("%s: %s", libname(), err_string)
	}
	if session == nil {
		return session, fmt.Errorf("%s: unable to create session", libname())
	}
	return session, nil
}
