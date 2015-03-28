package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	humanize "github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/dustin/go-humanize"

	cmds "github.com/jbenet/go-ipfs/commands"
	metrics "github.com/jbenet/go-ipfs/metrics"
	peer "github.com/jbenet/go-ipfs/p2p/peer"
	protocol "github.com/jbenet/go-ipfs/p2p/protocol"
	u "github.com/jbenet/go-ipfs/util"
)

var StatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Query IPFS daemon statistics",
		ShortDescription: ``,
	},

	Subcommands: map[string]*cmds.Command{
		"bw":      statBwCmd,
		"bw-poll": statBwPollCmd,
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

var statBwPollCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Print ipfs bandwidth information continuously",
		ShortDescription: ``,
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

		out := make(chan interface{})
		go func() {
			defer close(out)
			tick := time.NewTicker(time.Second)
			for {
				select {
				case <-req.Context().Context.Done():
					return
				case <-tick.C:
					totals := nd.Reporter.GetBandwidthTotals()
					out <- &totals
				}
			}
		}()

		res.SetOutput((<-chan interface{})(out))
	},
	Type: metrics.Stats{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			outCh, ok := res.Output().(<-chan interface{})
			if !ok {
				return nil, u.ErrCast()
			}

			first := true
			marshal := func(v interface{}) (io.Reader, error) {
				bs, ok := v.(*metrics.Stats)
				if !ok {
					return nil, u.ErrCast()
				}
				out := new(bytes.Buffer)
				if first {
					fmt.Fprintln(out, "Total Up\t Total Down\t Rate Up\t Rate Down")
					first = false
				}
				fmt.Fprint(out, "\r")
				fmt.Fprintf(out, "%s \t\t", humanize.Bytes(uint64(bs.TotalOut)))
				fmt.Fprintf(out, " %s \t\t", humanize.Bytes(uint64(bs.TotalIn)))
				fmt.Fprintf(out, " %s/s   \t", humanize.Bytes(uint64(bs.RateOut)))
				fmt.Fprintf(out, " %s/s     ", humanize.Bytes(uint64(bs.RateIn)))
				return out, nil

			}

			return &cmds.ChannelMarshaler{
				Channel:   outCh,
				Marshaler: marshal,
			}, nil
		},
	},
}
