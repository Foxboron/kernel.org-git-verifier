package main

import (
	"fmt"
	"io"
	"log"
	"time"
)

var (
	Signed   = 0
	Unsigned = 0
)

func main() {
	keyring, err := GetKeyring()
	if err != nil {
		log.Fatal(err)
	}

	r, err := GetRepo()
	if err != nil {
		log.Fatal(err)
	}

	toRound := time.Now()
	rounded := time.Date(toRound.Year(), 4, 29, 0, 0, 0, 0, toRound.Location())
	r.LogSince(&rounded,
		func(r io.Reader) error {
			files, err := ParseMail(r)
			if err != nil {
				log.Fatal(err)
			}
			if content, ok := files["git-push-certificate.txt"]; ok {
				p, err := DecodePushCertbuf(content)
				if err != nil {
					log.Fatal(err)
				}
				_, err = p.Verify(keyring)
				if err != nil {
					fmt.Printf("Invalid signature from %s\n", p.Pusher)
					return nil
				}
				Signed += 1
			} else {
				Unsigned += 1
			}
			return nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of certs: %d\nSigned: %d\nUnsigned: %d\n", (Signed + Unsigned), Signed, Unsigned)
}
