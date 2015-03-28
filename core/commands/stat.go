package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	humanize "github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/dustin/go-humanize"

	cmds "github.com/jbenet/go-ipfs/commands"
	metrics "github.com/jbenet/go-ipfs/metrics"
	peer "github.com/jbenet/go-ipfs/p2p/peer"
	protocol "github.com/jbenet/go-ipfs/p2p/protocol"
)

var StatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Query IPFS daemon statistics",
		ShortDescription: ``,
	},

	Subcommands: map[string]*cmds.Command{
		"bw": statBwCmd,
	},
}

var statBwCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Print ipfs bandwidth information",
		ShortDescription: ``,
	},
	Options: []cmds.Option{
		cmds.StringOption("peer", "p", "specify a peer to print bandwidth for"),
		cmds.StringOption("proto", "t", "specify a protocol to print bandwidth for"),
	},

	Run: func(req cmds.Request, res cmds.Response) {
		nd, err := req.Context().GetNode()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		// Must be online!
		if !nd.OnlineMode() {
			res.SetError(errNotOnline, cmds.ErrClient)
			return
		}

		pstr, pfound, err := req.Option("peer").String()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		tstr, tfound, err := req.Option("proto").String()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}
		if pfound && tfound {
			res.SetError(errors.New("please only specify peer OR protocol"), cmds.ErrClient)
			return
		}

		if pfound {
			pid, err := peer.IDB58Decode(pstr)
			if err != nil {
				res.SetError(err, cmds.ErrNormal)
				return
			}

			stats := nd.Reporter.GetBandwidthForPeer(pid)
			res.SetOutput(&stats)
		} else if tfound {
			pid := protocol.ID(tstr)
			stats := nd.Reporter.GetBandwidthForProtocol(pid)
			res.SetOutput(&stats)
		} else {
			totals := nd.Reporter.GetBandwidthTotals()
			res.SetOutput(&totals)
		}
	},
	Type: metrics.Stats{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			bs := res.Output().(*metrics.Stats)
			out := new(bytes.Buffer)
			fmt.Fprintln(out, "Bandwidth")
			fmt.Fprintf(out, "TotalIn: %s\n", humanize.Bytes(uint64(bs.TotalIn)))
			fmt.Fprintf(out, "TotalOut: %s\n", humanize.Bytes(uint64(bs.TotalOut)))
			fmt.Fprintf(out, "RateIn: %s/s\n", humanize.Bytes(uint64(bs.RateIn)))
			fmt.Fprintf(out, "RateOut: %s/s\n", humanize.Bytes(uint64(bs.RateOut)))
			return out, nil
		},
	},
}
