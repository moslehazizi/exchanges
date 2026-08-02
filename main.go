package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	accountbitget "order-maker/account/bitget"
	authbitpin "order-maker/auth/bitpin"
	authnobitex "order-maker/auth/nobitex"
	balancebitget "order-maker/balance/bitget"
	balancenobitex "order-maker/balance/nobitex"
	orderbitget "order-maker/order/bitget"
)

func main() {
	root := &cobra.Command{Use: "order-maker"}

	accountCmd := &cobra.Command{Use: "account"}
	accountCmd.AddCommand(&cobra.Command{
		Use: "bitget",
		Run: func(cmd *cobra.Command, args []string) { accountbitget.Run() },
	})

	balanceCmd := &cobra.Command{Use: "balance"}
	balanceCmd.AddCommand(&cobra.Command{
		Use: "bitget",
		Run: func(cmd *cobra.Command, args []string) { balancebitget.Run() },
	})
	balanceCmd.AddCommand(&cobra.Command{
		Use: "nobitex",
		Run: func(cmd *cobra.Command, args []string) { balancenobitex.Run() },
	})

	orderCmd := &cobra.Command{Use: "createorder"}
	orderCmd.AddCommand(&cobra.Command{
		Use: "bitget",
		Run: func(cmd *cobra.Command, args []string) { orderbitget.Run() },
	})

	authCmd := &cobra.Command{Use: "auth"}
	authCmd.AddCommand(&cobra.Command{
		Use: "bitpin",
		Run: func(cmd *cobra.Command, args []string) { authbitpin.Run() },
	})
	authCmd.AddCommand(&cobra.Command{
		Use: "nobitex",
		Run: func(cmd *cobra.Command, args []string) { authnobitex.Run() },
	})

	root.AddCommand(accountCmd, balanceCmd, orderCmd, authCmd)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
