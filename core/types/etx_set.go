package types

import (
	"github.com/dominant-strategies/go-quai/common"
)

const (
	EtxExpirationAge = 100 // With 10s blocks, ETX expire after ~24hrs
)

// The EtxSet maps an ETX hash to the ETX and block number in which it became available.
// If no entry exists for a given ETX hash, then that ETX is not available.
type EtxSet map[common.Hash]EtxSetEntry

type EtxSetEntry struct {
	Height uint64
	ETX    Transaction
}

func NewEtxSet() EtxSet {
	return make(EtxSet)
}

// updateInboundEtxs updates the set of inbound ETXs available to be mined into
// a block in this location. This method adds any new ETXs to the set and
// removes expired ETXs.
func (set *EtxSet) Update(newInboundEtxs Transactions, currentHeight uint64) {
	// Add new ETX entries to the inbound set
	for _, etx := range newInboundEtxs {
		if etx.To().Location().Equal(common.NodeLocation) {
			(*set)[etx.Hash()] = EtxSetEntry{currentHeight, *etx}
		} else {
			panic("cannot add ETX destined to other chain to our ETX set")
		}
	}

	// Remove expired ETXs
	for txHash, entry := range *set {
		availableAtBlock := entry.Height
		etxExpirationHeight := availableAtBlock + EtxExpirationAge
		if currentHeight > etxExpirationHeight {
			delete(*set, txHash)
		}
	}
}
