package goattypes

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	RelayerExecutor = common.HexToAddress("0xBc10000000000000000000000000000000001000")
	LockingExecutor = common.HexToAddress("0xBC10000000000000000000000000000000001001")
)

var (
	GoatTokenContract      = common.HexToAddress("0xbC10000000000000000000000000000000000001")
	GoatFoundationContract = common.HexToAddress("0xBc10000000000000000000000000000000000002")
	BridgeContract         = common.HexToAddress("0xBC10000000000000000000000000000000000003")
	LockingContract        = common.HexToAddress("0xbC10000000000000000000000000000000000004")
	BitcoinContract        = common.HexToAddress("0xbc10000000000000000000000000000000000005")
	RelayerContract        = common.HexToAddress("0xBC10000000000000000000000000000000000006")
)

var (
	// EmptyRequestsHash is the known hash of the empty requests set.
	EmptyRequestsHash = common.HexToHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	// EmptyGoatTxRoot is the known hash root of the empty goat txes list(with count 0 prefix)
	EmptyGoatTxRoot = common.Hex2Bytes("0056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)
