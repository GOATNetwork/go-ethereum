package goattypes

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type RelayerRequests struct {
	Adds    []*AddVoterRequest
	Removes []*RemoveVoterRequest
}

type AddVoterRequest struct {
	Voter  common.Address
	Pubkey common.Hash
}

func UnpackIntoAddVoterRequest(topics []common.Hash, data []byte) (*AddVoterRequest, error) {
	if len(topics) != 2 {
		return nil, fmt.Errorf("invalid AddVoter event topics length: expect 2 got %d", len(topics))
	}

	if len(data) != 32 {
		return nil, fmt.Errorf("invalid AddVoter event data length: want 32 have %d", len(data))
	}

	return &AddVoterRequest{
		Voter:  common.BytesToAddress(topics[1][:]),
		Pubkey: common.BytesToHash(data[:]),
	}, nil
}

func (req *AddVoterRequest) RequestType() byte { return AddVoterRequestType }
func (req *AddVoterRequest) Encode() []byte {
	res := make([]byte, 0, 53)
	res = append(res, req.RequestType())
	res = append(res, req.Voter.Bytes()...)
	res = append(res, req.Pubkey.Bytes()...)
	return res
}

func (req *AddVoterRequest) Decode(input []byte) error {
	if len(input) != 53 {
		return errors.New("invalid AddVoterRequest length")
	}
	if input[0] != req.RequestType() {
		return errors.New("not AddVoterRequest bytes")
	}
	req.Voter = common.BytesToAddress(input[1:21])
	req.Pubkey = common.BytesToHash(input[21:])
	return nil
}

func (req *AddVoterRequest) Copy() Request {
	return &AddVoterRequest{
		Voter:  req.Voter,
		Pubkey: req.Pubkey,
	}
}

type RemoveVoterRequest struct {
	Voter common.Address
}

func UnpackIntoRemoveVoterRequest(topics []common.Hash, data []byte) (*RemoveVoterRequest, error) {
	if len(topics) != 2 {
		return nil, fmt.Errorf("invalid RemoveVoter event topics length: expect 2 got %d", len(topics))
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("invalid RemoveVoter event data length: want 0, have %d", len(data))
	}
	return &RemoveVoterRequest{Voter: common.BytesToAddress(topics[1][:])}, nil
}

func (req *RemoveVoterRequest) RequestType() byte { return RemoveVoterRequestType }

func (req *RemoveVoterRequest) Encode() []byte {
	res := make([]byte, 0, 21)
	res = append(res, req.RequestType())
	res = append(res, req.Voter.Bytes()...)
	return res
}

func (req *RemoveVoterRequest) Decode(input []byte) error {
	if len(input) != 21 {
		return errors.New("invalid RemoveVoterRequest length")
	}
	if input[0] != req.RequestType() {
		return errors.New("not RemoveVoterRequest bytes")
	}
	req.Voter = common.BytesToAddress(input[1:])
	return nil
}

func (req *RemoveVoterRequest) Copy() Request {
	return &RemoveVoterRequest{
		Voter: req.Voter,
	}
}
