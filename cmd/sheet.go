/*
Copyright © 2024 Lukas Werner <me@lukaswerner.com>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/lukasmwerner/timecard/store"
	"github.com/spf13/cobra"
)

// sheetCmd represents the sheet command
var sheetCmd = &cobra.Command{
	Use:   "sheet",
	Short: "Presents your timesheet",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := store.Open()
		if err != nil {
			log.Fatalln(err.Error())
		}

		t := table.New().Border(lipgloss.NormalBorder()).Headers("DAY", "HOURS", "PUNCHES")

		days, totalHours, err := db.Report()

		fmt.Println("This Week's Hours:", totalHours)

		for _, day := range days {
			t.Row(day.Date.Format("Mon Jan _2"), fmt.Sprint(day.Hours), string(day.ClockInCount))
		}

		fmt.Println(t.Render())
	},
}

func init() {
	rootCmd.AddCommand(sheetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sheetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sheetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
