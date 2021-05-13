package cmd

import (
	"path/filepath"

	liteserver "github.com/deep2chain/sscq/lite/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

func StartLiteNodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lite",
		Short: "Run lite-client proxy server, verifying tendermint rpc",
		Long: `This node will run a secure proxy to a tendermint rpc server.

		All calls that can be tracked back to a block header by a proof
		will be verified before passing them back to the caller. Other that
		that it will present the same interface as a full tendermint node,
		just with added trust and running locally.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootDir := viper.GetString(cli.HomeFlag)
			defaultKeyStoreHome := filepath.Join(rootDir, "litenode-data")
			var err error
			home := defaultKeyStoreHome
			maxOpenConnections := 100
			cacheSize := 10
			srv := liteserver.NewServer(args[0], args[1], args[2], home, maxOpenConnections, cacheSize)
			err = srv.RunProxy()
			return err
		},
	}
}
