package redt

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethertypes "github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type StatisticsRedT struct {
	allValidators      map[common.Address]*NodeInfo
	valSet             []common.Address
	asProposer         map[common.Address]int
	asSigner           map[common.Address]int
	lastBlockProcessed int64
	previousTimestamp  uint64
	elapsed            uint64
	allTxsCount        int
}

// NewStatistics creates and initializes a statistics object that will accumulate counters.
func NewStatistics(allValidators map[common.Address]*NodeInfo, valSet []common.Address) *StatisticsRedT {
	st := &StatisticsRedT{}

	// All configured Validators
	st.allValidators = allValidators

	// The subset currently in the consensus set
	st.valSet = valSet

	// Initialise the counters for validators/signers
	st.asProposer = map[common.Address]int{}
	st.asSigner = map[common.Address]int{}

	return st
}

type StatisticsSummary struct {
	BlockNumber      int64
	Elapsed          uint64
	Timestamp        time.Time
	ProposerName     string
	ProposerCount    int
	ProposerAddress  common.Address
	NextProposerName string
	GasLimit         uint64
	GasUsed          uint64
	GasLimitH        string
	GasUsedH         string
	BlockNumTxs      int
	AllNumTxs        int
	Signers          []Signer
}

type Signer struct {
	NextProposer    bool
	CurrentProposer bool
	AsProposer      int
	CurrentSigner   bool
	AsSigner        int
	Name            string
	Address         common.Address
}

// StatisticsForBlock returns a map with statistical data prepared for HTML templates
func (st *StatisticsRedT) StatisticsForBlock(fullBlock *ethertypes.Block) (*StatisticsSummary, error) {

	fmt2 := message.NewPrinter(language.EuropeanSpanish)
	data := &StatisticsSummary{}

	// Get the header
	header := fullBlock.Header()

	// Update the statistics counters
	author, signers, err := st.UpdateStatisticsForBlock(fullBlock)
	if err != nil {
		return nil, err
	}

	// Get the info of the node operator
	oper := st.ValidatorInfo(author)

	// Determine the next node that should be proposer, according to the round-robin
	// selection algorithm
	var nextProposer common.Address
	for i := 0; i < len(st.valSet); i++ {
		if author == st.valSet[i] {
			nextProposer = st.valSet[(i+1)%len(st.valSet)]
			break
		}
	}

	// Calculate values for the header
	data.BlockNumber = header.Number.Int64()
	data.Elapsed = st.elapsed
	data.Timestamp = time.Unix(int64(header.Time), 0)
	data.ProposerName = oper.Operator
	data.ProposerCount = st.asProposer[author]
	data.ProposerAddress = author
	data.NextProposerName = st.ValidatorInfo(nextProposer).Operator
	data.GasLimit = header.GasLimit
	data.GasUsed = header.GasUsed
	data.BlockNumTxs = fullBlock.Transactions().Len()
	data.AllNumTxs = st.allTxsCount
	data.GasLimitH = fmt2.Sprintf("%d", header.GasLimit)
	data.GasUsedH = fmt2.Sprintf("%d", header.GasUsed)

	// Create the map with signers of this block
	var currentSigners = map[common.Address]bool{}
	for _, seal := range signers {
		currentSigners[seal] = true
	}

	// Calculate the table of validators
	validatorTable := make([]Signer, len(st.ValidatorSet()))

	for i, addr := range st.ValidatorSet() {

		// Mark if this is the next proposer
		if nextProposer.Hex() == addr.Hex() {
			validatorTable[i].NextProposer = true
		}

		// Mark if this is the proposer of this block
		if author.Hex() == addr.Hex() {
			validatorTable[i].CurrentProposer = true
		}

		// Number of times this has been a proposer
		validatorTable[i].AsProposer = st.asProposer[addr]

		// Mark is this is a signer of current block
		if currentSigners[addr] {
			validatorTable[i].CurrentSigner = true
		}

		// Number of times this has been a signer
		validatorTable[i].AsSigner = st.asSigner[addr]

		// Name of the node
		validatorTable[i].Name = st.ValidatorInfo(addr).Operator

		// Address of the node
		validatorTable[i].Address = addr
	}

	data.Signers = validatorTable

	return data, nil

}

// UpdateStatisticsForBlock uses the info of the block to update the statistics.
// It receives a full Block and returns the address of the author and an array with the signers of the block
func (st *StatisticsRedT) UpdateStatisticsForBlock(fullblock *ethertypes.Block) (author common.Address, signers []common.Address, err error) {

	// Get the header
	header := fullblock.Header()

	// Get the author (proposer) and signers for this block
	author, signers, err = SignersFromBlock(header)
	if err != nil {
		return author, signers, err
	}

	// Do nothing if the block was already processed
	thisBlockNumber := header.Number.Int64()
	if thisBlockNumber <= st.lastBlockProcessed {
		return author, signers, nil
	}

	// Update the last block processed
	st.lastBlockProcessed = thisBlockNumber

	// Calculate elapsed time since last block processed
	st.elapsed = header.Time - st.previousTimestamp

	// And update the last timestamp
	st.previousTimestamp = header.Time

	// Update the number of transactions since the initialization of statistics
	st.allTxsCount = st.allTxsCount + fullblock.Transactions().Len()

	// Increment the counter for authors
	st.asProposer[author] += 1

	// Increment counters for signers
	for _, seal := range signers {
		// Increment the counter of signatures
		st.asSigner[seal] += 1
	}

	return author, signers, nil

}

// ValidatorSet returns the current set of validator addresses.
// It does not make a request to the blockchain node but instead uses the set calculated
// when the statistics object was created.
func (st *StatisticsRedT) ValidatorSet() []common.Address {
	return st.valSet
}

// ValidatorInfo returns the validator info object associated to the specified address.
func (st *StatisticsRedT) ValidatorInfo(validator common.Address) *NodeInfo {
	return st.allValidators[validator]
}
