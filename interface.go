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

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Varbinder interface {
	// Integer() needs to handle both signed numbers (int32), as well as
	// unsigned int 64 (uint64). Therefore it returns a *big.Int.
	Integer() *big.Int
	fmt.Stringer
}

// GNET_SNMP_VARBIND_TYPE_NULL
type VBT_Null struct{}

func (r VBT_Null) Integer() *big.Int {
	return big.NewInt(0)
}

func (r VBT_Null) String() string {
	return "NULL"
}

// GNET_SNMP_VARBIND_TYPE_OCTETSTRING
type VBT_OctetString string

func (r VBT_OctetString) Integer() *big.Int {
	return big.NewInt(0)
}

func (r VBT_OctetString) String() string {
	return fmt.Sprintf("%s", string(r))
}

// GNET_SNMP_VARBIND_TYPE_OBJECTID
type VBT_ObjectID string

func (r VBT_ObjectID) Integer() *big.Int {
	return big.NewInt(0)
}

func (r VBT_ObjectID) String() string {
	return fmt.Sprintf("%s", string(r))
}

// GNET_SNMP_VARBIND_TYPE_IPADDRESS
type VBT_IPAddress string

func (r VBT_IPAddress) Integer() *big.Int {
	// convert ip address to a uint32 ie it's numeric form
	// alternatively, could return a net/IP from union_ui8v_ipaddress
	var result uint32
	for _, octet := range strings.Split(string(r), ".") {
		result <<= 8
		n, _ := strconv.ParseUint(octet, 10, 8)
		result += uint32(n)

	}
	return big.NewInt(int64(result))
}

func (r VBT_IPAddress) String() string {
	return fmt.Sprintf("%s", string(r))
}

// GNET_SNMP_VARBIND_TYPE_INTEGER32
type VBT_Integer32 int32

func (r VBT_Integer32) Integer() *big.Int {
	return big.NewInt(int64(r))
}

func (r VBT_Integer32) String() string {
	return fmt.Sprintf("%d", r)
}

// GNET_SNMP_VARBIND_TYPE_UNSIGNED32
type VBT_Unsigned32 uint32

func (r VBT_Unsigned32) Integer() *big.Int {
	return big.NewInt(int64(r))
}

func (r VBT_Unsigned32) String() string {
	return fmt.Sprintf("%d", r)
}

// GNET_SNMP_VARBIND_TYPE_COUNTER32
type VBT_Counter32 uint32

func (r VBT_Counter32) Integer() *big.Int {
	return big.NewInt(int64(r))
}

func (r VBT_Counter32) String() string {
	return fmt.Sprintf("%d", r)
}

// GNET_SNMP_VARBIND_TYPE_TIMETICKS
type VBT_Timeticks uint32

func (r VBT_Timeticks) Integer() *big.Int {
	return big.NewInt(int64(r))
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

func (r VBT_Opaque) Integer() *big.Int {
	return big.NewInt(0)
}

func (r VBT_Opaque) String() string {
	return fmt.Sprintf("%s", string(r))
}

// GNET_SNMP_VARBIND_TYPE_COUNTER64
type VBT_Counter64 uint64

func (r VBT_Counter64) Integer() *big.Int {
	// math.Big constructor only accepts int64; see "Issue 4389" below
	return uint64ToBigInt(uint64(r))
}

func (r VBT_Counter64) String() string {
	return fmt.Sprintf("%d", r)
}

// GNET_SNMP_VARBIND_TYPE_NOSUCHOBJECT
type VBT_NoSuchObject struct{}

func (r VBT_NoSuchObject) Integer() *big.Int {
	return big.NewInt(0)
}

func (r VBT_NoSuchObject) String() string {
	return "No Such Object available on this agent at this OID" // same as netsnmp
}

// GNET_SNMP_VARBIND_TYPE_NOSUCHINSTANCE
type VBT_NoSuchInstance struct{}

func (r VBT_NoSuchInstance) Integer() *big.Int {
	return big.NewInt(0)
}

func (r VBT_NoSuchInstance) String() string {
	return "No Such Instance"
}

// GNET_SNMP_VARBIND_TYPE_ENDOFMIBVIEW
type VBT_EndOfMibView struct{}

func (r VBT_EndOfMibView) Integer() *big.Int {
	return big.NewInt(0)
}

func (r VBT_EndOfMibView) String() string {
	return "End of MIB View"
}

// Issue 4389: math/big: add SetUint64 and Uint64 functions to *Int
//
// uint64ToBigInt copied from: http://github.com/cznic/mathutil/blob/master/mathutil.go#L341
//
// replace with Uint64ToBigInt or equivalent when using Go 1.1

var uint64ToBigIntDelta big.Int

func init() {
	uint64ToBigIntDelta.SetBit(&uint64ToBigIntDelta, 63, 1)
}

func uint64ToBigInt(n uint64) *big.Int {
	if n <= math.MaxInt64 {
		return big.NewInt(int64(n))
	}

	y := big.NewInt(int64(n - uint64(math.MaxInt64) - 1))
	return y.Add(y, &uint64ToBigIntDelta)
}
