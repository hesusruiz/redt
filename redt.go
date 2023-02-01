package redt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	ethertypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rpc"

	ibftengine "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/engine"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	lrucache "github.com/hashicorp/golang-lru"

	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

var logger *zap.SugaredLogger

func init() {
	// Set the common logger
	logger = zap.Must(zap.NewProduction()).Sugar()
}

type NodeInfo struct {
	Operator string `json:"operator"`
	Enode    string `json:"enode"`
	Address  common.Address
}

// sigHash returns the hash which is used as input for the Istanbul
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
func sigHash(header *ethertypes.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	rlp.Encode(hasher, ethertypes.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

func toBlockNumArg(number int64) string {
	if number == -1 {
		return "latest"
	}
	return fmt.Sprintf("0x%X", number)
}

type RedTNode struct {
	cli           *ethclient.Client
	headerCache   *lrucache.Cache
	rpccli        *rpc.Client
	valSet        []common.Address
	allValidators map[common.Address]*NodeInfo
	timeout       time.Duration
}

// NewRedTNode creates a connection to a blockchain node in the RedT network
// url can be an HTTP or WebSocket schema
func NewRedTNode(url string) (rt *RedTNode, err error) {

	// Allocate the structure
	rt = &RedTNode{}

	// Set default timeout
	rt.timeout = 30 * time.Second

	// Set the timeout for the next call
	ctx, cancel := context.WithTimeout(context.Background(), rt.timeout)
	defer cancel()

	// Connect to Client using an RPC client and then wrap it with a standard ethclient
	// We do this so we have access to the underlying raw RPC functionalities in cases we do not have
	// standard functionalities from ethclient.
	rt.rpccli, err = rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	rt.cli = ethclient.NewClient(rt.rpccli)

	// Get the Validators information from the node
	err = rt.RefreshValidators()
	if err != nil {
		return nil, err
	}

	// Initialise the header cache
	rt.headerCache, err = lrucache.New(100)
	if err != nil {
		panic(err)
	}

	return rt, nil
}

func (rt *RedTNode) Close() {
	rt.cli.Close()
}

func (rt *RedTNode) EthClient() *ethclient.Client {
	return rt.cli
}

func (rt *RedTNode) RpcClient() *rpc.Client {
	return rt.rpccli
}

// BlockByHash returns the given full block.
func (rt *RedTNode) BlockByHash(blockhash common.Hash) (*ethertypes.Block, error) {

	// Set the timeout for the next call
	ctx, cancel := context.WithTimeout(context.Background(), rt.timeout)
	defer cancel()

	// Retrieve the block from the node (may be remote)
	block, err := rt.cli.BlockByHash(ctx, blockhash)
	if err != nil {
		return nil, err
	}

	return block, nil

}

// HeaderByNumber retrieves a header from the internal cache of from the blockchain node if cache miss
func (rt *RedTNode) HeaderByNumber(number int64) (*ethertypes.Header, error) {

	var blockNumber *big.Int

	// If number is negative, bypass the cache and retrieve the latest block from the blockchain node
	if number >= 0 {
		// Try to get the header from the cache
		cachedHeader, _ := rt.headerCache.Get(number)
		if cachedHeader != nil {
			return cachedHeader.(*ethertypes.Header), nil
		}

		blockNumber = big.NewInt(number)
	}

	// Set the timeout for the next call
	ctx, cancel := context.WithTimeout(context.Background(), rt.timeout)
	defer cancel()

	// Retrieve the header from the node (may be remote)
	head, err := rt.cli.HeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	// Add it to the cache
	rt.headerCache.Add(head.Number.Int64(), head)

	return head, err
}

// CurrentBlockNumber asks the blockchain node for the number of the latest current block
func (rt *RedTNode) CurrentBlockNumber() (int64, error) {
	header, err := rt.HeaderByNumber(-1)
	if err != nil {
		return 0, err
	}
	number := header.Number.Int64()
	return number, nil
}

func (rt *RedTNode) NodeInfo() (*p2p.NodeInfo, error) {
	var ni *p2p.NodeInfo

	// Set the timeout for the next call
	ctx, cancel := context.WithTimeout(context.Background(), rt.timeout)
	defer cancel()

	// Use the 'nodeInfo' API in the 'admin' namespace
	err := rt.rpccli.CallContext(ctx, &ni, "admin_nodeInfo")
	if err == nil && ni == nil {
		err = ethereum.NotFound
		return nil, err
	}
	return ni, err

}

func (rt *RedTNode) Peers() ([]*p2p.PeerInfo, error) {
	var peers []*p2p.PeerInfo

	// Set the timeout for the next call
	ctx, cancel := context.WithTimeout(context.Background(), rt.timeout)
	defer cancel()

	// Use the 'nodeInfo' API in the 'admin' namespace
	err := rt.rpccli.CallContext(ctx, &peers, "admin_peers")
	if err == nil && peers == nil {
		err = ethereum.NotFound
		return nil, err
	}
	return peers, err
}

func (rt *RedTNode) RefreshValidators() (err error) {

	// Load the current Validator set from the blockchain node as it is used in the execution of the consensus algorithm.
	// It has to be refreshed in case a new Validator is added or removed (an infrequent event)
	// Restarting the program loads the most recent Validator set
	rt.valSet, err = rt.getValSet()
	if err != nil {
		return err
	}

	// Retrieve dynamically the official Validator list from Alastria Github
	// Please note that some of them may not be in the current consensus validator set above
	//	enodes := readValidatorListFromGithub()
	enodes, err := readNodeListFromGithub("validator")
	if err != nil {
		return err
	}

	// Calculate the ethereum address from the enode string
	rt.allValidators = make(map[common.Address]*NodeInfo, len(enodes))
	for _, item := range enodes {

		// Parse the enode
		en := enode.MustParse(item.Enode)

		// The address is derived from the public key of the node
		address := crypto.PubkeyToAddress(*en.Pubkey())

		// Set the address in the corresponding entry
		item.Address = address

		// Update the array of Validators
		rt.allValidators[address] = item
	}

	return err
}

func (rt *RedTNode) Validators() []common.Address {
	return rt.valSet
}

func (rt *RedTNode) AllValidators() map[common.Address]*NodeInfo {
	return rt.allValidators
}

func (rt *RedTNode) ValidatorInfo(validator common.Address) *NodeInfo {
	return rt.allValidators[validator]
}

func (rt *RedTNode) DisplayMyInfo() {

	ni, err := rt.NodeInfo()
	if err != nil {
		logger.Error(err)
		return
	}

	fmt.Printf("About my node:\n")
	out, err := json.MarshalIndent(ni, "", "  ")
	if err != nil {
		logger.Error(err)
		return
	}
	fmt.Printf("%v\n\n", string(out))

}

func (rt *RedTNode) DisplayPeersInfo() {

	peers, err := rt.Peers()
	if err != nil {
		logger.Error(err)
		return
	}
	fmt.Printf("About my peers:\n")
	out, err := json.MarshalIndent(peers, "", "  ")
	if err != nil {
		logger.Error(err)
		return
	}
	fmt.Printf("%v\n\n", string(out))

}

// getValSet retrieves the current validator set executing the consensus algorithm in the network
func (rt *RedTNode) getValSet() ([]common.Address, error) {

	// The set of validators, no necessarily ordered in the round-robin order used in the consensus
	var unorderedVals []string

	// We are going to call the Geth API, with a timeout of 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := rt.rpccli.CallContext(ctx, &unorderedVals, "istanbul_getValidators", toBlockNumArg(-1))
	if err != nil {
		return nil, err
	}

	// In order to have the same order as in the IBFT consensus algorithm,
	// we have to sort addresses in string format by lexicographic order.
	// But the hex strings to order should be in the checked address Ethereum format
	// where some hex letters are uppercase and some are lowercase.
	// This is important because the letter "E" goes before letter "a", for example

	// First we convert the strings into Ethereum Addresses
	valSet := make([]common.Address, len(unorderedVals))
	for i, addrStr := range unorderedVals {
		valSet[i] = common.HexToAddress(addrStr)
	}

	// Now convert them back into strings but in Ethereum format (checked address)
	for i := range unorderedVals {
		unorderedVals[i] = valSet[i].String()
	}

	// Sort the resulting slice lexicographically
	sort.Strings(unorderedVals)

	// And finally get the slice with Addresses in the right order
	for i, addrStr := range unorderedVals {
		valSet[i] = common.HexToAddress(addrStr)
	}

	return valSet, nil

}

// SignersFromBlock returns the author of the block (proposer) and the list of signers
func SignersFromBlock(header *ethertypes.Header) (author common.Address, signers []common.Address, err error) {

	// Retrieve the signature from the header extra-data
	extra, err := ethertypes.ExtractIstanbulExtra(header)
	if err != nil {
		return author, signers, err
	}

	// Recover the address that signed the block
	author, err = istanbul.GetSignatureAddress(sigHash(header).Bytes(), extra.Seal)
	if err != nil {
		return author, signers, err
	}

	// committedSeal is the array of signatures, one per each node that signed the block
	committedSeal := extra.CommittedSeal

	// proposalSeal is the data that was signed by each and all nodes. The data is the hash of the block
	proposalSeal := ibftengine.PrepareCommittedSeal(header.Hash())

	// To recover each address we need the data that was signed (hash of header) and the signature
	// We loop the array of signatures to recover the corresponding address for each entry

	// Get committed seals from current header
	for _, seal := range committedSeal {
		// Get the original address by seal and parent block hash
		addr, err := istanbul.GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			return author, signers, err
		}
		signers = append(signers, addr)
	}

	return
}

func DisplayPeersInfo(url string) {

	// Connect to the RedT node
	rt, err := NewRedTNode(url)
	if err != nil {
		logger.Error(err)
		return
	}

	rt.DisplayPeersInfo()

}

func readNodeListFromGithub(name string) ([]*NodeInfo, error) {
	var targetUrl string = "https://raw.githubusercontent.com/alastria/alastria-node-quorum-directory/main/"

	// Build the full url
	switch name {
	case "validator":
		targetUrl = targetUrl + "directory.val"
	case "regular":
		targetUrl = targetUrl + "directory.reg"
	case "boot":
		targetUrl = targetUrl + "directory.bot"
	default:
		return nil, fmt.Errorf("invalid node list name")
	}

	// Read the raw file, with one line per Validator
	// Each line is like: VAL_IN2 enode://0ede782b7ce6...2b8e64059adc03@15.236.56.133:21000?discport=0
	// Where the dots are an ellipsis to make the line shorter
	resp, err := http.Get(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("error getting list from Github: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error getting list from Github: %v", err)
	}

	// Process each line
	tokens := bytes.Fields(body)
	num_tokens := len(tokens)

	// Some sanity checks
	// The number of tokens should be even (and greater than zero)
	if num_tokens <= 0 || num_tokens%2 != 0 {
		return nil, fmt.Errorf("invalid number of tokens: %v", len(tokens))
	}

	// Generate the enodes array
	var enodes = make([]*NodeInfo, num_tokens/2)
	for i, t := range tokens {
		if i%2 == 0 {
			// The even entries have the name of the entity operating the node

			// Create the enode item
			enode := &NodeInfo{}
			enodes[i/2] = enode

			// If the name starts with "VAL_" and similar, trim it to display a simpler
			t = bytes.TrimPrefix(t, []byte("VAL_"))
			t = bytes.TrimPrefix(t, []byte("BOT_"))
			t = bytes.TrimPrefix(t, []byte("REG_"))

			enodes[i/2].Operator = string(t)

		} else {
			// The odd tokens should start with "enode://"
			if bytes.HasPrefix(t, []byte("enode://")) {
				enodes[i/2].Enode = string(t)
			} else {
				return nil, fmt.Errorf("invalid enode: %v", string(t))
			}
		}
	}

	return enodes, nil

}
