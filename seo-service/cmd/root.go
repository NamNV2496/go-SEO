package cmd

import "github.com/spf13/cobra"

func Execute() error {
	rootCmd := &cobra.Command{
		Short: "A simple web seo",
	}
	rootCmd.AddCommand(serverCmd)
	// rootCmd.AddCommand(CrawlerWorkerCmd)
	return rootCmd.Execute()
}
