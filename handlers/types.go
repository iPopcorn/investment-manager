package handlers

import "github.com/spf13/cobra"

type CobraCommandHandler func(cmd *cobra.Command, args []string) error
