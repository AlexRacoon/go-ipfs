// +build !nofuse

package incfusever

import (
	fuseversion "github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-fuse-version"
)

var _ = fuseversion.LocalFuseSystems
