package goattypes

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BridgeRequests struct {
	Withdraws     []*WithdrawalRequest
	ReplaceByFees []*ReplaceByFeeRequest
	Cancel1s      []*Cancel1Request
}

type WithdrawalRequest struct {
	Id      uint64
	Amount  uint64
	TxPrice uint64
	Address string
}

var (
	withdrawalReqAddrLoc = big.NewInt(128)
	satoshi              = big.NewInt(1e10)
)

func UnpackIntoWithdrawRequest(topics []common.Hash, data []byte) (*WithdrawalRequest, error) {
	if len(topics) != 3 {
		return nil, fmt.Errorf("invalid Withdraw event topics length: expect 3 got %d", len(topics))
	}

	if dl := len(data); dl < 192 || dl%32 != 0 {
		return nil, fmt.Errorf("invalid Withdraw event data length: %d", len(data))
	}

	id := new(big.Int).SetBytes(topics[1][:])
	if !id.IsUint64() {
		return nil, fmt.Errorf("withdrawal id is too large")
	}

	amount := new(big.Int).SetBytes(data[:32]) // amount
	_, dust := amount.DivMod(amount, satoshi, new(big.Int))
	if !amount.IsUint64() {
		return nil, fmt.Errorf("withdrawal amount is too large: %d", amount)
	}

	if dust.BitLen() != 0 {
		return nil, fmt.Errorf("withdrawal amount has dust: %d", dust)
	}

	maxTxPrice := new(big.Int).SetBytes(data[64:96])
	if !maxTxPrice.IsUint64() {
		return nil, fmt.Errorf("max tx price is too large: %d", maxTxPrice)
	}

	// receiver
	if addrLoc := new(big.Int).SetBytes(data[96:128]); addrLoc.Cmp(withdrawalReqAddrLoc) != 0 {
		return nil, fmt.Errorf("address location in the withdraw event should be 128 but goat %d", addrLoc)
	}

	addrLen := new(big.Int).SetBytes(data[128:160]) // length
	addrLenInt64 := addrLen.Int64()
	if addrLenInt64 > 90 {
		return nil, errors.New("address length too large")
	}
	if int64(len(data[160:])) < addrLenInt64 {
		return nil, errors.New("address slice is out of range")
	}

	return &WithdrawalRequest{
		Id:      id.Uint64(),
		Amount:  amount.Uint64(),
		TxPrice: maxTxPrice.Uint64(),
		Address: string(data[160 : 160+addrLenInt64]),
	}, nil
}

func (req *WithdrawalRequest) RequestType() byte {
	return WithdrawalRequestType
}

func (req *WithdrawalRequest) Encode() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(req.RequestType())
	buf.Write(EncodeUint64(req.Id, req.Amount, req.TxPrice))
	buf.WriteByte(byte(len(req.Address)))
	buf.WriteString(req.Address)
	return buf.Bytes()
}

func (req *WithdrawalRequest) Decode(input []byte) error {
	if len(input) < 26 {
		return errors.New("WithdrawalRequest bytes length too short")
	}

	if input[0] != req.RequestType() {
		return errors.New("not WithdrawalRequest bytes")
	}

	res, err := DecodeUint64(input[1:25], 3)
	if err != nil {
		return err
	}
	req.Id, req.Amount, req.TxPrice = res[0], res[1], res[2]

	if addrLength := int(input[25]); len(input[26:]) != addrLength {
		return errors.New("invalid WithdrawalRequest length")
	}
	req.Address = string(input[26:])
	return nil
}

func (req *WithdrawalRequest) Copy() Request {
	return &WithdrawalRequest{
		Id:      req.Id,
		Amount:  req.Amount,
		TxPrice: req.TxPrice,
		Address: req.Address,
	}
}

type ReplaceByFeeRequest struct {
	Id      uint64
	TxPrice uint64
}

func UnpackIntoReplaceByFeeRequest(topics []common.Hash, data []byte) (*ReplaceByFeeRequest, error) {
	if len(topics) != 2 {
		return nil, fmt.Errorf("invalid ReplaceByFee event topics length: expect 3 got %d", len(topics))
	}

	if len(data) != 32 {
		return nil, fmt.Errorf("invalid ReplaceByFee event data length: %d", len(data))
	}

	id := new(big.Int).SetBytes(topics[1][:])
	if !id.IsUint64() {
		return nil, fmt.Errorf("withdrawal id is too large")
	}

	txPrice := new(big.Int).SetBytes(data) // maxTxPrice
	if !txPrice.IsUint64() {
		return nil, fmt.Errorf("max tx price is too large")
	}
	return &ReplaceByFeeRequest{Id: id.Uint64(), TxPrice: txPrice.Uint64()}, nil
}

func (req *ReplaceByFeeRequest) RequestType() byte { return ReplaceByFeeRequestType }
func (req *ReplaceByFeeRequest) Encode() []byte {
	return append([]byte{ReplaceByFeeRequestType}, EncodeUint64(req.Id, req.TxPrice)...)
}

func (req *ReplaceByFeeRequest) Decode(input []byte) error {
	if len(input) != 17 {
		return errors.New("invalid ReplaceByFeeRequest bytes length")
	}

	if input[0] != req.RequestType() {
		return errors.New("not ReplaceByFeeRequest request")
	}

	res, err := DecodeUint64(input[1:], 2)
	if err != nil {
		return err
	}
	req.Id, req.TxPrice = res[0], res[1]
	return nil
}

func (req *ReplaceByFeeRequest) Copy() Request {
	return &ReplaceByFeeRequest{
		Id:      req.Id,
		TxPrice: req.TxPrice,
	}
}

type Cancel1Request struct {
	Id uint64
}

func UnpackIntoCancel1Request(topics []common.Hash, data []byte) (*Cancel1Request, error) {
	if len(topics) != 2 {
		return nil, fmt.Errorf("invalid Cancel1 event topics length: expect 2 got %d", len(topics))
	}

	if len(data) != 0 {
		return nil, fmt.Errorf("invalid Cancel1 event data length, expect 0 got %d", len(data))
	}

	id := new(big.Int).SetBytes(topics[1][:])
	if !id.IsUint64() {
		return nil, fmt.Errorf("withdrawal id is too large")
	}
	return &Cancel1Request{Id: id.Uint64()}, nil
}

func (req *Cancel1Request) RequestType() byte { return Cancel1RequestType }
func (req *Cancel1Request) Encode() []byte {
	return append([]byte{req.RequestType()}, EncodeUint64(req.Id)...)
}
func (req *Cancel1Request) Decode(input []byte) error {
	if len(input) != 9 {
		return errors.New("invalid Cancel1 bytes length")
	}

	if input[0] != req.RequestType() {
		return errors.New("not Cancel1 request")
	}

	res, err := DecodeUint64(input[1:], 1)
	if err != nil {
		return err
	}
	req.Id = res[0]
	return nil
}

func (req *Cancel1Request) Copy() Request {
	return &Cancel1Request{
		Id: req.Id,
	}
}
