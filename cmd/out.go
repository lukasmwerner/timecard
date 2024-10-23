/*
Copyright © 2024 Lukas Werner <me@lukaswerner.com>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/lukasmwerner/timecard/store"
	"github.com/spf13/cobra"
)

// outCmd represents the out command
var outCmd = &cobra.Command{
	Use:   "out",
	Short: "Clocks into work",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := store.Open()
		if err != nil {
			log.Fatalln(err.Error())
		}

		db.PunchOut("")

		fmt.Println("punched out")

	},
}

func init() {
	rootCmd.AddCommand(outCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// outCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// outCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
