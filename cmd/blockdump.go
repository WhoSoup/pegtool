package cmd

import (
	"fmt"

	"github.com/pegnet/pegnet/modules/opr"

	"github.com/Factom-Asset-Tokens/factom"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(blockdump)
}

var blockdump = &cobra.Command{
	Use:   "blockdump",
	Short: "Analyzes an EBlock",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cli := &factom.Client{FactomdServer: "http://spoon:8088"}

		var eb factom.EBlock
		eb.KeyMR = factom.NewBytes32FromString(args[0])
		err := eb.Get(cli)
		if err != nil {
			panic(err)
		}
		err = eb.GetEntries(cli)
		if err != nil {
			panic(err)
		}

		total := 0
		byaddr := make(map[string]map[string]int)
		for _, e := range eb.Entries {
			o, err := opr.Parse(e.Content)
			if err != nil {
				continue
			}
			total++

			id := o.GetID()
			addr := o.GetAddress()

			if _, b := byaddr[addr]; !b {
				byaddr[addr] = make(map[string]int)
			}

			byaddr[addr][id]++
		}

		fmt.Println("Total OPRs:", total)
		fmt.Println("Total Addresses:", len(byaddr))

		for addr, more := range byaddr {
			fmt.Printf("%s:\n", addr)
			for id, count := range more {
				fmt.Printf("\t%3d %s\n", count, id)
			}
		}
	},
}
