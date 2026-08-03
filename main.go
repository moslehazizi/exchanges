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
	orderbitgetcancel "order-maker/order/bitget/cancel"
	orderbitgethistory "order-maker/order/bitget/history"
	orderbitgetlist "order-maker/order/bitget/list"
	orderbitgetplace "order-maker/order/bitget/place"
	websocketbitgetaccount "order-maker/web-socket/bitget/account"
	websocketbitgetorder "order-maker/web-socket/bitget/order"

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

	orderCmd := &cobra.Command{Use: "order"}
	orderCmd.AddCommand(&cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) { orderbitgetplace.Run() },
	})

	orderCmd.AddCommand(&cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) { orderbitgetlist.Run() },
	})

	orderCmd.AddCommand(&cobra.Command{
		Use: "history",
		Run: func(cmd *cobra.Command, args []string) { orderbitgethistory.Run() },
	})

	orderCmd.AddCommand(&cobra.Command{
		Use: "cancel",
		Run: func(cmd *cobra.Command, args []string) { orderbitgetcancel.Run() },
	})

	websocketCmd := &cobra.Command{Use: "websocket"}
	websocketCmd.AddCommand(&cobra.Command{
		Use: "bitgetaccount",
		Run: func(cmd *cobra.Command, args []string) { websocketbitgetaccount.Run() },
	})
	websocketCmd.AddCommand(&cobra.Command{
		Use: "bitgetorder",
		Run: func(cmd *cobra.Command, args []string) { websocketbitgetorder.Run() },
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

	root.AddCommand(accountCmd, balanceCmd, orderCmd, websocketCmd, authCmd)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
