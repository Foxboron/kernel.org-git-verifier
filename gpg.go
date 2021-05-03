package main

import (
	"os"
	"path/filepath"

	"github.com/ProtonMail/go-crypto/openpgp"
)

var Keyring openpgp.EntityList

// Get a single keyring from file
func GetKey(path string) (openpgp.EntityList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	ents, err := openpgp.ReadArmoredKeyRing(f)
	if err != nil {
		return nil, err
	}
	return ents, err
}

// This walks over a keyring directory and fetches our known kernel developer keys
func GetKeyring() (openpgp.EntityList, error) {
	var ent []*openpgp.Entity
	if err := filepath.Walk("pgpkeys/keys", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		ents, err := GetKey(path)
		if err != nil {
			return err
		}
		ent = append(ent, ents...)
		return nil
	}); err != nil {
		return ent, err
	}
	return ent, nil
}

func init() {
	Keyring, _ = GetKeyring()
}
