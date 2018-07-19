package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/examples/democoin/x/cool"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
)

// QuizTxCmd invokes the coolness quiz transaction.
func QuizTxCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cool [answer]",
		Short: "What's cooler than being cool?",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			queryCtx := context.NewQueryContextFromCLI().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			from, err := queryCtx.GetFromAddress()
			if err != nil {
				return err
			}

			account, err := queryCtx.GetAccount(from)
			if err != nil {
				return fmt.Errorf(`Failed to find our decode account with address: %s.
Are you sure there has been a transaction involving it?`, from)
			}

			msg := cool.NewMsgQuiz(from, args[0])
			name := viper.GetString(client.FlagName)

			utils.SendTx(txCtx, queryCtx, account, name, []sdk.Msg{msg})
			return nil
		},
	}
}

// SetTrendTxCmd sends a new cool trend transaction.
func SetTrendTxCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "setcool [answer]",
		Short: "You're so cool, tell us what is cool!",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			queryCtx := context.NewQueryContextFromCLI().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from address from the name flag
			from, err := queryCtx.GetFromAddress()
			if err != nil {
				return err
			}

			name := viper.GetString(client.FlagName)
			msg := cool.NewMsgSetTrend(from, args[0])

			utils.SendTx(txCtx, queryCtx, from, name, []sdk.Msg{msg})
			return nil
		},
	}
}
