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
	flagsLong             = "long"
	hscliClientIdentifier = "hscli"
)

var (
	// VersionCmd prints out the current sdk version
	VersionHscliCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the hscli version",
		RunE: func(_ *cobra.Command, _ []string) error {
			return versionHscli()
		},
	}
)

func init() {
	VersionHscliCmd.Flags().Bool(flagsLong, false, "Print long version information")
}

func versionHscli() error {
	fmt.Println(strings.Title(hscliClientIdentifier))
	fmt.Println("Version:", params.VersionWithMeta)
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}
