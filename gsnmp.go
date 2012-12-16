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
*/
import "C"

import (
	"strconv"
	"unsafe"
)

/*
glib typedefs - see http://developer.gnome.org/glib/2.35/glib-Basic-Types.html
glib tutorial - see http://www.dlhoffman.com/publiclibrary/software/gtk+-html-docs/gtk_tut-17.html
*/

// ParseURI: gnet_snmp_parse_uri
// GURI*
// gnet_snmp_parse_uri(const gchar *uri_string, GError **error)
func ParseURI(uri string) (parsed_uri *_Ctype_GURI) {
	curi := (*C.gchar)(C.CString(uri))
	defer C.free(unsafe.Pointer(curi))

	var gerror *C.GError
	parsed_uri = C.gnet_snmp_parse_uri(curi, &gerror) // TODO handle error
	return
}

// Stringer for *_Ctype_GURI
func (parsed_uri *_Ctype_GURI) String() string {
	/*
		see /usr/include/gnet-2.0/uri.h for GURI
		struct _GURI
		{
		  gchar* scheme;
		  gchar* userinfo;
		  gchar* hostname;
		  gint   port;
		  gchar* path;
		  gchar* query;
		  gchar* fragment;
		};
	*/
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
	result += "fragment:" + fragment + "}"

	return result
}

// ParsePath: gnet_snmp_parse_path
//    gboolean
//    gnet_snmp_parse_path(const gchar *path,
//    		     GList **vbl,
//    		     GNetSnmpUriType *type,
//    		     GError **error)
func ParsePath(parsed_uri *_Ctype_GURI) (vbl *_Ctype_GList, _type _Ctype_GNetSnmpUriType, result bool) {
	path := parsed_uri.path
	var gerror *C.GError
	rv := C.gnet_snmp_parse_path(path, &vbl, &_type, &gerror) // TODO handle error
	if rv != 0 {
		result = true
	}
	return
}

// Stringer for *_Ctype_GList
/*
see /usr/include/glib-2.0/glib/glist.h for GList

typedef struct _GList GList;

struct _GList
{
  gpointer data;
  GList *next;
  GList *prev;
};

typedef void* gpointer;
*/
func (glist *_Ctype_GList) String() string {
	// what does it means to print a "void* data"??
	// follow down next's and print each data
	return "TODO"
}

type UriType int

const (
	GNET_SNMP_URI_GET UriType = iota
	GNET_SNMP_URI_NEXT
	GNET_SNMP_URI_WALK
)

// Stringer for _Ctype_GNetSnmpUriType
//
// /usr/include/gsnmp/utils.h
//
// typedef enum
// {
//     GNET_SNMP_URI_GET,
//     GNET_SNMP_URI_NEXT,
//     GNET_SNMP_URI_WALK
// } GNetSnmpUriType;
func (uritype _Ctype_GNetSnmpUriType) String() (result string) {
	switch UriType(uritype) {
	case GNET_SNMP_URI_GET:
		return "GNET_SNMP_URI_GET"
	case GNET_SNMP_URI_NEXT:
		return "GNET_SNMP_URI_NEXT"
	case GNET_SNMP_URI_WALK:
		return "GNET_SNMP_URI_WALK"
	}
	return "UNKNOWN GNetSnmpUriType"
}
