package airscan_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/brutella/dnssd"
	"github.com/stapelberg/airscan"
)

func TestDialer(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go http.Serve(ln, mockScanner())

	addr := ln.Addr().(*net.TCPAddr)

	// use a dnssd service struct
	svc := dnssd.Service{
		// Likely unreachable:
		Host:   "unreachable.invalid",
		Domain: "local",
		Port:   addr.Port,
		IPs: []net.IP{
			// Likely unreachable:
			net.ParseIP("255.255.255.255"),
			// Actually reachable:
			addr.IP,
		},
	}
	cl := airscan.NewClientForService(&svc)
	if _, err := cl.ScannerStatus(); err != nil {
		t.Fatal(err)
	}
}
