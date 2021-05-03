package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	pgperrors "github.com/ProtonMail/go-crypto/openpgp/errors"
)

func VerifyPushCert(content []byte) (string, error) {
	p, err := DecodePushCertbuf(content)
	if err != nil {
		return "", err
	}
	ent, err := p.Verify(Keyring)
	if err != nil {
		return "", err
	}
	for s := range ent.Identities {
		return s, nil
	}
	return "", nil
}

// TODO: Generalize the revlist parsing
// Since the revlist is either in a file or part of the
// tlog entry we need figure out how to nicely abstract it
func RevlistParser(revlist io.Reader) [][]byte {
	r := bufio.NewReader(revlist)
	var ret [][]byte
	for {
		line, err := r.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil
		}
		split := bytes.SplitN(line, []byte{' '}, 2)
		ret = append(ret, split[0])
	}
	return ret
}

func WorkRevisions(files MimeFiles) (string, []Revision) {
	var wc []Revision
	entry, _ := DecodeTLogEntry(bytes.NewReader(files[0].Content))
	var r io.Reader
	for _, change := range entry.Changes {
		// If there is a revision list we find that instead
		if strings.HasPrefix(change.Log, "revlist") {
			for _, f := range files {
				if f.Filename == change.Log {
					r = bytes.NewBuffer(f.Content)
				}
			}
		} else {
			r = strings.NewReader(change.Log)
		}
	}
	revlist := RevlistParser(r)
	for _, rev := range revlist {
		wc = append(wc, Revision{Who: entry.User, Repository: entry.Repo, Commit: fmt.Sprintf("%s", rev)})
	}
	return entry.User, wc
}

func WorkTLog(h string, files MimeFiles) *TLogCommit {
	var tlog TLogCommit
	tlog.Commit = h
	user, wc := WorkRevisions(files)
	tlog.User = user
	tlog.Revisions = wc
	// Some commits are just plain single file entries
	if len(files) < 2 {
		return &tlog
	}
	if files[1].Filename != "git-push-certificate.txt" {
		return &tlog
	}
	tlog.Signature = true
	cert := files[1].Content
	issuer, err := VerifyPushCert(cert)
	if errors.Is(err, pgperrors.ErrUnknownIssuer) {
		tlog.Unknown = true
	} else if err == nil {
		tlog.Valid = true
	}
	tlog.SigIssuer = issuer
	return &tlog
}
