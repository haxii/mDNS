package dns

import (
	"errors"
	"net"

	"github.com/miekg/dns"
)

func LookupIPOnConn(conn net.Conn, host string) ([]net.IPAddr, error) {
	msg := packDnsMsg(host)
	co := &dns.Conn{Conn: conn}
	err := co.WriteMsg(msg)
	if err != nil {
		return nil, err
	}
	r, err := co.ReadMsg()
	if err != nil {
		return nil, err
	}
	if r.Id != msg.Id {
		return nil, errors.New("Id mismatch")
	}

	ips := []net.IPAddr{}
	for _, answer := range r.Answer {
		header := answer.Header()
		if header.Rrtype != dns.TypeA {
			continue
		}
		a := answer.(*dns.A)
		ips = append(ips, net.IPAddr{IP: a.A})
	}

	return ips, err
}

//pack dns msg object
func packDnsMsg(host string) *dns.Msg {
	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative:     false,
			AuthenticatedData: false,
			CheckingDisabled:  false,
			RecursionDesired:  true,
			Opcode:            dns.OpcodeQuery,
		},
		Question: make([]dns.Question, 1),
	}
	m.Rcode = dns.RcodeSuccess
	qt := dns.TypeA
	qc := uint16(dns.ClassINET)
	m.Question[0] = dns.Question{Name: dns.Fqdn(host), Qtype: qt, Qclass: qc}
	m.Id = dns.Id()

	return m
}
