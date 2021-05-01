package main

import (
	"bytes"
	"testing"
)

const TestString = `certificate version 0.1
pusher B6C41CE35664996C! 1619726851 -0400
pushee gitolite.kernel.org:pub/scm/linux/kernel/git/mricon/patch-attestation-poc
nonce 1619726851-58bb2e188376b269793ea0afac98f8d894dbab80

93628b9ad601b9b39c9c3d6829f87fb7b6eacc64 b0e98279568c2339d5391716744cfb7bcd46b31d refs/heads/master
-----BEGIN PGP SIGNATURE-----

iHUEABYIAB0WIQR2vl2yUnHhSB5njDW2xBzjVmSZbAUCYIsSAwAKCRC2xBzjVmSZ
bF0yAP4wlxeYaBEblhI2tp6QK7cu7fWEHGHQpJ3SonwrjwORnwEAhHAo8zZtZpCx
CZO7Wu5HCY1MrVqJB/L2OhgfOT0Siw0=
=8iN1
-----END PGP SIGNATURE-----
`

func TestParsePushCert(t *testing.T) {
	_, err := DecodePushCertbuf([]byte(TestString))
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncodePushCert(t *testing.T) {
	p, _ := DecodePushCertbuf([]byte(TestString))
	b := bytes.NewBuffer(nil)
	p.Encode(b)
	if !bytes.Equal(b.Bytes(), []byte(TestString)) {
		t.Fatal("Failed to encode push-cert")
	}
}

func TestVerifyPushCert(t *testing.T) {
	p, _ := DecodePushCertbuf([]byte(TestString))
	keyring, _ := GetKey("keyring/mricon.gpg")
	_, err := p.Verify(keyring)
	if err != nil {
		t.Fatalf("Can't verify content: %s", err)
	}
}
