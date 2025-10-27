package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Option struct {
	Address    string
	TrialCount int
	Identifier int
	Timeout    time.Duration
}

type Pinger struct {
	Conn     *icmp.PacketConn
	TargetIP net.IP
	ID       int
	Seq      int
	Timeout  time.Duration
}

func main() {
	opt := parseOptions()
	pinger, err := NewPinger(opt)
	if err != nil {
		panic(err)
	}
	defer pinger.Close()

	fmt.Printf("PING %s (id=%d):\n", opt.Address, opt.Identifier)

	for i := 0; i < opt.TrialCount; i++ {
		rtt, err := pinger.PingOnce()
		if err != nil {
			fmt.Printf("Request timeout (seq=%d): %v\n", pinger.Seq, err)
		} else {
			fmt.Printf("%s: icmp_seq=%d id=%d time=%v\n",
				pinger.TargetIP, pinger.Seq, pinger.ID, rtt)
		}
		time.Sleep(time.Second)
	}
}

func parseOptions() Option {
	var opt Option
	flag.StringVar(&opt.Address, "a", "", "Destination address (required)")
	flag.IntVar(&opt.TrialCount, "n", 3, "Number of echo requests (0 for infinite)")
	flag.IntVar(&opt.Identifier, "id", -1, "ICMP identifier (0-65535). Default: PID & 0xffff")
	flag.DurationVar(&opt.Timeout, "t", 3*time.Second, "Timeout for each ping")
	flag.Parse()

	if opt.Address == "" {
		fmt.Println("Required address (-a).")
		os.Exit(1)
	}
	if opt.Identifier == -1 {
		opt.Identifier = os.Getpid() & 0xffff
	}
	if opt.Identifier < 0 || opt.Identifier > 0xffff {
		fmt.Println("Identifier must be 0-65535")
		os.Exit(1)
	}
	if opt.TrialCount <= 0 {
		fmt.Println("Trial count must be greater than 0")
		os.Exit(1)
	}
	return opt
}

func NewPinger(opt Option) (*Pinger, error) {
	ip, err := net.ResolveIPAddr("ip4", opt.Address)
	if err != nil {
		return nil, fmt.Errorf("resolve error: %w", err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("listen error: %w", err)
	}

	return &Pinger{
		Conn:     conn,
		TargetIP: ip.IP,
		ID:       opt.Identifier,
		Seq:      -1,
		Timeout:  opt.Timeout,
	}, nil
}

func (p *Pinger) PingOnce() (time.Duration, error) {
	p.Seq++

	now := time.Now().UnixMilli()
	data := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(data, now)

	req := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   p.ID,
			Seq:  p.Seq,
			Data: data,
		},
	}

	reqBytes, err := req.Marshal(nil)
	if err != nil {
		return 0, fmt.Errorf("marshal error: %w", err)
	}

	if _, err := p.Conn.WriteTo(reqBytes, &net.IPAddr{IP: p.TargetIP}); err != nil {
		return 0, fmt.Errorf("send error: %w", err)
	}

	p.Conn.SetReadDeadline(time.Now().Add(p.Timeout))
	reply := make([]byte, 1500)

	n, _, err := p.Conn.ReadFrom(reply)
	if err != nil {
		return 0, fmt.Errorf("recv timeout: %w", err)
	}

	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), reply[:n])
	if err != nil {
		return 0, fmt.Errorf("parse error: %w", err)
	}

	if rm.Type != ipv4.ICMPTypeEchoReply {
		return 0, fmt.Errorf("unexpected ICMP type: %v", rm.Type)
	}

	body, ok := rm.Body.(*icmp.Echo)
	if !ok {
		return 0, fmt.Errorf("invalid echo body")
	}

	if body.ID != p.ID {
		return 0, fmt.Errorf("mismatched ID (got %d)", body.ID)
	}

	t, _ := binary.Varint(body.Data)
	rtt := time.Duration(time.Now().UnixMilli()-t) * time.Millisecond
	return rtt, nil
}

func (p *Pinger) Close() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}
