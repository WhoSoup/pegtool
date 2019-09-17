package cmd

import (
	"encoding/binary"
	"fmt"

	"github.com/pegnet/pegnet/modules/opr"

	"github.com/Factom-Asset-Tokens/factom"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decode)
}

var decode = &cobra.Command{
	Use:   "decode",
	Short: "Decodes an OPR entry",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cli := &factom.Client{FactomdServer: "https://api.factomd.net"}
		var e factom.Entry
		e.ChainID = factom.NewBytes32FromString("a642a8674f46696cc47fdb6b65f9c87b2a19c5ea8123b3d2f0c13b6f33a9d5ef")
		e.Hash = factom.NewBytes32FromString(args[0])
		err := e.Get(cli)
		if err != nil {
			cmd.PrintErr(err)
		}

		fmt.Println("===== HEADER ======")
		fmt.Printf("  Entry: %s\n", e.Hash)
		fmt.Printf("Version: %d\n", e.ExtIDs[2][0])
		fmt.Printf("  Nonce: %x\n", e.ExtIDs[0])

		srd := binary.BigEndian.Uint64(e.ExtIDs[1])
		fmt.Printf("    SRD: %d\n", srd)
		fmt.Println("===== CONTENT =====")

		o, err := opr.Parse(e.Content)
		if err != nil {
			cmd.PrintErr(err)
		}

		fmt.Printf("Address: %s\n", o.GetAddress())
		fmt.Printf("ID: %s\n", o.GetID())
		fmt.Printf("Height: %d\n", o.GetHeight())
		fmt.Printf("PrevWinners:\n")
		prev := o.GetPreviousWinners()

		for i := 0; i < len(prev); i += 5 {
			fmt.Printf("\t")
			for j := 0; j < 5; j++ {
				fmt.Printf("%s ", prev[i+j])
			}
			fmt.Println()
		}
		fmt.Printf("Assets:\n")

		for _, a := range o.GetOrderedAssetsFloat() {
			fmt.Printf("\t%4s = %f\n", a.Name, a.Value)
		}

	},
}
