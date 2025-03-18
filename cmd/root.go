package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	api      string
	db       string
	join     *[]string
	dir      string
	verbose  bool
	diskMode bool
	crt      string
	key      string
)

var rootCmd = &cobra.Command{
	Use:   "dqldemo",
	Short: "dqldemo showcases distributed systems concepts with dqlite",
	Long:  "A small API that allows setting key value pairs in a distributed sqlite database.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		fmt.Println("Hello, World!")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&api, "api", "", "api endpoint to use for serving")
	rootCmd.PersistentFlags().StringVar(&db, "db", "", "database name")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", "", "path for storing sqlite data")
	rootCmd.PersistentFlags().StringVar(&crt, "crt", "", "certificate file for TLS")
	rootCmd.PersistentFlags().StringVar(&key, "key", "", "key file for certificate")

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&diskMode, "disk-mode", false, "use diskmode")

	join = rootCmd.PersistentFlags().StringSlice("join", []string{}, "List of nodes to form a cluster")

	if err := rootCmd.MarkPersistentFlagRequired("api"); err != nil {
		panic(err)
	}

	if err := rootCmd.MarkPersistentFlagRequired("db"); err != nil {
		panic(err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
