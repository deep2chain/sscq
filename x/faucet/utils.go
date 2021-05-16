package faucet

import (
	"encoding/hex"

	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/deep2chain/sscq/x/auth"
	"github.com/tendermint/go-amino"

	"strings"

	"github.com/deep2chain/sscq/codec"
	newevmtypes "github.com/deep2chain/sscq/evm/types"
	"github.com/deep2chain/sscq/server"
	sdk "github.com/deep2chain/sscq/types"
	"github.com/deep2chain/sscq/x/params"
	"github.com/spf13/viper"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"

	"github.com/deep2chain/sscq/x/crisis"
	distr "github.com/deep2chain/sscq/x/distribution"
	"github.com/deep2chain/sscq/x/gov"
	"github.com/deep2chain/sscq/x/guardian"
	"github.com/deep2chain/sscq/x/mint"
	"github.com/deep2chain/sscq/x/service"
	"github.com/deep2chain/sscq/x/slashing"
	stake "github.com/deep2chain/sscq/x/staking"
	"github.com/deep2chain/sscq/x/upgrade"
)

//
func Decode_Hex(str string) ([]byte, error) {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		//panic(fmt.Sprintf("invalid hex string: %q", str))
		return nil, err
	}
	return b, nil
}

//
func Encode_Hex(str []byte) string {
	return hex.EncodeToString(str)
}

// Read and decode a StdTx from rawdata
func ReadStdTxFromRawData(cdc *amino.Codec, str string) (stdTx auth.StdTx, err error) {
	bytes, err := Decode_Hex(str)
	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return stdTx, err
	}
	return stdTx, err
}

// Read and decode a StdTx from rawdata
func ReadStdTxFromString(cdc *amino.Codec, str string) (stdTx auth.StdTx, err error) {
	bytes := []byte(str)
	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return stdTx, err
	}
	return stdTx, err
}

/////////////////////////////////////////////////////////////////
// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`

	// vesting account fields
	OriginalVesting  sdk.Coins `json:"original_vesting"`  // total vesting coins upon initialization
	DelegatedFree    sdk.Coins `json:"delegated_free"`    // delegated vested coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        int64     `json:"start_time"`        // vesting start time (UNIX Epoch time)
	EndTime          int64     `json:"end_time"`          // vesting end time (UNIX Epoch time)
}

type GenesisState struct {
	Accounts     []GenesisAccount      `json:"accounts"`
	AuthData     auth.GenesisState     `json:"auth"`
	StakeData    stake.GenesisState    `json:"staking"`
	MintData     mint.GenesisState     `json:"mint"`
	DistrData    distr.GenesisState    `json:"distr"`
	GovData      gov.GenesisState      `json:"gov"`
	UpgradeData  upgrade.GenesisState  `json:"upgrade"`
	CrisisData   crisis.GenesisState   `json:"crisis"`
	SlashingData slashing.GenesisState `json:"slashing"`
	ServiceData  service.GenesisState  `json:"service"`
	GuardianData guardian.GenesisState `json:"guardian"`
	GenTxs       []json.RawMessage     `json:"gentxs"`
}

func MakeLatestCodec() *codec.Codec {
	var cdc = codec.New()
	newevmtypes.RegisterCodec(cdc)
	RegisterCodec(cdc)
	params.RegisterCodec(cdc) // only used by querier
	mint.RegisterCodec(cdc)   // only used by querier
	// bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	upgrade.RegisterCodec(cdc)
	service.RegisterCodec(cdc)
	guardian.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	crisis.RegisterCodec(cdc)
	// RegisterCodec(cdc)
	return cdc
}

// // junying-todo-20190429
// // Unlock unlocks an account when address and password are given.
// func Unlock(encrypted, passphrase string) (tmcrypto.PrivKey, error) {
// 	account := accounts.Account{Address: encrypted}
// 	privkey, err := keystore.GetPrivKey(account, passphrase, "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return privkey, nil
// }

// // junying-todo-20190429
// // UnlockByStdIn needs user to type password when bechaddr is given.
// func UnlockByStdIn(bech32 string) (tmcrypto.PrivKey, error) {
// 	passphrase, err := keys.ReadPassphraseFromStdin(bech32)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return Unlock(bech32, passphrase)
// }

// junying-todo-20190429
// This election algorithm is based on the hyposis that all the genesis accounts are owned by
// htdf development team.
// Now it changes every day.
func ElectDefaultSystemIssuer(genstate GenesisState) sdk.AccAddress {
	accounts := genstate.Accounts
	length := len(accounts)
	if length == 0 {
		return nil
	}
	return accounts[0].Address //accounts[time.Now().Day()%length].Address
}

// junying-todo-20190429
// This function finds genesis.json in root directory,
// the read genesis accounts
// return current system issuer account
func GetSystemIssuerFromRoot() (sdk.AccAddress, error) {
	ctx := server.NewDefaultContext()
	config := ctx.Config
	config.SetRoot(viper.GetString(tmcli.HomeFlag))
	cdc := MakeLatestCodec()

	genDoc, err := LoadGenesisDoc(cdc, config.GenesisFile())
	if err != nil {
		return nil, err
	}

	genesisState := GenesisState{}
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return nil, err
	}

	systemissuer := ElectDefaultSystemIssuer(genesisState)
	return systemissuer, nil
}

// junying-todo-20190429
// This function finds genesis.json where user indicates,
// the read genesis accounts
// return current system issuer account
func GetSystemIssuerFromFile(genfilepath string) (sdk.AccAddress, error) {
	cdc := MakeLatestCodec()
	genDoc, err := LoadGenesisDoc(cdc, genfilepath)
	if err != nil {
		return nil, err
	}

	genesisState := GenesisState{}
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return nil, err
	}
	systemissuer := ElectDefaultSystemIssuer(genesisState)
	return systemissuer, nil
}

// LoadGenesisDoc reads and unmarshals GenesisDoc from the given file.
func LoadGenesisDoc(cdc *amino.Codec, genFile string) (genDoc types.GenesisDoc, err error) {
	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}

//
func WriteString(filepath string, msg string) error {
	err := ioutil.WriteFile(filepath, []byte(msg), 0644)
	return err
}

//
func ReadString(filepath string, lineNum int) (line string, lastLine int, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return line, lastLine, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

// junying-todo, 2020-01-15
// this is used to export accounts from accounts.text into genesis.json
// in: accounts text file path
// out: accounts, balances
func ReadAccounts(filepath string) (accounts []string, balances []int, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	linenum := 0
	for sc.Scan() {
		linenum++
		line := sc.Text()
		items := strings.Split(line, "	")
		// fmt.Print(items[0], ",", items[1], "\n")
		accounts = append(accounts, items[0])
		balance, err := strconv.Atoi(items[1])
		if err != nil {
			return accounts, balances, err
		}
		// fmt.Print(balance, "\n")
		balances = append(balances, balance)
	}
	return accounts, balances, sc.Err()
}
