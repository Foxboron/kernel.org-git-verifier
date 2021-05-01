package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
)

var data = `certificate version 0.1
pusher B6C41CE35664996C! 1619726851 -0400
pushee gitolite.kernel.org:pub/scm/linux/kernel/git/mricon/patch-attestation-poc
nonce 1619726851-58bb2e188376b269793ea0afac98f8d894dbab80

93628b9ad601b9b39c9c3d6829f87fb7b6eacc64 b0e98279568c2339d5391716744cfb7bcd46b31d refs/heads/master
`

var sig = `-----BEGIN PGP SIGNATURE-----

iHUEABYIAB0WIQR2vl2yUnHhSB5njDW2xBzjVmSZbAUCYIsSAwAKCRC2xBzjVmSZ
bF0yAP4wlxeYaBEblhI2tp6QK7cu7fWEHGHQpJ3SonwrjwORnwEAhHAo8zZtZpCx
CZO7Wu5HCY1MrVqJB/L2OhgfOT0Siw0=
=8iN1
-----END PGP SIGNATURE-----`

func maintest() {
	file := strings.NewReader(data)
	s := strings.NewReader(sig)
	keyring, err := GetKeyring()
	if err != nil {
		log.Fatal(err)
	}
	ent, err := openpgp.CheckArmoredDetachedSignature(keyring, file, s, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ent)
}
