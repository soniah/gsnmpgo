package gsnmpgo

// Copyright 2013 Sonia Hamilton <sonia@snowfrog.net>. All rights
// reserved.  Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

import (
	"fmt"
)

type Varbinder interface {
	Integer() int64
	fmt.Stringer
}

// GNET_SNMP_VARBIND_TYPE_NULL
type VBT_Null struct{}

func (r VBT_Null) Integer() int64 {
	return 0
}

func (r VBT_Null) String() string {
	return "NULL"
}

// GNET_SNMP_VARBIND_TYPE_OCTETSTRING
type VBT_OctetString string

func (r VBT_OctetString) Integer() int64 {
	return 0
}

func (r VBT_OctetString) String() string {
	return fmt.Sprintf("%s", string(r))
}

// GNET_SNMP_VARBIND_TYPE_OBJECTID
type VBT_ObjectID string

func (r VBT_ObjectID) Integer() int64 {
	return 0
}

func (r VBT_ObjectID) String() string {
	return fmt.Sprintf("%s", string(r))
}

// GNET_SNMP_VARBIND_TYPE_IPADDRESS
type VBT_IPAddress string

func (r VBT_IPAddress) Integer() int64 {
	// TODO convert ip address *back* to a number. Or store as
	// a number and convert here to dotted form here??
	return 0
}

func (r VBT_IPAddress) String() string {
	return fmt.Sprintf("%s", string(r))
}

// GNET_SNMP_VARBIND_TYPE_INTEGER32
type VBT_Integer32 int32

func (r VBT_Integer32) Integer() int64 {
	return int64(r)
}

func (r VBT_Integer32) String() string {
	return fmt.Sprintf("%d", r)
}

// GNET_SNMP_VARBIND_TYPE_UNSIGNED32
type VBT_Unsigned32 uint32

func (r VBT_Unsigned32) Integer() int64 {
	return int64(r)
}

func (r VBT_Unsigned32) String() string {
	return fmt.Sprintf("%d", r)
}

// GNET_SNMP_VARBIND_TYPE_COUNTER32
type VBT_Counter32 uint32

func (r VBT_Counter32) Integer() int64 {
	return int64(r)
}

func (r VBT_Counter32) String() string {
	return fmt.Sprintf("%d", r)
}

// GNET_SNMP_VARBIND_TYPE_TIMETICKS
type VBT_Timeticks uint32

func (r VBT_Timeticks) Integer() int64 {
	return int64(r)
}

func (r VBT_Timeticks) String() string {
	ticks := uint32(r)
	if ticks == uint32(0) {
		return "0:0:00:00.00"
	}

	days := int(ticks / (24 * 60 * 60 * 100))
	ticks %= (24 * 60 * 60 * 100)
	hours := int(ticks / (60 * 60 * 100))
	ticks %= (60 * 60 * 100)
	minutes := int(ticks / (60 * 100))
	ticks %= (60 * 100)
	seconds := float64(ticks / 100)
	return fmt.Sprintf("%d days, %d:%02d:%05.02f", days, hours, minutes, seconds)
}

// GNET_SNMP_VARBIND_TYPE_OPAQUE
type VBT_Opaque string

func (r VBT_Opaque) Integer() int64 {
	return 0
}

func (r VBT_Opaque) String() string {
	return fmt.Sprintf("%s", string(r))
}

// GNET_SNMP_VARBIND_TYPE_COUNTER64
type VBT_Counter64 uint64

func (r VBT_Counter64) Integer() int64 {
	// TODO bzzzt fail uint64 -> int64.
	// Should Integer() return a uint64? Or a longer value??
	return int64(r)
}

func (r VBT_Counter64) String() string {
	return fmt.Sprintf("%d", r)
}

// GNET_SNMP_VARBIND_TYPE_NOSUCHOBJECT
type VBT_NoSuchObject struct{}

func (r VBT_NoSuchObject) Integer() int64 {
	return 0
}

func (r VBT_NoSuchObject) String() string {
	return "No Such Object available on this agent at this OID" // same as netsnmp
}

// GNET_SNMP_VARBIND_TYPE_NOSUCHINSTANCE
type VBT_NoSuchInstance struct{}

func (r VBT_NoSuchInstance) Integer() int64 {
	return 0
}

func (r VBT_NoSuchInstance) String() string {
	return "No Such Instance"
}

// GNET_SNMP_VARBIND_TYPE_ENDOFMIBVIEW
type VBT_EndOfMibView struct{}

func (r VBT_EndOfMibView) Integer() int64 {
	return 0
}

func (r VBT_EndOfMibView) String() string {
	return "End of MIB View"
}
