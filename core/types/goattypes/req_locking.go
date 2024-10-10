package goattypes

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LockingRequests struct {
	Gas              []*GasRequest
	Creates          []*CreateRequest
	Locks            []*LockRequest
	Unlocks          []*UnlockRequest
	Claims           []*ClaimRequest
	Grants           []*GrantRequest
	UpdateWeights    []*UpdateTokenWeightRequest
	UpdateThresholds []*UpdateTokenThresholdRequest
}

type GasRequest struct {
	Height uint64
	Amount *big.Int
}

func NewGasRequest(height uint64, amount *big.Int) *GasRequest {
	return &GasRequest{Height: height, Amount: new(big.Int).Set(amount)}
}

func (req *GasRequest) RequestType() byte { return GasRequestType }

func (req *GasRequest) Encode() []byte {
	res := make([]byte, 0, 41)
	res = append(res, req.RequestType())
	res = append(res, EncodeUint64(req.Height)...)
	res = append(res, req.Amount.FillBytes(make([]byte, 32))...)
	return res
}

func (req *GasRequest) Decode(input []byte) error {
	if len(input) != 41 {
		return errors.New("invalid GasRequest bytes length")
	}

	if input[0] != req.RequestType() {
		return errors.New("not GasRequest")
	}

	res, err := DecodeUint64(input[1:9], 1)
	if err != nil {
		return err
	}
	req.Height = res[0]

	req.Amount = new(big.Int).SetBytes(input[9:])
	return nil
}

func (req *GasRequest) Copy() Request {
	return &GasRequest{
		Height: req.Height,
		Amount: new(big.Int).Set(req.Amount),
	}
}

type CreateRequest struct {
	Validator common.Address
	Pubkey    [64]byte
}

func UnpackIntoCreateRequest(data []byte) (*CreateRequest, error) {
	if len(data) != 128 {
		return nil, fmt.Errorf("invalid CreateValidator event data length: want 128, have %d", len(data))
	}
	return &CreateRequest{Validator: common.BytesToAddress(data[:32]), Pubkey: [64]byte(data[64:])}, nil
}

func (req *CreateRequest) RequestType() byte { return CreateRequestType }
func (req *CreateRequest) Encode() []byte {
	res := make([]byte, 0, 85)
	res = append(res, req.RequestType())
	res = append(res, req.Validator.Bytes()...)
	res = append(res, req.Pubkey[:]...)
	return res
}

func (req *CreateRequest) Decode(input []byte) error {
	if len(input) != 85 {
		return errors.New("invalid CreateRequest bytes length")
	}

	if input[0] != req.RequestType() {
		return errors.New("not CreateRequest")
	}
	req.Validator = common.BytesToAddress(input[1:21])
	req.Pubkey = [64]byte(input[21:])
	return nil
}

func (req *CreateRequest) Copy() Request {
	return &CreateRequest{
		Validator: req.Validator,
		Pubkey:    req.Pubkey,
	}
}

type LockRequest struct {
	Validator common.Address
	Token     common.Address
	Amount    *big.Int
}

func (req *LockRequest) RequestType() byte { return LockRequestType }
func (req *LockRequest) Encode() []byte {
	res := make([]byte, 0, 73)
	res = append(res, req.RequestType())
	res = append(res, req.Validator.Bytes()...)
	res = append(res, req.Token.Bytes()...)
	res = append(res, req.Amount.FillBytes(make([]byte, 32))...)
	return res
}

func (req *LockRequest) Decode(input []byte) error {
	if len(input) != 73 {
		return errors.New("invalid LockRequest bytes length")
	}

	if input[0] != req.RequestType() {
		return errors.New("not LockRequest")
	}

	req.Validator = common.BytesToAddress(input[1:21])
	input = input[21:]
	req.Token = common.BytesToAddress(input[:20])
	input = input[20:]
	req.Amount = new(big.Int).SetBytes(input)
	return nil
}

func (req *LockRequest) Copy() Request {
	return &LockRequest{
		Validator: req.Validator,
		Token:     req.Token,
		Amount:    new(big.Int).Set(req.Amount),
	}
}

func UnpackIntoLockRequest(data []byte) (*LockRequest, error) {
	if len(data) != 96 {
		return nil, fmt.Errorf("invalid Lock event data length: want 96, have %d", len(data))
	}
	return &LockRequest{
		Validator: common.BytesToAddress(data[:32]),
		Token:     common.BytesToAddress(data[32:64]),
		Amount:    new(big.Int).SetBytes(data[64:]),
	}, nil
}

type UnlockRequest struct {
	Id        uint64
	Validator common.Address
	Recipient common.Address
	Token     common.Address
	Amount    *big.Int
}

func UnpackIntoUnlockRequest(data []byte) (*UnlockRequest, error) {
	if len(data) != 160 {
		return nil, fmt.Errorf("invalid Unlock event data length: want 160, have %d", len(data))
	}
	return &UnlockRequest{
		Id:        new(big.Int).SetBytes(data[:32]).Uint64(),
		Validator: common.BytesToAddress(data[32:64]),
		Recipient: common.BytesToAddress(data[64:96]),
		Token:     common.BytesToAddress(data[96:128]),
		Amount:    new(big.Int).SetBytes(data[128:160]),
	}, nil
}

func (req *UnlockRequest) RequestType() byte { return UnlockRequestType }
func (req *UnlockRequest) Encode() []byte {
	res := make([]byte, 0, 101)
	res = append(res, req.RequestType())
	res = append(res, EncodeUint64(req.Id)...)
	res = append(res, req.Validator.Bytes()...)
	res = append(res, req.Recipient.Bytes()...)
	res = append(res, req.Token.Bytes()...)
	res = append(res, req.Amount.FillBytes(make([]byte, 32))...)
	return res
}

func (req *UnlockRequest) Decode(input []byte) error {
	if len(input) != 101 {
		return errors.New("invalid UnlockRequest bytes length")
	}

	if input[0] != req.RequestType() {
		return errors.New("not UnlockRequest")
	}

	res, err := DecodeUint64(input[1:9], 1)
	if err != nil {
		return err
	}
	req.Id = res[0]

	input = input[9:]
	req.Validator = common.BytesToAddress(input[:20])
	input = input[20:]
	req.Recipient = common.BytesToAddress(input[:20])
	input = input[20:]
	req.Token = common.BytesToAddress(input[:20])
	input = input[20:]
	req.Amount = new(big.Int).SetBytes(input)
	return nil
}

func (req *UnlockRequest) Copy() Request {
	return &UnlockRequest{
		Id:        req.Id,
		Validator: req.Validator,
		Token:     req.Token,
		Recipient: req.Recipient,
		Amount:    new(big.Int).Set(req.Amount),
	}
}

type ClaimRequest struct {
	Id        uint64
	Validator common.Address
	Recipient common.Address
}

func (req *ClaimRequest) RequestType() byte { return ClaimRequestType }
func (req *ClaimRequest) Encode() []byte {
	res := make([]byte, 0, 49)
	res = append(res, req.RequestType())
	res = append(res, EncodeUint64(req.Id)...)
	res = append(res, req.Validator.Bytes()...)
	res = append(res, req.Recipient.Bytes()...)
	return res
}

func (req *ClaimRequest) Decode(input []byte) error {
	if len(input) != 49 {
		return errors.New("invalid UnlockRequest bytes length")
	}
	if input[0] != req.RequestType() {
		return errors.New("not ClaimRequest")
	}
	res, err := DecodeUint64(input[1:9], 1)
	if err != nil {
		return err
	}
	req.Id = res[0]

	input = input[9:]
	req.Validator = common.BytesToAddress(input[:20])
	req.Recipient = common.BytesToAddress(input[20:])
	return nil
}

func (req *ClaimRequest) Copy() Request {
	return &ClaimRequest{
		Id:        req.Id,
		Validator: req.Validator,
		Recipient: req.Recipient,
	}
}

func UnpackIntoClaimRequest(data []byte) (*ClaimRequest, error) {
	if len(data) != 96 {
		return nil, fmt.Errorf("GoatRewardClaim wrong length: want 96, have %d", len(data))
	}
	return &ClaimRequest{
		Id:        new(big.Int).SetBytes(data[:32]).Uint64(),
		Validator: common.BytesToAddress(data[32:64]),
		Recipient: common.BytesToAddress(data[64:96]),
	}, nil
}

type UpdateTokenWeightRequest struct {
	Token  common.Address
	Weight uint64
}

func UnpackIntoUpdateTokenWeightRequest(data []byte) (*UpdateTokenWeightRequest, error) {
	if len(data) != 64 {
		return nil, fmt.Errorf("UpdateTokenWeight wrong length: want 64, have %d", len(data))
	}
	return &UpdateTokenWeightRequest{
		Token:  common.BytesToAddress(data[:32]),
		Weight: new(big.Int).SetBytes(data[32:64]).Uint64(),
	}, nil
}

func (req *UpdateTokenWeightRequest) RequestType() byte { return UpdateTokenWeightRequestType }
func (req *UpdateTokenWeightRequest) Encode() []byte {
	res := make([]byte, 0, 29)
	res = append(res, req.RequestType())
	res = append(res, req.Token.Bytes()...)
	res = append(res, EncodeUint64(req.Weight)...)
	return res
}

func (req *UpdateTokenWeightRequest) Decode(input []byte) error {
	if len(input) != 29 {
		return errors.New("invalid UpdateTokenWeightRequest bytes length")
	}

	if input[0] != req.RequestType() {
		return errors.New("not UpdateTokenWeightRequest")
	}
	req.Token = common.BytesToAddress(input[1:21])

	res, err := DecodeUint64(input[21:], 1)
	if err != nil {
		return err
	}
	req.Weight = res[0]
	return nil
}

func (req *UpdateTokenWeightRequest) Copy() Request {
	return &UpdateTokenWeightRequest{
		Token:  req.Token,
		Weight: req.Weight,
	}
}

type UpdateTokenThresholdRequest struct {
	Token     common.Address
	Threshold *big.Int
}

func UnpackIntoUpdateTokenThresholdRequest(data []byte) (*UpdateTokenThresholdRequest, error) {
	if len(data) != 64 {
		return nil, fmt.Errorf("invalid UpdateTokenThreshold event data length: want 64, have %d", len(data))
	}
	return &UpdateTokenThresholdRequest{
		Token:     common.BytesToAddress(data[:32]),
		Threshold: new(big.Int).SetBytes(data[32:64]),
	}, nil
}

func (req *UpdateTokenThresholdRequest) RequestType() byte { return UpdateTokenThresholdRequestType }
func (req *UpdateTokenThresholdRequest) Encode() []byte {
	res := make([]byte, 0, 53)
	res = append(res, req.RequestType())
	res = append(res, req.Token.Bytes()...)
	res = append(res, req.Threshold.FillBytes(make([]byte, 32))...)
	return res
}

func (req *UpdateTokenThresholdRequest) Decode(input []byte) error {
	if len(input) != 53 {
		return errors.New("invalid UpdateTokenThresholdRequest bytes length")
	}

	if input[0] != req.RequestType() {
		return errors.New("not UpdateTokenThresholdRequest")
	}
	req.Token = common.BytesToAddress(input[1:21])
	req.Threshold = new(big.Int).SetBytes(input[21:])
	return nil
}

func (req *UpdateTokenThresholdRequest) Copy() Request {
	return &UpdateTokenThresholdRequest{
		Token:     req.Token,
		Threshold: new(big.Int).Set(req.Threshold),
	}
}

type GrantRequest struct {
	Amount *big.Int
}

func UnpackIntoGrantRequest(data []byte) (*GrantRequest, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("invalid GoatGrant event data length: want 32, have %d", len(data))
	}
	return &GrantRequest{Amount: new(big.Int).SetBytes(data[:])}, nil
}

func (req *GrantRequest) RequestType() byte { return GrantRequestType }
func (req *GrantRequest) Encode() []byte {
	return append([]byte{req.RequestType()}, req.Amount.FillBytes(make([]byte, 32))...)
}
func (req *GrantRequest) Decode(input []byte) error {
	if len(input) != 33 {
		return errors.New("invalid UpdateTokenThresholdRequest bytes length")
	}
	if input[0] != req.RequestType() {
		return errors.New("not UpdateTokenThresholdRequest")
	}
	req.Amount = new(big.Int).SetBytes(input[1:])
	return nil
}

func (req *GrantRequest) Copy() Request {
	return &GrantRequest{
		Amount: new(big.Int).Set(req.Amount),
	}
}
