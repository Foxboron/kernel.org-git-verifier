package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// This should be upstreamed to go-git
// Thus parts of the API is mirrored

const (
	beginpgp  string = "-----BEGIN PGP SIGNATURE-----"
	endpgp    string = "-----END PGP SIGNATURE-----"
	headerpgp string = "gpgsig"
)

type PushCert struct {
	Version      []byte
	Pusher       []byte
	Timestamp    time.Time
	Pushee       []byte
	Nonce        []byte
	Protocol     []byte
	Content      string
	GPGSignature []byte
}

func DecodePushCertbuf(buf []byte) (*PushCert, error) {
	pc := new(PushCert)
	pc.Content = string(buf)
	byteBuffer := bytes.NewBuffer(buf)
	r := bufio.NewReader(byteBuffer)
	isGPG := false
	var sig [][]byte
	for {
		var line []byte
		line, err := r.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		line = bytes.TrimSpace(line)
		if bytes.Equal(line, []byte(beginpgp)) {
			// start parsing GPG signature until we see the end
			isGPG = true
		}
		if bytes.Equal(line, []byte(endpgp)) {
			isGPG = false
			sig = append(sig, line)
			pc.GPGSignature = bytes.Join(sig, []byte("\n"))
		}
		if isGPG {
			sig = append(sig, line)
			continue
		}

		split := bytes.SplitN(line, []byte{' '}, 2)
		switch string(split[0]) {
		case "certificate":
			version := bytes.SplitN(split[1], []byte{' '}, 2)
			pc.Version = version[1]
		case "pusher":
			pusher := bytes.Split(split[1], []byte{' '})
			pc.Pusher = bytes.Join(pusher[:len(pusher)-2], []byte(" "))
			unix, err := strconv.ParseInt(string(pusher[len(pusher)-2]), 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			t := time.Unix(unix, 0)
			// Super hacky but we need to retain the timezone when encoding the cert
			tz := string(pusher[len(pusher)-1])
			test, err := time.Parse("-0700", tz)
			if err != nil {
				log.Fatal(err)
			}
			t = t.In(test.Location())
			pc.Timestamp = t
		case "pushee":
			pc.Pushee = split[1]
		case "nonce":
			pc.Nonce = split[1]
		case "":
			var buf []byte
			for {
				b, err := r.Peek(5)
				if err != nil {
					break
				}
				// Breka if we see a signature
				if bytes.Equal(b, []byte("-----")) {
					break
				}
				line, err = r.ReadBytes('\n')
				if errors.Is(err, io.EOF) {
					break
				}
				buf = append(buf, line...)
			}
			pc.Protocol = buf
		}
	}
	return pc, nil
}

func (p *PushCert) encode(w io.Writer, includeSig bool) error {
	t := fmt.Sprintf("%d %s", p.Timestamp.Unix(), p.Timestamp.Format("-0700"))
	tmpl := "certificate version %s\npusher %s %s\npushee %s\nnonce %s\n\n%s"
	if _, err := fmt.Fprintf(w, tmpl,
		p.Version, p.Pusher, t, p.Pushee, p.Nonce, p.Protocol); err != nil {
		return err
	}
	if includeSig {
		if _, err := fmt.Fprintf(w, "%s\n", p.GPGSignature); err != nil {
			return err
		}
	}
	return nil
}

func (p *PushCert) Encode(w io.Writer) error {
	return p.encode(w, true)
}

func (p *PushCert) EncodeWithoutCert(w io.Writer) error {
	return p.encode(w, false)
}

// Verify the content of the cert towards a pgp keyring
func (p *PushCert) Verify(keyring openpgp.KeyRing) (ent *openpgp.Entity, err error) {
	file := bytes.NewBuffer(nil)
	p.EncodeWithoutCert(file)
	s := bytes.NewReader(p.GPGSignature)
	// Because CLOCKS ARE AMAZING:
	// The signature *might* be in the future, thus we need to add 5 minutes to the date
	// TODO: If signatures fail, check this.
	jitter := 5
	roundedUp := time.Date(p.Timestamp.Year(),
		p.Timestamp.Month(), p.Timestamp.Day(), p.Timestamp.Hour(), p.Timestamp.Minute()+jitter, 0, 0, p.Timestamp.Location())
	// The extra code is incase we need to add jitter to round up and down :/
	// Might not need it!
	// roundedDown := time.Date(p.Timestamp.Year(),
	// 	p.Timestamp.Month(), p.Timestamp.Day(), p.Timestamp.Hour(), p.Timestamp.Minute()-jitter, 0, 0, p.Timestamp.Location())
	// for _, _ = range []time.Time{p.Timestamp, roundedUp, roundedDown} {
	ent, err = openpgp.CheckArmoredDetachedSignature(keyring, file, s, &packet.Config{Time: roundedUp.Local})
	if err != nil {
		return nil, err
	}
	return ent, nil
}
