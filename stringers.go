package gsnmp

// Copyright 2012 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// stringers.go contains stringers for C enums and other types. To help with
// the generation of the boilerplate code for the C enums,
// github.com/natefinch/gocog is used. AFTER EDITING any gocog sections
// (between gocog open and close square brackets), you MUST run:
//
// rm -f stringers.go_cog; $GOPATH/bin/gocog stringers.go
//

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
#define MAX_OIDS_STR_LEN 1000

static void
oid_to_str(GList *list, char result[MAX_OIDS_STR_LEN]) {
	result[0] = '\0';
	while (list) {
		// assume an oid isn't longer than 200 characters
		if (strlen(result) > (MAX_OIDS_STR_LEN - 200)) {
			// run out of space, just append ...
			strcat(result, "...");
			return;
		}

		GList *next = list->next;
		GNetSnmpVarBind *vb = list->data;

		gint i;
		// assume an oid octet isn't longer than 20 characters
		char some_digits[20];
		for (i = 0; i < vb->oid_len; i++) {
			strcat(result, ".");
			sprintf(some_digits, "%i", vb->oid[i]);
			strcat(result, some_digits);
			some_digits[0] = '\0';
		}
		if (next != NULL) {
			strcat(result, ":");
		}
		list = next;
	}
}
*/
import "C"

import (
	"fmt"
	"strconv"
	"unsafe"
)

// returns a string represention of OIDs in vbl (var bind list)
//
// C:
//     /usr/include/glib-2.0/glib/glist.h
//     typedef struct _GList GList;
//     struct _GList
//     {
// 	       gpointer data;
// 	       GList *next;
// 	       GList *prev;
//     };
func OidToString(vbl *_Ctype_GList) string {
	// allocate "char result[MAX_OIDS_STR_LEN]"
	const MAX_OIDS_STR_LEN = 1000 // same as C code define
	result_go := fmt.Sprintf("%"+strconv.Itoa(MAX_OIDS_STR_LEN)+"s", " ")
	var result_c *C.char = C.CString(result_go)
	defer C.free(unsafe.Pointer(result_c))

	C.oid_to_str(vbl, result_c)
	return C.GoString(result_c)
}

// Stringer for *_Ctype_GURI
//
// Example:
//    fmt.Println("ParseURI():", parsed_uri)
//
// C:
//    /usr/include/gnet-2.0/uri.h
//    struct _GURI
//    {
//      gchar* scheme;
//      gchar* userinfo;
//      gchar* hostname;
//      gint   port;
//      gchar* path;
//      gchar* query;
//      gchar* fragment;
//    };
//
func (parsed_uri *_Ctype_GURI) String() string {
	scheme := C.GoString((*C.char)(parsed_uri.scheme))
	userinfo := C.GoString((*C.char)(parsed_uri.userinfo))
	hostname := C.GoString((*C.char)(parsed_uri.hostname))
	port := int(parsed_uri.port)
	path := C.GoString((*C.char)(parsed_uri.path))
	query := C.GoString((*C.char)(parsed_uri.query))
	fragment := C.GoString((*C.char)(parsed_uri.fragment))

	result := "{"
	result += "scheme:" + scheme + " "
	result += "userinfo:" + userinfo + " "
	result += "hostname:" + hostname + " "
	result += "port:" + strconv.Itoa(port) + " "
	result += "path:" + path + " "
	result += "query:" + query + " "
	result += "fragment:" + fragment
	result += "}"

	return result
}

// Stringer for GString
//
//     glib/gstring.h
//     typedef struct _GString GString;
//     struct _GString {
//         gchar  *str;
//         gsize len;
//         gsize allocated_len;
//     };
//
// http://developer.gnome.org/glib/2.34/glib-Strings.html
// gchar *str - points to the character data. It may move as text is added. The
// str field is null-terminated and so can be used as an ordinary C string.
//
func (s _Ctype_GString) String() string {
	return C.GoString((*_Ctype_char)(s.str))
}

// Stringer for *_Ctype_GNetSnmp (a session)
//
// C:
//     gsnmp-0.3.0/src/session.h
//     typedef struct _GNetSnmp GNetSnmp;
//     struct _GNetSnmp {
//         GNetSnmpTAddress *taddress;
//         GURI             *uri;
//         gint32           error_status;
//         guint32          error_index;
//         guint            retries;        /* number of retries */
//         guint            timeout;        /* timeout in milliseconds */
//         GNetSnmpVersion  version;        /* message version */
//         GString          *ctxt_name;     /* context name */
//         GString          *sec_name;      /* security name */
//         GNetSnmpSecModel sec_model;      /* security model */
//         GNetSnmpSecLevel sec_level;      /* security level */
//         GNetSnmpDoneFunc done_callback;  /* what to call when complete */
//     }
//
func (s *_Ctype_GNetSnmp) String() string {

	error_status := strconv.Itoa(int(s.error_status))
	error_index := strconv.Itoa(int(s.error_index))
	retries := strconv.Itoa(int(s.retries))
	version := strconv.Itoa(int(s.version))

	result := "{"
	result += "taddress:" + fmt.Sprintf("%s", s.taddress) + " "
	result += "uri:" + fmt.Sprintf("%s", s.uri) + " "
	result += "error_status:" + error_status + " "
	result += "error_index:" + error_index + " "
	result += "retries:" + retries + " "
	// TODO need function to set version on session
	result += "version:" + version + " "
	result += "context_name:" + fmt.Sprintf("%s", s.ctxt_name) + " "
	result += "security_name:" + fmt.Sprintf("%s", s.sec_name) + " "
	result += "security_model:" + fmt.Sprintf("%s", s.sec_model) + " "
	result += "security_level:" + fmt.Sprintf("%s", s.sec_level) + " "
	result += "done_callback:UNIMPLEMENTED"
	result += "}"
	return result
}

// Stringer for *_Ctype_GNetSnmpTAddress
//
// Example:
//    fmt.Printf("GNetSnmpTAddress: %s", taddress)
//
// C:
//     gsnmp-0.3.0/src/transport.h
//     typedef struct {
//         GNetSnmpTDomain domain;
//         union {
//             GInetAddr *inetaddr;
//             gchar     *path;
//         };
//     } GNetSnmpTAddress;
func (t *_Ctype_GNetSnmpTAddress) String() string {
	name := C.gnet_snmp_taddress_get_short_name(t)
	return C.GoString((*_Ctype_char)(name))
}

///////////
// enums //
///////////

/*[[[gocog
package main
import ("github.com/soniah/gsnmp/enumconv")
func main() {
	ccode := "/usr/include/gsnmp/utils.h"
	vals := []string{"GNET_SNMP_URI_GET", "GNET_SNMP_URI_NEXT", "GNET_SNMP_URI_WALK"}
	enumconv.Write("UriType", "_Ctype_GNetSnmpUriType", vals, ccode, 0)
}
gocog]]]*/

// type and values for _Ctype_GNetSnmpUriType
//
type UriType int

const (
	GNET_SNMP_URI_GET UriType = iota
	GNET_SNMP_URI_NEXT
	GNET_SNMP_URI_WALK
)

// Stringer for _Ctype_GNetSnmpUriType
//
// C:
//    /usr/include/gsnmp/utils.h
//
func (uritype _Ctype_GNetSnmpUriType) String() string {
	switch UriType(uritype) {
	case GNET_SNMP_URI_GET:
		return "GNET_SNMP_URI_GET"
	case GNET_SNMP_URI_NEXT:
		return "GNET_SNMP_URI_NEXT"
	case GNET_SNMP_URI_WALK:
		return "GNET_SNMP_URI_WALK"
	}
	return "UNKNOWN _Ctype_GNetSnmpUriType"
}

//[[[end]]]

/*[[[gocog
package main
import ("github.com/soniah/gsnmp/enumconv")
func main() {
	ccode := "gsnmp-0.3.0/src/security.h"
	vals := []string{"GNET_SNMP_SECMODEL_ANY", "GNET_SNMP_SECMODEL_SNMPV1", "GNET_SNMP_SECMODEL_SNMPV2C", "GNET_SNMP_SECMODEL_SNMPV3"}
	enumconv.Write("SecModel", "_Ctype_GNetSnmpSecModel", vals, ccode, 0)
}
gocog]]]*/

// type and values for _Ctype_GNetSnmpSecModel
//
type SecModel int

const (
	GNET_SNMP_SECMODEL_ANY SecModel = iota
	GNET_SNMP_SECMODEL_SNMPV1
	GNET_SNMP_SECMODEL_SNMPV2C
	GNET_SNMP_SECMODEL_SNMPV3
)

// Stringer for _Ctype_GNetSnmpSecModel
//
// C:
//    gsnmp-0.3.0/src/security.h
//
func (secmodel _Ctype_GNetSnmpSecModel) String() string {
	switch SecModel(secmodel) {
	case GNET_SNMP_SECMODEL_ANY:
		return "GNET_SNMP_SECMODEL_ANY"
	case GNET_SNMP_SECMODEL_SNMPV1:
		return "GNET_SNMP_SECMODEL_SNMPV1"
	case GNET_SNMP_SECMODEL_SNMPV2C:
		return "GNET_SNMP_SECMODEL_SNMPV2C"
	case GNET_SNMP_SECMODEL_SNMPV3:
		return "GNET_SNMP_SECMODEL_SNMPV3"
	}
	return "UNKNOWN _Ctype_GNetSnmpSecModel"
}

//[[[end]]]

/*[[[gocog
package main
import ("github.com/soniah/gsnmp/enumconv")
func main() {
	ccode := "gsnmp-0.3.0/src/security.h"
	vals := []string{"GNET_SNMP_SECLEVEL_NANP", "GNET_SNMP_SECLEVEL_ANP", "GNET_SNMP_SECLEVEL_AP"}
	enumconv.Write("SecLevel", "_Ctype_GNetSnmpSecLevel", vals, ccode, 0)
}
gocog]]]*/

// type and values for _Ctype_GNetSnmpSecLevel
//
type SecLevel int

const (
	GNET_SNMP_SECLEVEL_NANP SecLevel = iota
	GNET_SNMP_SECLEVEL_ANP
	GNET_SNMP_SECLEVEL_AP
)

// Stringer for _Ctype_GNetSnmpSecLevel
//
// C:
//    gsnmp-0.3.0/src/security.h
//
func (seclevel _Ctype_GNetSnmpSecLevel) String() string {
	switch SecLevel(seclevel) {
	case GNET_SNMP_SECLEVEL_NANP:
		return "GNET_SNMP_SECLEVEL_NANP"
	case GNET_SNMP_SECLEVEL_ANP:
		return "GNET_SNMP_SECLEVEL_ANP"
	case GNET_SNMP_SECLEVEL_AP:
		return "GNET_SNMP_SECLEVEL_AP"
	}
	return "UNKNOWN _Ctype_GNetSnmpSecLevel"
}

//[[[end]]]

//
// GNetSnmp.error_status has type gint32, not GNetSnmpPduError <sigh>
//

/*[[[gocog
package main
import ("github.com/soniah/gsnmp/enumconv")
uunc main() {
	ccode := "gsnmp-0.3.0/src/pdu.h"
	vals := []string{"GNET_SNMP_PDU_ERR_DONE", "GNET_SNMP_PDU_ERR_PROCEDURE", "GNET_SNMP_PDU_ERR_INTERNAL", "GNET_SNMP_PDU_ERR_NORESPONSE", "GNET_SNMP_PDU_ERR_NOERROR", "GNET_SNMP_PDU_ERR_TOOBIG", "GNET_SNMP_PDU_ERR_NOSUCHNAME", "GNET_SNMP_PDU_ERR_BADVALUE", "GNET_SNMP_PDU_ERR_READONLY", "GNET_SNMP_PDU_ERR_GENERROR", "GNET_SNMP_PDU_ERR_NOACCESS", "GNET_SNMP_PDU_ERR_WRONGTYPE", "GNET_SNMP_PDU_ERR_WRONGLENGTH", "GNET_SNMP_PDU_ERR_WRONGENCODING", "GNET_SNMP_PDU_ERR_WRONGVALUE", "GNET_SNMP_PDU_ERR_NOCREATION", "GNET_SNMP_PDU_ERR_INCONSISTENTVALUE", "GNET_SNMP_PDU_ERR_RESOURCEUNAVAILABLE", "GNET_SNMP_PDU_ERR_COMMITFAILED", "GNET_SNMP_PDU_ERR_UNDOFAILED", "GNET_SNMP_PDU_ERR_AUTHORIZATIONERROR", "GNET_SNMP_PDU_ERR_NOTWRITABLE", "GNET_SNMP_PDU_ERR_INCONSISTENTNAME"}
	enumconv.Write("PduError", "_Ctype_gint32", vals, ccode, -4)
}
gocog]]]*/

// type and values for _Ctype_gint32
//
type PduError int

const (
	GNET_SNMP_PDU_ERR_DONE PduError = iota - 4
	GNET_SNMP_PDU_ERR_PROCEDURE
	GNET_SNMP_PDU_ERR_INTERNAL
	GNET_SNMP_PDU_ERR_NORESPONSE
	GNET_SNMP_PDU_ERR_NOERROR
	GNET_SNMP_PDU_ERR_TOOBIG
	GNET_SNMP_PDU_ERR_NOSUCHNAME
	GNET_SNMP_PDU_ERR_BADVALUE
	GNET_SNMP_PDU_ERR_READONLY
	GNET_SNMP_PDU_ERR_GENERROR
	GNET_SNMP_PDU_ERR_NOACCESS
	GNET_SNMP_PDU_ERR_WRONGTYPE
	GNET_SNMP_PDU_ERR_WRONGLENGTH
	GNET_SNMP_PDU_ERR_WRONGENCODING
	GNET_SNMP_PDU_ERR_WRONGVALUE
	GNET_SNMP_PDU_ERR_NOCREATION
	GNET_SNMP_PDU_ERR_INCONSISTENTVALUE
	GNET_SNMP_PDU_ERR_RESOURCEUNAVAILABLE
	GNET_SNMP_PDU_ERR_COMMITFAILED
	GNET_SNMP_PDU_ERR_UNDOFAILED
	GNET_SNMP_PDU_ERR_AUTHORIZATIONERROR
	GNET_SNMP_PDU_ERR_NOTWRITABLE
	GNET_SNMP_PDU_ERR_INCONSISTENTNAME
)

// Stringer for _Ctype_gint32
//
// C:
//    gsnmp-0.3.0/src/pdu.h
//
func (pduerror _Ctype_gint32) String() string {
	switch PduError(pduerror) {
	case GNET_SNMP_PDU_ERR_DONE:
		return "GNET_SNMP_PDU_ERR_DONE"
	case GNET_SNMP_PDU_ERR_PROCEDURE:
		return "GNET_SNMP_PDU_ERR_PROCEDURE"
	case GNET_SNMP_PDU_ERR_INTERNAL:
		return "GNET_SNMP_PDU_ERR_INTERNAL"
	case GNET_SNMP_PDU_ERR_NORESPONSE:
		return "GNET_SNMP_PDU_ERR_NORESPONSE"
	case GNET_SNMP_PDU_ERR_NOERROR:
		return "GNET_SNMP_PDU_ERR_NOERROR"
	case GNET_SNMP_PDU_ERR_TOOBIG:
		return "GNET_SNMP_PDU_ERR_TOOBIG"
	case GNET_SNMP_PDU_ERR_NOSUCHNAME:
		return "GNET_SNMP_PDU_ERR_NOSUCHNAME"
	case GNET_SNMP_PDU_ERR_BADVALUE:
		return "GNET_SNMP_PDU_ERR_BADVALUE"
	case GNET_SNMP_PDU_ERR_READONLY:
		return "GNET_SNMP_PDU_ERR_READONLY"
	case GNET_SNMP_PDU_ERR_GENERROR:
		return "GNET_SNMP_PDU_ERR_GENERROR"
	case GNET_SNMP_PDU_ERR_NOACCESS:
		return "GNET_SNMP_PDU_ERR_NOACCESS"
	case GNET_SNMP_PDU_ERR_WRONGTYPE:
		return "GNET_SNMP_PDU_ERR_WRONGTYPE"
	case GNET_SNMP_PDU_ERR_WRONGLENGTH:
		return "GNET_SNMP_PDU_ERR_WRONGLENGTH"
	case GNET_SNMP_PDU_ERR_WRONGENCODING:
		return "GNET_SNMP_PDU_ERR_WRONGENCODING"
	case GNET_SNMP_PDU_ERR_WRONGVALUE:
		return "GNET_SNMP_PDU_ERR_WRONGVALUE"
	case GNET_SNMP_PDU_ERR_NOCREATION:
		return "GNET_SNMP_PDU_ERR_NOCREATION"
	case GNET_SNMP_PDU_ERR_INCONSISTENTVALUE:
		return "GNET_SNMP_PDU_ERR_INCONSISTENTVALUE"
	case GNET_SNMP_PDU_ERR_RESOURCEUNAVAILABLE:
		return "GNET_SNMP_PDU_ERR_RESOURCEUNAVAILABLE"
	case GNET_SNMP_PDU_ERR_COMMITFAILED:
		return "GNET_SNMP_PDU_ERR_COMMITFAILED"
	case GNET_SNMP_PDU_ERR_UNDOFAILED:
		return "GNET_SNMP_PDU_ERR_UNDOFAILED"
	case GNET_SNMP_PDU_ERR_AUTHORIZATIONERROR:
		return "GNET_SNMP_PDU_ERR_AUTHORIZATIONERROR"
	case GNET_SNMP_PDU_ERR_NOTWRITABLE:
		return "GNET_SNMP_PDU_ERR_NOTWRITABLE"
	case GNET_SNMP_PDU_ERR_INCONSISTENTNAME:
		return "GNET_SNMP_PDU_ERR_INCONSISTENTNAME"
	}
	return "UNKNOWN _Ctype_gint32"
}

//[[[end]]]
