package keystore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/deep2chain/sscq/accounts"
)

const KeyStoreScheme = "keystores"

type accountsByURL []accounts.Account

type scaner struct {
	keydir string
	mu     sync.Mutex
	set    mapset.Set
	all    accountsByURL
	byAddr map[string][]accounts.Account
}

func newScaner(keydir string) *scaner {
	sc := &scaner{
		keydir: keydir,
		set:    mapset.NewThreadUnsafeSet(),
		byAddr: make(map[string][]accounts.Account),
	}

	return sc
}

func (sc *scaner) getSigner(addr string) (*Key, error) {
	found := sc.hasAddress(addr)
	if found {
		account := accounts.Account{Address: addr}
		acc, err := sc.find(account)
		if err != nil {
			return nil, err
		}

		key, err := getKey(addr, acc.URL.Path)
		if err != nil {
			return nil, err
		}
		return key, err
	}

	return nil, ErrNoMatch
}

func (sc *scaner) find(a accounts.Account) (accounts.Account, error) {
	// Limit search to address candidates if possible.
	matches := sc.all
	if a.Address != "" {
		matches = sc.byAddr[a.Address]
	}
	if a.URL.Path != "" {
		// If only the basename is specified, complete the path.
		if !strings.ContainsRune(a.URL.Path, filepath.Separator) {
			a.URL.Path = filepath.Join(sc.keydir, a.URL.Path)
		}
		for i := range matches {
			if matches[i].URL == a.URL {
				return matches[i], nil
			}
		}
		if a.Address == "" {
			return accounts.Account{}, ErrNoMatch
		}
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return accounts.Account{}, ErrNoMatch
	default:
		return accounts.Account{}, nil
	}
}

func (sc *scaner) accounts() ([]accounts.Account, error) {
	err := sc.scanAccounts()
	if err != nil {
		return nil, err
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()
	cpy := make([]accounts.Account, len(sc.all))
	copy(cpy, sc.all)

	return cpy, err
}

func (sc *scaner) scanAccounts() error {
	// Create a helper method to scan the contents of the key files
	var (
		buf = new(bufio.Reader)
		key struct {
			Address string `json:"address"`
		}
	)

	// Scan the entire folder metadata for file changes
	creates, deletes, updates, err := sc.scan(sc.keydir)
	if err != nil {
		fmt.Print("Failed to reload keystore contents err: ", err, "\n")
		return err
	}
	if creates.Cardinality() == 0 && deletes.Cardinality() == 0 && updates.Cardinality() == 0 {
		return nil
	}

	readAccount := func(path string) *accounts.Account {
		fd, err := os.Open(path)
		if err != nil {
			fmt.Print("Failed to open keystore file", "path", path, "err", err)
			return nil
		}
		defer fd.Close()
		buf.Reset(fd)
		// Parse the address.
		key.Address = ""
		err = json.NewDecoder(buf).Decode(&key)
		addr := key.Address
		switch {
		case err != nil:
			fmt.Print("Failed to decode keystore key", "path", path, "err", err)
		case (addr == ""):
			fmt.Print("Failed to decode keystore key", "path", path, "err", "missing or zero address")
		default:
			return &accounts.Account{Address: addr, URL: accounts.URL{Scheme: KeyStoreScheme, Path: path}}
		}
		return nil
	}

	for _, p := range creates.ToSlice() {
		if a := readAccount(p.(string)); a != nil {
			sc.add(*a)
		}
	}
	return err
}

func (sc *scaner) scan(keyDir string) (mapset.Set, mapset.Set, mapset.Set, error) {
	// List all the failes from the keystore folder
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return nil, nil, nil, err
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Iterate all the files and gather their metadata
	all := mapset.NewThreadUnsafeSet()
	mods := mapset.NewThreadUnsafeSet()

	for _, fi := range files {
		path := filepath.Join(keyDir, fi.Name())
		// Skip any non-key files from the folder
		if nonKeyFile(fi) {
			//log.Trace("Ignoring file on account scan", "path", path)
			continue
		}
		// Gather the set of all and fresly modified files
		all.Add(path)

	}

	// Update the tracked files and return the three sets
	deletes := sc.set.Difference(all)   // Deletes = previous - current
	creates := all.Difference(sc.set)   // Creates = current - previous
	updates := mods.Difference(creates) // Updates = modified - creates

	sc.set = all

	// Report on the scanning stats and return
	return creates, deletes, updates, nil
}

// nonKeyFile ignores editor backups, hidden files and folders/symlinks.
func nonKeyFile(fi os.FileInfo) bool {
	// Skip editor backups and UNIX-style hidden files.
	if strings.HasSuffix(fi.Name(), "~") || strings.HasPrefix(fi.Name(), ".") {
		return true
	}
	// Skip misc special files, directories (yes, symlinks too).
	if fi.IsDir() || fi.Mode()&os.ModeType != 0 {
		return true
	}
	return false
}

func (sc *scaner) add(newAccount accounts.Account) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	i := sort.Search(len(sc.all), func(i int) bool { return sc.all[i].URL.Cmp(newAccount.URL) >= 0 })
	if i < len(sc.all) && sc.all[i] == newAccount {
		return
	}

	// newAccount is not in the cache.
	sc.all = append(sc.all, accounts.Account{})
	copy(sc.all[i+1:], sc.all[i:])
	sc.all[i] = newAccount
	sc.byAddr[newAccount.Address] = append(sc.byAddr[newAccount.Address], newAccount)

}

func (sc *scaner) hasAddress(addr string) bool {
	err := sc.scanAccounts()
	if err == nil {
		sc.mu.Lock()
		defer sc.mu.Unlock()
		return len(sc.byAddr[addr]) > 0
	}
	return false
}
