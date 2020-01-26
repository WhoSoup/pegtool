package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pegnet/pegnet/modules/opr"

	"github.com/Factom-Asset-Tokens/factom"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(prevWinners)
}

var prevWinners = &cobra.Command{
	Use:   "prev-winners",
	Short: "Attempts to guess the previous winners of an EBlock",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cli := &factom.Client{FactomdServer: "http://spoon:8088"}

		height, err := strconv.Atoi(args[0])
		if err == nil {
			var dblock factom.DBlock
			dblock.Header.Height = uint32(height)
			err = dblock.Get(cli)
			if err != nil {
				panic(err)
			}

			for _, eb := range dblock.EBlocks {
				if eb.ChainID.String() == "a642a8674f46696cc47fdb6b65f9c87b2a19c5ea8123b3d2f0c13b6f33a9d5ef" {
					args[0] = eb.KeyMR.String()

					break
				}
			}
		}

		var eb factom.EBlock
		eb.KeyMR = factom.NewBytes32FromString(args[0])
		err = eb.Get(cli)
		if err != nil {
			panic(err)
		}
		err = eb.GetEntries(cli)
		if err != nil {
			panic(err)
		}

		total := 0
		prev := make(map[string]int)
		ids := make(map[string]map[string]int)
		for _, e := range eb.Entries {
			o, err := opr.Parse(e.Content)
			if err != nil {
				continue
			}
			total++

			prevs := "\"" + strings.Join(o.GetPreviousWinners(), "\",\"") + "\""
			prev[prevs]++
			if _, ok := ids[prevs]; !ok {
				ids[prevs] = make(map[string]int)
			}
			ids[prevs][o.GetID()]++
		}

		sorted := make([]string, 0)
		for s := range prev {
			sorted = append(sorted, s)
		}
		sort.Slice(sorted, func(i, j int) bool {
			return prev[sorted[j]] < prev[sorted[i]]
		})

		for _, s := range sorted {
			fmt.Printf("%.2f%% %s\n", float64(prev[s])/float64(total)*100, s)
			fmt.Printf("miners %s\n\n", idfs(ids[s]))
		}
	},
}

func idfs(data map[string]int) string {
	sortIDs := make([]string, 0)
	for s := range data {
		sortIDs = append(sortIDs, s)
	}

	sort.Slice(sortIDs, func(i, j int) bool {
		return data[sortIDs[j]] < data[sortIDs[i]]
	})
	return strings.Join(sortIDs, ",")
}
