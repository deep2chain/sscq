package cli

import (
	"github.com/deep2chain/sscq/accounts/keystore"
	"github.com/spf13/cobra"
)

func GetDelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "del",
		Short: "delete account",
		Long:  "delete account from keystores",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())
			address := args[0]
			err := ksw.Drop(address)
			if err != nil {
				return err
			}

			println("Delete success: ", address)
			return nil
		},
	}
}
