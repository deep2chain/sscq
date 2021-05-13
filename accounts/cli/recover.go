package cli

import (
	"fmt"

	"github.com/deep2chain/sscq/accounts/keystore"
	"github.com/spf13/cobra"
)

//recover a account from a private key
func GetRecoverAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recover [private-key] [password]",
		Short: "recover a account from a private key.",
		Long:  "recover a account from a private key. when success,you can find the recovered account keyfile in keystore",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			//var acc accounts.Account
			//var encryptPassword string
			var err error
			strPrivateKey := args[0]
			strPasswd := args[1]

			ks := keystore.NewKeyStore(keystore.DefaultKeyStoreHome())

			err = ks.RecoverKey(strPrivateKey, strPasswd)
			if err != nil {
				fmt.Printf("RecoverAccount error|err=%s\n", err)
				return err
			}

			fmt.Printf("strPubkey=%v\n", ks.Key().PubKey)

			fmt.Printf("RecoverAccount success|address=%s\n", ks.Key().Address)

			return err
		},
	}
}
