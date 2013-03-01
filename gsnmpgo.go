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
#include "c_bridge.h"
*/
import "C"

import (
	"code.google.com/p/tcgl/applog"
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"strconv"
	"strings"
	"unsafe"
)

// the maximum number of paths that can be in a single uri
const MAX_URI_COUNT = 50

var Debug bool // global debugging flag

// Struct of parameters to pass to Query
type QueryParams struct {
	Uri     string
	Version SnmpVersion
	Timeout int // timeout in milliseconds
	Retries int // number of retries
	// Nonrep and Maxrep will be used by v2c BULK GETs.
	// From O'Reilly "Essential SNMP": "nonrep is the number of scalar
	// objects that this command will return; rep is the number of
	// instances of each nonscalar object that the command will return."
	Nonrep int
	Maxrep int
	// if Tree is non-nil, it will be used for appending Query()
	// results eg when doing two GETs in a row
	Tree *llrb.Tree
	//
	// GetBulk()
	//
	Community string
	IP        net.IP
	Port      uint16
	Oids      []string
}

// A single result, used as an Item in the llrb tree
type QueryResult struct {
	Oid   string
	Value Varbinder
}

// Query takes a URI in RFC 4088 format, does an SNMP query and returns the results.
func Query(params *QueryParams) (results *llrb.Tree, err error) {

	parsed_uri, err := parseURI(params.Uri)
	if Debug {
		applog.Debugf("parsed_uri: %s\n\n", parsed_uri)
	}
	if err != nil {
		return nil, err
	}

	path := C.GoString((*C.char)(parsed_uri.path))
	if Debug {
		applog.Warningf("number of incoming uris: %d", uriCount(path))
	}
	if err := uriCountMaxed(path, MAX_URI_COUNT); err != nil {
		return nil, err
	}

	vbl, uritype, err := parsePath(params.Uri, parsed_uri)
	defer uriDelete(parsed_uri)
	if Debug {
		applog.Debugf("vbl, uritype: %s, %s", gListOidsString(vbl), uritype)
	}
	if err != nil {
		return nil, err
	}

	session, err := newUri(params, parsed_uri)
	/*
		causing <undefined symbol: gnet_snmp_taddress_get_short_name>
		if Debug {
			applog.Warningf("session: %s\n\n", session)
		}
	*/
	if err != nil {
		return nil, err
	}

	vbl_results, err := querySync(session, vbl, uritype, params.Version)
	defer vblDelete(vbl_results)
	if err != nil {
		return nil, err
	}
	return convertResults(params, vbl_results), nil
}

// ------------------- other functions in alphabetical order --------------------

// convertResults converts C results to a Go struct.
func convertResults(params *QueryParams, out *_Ctype_GList) (results *llrb.Tree) {

	var out_count int

	// create or re-use an existing llrb Tree
	if params.Tree == nil {
		results = llrb.New(LessOID)
	} else {
		results = params.Tree
	}

	for {
		if out == nil {
			// finished
			if Debug {
				applog.Warningf("number of results converted: %d", out_count)
			}
			return results
		}

		// another result: initialise
		out_count++
		data := (*C.GNetSnmpVarBind)(out.data)
		oid := gIntArrayOidString(data.oid, data.oid_len)
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
			value = VBT_ObjectID("." + gIntArrayOidString(guint32_ptr, data.value_len))

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
		result := QueryResult{Oid: oid, Value: value}
		results.ReplaceOrInsert(result)

		// move on to next element in list
		out = out.next
	}
	panic(fmt.Sprintf("%s: convertResults(): fell out of for loop", libname()))
}

// Dump is a convenience function for printing the results of a Query.
func Dump(results *llrb.Tree) {
	if results == nil {
		fmt.Println("Dump: results are NIL")
		return
	}
	fmt.Println("Dump:")
	ch := results.IterAscend()
	for {
		r := <-ch
		if r == nil {
			break
		}
		result := r.(QueryResult)
		fmt.Printf("oid, type: %s, %T\n", result.Oid, result.Value)
		fmt.Printf("INTEGER: %d\n", result.Value.Integer())
		fmt.Printf("STRING : %s\n", result.Value)
		fmt.Println()
	}
}

// LessOID is the LessFunc for GoLLRB
//
// It returns true if oid a is less than oid b.
func LessOID(astruct, bstruct interface{}) bool {
	a := astruct.(QueryResult).Oid
	b := bstruct.(QueryResult).Oid

	if a == "" && b == "" {
		return false
	} else if a == "" {
		return true
	} else if b == "" {
		return false
	}

	a_splits := strings.Split(a, ".")
	b_splits := strings.Split(b, ".")
	len_b := len(b_splits)

	for i, a_digit := range a_splits {
		if i > len_b-1 {
			return false
		}
		a_num, _ := strconv.Atoi(a_digit)
		b_num, _ := strconv.Atoi(b_splits[i])
		if a_num < b_num {
			return true
		} else if i == len_b-1 {
			return false
		}
	}
	return true
}

// libname returns the name of this library, for generating error messages.
func libname() string {
	return "gsnmpgo"
}

// NewDefaultParams returns QueryParams with sensible default values
func NewDefaultParams(uri string) *QueryParams {
	return &QueryParams{
		Uri:     uri,
		Version: GNET_SNMP_V2C,
		Timeout: 200,
		Retries: 3,
		// From O'Reilly "Essential SNMP": "nonrep is the number of scalar
		// objects that this command will return; rep is the number of
		// instances of each nonscalar object that the command will return. If
		// you omit this option the default values of nonrep and rep, 1 and
		// 100, respectively, will be used." So use these defaults for the
		// moment.
		Nonrep: 1,
		Maxrep: 100,
	}
}

// newUri creates a session from a parsed uri.
func newUri(params *QueryParams, parsed_uri *_Ctype_GURI) (session *_Ctype_GNetSnmp, err error) {

	var gerror *C.GError
	session = C.gnet_snmp_new_uri(parsed_uri, &gerror)

	// error handling
	if gerror != nil {
		err_string := C.GoString((*_Ctype_char)(gerror.message))
		C.g_clear_error(&gerror)
		return session, fmt.Errorf("%s: newUri(): %s", libname(), err_string)
	}
	if session == nil {
		return session, fmt.Errorf("%s: newUri(): unable to create session", libname())
	}

	if params.Version == GNET_SNMP_V1 { // default in library is v2c
		C.gnet_snmp_set_version(session, 0)
	}
	C.gnet_snmp_set_timeout(session, (_Ctype_guint)(params.Timeout))
	C.gnet_snmp_set_retries(session, (_Ctype_guint)(params.Timeout))

	return session, nil
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
		return vbl, uritype, fmt.Errorf("%s: parsePath(): %s: <%s>", libname(), err_string, uri)
	}
	return vbl, uritype, nil
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
		return nil, fmt.Errorf("%s: parseURI(): invalid snmp uri: %s", libname(), uri)
	}
	return parsed_uri, nil
}

// PartitionAllP - returns true when dividing a slice into
// partition_size lengths, including last partition which may be smaller
// than partition_size.
//
// For example for a slice of 8 items to be broken into partitions of
// length 3, PartitionAllP returns true for the current_position having
// the following values:
//
// 0  1  2  3  4  5  6  7
//       T        T     T
//
// 'P' stands for Predicate (like foo? in Ruby, foop in Lisp)
//
func PartitionAllP(current_position, partition_size, slice_length int) bool {
	// TODO should handle partition_size > slice_length, slice_length < 0
	if current_position < 0 || current_position >= slice_length {
		return false
	}
	if partition_size == 1 { // redundant, but an obvious optimisation
		return true
	}
	if current_position%partition_size == partition_size-1 {
		return true
	}
	if current_position == slice_length-1 {
		return true
	}
	return false
}

// querySync - do an gsnmp library sync_* query
//
// Results are returned in C form, use convertResults() to convert to a Go struct.
func querySync(session *_Ctype_GNetSnmp, vbl *_Ctype_GList, uritype _Ctype_GNetSnmpUriType,
	version SnmpVersion) (*_Ctype_GList, error) {
	var gerror *C.GError
	var out *_Ctype_GList

	if Debug {
		applog.Debugf("Starting a %s", uritype)
	}
	switch UriType(uritype) {
	case GNET_SNMP_URI_GET:
		out = C.gnet_snmp_sync_get(session, vbl, &gerror)
	case GNET_SNMP_URI_NEXT:
		out = C.gnet_snmp_sync_getnext(session, vbl, &gerror)
	case GNET_SNMP_URI_WALK:
		out = C.gnet_snmp_sync_walk(session, vbl, &gerror)
		/* TODO gnet_snmp_sync_walk is just a series of 'getnexts'
		if version == GNET_SNMP_V1 {
			out = C.gnet_snmp_sync_walk(session, vbl, &gerror)
		} else {
			// do a proper bulkwalk
		}
		*/
	default:
		return nil, fmt.Errorf("%s: querySync(): unknown uritype", libname())
	}

	/*
		Originally error handling was done at this point, like
		gsnmp-0.3.0/examples/gsnmp-get.c. However in production too many results
		were being discarded. Hence just return out, and convertResults() will
		convert any errors in out to nil values.
	*/

	return out, nil
}

// uriCount returns a count of the number of uri's in the path
func uriCount(path string) int {
	left_paren := strings.Index(path, "(")
	right_paren := strings.Index(path, ")")
	if left_paren < 0 || right_paren < 0 {
		return -1
	}
	uris := path[left_paren+1 : right_paren]
	return len(strings.Split(uris, ","))
}

// uriCountMaxed returns an error if there are more uri's in path than max
func uriCountMaxed(path string, max int) (err error) {
	if uri_count := uriCount(path); uri_count > max {
		return fmt.Errorf("number of uris is greater than max (%d/%d) in path %s", uri_count, max, path)
	}
	return nil
}

// uriDelete frees the memory used by a parsed_uri.
//
// A defered call to uriDelete should be made after parsePath().
func uriDelete(parsed_uri *_Ctype_GURI) {
	C.gnet_uri_delete(parsed_uri)
}

// vblDelete frees the memory used by a var bind list.
//
// A deferred call to vblDelete should be made after call to
// gnet_snmp_sync_get (or similar).
func vblDelete(vbl *_Ctype_GList) {
	C.vbl_delete(vbl)
}
