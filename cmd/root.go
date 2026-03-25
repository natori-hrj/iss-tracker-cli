package cmd

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/natori-hrj/iss-tracker-cli/internal/tui"
)

var interval int

var rootCmd = &cobra.Command{
	Use:   "iss-tracker",
	Short: "Track the International Space Station in your terminal",
	Long: `ISS Tracker - A CLI tool to track the International Space Station.

Displays the ISS position on an ASCII world map, shows astronauts
currently in space, calculates distance from your location, and
estimates the next pass time over your area.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		duration := time.Duration(interval) * time.Second
		model := tui.NewModel(duration)

		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&interval, "interval", "i", 5, "Update interval in seconds")
}
