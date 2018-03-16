package tdns

import (
	"net"
	"testing"

	"github.com/miekg/dns"
)

func TestLookupIPOnConn(t *testing.T) {
	udpConn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		t.Error(err)
	}
	defer udpConn.Close()

	ips, err := LookupIPOnConn(udpConn, "www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	if len(ips) == 0 {
		t.Fail()
	}

	tcpConn, err := net.Dial("tcp", "8.8.8.8:53")
	if err != nil {
		t.Error(err)
	}
	defer tcpConn.Close()
	ips2, err := LookupIPOnConn(tcpConn, "www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	if len(ips2) == 0 {
		t.Fail()
	}
}

func TestPackDNSMsg(t *testing.T) {
	m := packDNSMsg("www.baidu.com")
	if m.Id == 0 {
		t.Fatal("pack dns msg, id is 0")
	}
	if m.Rcode != dns.RcodeSuccess {
		t.Fatalf("pack dns msg, Rcode != %d", dns.RcodeSuccess)
	}

	if len(m.Question) != 1 {
		t.Fatal("pack dns msg, Question length != 1")
	}

	if m.Question[0].Name != "www.baidu.com." {
		t.Fatal("pack dns msg, Question Name != www.baidu.com.")
	}

	if m.Question[0].Qtype != dns.TypeA {
		t.Fatalf("pack dns msg, Question Qtype != %d", dns.TypeA)
	}

	if m.Question[0].Qclass != uint16(dns.ClassINET) {
		t.Fatalf("pack dns msg, Question Qclass != %d", dns.ClassINET)
	}
}
