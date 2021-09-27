package main

import (
	"strings"
	"testing"
)

const EntryTest = `---
service: git-receive-pack
repo: pub/scm/linux/kernel/git/mricon/patch-attestation-poc
user: mricon
git_push_cert_status: G
changes:
  - ref: refs/heads/master
    old: 93628b9ad601b9b39c9c3d6829f87fb7b6eacc64
    new: b0e98279568c2339d5391716744cfb7bcd46b31d
    log: |
         b0e98279568c2339d5391716744cfb7bcd46b31d Remove some redundant code
         
`

func TestDecodeEntry(t *testing.T) {
	r := strings.NewReader(EntryTest)
	_, err := DecodeTLogEntry(r)
	if err != nil {
		t.Fatal("failed to decode")
	}
}

func TestEncodeEntry(t *testing.T) {
	r := strings.NewReader(EntryTest)
	tlog, _ := DecodeTLogEntry(r)
	var s strings.Builder
	err := tlog.Encode(&s)
	if err != nil {
		t.Fatalf("Failed to encode")
	}
	// fmt.Println(s.String())
	// if EntryTest != s.String() {
	// 	t.Fatalf("failed to encode string")
	// }
}
