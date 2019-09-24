package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pegnet/pegnet/modules/opr"

	"github.com/Factom-Asset-Tokens/factom"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyze)
}

var analyze = &cobra.Command{
	Use:   "analyze",
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

		gencount := make(map[string]int)
		idcount := make(map[string]int)
		notidcount := make(map[string]int)
		notvalid := 0
		for _, e := range eb.Entries {
			o, err := opr.Parse(e.Content)
			if err != nil {
				notvalid++
				continue
			}

			id := o.GetID()
			baseid := strings.TrimRight(id, "0123456789")
			gencount[baseid]++
			idcount[id]++

			if o.GetHeight() != int32(eb.Height) {
				notidcount[id]++
				continue
			}
		}

		var data []Miner
		for k, v := range gencount {
			data = append(data, Miner{
				name:  k,
				valid: v,
				//notvalid: notidcount[k],
			})
		}

		sort.Slice(data, func(i, j int) bool {
			if data[i].valid == data[j].valid {
				return data[i].name < data[j].name
			}
			return data[i].valid > data[j].valid
		})

		fmt.Printf("## Height %d\nHash: %s\nEntries: %d\n\n", eb.Height, eb.KeyMR, len(eb.Entries))
		fmt.Println("|Prefix|Entries|Miners|Max Individual|\n|---|---|---|---|")
		for _, d := range data {
			sub := get(d.name, idcount)
			if len(sub) > 1 {
				fmt.Printf("|%s|%d|%d|%d", d.name, d.valid, len(sub), sub[0].valid)
			} else {
				fmt.Printf("|%s|%d|-|-", d.name, d.valid)
			}

			/*			if len(sub) > 1 {
						br := ""
						for _, x := range sub {
							fmt.Printf("`%d %s`%s", x.valid, x.name, br)
							br = "<br>"
						}

					}*/

			fmt.Println("|")
		}
	},
}

type Miner struct {
	name     string
	valid    int
	notvalid int
}

func get(prefix string, data map[string]int) []Miner {
	var m []Miner

	for k, v := range data {
		if strings.HasPrefix(k, prefix) && k != strings.TrimRight(k, "0123456789") {
			m = append(m, Miner{
				name:  k,
				valid: v,
				//notvalid: notidcount[k],
			})
		}
	}

	sort.Slice(m, func(i, j int) bool {
		if m[i].valid == m[j].valid {
			return m[i].name < m[j].name
		}
		return m[i].valid > m[j].valid
	})

	return m
}
