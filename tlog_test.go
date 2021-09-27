package main

import (
	"io"
	"os"
	"testing"
)

var (
	tests = []struct {
		Filename      string
		Reader        io.Reader
		Revisions     int
		Files         int
		CheckRevision string
		Who           string
		Skip          bool
	}{
		{
			Filename:      "tests/multipart1",
			Revisions:     55,
			Files:         3,
			CheckRevision: "e5b622e57b2b08f20828ead059db36338c6a54b0",
			Who:           "gregkh",
			Skip:          true,
		},
		{
			Filename:      "tests/multipart2",
			Revisions:     2982,
			Files:         5,
			CheckRevision: "14a832498c23cf480243222189066a8006182b9d",
			Who:           "sfr",
			Skip:          false,
		},
	}
)

func TestRevisions(t *testing.T) {
	for _, test := range tests {
		if test.Skip {
			continue
		}
		f, _ := os.Open(test.Filename)
		files, _ := ParseMail(f)
		if len(files) != test.Files {
			t.Fatalf("%s: Expected %d files, got %d", test.Filename, test.Files, len(files))
		}
		w, revs := WorkRevisions(files)
		if len(revs) != test.Revisions {
			t.Fatalf("%s: Expected %d revisions, got %d", test.Filename, test.Revisions, len(revs))
		}
		if w != test.Who {
			t.Fatalf("%s: Expected %s, got %s", test.Filename, test.Who, w)
		}
	}
}

func TestWorkDB(t *testing.T) {
	os.Remove("test.db")
	InitDB("test.db")
	count := int64(0)
	for _, test := range tests {
		if test.Skip {
			continue
		}
		f, _ := os.Open(test.Filename)
		files, err := ParseMail(f)
		if err != nil {
			t.Fatal(err)
		}
		tlogEntry := WorkTLog(test.Filename, files)
		AddCommit(tlogEntry)
		err = db.Table("revisions").Where("\"commit\" = ?", test.CheckRevision).Count(&count).Error
		if err != nil {
			t.Fatal(err)
		}
		if count == 0 {
			t.Fatalf("didn't find revision after insertion")
		}
	}
}
