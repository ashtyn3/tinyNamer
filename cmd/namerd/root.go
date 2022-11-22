package cmd

import (
	"os"
	"strings"

	"github.com/ashtyn3/tinynamer/p2p"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func level(lvl interface{}) string {
	switch lvl.(string) {
	case "info":
		return color.GreenString("INFO")
	case "warn":
		return color.YellowString("WARN")
	case "error":
		return color.RedString("ERROR")
	case "fatal":
		return color.RedString("FATAL")
	default:
		return strings.ToUpper(lvl.(string))
	}
}

var rootCmd = &cobra.Command{
	Use:   "tinynamer",
	Short: "A TinyNamer client",
	Long:  `A client written in Golang for the TinyNamer blockchain.`,
	Run: func(cmd *cobra.Command, args []string) {
		l, _ := cmd.Flags().GetBool("log")
		if l == false {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, FormatLevel: level, TimeFormat: "[01-02|03:04:05PM]"})
		} else {
			home, _ := os.UserHomeDir()
			log_file, _ := os.Create(home + "/.tinyNamer/log")
			log.Logger = log.Output(log_file)
		}
		n := p2p.NewNode()
		port, _ := cmd.Flags().GetString("port")
		n.Run(port)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("port", "p", "5770", "The network port of the client.")
	rootCmd.Flags().BoolP("log", "L", false, "Log to file")
}
