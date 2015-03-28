package meterstream

import (
	metrics "github.com/jbenet/go-ipfs/metrics"
	inet "github.com/jbenet/go-ipfs/p2p/net"
	peer "github.com/jbenet/go-ipfs/p2p/peer"
	protocol "github.com/jbenet/go-ipfs/p2p/protocol"
)

type meteredStream struct {
	// keys for accessing metrics data
	protoKey protocol.ID
	peerKey  peer.ID

	stream inet.Stream

	// callbacks for reporting bandwidth usage
	mesSent metrics.MeterProtoCallback
	mesRecv metrics.MeterProtoCallback
}

func NewMeteredStream(base inet.Stream, pid protocol.ID, sentCB, recvCB metrics.MeterProtoCallback) inet.Stream {
	return &meteredStream{
		stream:   base,
		mesSent:  sentCB,
		mesRecv:  recvCB,
		protoKey: pid,
	}
}

func (s *meteredStream) Read(b []byte) (int, error) {
	n, err := s.stream.Read(b)

	// Log bytes read
	s.mesRecv(int64(n), s.protoKey)

	return n, err
}

func (s *meteredStream) Write(b []byte) (int, error) {
	n, err := s.stream.Write(b)

	// Log bytes written
	s.mesSent(int64(n), s.protoKey)

	return n, err
}

func (s *meteredStream) Close() error {
	return s.stream.Close()
}

func (s *meteredStream) Conn() inet.Conn {
	return s.stream.Conn()
}
