package main

import (
	"os"
	"time"

	"github.com/hesusruiz/redt"
	"github.com/urfave/cli/v2"
)

// The default node address used is a local one
var localNodeHTTP = "http://127.0.0.1:22000"
var remoteNodeHTTP = "http://15.236.0.91:22000"
var localNodeWS = "ws://127.0.0.1:22001"

func main() {

	// Define commands, parse command line arguments and start execution
	app := &cli.App{
		Usage: "Monitoring of block signers activity for the Alastria RedT blockchain network",
		UsageText: `signers [global options] command [command options]
			where 'nodeURL' is the address of the Quorum node.
			It supports both HTTP and WebSockets endpoints.
			By default it uses WebSockets and for HTTP you have to use the 'poll' subcommand.`,

		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Version:                "v0.1",
		Compiled:               time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Jesus Ruiz",
				Email: "hesus.ruiz@gmail.com",
			},
		},

		Action: func(c *cli.Context) error {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		},
	}

	monitorCMD := &cli.Command{
		Name:      "poll",
		Usage:     "monitor the signers activity via HTTP polling",
		UsageText: "signers poll [options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Value:    remoteNodeHTTP,
				Usage:    "url of the endpoint of blockchain node",
				Aliases:  []string{"u"},
				Required: false,
			},
			&cli.Int64Flag{
				Name:    "refresh",
				Value:   2,
				Usage:   "refresh interval for presentation. All blocks are processed independent of this value",
				Aliases: []string{"r"},
			},
		},

		Action: func(c *cli.Context) error {
			url := c.String("url")
			refresh := c.Int64("refresh")
			err := MonitorSigners(url, refresh)
			return err
		},
	}

	displayPeersCMD := &cli.Command{
		Name:  "peers",
		Usage: "display peers information",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Value:    localNodeHTTP,
				Usage:    "url of the endpoint of blockchain node",
				Aliases:  []string{"u"},
				Required: false,
			},
		},

		Action: func(c *cli.Context) error {
			url := c.String("url")
			redt.DisplayPeersInfo(url)
			return nil
		},
	}

	app.Commands = []*cli.Command{
		monitorCMD,
		displayPeersCMD,
	}

	// Run the application
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}
