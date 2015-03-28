package metrics

import (
	peer "github.com/jbenet/go-ipfs/p2p/peer"
	protocol "github.com/jbenet/go-ipfs/p2p/protocol"
)

type MeterPeerCallback func(int64, peer.ID)
type MeterProtoCallback func(int64, protocol.ID)
type MeterCallback func(int64)

type Reporter interface {
	LogSentMessage(int64)
	LogRecvMessage(int64)
	LogSentMessagePeer(int64, peer.ID)
	LogRecvMessagePeer(int64, peer.ID)
	LogSentMessageProto(int64, protocol.ID)
	LogRecvMessageProto(int64, protocol.ID)
	GetBandwidthForPeer(peer.ID) Stats
	GetBandwidthForProtocol(protocol.ID) Stats
	GetBandwidthTotals() Stats
}
