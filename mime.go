package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

type MimeFile struct {
	Filename string
	Content  []byte
}

type MimeFiles []MimeFile

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

func ParsePart(out *MimeFiles, mime_data io.Reader, boundary string, index int) {
	if boundary == "" {
		bytes, _ := io.ReadAll(mime_data)
		*out = append(*out, MimeFile{Filename: "file", Content: bytes})
		return
	}
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
			*out = append(*out, MimeFile{Filename: filename, Content: WritePart(new_part)})
		}
	}
}

func ParseMail(r io.Reader) (MimeFiles, error) {
	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Fatal(err)
	}
	_, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	var files MimeFiles
	ParsePart(&files, m.Body, params["boundary"], 1)
	return files, nil
}
