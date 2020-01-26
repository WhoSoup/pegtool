package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Factom-Asset-Tokens/factom"
	"github.com/pegnet/pegnet/modules/opr"
	"github.com/spf13/cobra"
)

// creates a unit test output for pegnet/module/grader/testdata unit tests

func init() {
	rootCmd.AddCommand(unitTestDump)
}

var unitTestDump = &cobra.Command{
	Use:   "unit-test",
	Short: "Dumps an EBlock",
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

		data := make(map[string]interface{})
		data["Height"] = eb.Height

		data["Winners"] = make([]string, 25)

		entries := make([]map[string]interface{}, 0)

		total := 0
		prev := make(map[string]int)
		for _, e := range eb.Entries {
			entry := make(map[string]interface{})
			entry["Hash"] = e.Hash.String()
			extids := make([][]byte, 0)
			for _, xt := range e.ExtIDs {
				extids = append(extids, xt)
			}
			entry["ExtIDs"] = extids
			entry["Content"] = []byte(e.Content)

			entries = append(entries, entry)

			o, err := opr.Parse(e.Content)
			if err != nil {
				continue
			}
			total++

			prevs := strings.Join(o.GetPreviousWinners(), ",")
			prev[prevs]++
		}
		data["Entries"] = entries

		sorted := make([]string, 0)
		for s := range prev {
			sorted = append(sorted, s)
		}

		sort.Slice(sorted, func(i, j int) bool {
			return prev[sorted[j]] < prev[sorted[i]]
		})

		if len(sorted) > 0 {
			wins := strings.Split(sorted[0], ",")
			data["PreviousWinners"] = wins
		} else {
			data["PreviousWinners"] = make([]string, 10)
		}

		js, _ := json.Marshal(data)
		fmt.Println(string(js))
	},
}
