package main

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
	"fmt"
	"unsafe"
)

func Get() {
	uri := `snmp://public@device//1.3.6.1.2.1.1.*`
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))
	//uri2 := C.gnet_snmp_parse_uri(curi, &error); // TODO &error
}

func main() {
	fmt.Println("hello world")
}

/// cgo CFLAGS: -I/usr/include/glib-2.0
/// cgo pkg-config: --cflags glib-2.0 gsnmp
/// cgo pkg-config: gsnmp glib-2.0
/// cgo CFLAGS: -pthread -I/usr/include/glib-2.0 -I/usr/lib/x86_64-linux-gnu/glib-2.0/include -I/usr/include/gsnmp -I/usr/include/gnet-2.0 -I/usr/lib/gnet-2.0/include/
/// cgo LDFLAGS: -lgsnmp
