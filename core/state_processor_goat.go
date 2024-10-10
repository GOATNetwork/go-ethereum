package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/holiman/uint256"
)

var (
	gfBasePoint    = big.NewInt(200)
	gfMaxBasePoint = big.NewInt(1e4)
)

func ProcessGoatGasFee(statedb *state.StateDB, gasFees *big.Int) *big.Int {
	if gasFees.BitLen() == 0 {
		return new(big.Int)
	}

	// foundation tax 2%
	tax := new(big.Int).Mul(gasFees, gfBasePoint)
	tax.Div(tax, gfMaxBasePoint)

	if tax.BitLen() != 0 {
		f, _ := uint256.FromBig(tax)
		statedb.AddBalance(goattypes.GoatFoundationContract, f, tracing.BalanceIncreaseRewardTransactionFee)
	}

	// add gas revenue to locking contract
	// if the validator withdraws the gas reward, we will subtract it from locking contract then
	gas := new(big.Int).Sub(gasFees, tax)
	if gas.BitLen() != 0 {
		f, _ := uint256.FromBig(gas)
		statedb.AddBalance(goattypes.LockingContract, f, tracing.BalanceIncreaseRewardTransactionFee)
	}
	return gas
}

// ProcessGoatRequests processes goat requests
// It's not same with the eip-7685, the order is by it's emitted in the block
// and every request has its type prefix
func ProcessGoatRequests(height uint64, reward *big.Int, allLogs []*types.Log) ([][]byte, error) {
	requests := [][]byte{goattypes.NewGasRequest(height, reward).Encode()}

	for _, log := range allLogs {
		switch log.Address {
		case goattypes.BridgeContract:
			if len(log.Topics) < 2 {
				continue
			}
			switch log.Topics[0] {
			case goattypes.WithdrawEventTopic:
				req, err := goattypes.UnpackIntoWithdrawRequest(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.ReplaceByFeeEventTopic:
				req, err := goattypes.UnpackIntoReplaceByFeeRequest(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.Cancel1EventTopic:
				req, err := goattypes.UnpackIntoCancel1Request(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			}
		case goattypes.LockingContract:
			if len(log.Topics) != 1 {
				continue
			}
			switch log.Topics[0] {
			case goattypes.CreateEventTopic:
				req, err := goattypes.UnpackIntoCreateRequest(log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.LockEventTopic:
				req, err := goattypes.UnpackIntoLockRequest(log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.UnlockEventTopic:
				req, err := goattypes.UnpackIntoUnlockRequest(log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.ClaimEventTopic:
				req, err := goattypes.UnpackIntoClaimRequest(log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.GrantEventTopic:
				req, err := goattypes.UnpackIntoGrantRequest(log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.UpdateTokenWeightEventTopic:
				req, err := goattypes.UnpackIntoUpdateTokenWeightRequest(log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.UpdateTokenThresholdEventTopic:
				req, err := goattypes.UnpackIntoUpdateTokenThresholdRequest(log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			}
		case goattypes.RelayerContract:
			if len(log.Topics) != 2 {
				continue
			}
			switch log.Topics[0] {
			case goattypes.AddVoterEventTopoic:
				req, err := goattypes.UnpackIntoAddVoterRequest(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			case goattypes.RemoveVoterEventTopic:
				req, err := goattypes.UnpackIntoRemoveVoterRequest(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				requests = append(requests, req.Encode())
			}
		}
	}

	return requests, nil
}
