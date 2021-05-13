package cli

import (
	"github.com/deep2chain/sscq/accounts/keystore"
	"github.com/deep2chain/sscq/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"fmt"
)

const FlagPublicKey string = "pubkey"

func GetNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "create a new account.",
		Long:  "new_Account: If OFF to use standard client tool ,else use common stdin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//var acc accounts.Account
			var encryptPassword string
			var err error
			if args[0] == "OFF" {

				buf := client.BufferStdin()
				//interactive := viper.GetBool(flagInteractive)
				encryptPassword = "12345678"
				if viper.GetString(FlagPublicKey) == "" && !viper.GetBool(client.FlagUseLedger) {
					encryptPassword, err = client.GetCheckPassword(
						"Enter a passphrase to encrypt your key to disk:",
						"Repeat the passphrase:", buf)
					if err != nil {
						return err
					}
				}

			} else {
				encryptPassword = args[0]
			}

			var ks = new(keystore.KeyStore)
			ks = keystore.NewKeyStore(keystore.DefaultKeyStoreHome())
			_, err = ks.NewKey(encryptPassword)
			if err != nil {
				return err
			}

			// println("Create new account successful!")
			// println("Address: ",ks.Key().Address)
			fmt.Println(ks.Key().Address)

			return err
		},
	}
}
