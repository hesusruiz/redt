package main

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hesusruiz/redt"
	"github.com/pterm/pterm"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	spinner *pterm.SpinnerPrinter
	rt      *redt.RedTNode
	stats   *redt.StatisticsRedT
	fmt2    *message.Printer
)

func MonitorSigners(url string, refresh int64) error {
	var err error

	fmt2 = message.NewPrinter(language.EuropeanSpanish)

	// Connect to the RedT node
	rt, err = redt.NewRedTNode(url)
	if err != nil {
		panic(err)
	}

	spinner, _ = pterm.DefaultSpinner.Start("Connecting to ", url, " ...")
	spinner.RemoveWhenDone = true

	// Check that we can get info about the node
	thisNode, err := rt.NodeInfo()
	if err != nil {
		redt.Logger.Fatal(err)
	}
	pterm.Println("Connected to: " + pterm.Green(thisNode.Name))

	// Create new statistics
	stats = redt.NewStatistics(rt.AllValidators(), rt.Validators())

	// Start from an old block
	latestBlock, err := rt.BlockByNumber(-1)
	if err != nil {
		redt.Logger.Fatal(err)
	}

	latestNumber := int64(latestBlock.NumberU64())

	_, err = stats.StatisticsForBlock(latestBlock)
	if err != nil {
		panic(err)
	}

	spinner.Stop()

	for {

		// Sleep before getting the next one
		time.Sleep(time.Duration(refresh) * time.Second)

		// Get the current block number
		currentHeader, err := rt.HeaderByNumber(-1)
		if err != nil {
			redt.Logger.Fatal(err)
		}
		currentNumber := currentHeader.Number.Int64()

		// Make sure we have advanced at least one block
		if currentNumber == latestNumber {
			continue
		}

		// Display all blocks from the latest one until the current one
		for i := latestNumber + 1; i <= currentNumber; i++ {
			b, err := rt.BlockByNumber(i)
			if err != nil {
				return err
			}

			summary, err := stats.StatisticsForBlock(b)
			if err != nil {
				panic(err)
			}

			// Check if there are any big transactions
			var bigTxs types.Transactions
			transactions := b.Transactions()
			for _, tx := range transactions {
				if tx.Gas() >= 1_000_000 {
					bigTxs = append(bigTxs, tx)
				}
			}

			DisplaySigners(rt, summary, bigTxs)

		}

		// Update the latest block number that we processed
		latestNumber = currentNumber

	}

}

func DisplaySigners(rt *redt.RedTNode, sum *redt.StatisticsSummary, bigTxs types.Transactions) {

	if spinner != nil && spinner.IsActive {
		spinner.Stop()
	}

	tableMsg := ""

	// Print the title of the table
	tableMsg += pterm.Sprintf("  Author |  Signer  |       Name      Address\n")

	for _, val := range sum.Signers {

		var authorCountStr string
		if val.AsProposer == 0 {
			authorCountStr = pterm.FgRed.Sprintf("%6v", val.AsProposer)
		} else {
			authorCountStr = pterm.Sprintf("%6v", val.AsProposer)
		}

		if val.CurrentProposer {
			authorCountStr = pterm.BgLightBlue.Sprint(pterm.Bold.Sprintf("%v %1v", authorCountStr, "X"))
		} else {
			authorCountStr = pterm.Bold.Sprintf("%v %1v", authorCountStr, " ")
		}

		var signerCountStr string
		if val.AsSigner == 0 {
			signerCountStr = pterm.FgRed.Sprintf("%6v", val.AsSigner)
		} else {
			signerCountStr = pterm.Sprintf("%6v", val.AsSigner)
		}

		if val.CurrentSigner {
			signerCountStr = pterm.BgLightBlue.Sprint(pterm.Bold.Sprintf("%v %1v", signerCountStr, "X"))
		} else {
			signerCountStr = pterm.Bold.Sprintf("%v %1v", signerCountStr, " ")
		}

		name := val.Name
		if len(name) > 12 {
			name = name[:12]
		}
		tableMsg += pterm.Sprintf("%v | %v | %12v %v\n", authorCountStr, signerCountStr, name, val.Address)

	}

	// Create the box for display of the info inside
	var blockNum string
	if sum.Elapsed < 6 {
		blockNum = pterm.Green(sum.BlockNumber)
	} else {
		blockNum = pterm.Red(sum.BlockNumber)
	}
	blockInfo := pterm.DefaultBox.WithTitle("Block " + blockNum)

	// Build the header message, in red if block time was bad
	headerMsg1 := pterm.Sprintf("NumTxs %v Elapsed %vs %v\n", pterm.Green(sum.BlockNumTxs), pterm.Green(sum.Elapsed), sum.Timestamp)

	// The author info
	headerMsg2 := pterm.Sprintf("Author: %v (%v) (%v)\n", sum.ProposerName, sum.ProposerCount, sum.ProposerAddress)

	// Gas limit and number of txs
	headerMsg3 := pterm.Sprintf("GasUsed: %v from %v\n", pterm.Green(sum.GasUsedH), sum.GasLimitH)
	if sum.GasUsed > 5_000_000 {
		headerMsg3 = pterm.FgRed.Sprintf("GasUsed: %v from %v\n\n", sum.GasUsedH, sum.GasLimitH)
	} else if sum.GasUsed > 1_000_000 {
		headerMsg3 = pterm.FgYellow.Sprintf("GasUsed: %v from %v\n\n", sum.GasUsedH, sum.GasLimitH)
	}

	blockInfo.Println(headerMsg1 + headerMsg2 + headerMsg3 + tableMsg + DisplayBigTxs(rt, bigTxs))

	spinner, _ = pterm.DefaultSpinner.Start("Waiting for ", sum.NextProposerName, " to create next block ...")
	spinner.RemoveWhenDone = true

}

func DisplayBigTxs(rt *redt.RedTNode, bigTxs types.Transactions) string {

	var bigTxsStr string
	for _, tx := range bigTxs {
		bigTxsStr += pterm.FgRed.Sprintf("\nTx with Gas %v and Size %v\n", fmt2.Sprint(tx.Gas()), tx.Size())

		txReceipt, err := rt.EthClient().TransactionReceipt(context.Background(), tx.Hash())
		if err == nil {
			var gasUsed string
			if txReceipt.GasUsed > 10_000_000 {
				gasUsed = pterm.Red(fmt2.Sprint(txReceipt.GasUsed))
			} else {
				gasUsed = pterm.Green(fmt2.Sprint(txReceipt.GasUsed))
			}
			if txReceipt.Status == 0 {
				bigTxsStr += pterm.Yellow("  GasUsed ") + gasUsed + pterm.Yellow(" but has been Reverted\n")
			} else {
				bigTxsStr += pterm.Sprintf("  GasUsed %v\n", gasUsed)
			}
		}

		bigTxsStr += pterm.Sprintf("  hash %v\n", tx.Hash())
		bigTxsStr += pterm.Sprintf("  from %v to %v\n", tx.From(), tx.To())

	}

	return bigTxsStr

}
