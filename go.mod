module github.org/foxboron/kernel.org-git-verifier

go 1.16

require (
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/ProtonMail/gopenpgp/v2 v2.1.8 // indirect
	github.com/go-git/go-git/v5 v5.3.0
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/ProtonMail/go-crypto => ../go-crypto
