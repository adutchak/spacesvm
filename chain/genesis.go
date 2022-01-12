// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package chain

import (
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/units"
)

type Genesis struct {
	// Tx params
	BaseTxUnits uint64 `serialize:"true" json:"baseTxUnits"`

	// SetTx params
	ValueUnitSize int `serialize:"true" json:"valueUnitSize"`
	MaxValueSize  int `serialize:"true" json:"maxValueSize"`

	// Claim Params
	ClaimFeeMultiplier   uint64 `serialize:"true" json:"claimFeeMultiplier"`
	ExpiryTime           uint64 `serialize:"true" json:"expiryTime"`
	ClaimTier3Multiplier uint64 `serialize:"true" json:"claimTier3Multiplier"`
	ClaimTier2Size       int    `serialize:"true" json:"claimTier2Size"`
	ClaimTier2Multiplier uint64 `serialize:"true" json:"claimTier2Multiplier"`
	ClaimTier1Size       int    `serialize:"true" json:"claimTier1Size"`
	ClaimTier1Multiplier uint64 `serialize:"true" json:"claimTier1Multiplier"`

	// Lifeline Params
	PrefixRenewalDiscount uint64 `serialize:"true" json:"prefixRenewalDiscount"`

	// Fee Mechanism Params
	LookbackWindow int64  `serialize:"true" json:"lookbackWindow"`
	BlockTarget    int64  `serialize:"true" json:"blockTarget"`
	TargetUnits    uint64 `serialize:"true" json:"targetUnits"`
	MinDifficulty  uint64 `serialize:"true" json:"minDifficulty"`
	MinBlockCost   uint64 `serialize:"true" json:"minBlockCost"`
}

func DefaultGenesis() *Genesis {
	return &Genesis{
		// Tx params
		BaseTxUnits: 10,

		// SetTx params
		ValueUnitSize: 256,             // 256B
		MaxValueSize:  128 * units.KiB, // (500 Units)

		// Claim Params
		ClaimFeeMultiplier:   5,
		ExpiryTime:           60 * 60 * 24 * 30, // 30 Days
		ClaimTier3Multiplier: 1,
		ClaimTier2Size:       36,
		ClaimTier2Multiplier: 5,
		ClaimTier1Size:       12,
		ClaimTier1Multiplier: 25,

		// Lifeline Params
		PrefixRenewalDiscount: 5,

		// Fee Mechanism Params
		LookbackWindow: 60,            // 60 Seconds
		BlockTarget:    1,             // 1 Block per Second
		TargetUnits:    10 * 512 * 60, // 5012 Units Per Block (~1.2MB of SetTx)
		MinDifficulty:  100,           // ~100ms per unit (~5s for easiest claim)
		MinBlockCost:   1,             // Minimum Unit Overhead
	}
}

func VerifyGenesis(b *StatelessBlock) error {
	if b.Prnt != ids.Empty {
		return ErrInvalidGenesisParent
	}
	if b.Hght != 0 {
		return ErrInvalidGenesisHeight
	}
	if b.Tmstmp == 0 || time.Now().Unix()-b.Tmstmp < 0 {
		return ErrInvalidGenesisTimestamp
	}
	if b.Genesis == nil {
		return ErrMissingGenesis
	}
	if b.Difficulty != b.Genesis.MinDifficulty {
		return ErrInvalidGenesisDifficulty
	}
	if b.Cost != b.Genesis.MinBlockCost {
		return ErrInvalidGenesisCost
	}
	if len(b.Txs) > 0 {
		return ErrInvalidGenesisTxs
	}
	if b.Beneficiary != nil {
		return ErrInvalidGenesisBeneficiary
	}
	return nil
}
