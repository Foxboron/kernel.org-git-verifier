package main

import (
	"io"
	"log"

	"gopkg.in/yaml.v3"
)

type Change struct {
	Ref string
	Old string
	New string
	Log string
}

type TLogEntry struct {
	Service           string
	Repo              string
	User              string
	GitPushCertStatus string `yaml:"git_push_cert_status"`
	Changes           []Change
}

func DecodeTLogEntry(r io.Reader) (*TLogEntry, error) {
	var tlog TLogEntry
	reader := yaml.NewDecoder(r)
	err := reader.Decode(&tlog)
	if err != nil {
		return nil, err
	}
	return &tlog, nil
}

func (t *TLogEntry) Encode(w io.Writer) error {
	writer := yaml.NewEncoder(w)
	// This doesn't work :c
	// The log blob should have 4 indents, but rest has 2
	writer.SetIndent(2)
	err := writer.Encode(t)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
