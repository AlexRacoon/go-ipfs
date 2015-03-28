package metrics

import (
	gm "github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/whyrusleeping/go-metrics"
	"sync"

	peer "github.com/jbenet/go-ipfs/p2p/peer"
	protocol "github.com/jbenet/go-ipfs/p2p/protocol"
)

type Stats struct {
	TotalIn  int64
	TotalOut int64
	RateIn   float64
	RateOut  float64
}

type BandwidthCounter struct {
	lock     sync.Mutex
	totalIn  gm.Meter
	totalOut gm.Meter
	reg      gm.Registry
}

func NewBandwidthCounter() *BandwidthCounter {
	reg := gm.NewRegistry()
	return &BandwidthCounter{
		totalIn:  gm.GetOrRegisterMeter("totalIn", reg),
		totalOut: gm.GetOrRegisterMeter("totalOut", reg),
		reg:      reg,
	}
}

func (bwc *BandwidthCounter) LogSentMessagePeer(size int64, p peer.ID) {
	bwc.totalOut.Mark(size)

	meter := gm.GetOrRegisterMeter("/peer/out/"+string(p), bwc.reg)
	meter.Mark(size)
}

func (bwc *BandwidthCounter) LogRecvMessagePeer(size int64, p peer.ID) {
	bwc.totalIn.Mark(size)

	meter := gm.GetOrRegisterMeter("/peer/in/"+string(p), bwc.reg)
	meter.Mark(size)
}

func (bwc *BandwidthCounter) LogSentMessageProto(size int64, proto protocol.ID) {
	meter := gm.GetOrRegisterMeter("/proto/out/"+string(proto), bwc.reg)
	meter.Mark(size)
}

func (bwc *BandwidthCounter) LogRecvMessageProto(size int64, proto protocol.ID) {
	meter := gm.GetOrRegisterMeter("/proto/in/"+string(proto), bwc.reg)
	meter.Mark(size)
}

func (bwc *BandwidthCounter) GetBandwidthForPeer(p peer.ID) (out Stats) {
	inMeter := gm.GetOrRegisterMeter("/peer/in/"+string(p), bwc.reg).Snapshot()
	outMeter := gm.GetOrRegisterMeter("/peer/out/"+string(p), bwc.reg).Snapshot()

	return Stats{
		TotalIn:  inMeter.Count(),
		TotalOut: outMeter.Count(),
		RateIn:   inMeter.RateFine(),
		RateOut:  outMeter.RateFine(),
	}
}

func (bwc *BandwidthCounter) GetBandwidthForProtocol(proto protocol.ID) (out Stats) {
	inMeter := gm.GetOrRegisterMeter(string("/proto/in/"+proto), bwc.reg).Snapshot()
	outMeter := gm.GetOrRegisterMeter(string("/proto/out/"+proto), bwc.reg).Snapshot()

	return Stats{
		TotalIn:  inMeter.Count(),
		TotalOut: outMeter.Count(),
		RateIn:   inMeter.RateFine(),
		RateOut:  outMeter.RateFine(),
	}
}

func (bwc *BandwidthCounter) GetBandwidthTotals() (out Stats) {
	return Stats{
		TotalIn:  bwc.totalIn.Count(),
		TotalOut: bwc.totalOut.Count(),
		RateIn:   bwc.totalIn.RateFine(),
		RateOut:  bwc.totalOut.RateFine(),
	}
}
