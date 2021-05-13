package cli

import (
	_ "path/filepath"

	"github.com/deep2chain/sscq/accounts/keystore"
	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"
	_ "github.com/tendermint/tmlibs/cli"
)

func GetUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update password",
		Long:  "update an account password",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())
			address := args[0]
			oldpass := args[1]
			newpass := args[2]
			err := ksw.Update(address, oldpass, newpass)
			if err != nil {
				return err
			}
			println("Update password success !")
			return nil
		},
	}
}
