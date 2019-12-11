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
	rootCmd.AddCommand(sortEB)
}

var sortEB = &cobra.Command{
	Use:   "sort",
	Short: "Sort an EBlock",
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

		type entry struct {
			diff uint64
			addr string
			id   string
		}

		entries := make([]entry, 0)
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

			var diff uint64
			diff = binary.BigEndian.Uint64(e.ExtIDs[1])

			entries = append(entries, entry{
				diff: diff,
				addr: addr,
				id:   id,
			})
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].diff > entries[j].diff
		})

		for i, e := range entries {
			fmt.Println(i, e)
		}
	},
}
