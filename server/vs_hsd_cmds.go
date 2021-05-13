package server

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/deep2chain/sscq/params"
)

const (
	flagLong         = "long"
	clientIdentifier = "ssd"
)

var (
	// VersionCmd prints out the current sdk version
	VersionHsdCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the ssd version",
		RunE: func(_ *cobra.Command, _ []string) error {
			return versionHsd()
		},
	}
)

func init() {
	VersionHsdCmd.Flags().Bool(flagLong, false, "Print long version information")
}

func versionHsd() error {
	fmt.Println(strings.Title(clientIdentifier))
	fmt.Println("Version:", params.VersionWithMeta)
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}
