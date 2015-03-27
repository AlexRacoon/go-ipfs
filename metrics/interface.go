package metrics

import (
	peer "github.com/jbenet/go-ipfs/p2p/peer"
	protocol "github.com/jbenet/go-ipfs/p2p/protocol"
)

type Reporter interface {
	LogSentMessage(size int64, proto protocol.ID, peer peer.ID)
	LogRecvMessage(size int64, proto protocol.ID, peer peer.ID)
	GetBandwidthForPeer(peer.ID) Stats
	GetBandwidthForProtocol(protocol.ID) Stats
	GetBandwidthTotals() Stats
}
