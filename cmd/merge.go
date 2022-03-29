/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	merge "github.com/stewartbutler/go-merge-imports/pkg"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Execute merge",
	Long: `Executes a git merge on a file. Unions the imports blocks then hands
off to 'git merge-file'.`,
	Args:       cobra.MinimumNArgs(3),
	ArgAliases: []string{"current", "base", "other"},
	Run:        runMerge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runMerge(cmd *cobra.Command, args []string) {
	m := merge.NewMerge(args[0], args[1], args[2])
	m.MergeFile()
	//m.CallNextBinary()
}
