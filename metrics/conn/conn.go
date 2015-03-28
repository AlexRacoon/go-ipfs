package meterconn

import (
	metrics "github.com/jbenet/go-ipfs/metrics"
	conn "github.com/jbenet/go-ipfs/p2p/net/conn"
)

type MeteredConn struct {
	mesRecv metrics.MeterPeerCallback
	mesSent metrics.MeterPeerCallback

	conn.Conn
}

func NewMeteredConn(base conn.Conn, rcb metrics.MeterPeerCallback, scb metrics.MeterPeerCallback) conn.Conn {
	return &MeteredConn{
		Conn:    base,
		mesRecv: rcb,
		mesSent: scb,
	}
}

func (mc *MeteredConn) Read(b []byte) (int, error) {
	n, err := mc.Conn.Read(b)

	mc.mesRecv(int64(n), mc.Conn.RemotePeer())
	return n, err
}

func (mc *MeteredConn) Write(b []byte) (int, error) {
	n, err := mc.Conn.Write(b)

	mc.mesSent(int64(n), mc.Conn.RemotePeer())
	return n, err
}
