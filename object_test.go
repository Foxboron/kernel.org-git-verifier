package main

import (
	"errors"
	"strings"
	"testing"

	pgperrors "github.com/ProtonMail/go-crypto/openpgp/errors"
)

const TestString = `certificate version 0.1
pusher Greg Kroah-Hartman <gregkh@linuxfoundation.org> 1619954742 +0200
pushee gitolite.kernel.org:/pub/scm/linux/kernel/git/stable/stable-queue.git
nonce 1619954741-8aa31f83d2c67834f75b05b6470e744bc2e3fc94

e667cd60cdd4c39c62d30ad0925bcfccaba25e3f fec5db375df45faac71048ea30e3827040668171 refs/heads/master
-----BEGIN PGP SIGNATURE-----

iQJPBAABCAA5FiEEZH8oZUiU471FcZm+ONu9yGCSaT4FAmCOjDcbHGdyZWdraEBs
aW51eGZvdW5kYXRpb24ub3JnAAoJEDjbvchgkmk+2QQQAK0R5G1yaql/XpArL7zf
m4wY9Wk/WzSmJM3iQe8daqu0Zkdy3Z2xev6XXua2bL1XCYwC+DNH4pExg3+fJdYm
raBYjQwsyIbCWCfJ1mwrbmESn7tY76ioYBwyrz+fgqr7WpNlTeBrMETVyMMjoKw8
e6AoC7Iw8E0QSc2IBS17i/c7LynUxh+dUIMNou9HZy7+XYpRq0OuIj27HP7/R0mB
NpDpYkJnu7yi/HmcwLTAPFUt89yemiHw8xyVd94eTKTUdFbJCIv5LKowTAO4etEK
gTi24x7BlSLZJAKQa2Rss+ufs59xtWj7NU1XE0yRPsxHO5bXUgCSVsAglfV8DN/C
1N+rr2nDFsmul8QQm/mxeGv/t1Kq7X8oGXFxpsBBtbOVajNzL6JpMgENwZPuwjNF
fXqLfD0MQdEcQXzAUgNNWXOxBRk88GYGz/zoJZ9EXPuvHzBVP/l5rHfGPMb/UHd2
izg6fkMGSdLJ24H6Cw2e8rX4uztwFVjr0YYYu+i0ILsp6xWoWpjrm8gretT7ZQoq
pXnasWupZ7gBEhM7GFDTb6ktte9mXYiBn85GMOIWE6Me358+EUjbW0Y8gU2k5PR3
GgSUJxsM4B3BJNFHl8JwvzkJ7sPK1N6VamxAtcSxD45rtbiTMZRwWzXzz6TZii7u
ZtZ+nNNEMoiBUpVPG9jOu71r
=qj3o
-----END PGP SIGNATURE-----
`

const TestString2 = `certificate version 0.1
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

// Expired cert. This should validate with a paste date set
const TestString3 = `certificate version 0.1
pusher 0x3D200E9CA6329909 1604298629 +0100
pushee ssh://gitolite.kernel.org/pub/scm/linux/kernel/git/ardb/linux.git
nonce 1604298628-c34b6aff6aa4d416a57c43c14e09b390ce74cb57

0000000000000000000000000000000000000000 3ebfd27aa35e4145a684e269aadf06e4062a6251 refs/heads/arm32-chacha
-----BEGIN PGP SIGNATURE-----

iQEzBAABCgAdFiEEnNKg2mrY9zMBdeK7wjcgfpV0+n0FAl+fp4UACgkQwjcgfpV0
+n3nlQgAn23fpOSdL5Q6CGUuhH5VMD+qpekWTd9fKGZTKMU+IG92/xlb8xpWmo01
brwO0gfa9jXlXZgQLyejROCvdE2QjRDSFT9ReSzJjSkmNfmNOSRmSTPDXDXxih1d
0t1hq7DcXoSnxyLeCejc6epsGXQ9MR92nQOsgopzoYhqFXb2FyXeVRpx75MQw0tR
UkWc5hyEZbbiWJgVYcAOIxIzQF6LeZC5jSH0TTfNOT2rck0BpvT1jI4x4NnPKY1k
1nChRzKjzzKwkoG0B54R78C+gtbah8yRk8fxyP60gDEm6tty9LltBLo7qBmhVtYx
7ncau7k5f4529pdFauwDk9fXuESpHg==
=LlUA
-----END PGP SIGNATURE-----
`

const TestString4 = `certificate version 0.1
pusher BE5675C37E818C8B5764241C254BCFC56BF6CE8D 1606558182 +0100
pushee gitolite.kernel.org:pub/scm/linux/kernel/git/mripard/linux.git
nonce 1606558181-893a31ce1398048e325b289fc3175fffeea72b5d

0000000000000000000000000000000000000000 fdb21b6152b545781b2c76cd290f5e2192be4fc7 refs/tags/sunxi/config-for-5.11-pull-request
-----BEGIN PGP SIGNATURE-----

iHUEABYIAB0WIQRcEzekXsqa64kGDp7j7w1vZxhRxQUCX8Ih5gAKCRDj7w1vZxhR
xXN8AP9r9WPz5A4yOcsNVqYKcy49Dga5bPChm5AT4/qQrMwaZgD+M7u6YcdpI/i1
6V4UEG4AwzE2XJJZ4aZh/ZHO15iAHgc=
=dqSt
-----END PGP SIGNATURE-----
`

var TestStrings = []string{TestString4, TestString, TestString2, TestString3}

func TestParsePushCert(t *testing.T) {
	for _, line := range TestStrings {
		_, err := DecodePushCertbuf([]byte(line))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestEncodePushCert(t *testing.T) {
	for _, line := range TestStrings {
		p, _ := DecodePushCertbuf([]byte(line))
		var s strings.Builder
		p.Encode(&s)
		if s.String() != line {
			t.Fatal("Failed to encode push-cert")
		}
	}
}

func TestVerifyPushCert(t *testing.T) {
	for _, line := range TestStrings {
		p, _ := DecodePushCertbuf([]byte(line))
		var s strings.Builder
		p.Encode(&s)
		keyring, _ := GetKeyring()
		_, err := p.Verify(keyring)
		if errors.Is(err, pgperrors.ErrUnknownIssuer) {
			t.Fatalf("Unknown issuer: %s", err)
		} else if err != nil {
			t.Fatalf("Can't verify content: %s", err)
		}
	}
}
