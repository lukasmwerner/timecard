/*
Copyright © 2024 Lukas Werner <me@lukaswerner.com>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
	"github.com/lukasmwerner/timecard/store"
	"github.com/spf13/cobra"
)

var (
	ClockedInStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#EEF572")).
			Align(lipgloss.Center).
			Width(7)
	ClockedOutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#DB45BE")).
			Align(lipgloss.Center).
			Width(7)
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "shows you your clock status",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := store.Open()
		if err != nil {
			log.Fatalln(err.Error())
		}

		clockedIn, timestamp, err := db.Status()
		if err != nil {
			log.Fatalln(err.Error())
		}

		if clockedIn {
			fmt.Println(ClockedInStyle.Render("In"), "Working since:", timestamp.Format("15:04:05"))
		} else {
			fmt.Println(ClockedOutStyle.Render("Out"), "Currently off the clock")
		}

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
