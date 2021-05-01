package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	StateDir = "state"
	GitDir   = "https://git.kernel.org/pub/scm/infra/transparency-logs/gitolite/git/1.git"
)

type MimeFile struct {
	Filename string
	Content  []byte
}

func BuildFileName(part *multipart.Part, radix string, index int) (filename string) {
	filename = part.FileName()
	if len(filename) > 0 {
		return
	}
	mediaType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
	if err == nil {
		mime_type, e := mime.ExtensionsByType(mediaType)
		if e == nil {
			return fmt.Sprintf("%s-%d%s", radix, index, mime_type[0])
		}
	}
	return
}

// WitePart decodes the data of MIME part and writes it to the file filename.
func WritePart(part *multipart.Part) []byte {
	part_data, err := ioutil.ReadAll(part)
	if err != nil {
		return nil
	}

	content_transfer_encoding := strings.ToUpper(part.Header.Get("Content-Transfer-Encoding"))
	switch {
	case strings.Compare(content_transfer_encoding, "BASE64") == 0:
		if decoded_content, err := base64.StdEncoding.DecodeString(string(part_data)); err != nil {
			return decoded_content
		}

	case strings.Compare(content_transfer_encoding, "QUOTED-PRINTABLE") == 0:
		if decoded_content, err := ioutil.ReadAll(quotedprintable.NewReader(bytes.NewReader(part_data))); err != nil {
			return decoded_content
		}

	default:
		return part_data
	}
	return nil
}

func ParsePart(out *[]MimeFile, mime_data io.Reader, boundary string, index int) {
	reader := multipart.NewReader(mime_data, boundary)
	if reader == nil {
		return
	}

	for {
		new_part, err := reader.NextPart()
		if err == io.EOF || err != nil {
			break
		}

		mediaType, params, err := mime.ParseMediaType(new_part.Header.Get("Content-Type"))
		if err == nil && strings.HasPrefix(mediaType, "multipart/") {
			ParsePart(out, new_part, params["boundary"], index+1)
		} else {
			filename := BuildFileName(new_part, boundary, 1)
			mime := MimeFile{
				Filename: filename,
				Content:  WritePart(new_part)}
			*out = append(*out, mime)
		}
	}
}

func main() {
	dir := path.Join(StateDir, path.Base(GitDir))
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL:      GitDir,
			Progress: os.Stdout,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	r, err := git.PlainOpen(dir)
	if err != nil {
		log.Fatal(err)
	}

	toRound := time.Now()
	rounded := time.Date(toRound.Year(), toRound.Month(), toRound.Day(), 0, 0, 0, 0, toRound.Location())
	cIter, err := r.Log(&git.LogOptions{Since: &rounded})
	if err != nil {
		log.Fatal(err)
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		f, err := c.File("m")
		if err != nil {
			log.Fatal(err)
		}
		blob, err := f.Contents()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(blob)
		r := strings.NewReader(blob)
		m, err := mail.ReadMessage(r)
		if err != nil {
			log.Fatal(err)
		}
		_, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
		if err != nil {
			log.Fatal(err)
		}
		var files []MimeFile
		ParsePart(&files, m.Body, params["boundary"], 1)
		fmt.Println(string(files[1].Content))
		os.Exit(1)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
