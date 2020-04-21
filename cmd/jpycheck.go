package cmd

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/pegnet/pegnet/modules/opr"

	"github.com/Factom-Asset-Tokens/factom"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(jpycheck)
}

type Tie struct {
	opr  *opr.V4Content
	hash string
	diff uint64
}

var jpycheck = &cobra.Command{
	Use:   "jpycheck",
	Short: "checks oprs for jpy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cli := &factom.Client{FactomdServer: "http://courtesy-node.factom.com:80"}

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

		oprs := make([]Tie, 0)

		for _, e := range eb.Entries {
			o, err := opr.ParseV2Content(e.Content)
			if err != nil {
				continue
			}
			total++
			o2 := &opr.V4Content{V2Content: *o}

			if o2.GetHeight() != int32(eb.Height) {
				continue
			}

			oprs = append(oprs, Tie{
				opr:  o2,
				hash: e.Hash.String(),
				diff: binary.BigEndian.Uint64(e.ExtIDs[1]),
			})
		}

		sort.Slice(oprs, func(i, j int) bool {
			return oprs[i].diff > oprs[j].diff
		})

		fmt.Println("===================================")
		fmt.Println("EBlock", eb.Height, eb.KeyMR.String())
		fmt.Println("===================================")
		fmt.Printf("ID  %-16s %-64s %s\n", "Miner ID", "Entry Hash", "JPY")
		attack := 0
		for p, tie := range oprs {
			jpy := 0.0
			for _, a := range tie.opr.GetOrderedAssetsFloat() {
				if a.Name == "JPY" {
					jpy = a.Value
					break
				}
			}
			if jpy > 1 {
				attack++
			}
			fmt.Printf("%2d. %-16s %s %f\n", p+1, tie.opr.GetID(), tie.hash, jpy)
		}
		fmt.Println("===================================")
		fmt.Printf("Malicious OPRs: %d / %d\n", attack, 50)
		fmt.Println("===================================")

	},
}
