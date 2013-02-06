gsnmpgo
======

gsnmpgo is a go/cgo wrapper around gsnmp; it currently provides support
for snmp v1 and v2c, and snmp get, snmp getnext, and snmp walk.

gsnmpgo is pre 1.0, therefore API's may change, and tests are minimal.

Sonia Hamilton, sonia@snowfrog.net, http://www.snowfrog.net.

Documentation
-------------

See http://godoc.org/github.com/soniah/gsnmpgo or your local
go doc server for full documentation:

    cd $GOPATH
    godoc -http=:6060 &
    $preferred_browser http://localhost:6060/pkg &

Installation
------------

See documentation.

Issues
------

* snmp walk results are returned as a map ie unordered. My intention is to
  replace the map with an "ordered map" equivalent eg
  github.com/petar/GoLLRB/llrb or code.google.com/p/biogo.llrb/

* snmp walks aren't being done as "GETBULK" when using snmp v2c, rather
  repeated getnexts are being done ie like snmp v1. The gsnmp C library doesn't
  implement GETBULK directly, my intention is to write the C code to do it.

* testing is still a WIP (work in progress)

License
-------

gsnmpgo is a go/cgo wrapper around gsnmp.

Copyright (C) 2012-2013 Sonia Hamilton sonia@snowfrog.net.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

Note on License
---------------

The preferred way to release Go code is under a BSD/MIT/Apache license.
However gsnmp is released under the GPL, therefore gsnmpgo must be
released under the GPL. See http://www.gnu.org/licenses/gpl-faq.html#IfLibraryIsGPL
