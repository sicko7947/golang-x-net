// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6_test

import (
	"net"
	"reflect"
	"runtime"
	"testing"

	"github.com/detunized/golang-x-net/ipv6"
	"github.com/detunized/golang-x-net/nettest"
)

var icmpStringTests = []struct {
	in  ipv6.ICMPType
	out string
}{
	{ipv6.ICMPTypeDestinationUnreachable, "destination unreachable"},

	{256, "<nil>"},
}

func TestICMPString(t *testing.T) {
	for _, tt := range icmpStringTests {
		s := tt.in.String()
		if s != tt.out {
			t.Errorf("got %s; want %s", s, tt.out)
		}
	}
}

func TestICMPFilter(t *testing.T) {
	switch runtime.GOOS {
	case "fuchsia", "hurd", "js", "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	}

	var f ipv6.ICMPFilter
	for _, toggle := range []bool{false, true} {
		f.SetAll(toggle)
		for _, typ := range []ipv6.ICMPType{
			ipv6.ICMPTypeDestinationUnreachable,
			ipv6.ICMPTypeEchoReply,
			ipv6.ICMPTypeNeighborSolicitation,
			ipv6.ICMPTypeDuplicateAddressConfirmation,
		} {
			f.Accept(typ)
			if f.WillBlock(typ) {
				t.Errorf("ipv6.ICMPFilter.Set(%v, false) failed", typ)
			}
			f.Block(typ)
			if !f.WillBlock(typ) {
				t.Errorf("ipv6.ICMPFilter.Set(%v, true) failed", typ)
			}
		}
	}
}

func TestSetICMPFilter(t *testing.T) {
	switch runtime.GOOS {
	case "fuchsia", "hurd", "js", "nacl", "plan9", "windows":
		t.Skipf("not supported on %s", runtime.GOOS)
	}
	if !nettest.SupportsIPv6() {
		t.Skip("ipv6 is not supported")
	}
	if !nettest.SupportsRawSocket() {
		t.Skipf("not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	c, err := net.ListenPacket("ip6:ipv6-icmp", "::1")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	p := ipv6.NewPacketConn(c)

	var f ipv6.ICMPFilter
	f.SetAll(true)
	f.Accept(ipv6.ICMPTypeEchoRequest)
	f.Accept(ipv6.ICMPTypeEchoReply)
	if err := p.SetICMPFilter(&f); err != nil {
		t.Fatal(err)
	}
	kf, err := p.ICMPFilter()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(kf, &f) {
		t.Fatalf("got %#v; want %#v", kf, f)
	}
}
