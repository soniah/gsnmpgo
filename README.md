gsnmpgo
======

Update 8/Mar/13
---------------

Work on gsnmpgo has halted. Making the gsnmp C library multithreaded was
proving too time consuming. Use http://github.com/soniah/gosnmp instead.

Update 4/Dec/13
---------------

Geert-Johan Riemer commented:

I saw your notice on gsnmpgo "Threading: either gsnmp isn't totally thread
safe, or I'm making errors with my C calls from Go (more likely). I'm getting
aborts with this message:"

I think this might be solved by locking the requesting goroutine to it's thread
using runtime.LockOSThread (http://golang.org/pkg/runtime/#LockOSThread). And
make sure requests are always made from that specific goroutine.

(I currently don't have time to work on gsnmpgo, but I thought I'd add this
comment for future reference. Sonia.)

About
-----

gsnmpgo is an snmp library for Go; it uses Go/CGo to wrap gsnmp.

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

See documentation.

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
