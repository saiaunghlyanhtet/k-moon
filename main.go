package main

import (
	"github.com/saiaunghlyanhtet/k-moon/cmd/mem"

	"github.com/saiaunghlyanhtet/k-moon/cmd/cpu"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "ebpfmon"}
	rootCmd.AddCommand(cpu.Cmd)
	rootCmd.AddCommand(mem.Cmd)
	rootCmd.Execute()
}
