module github.com/dappledger/ann-contracts

go 1.12

require (
	github.com/bitly/go-simplejson v0.0.0-20170206154632-da1a8928f709
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3 // indirect
	github.com/dappledger/AnnChain v1.4.1
	github.com/ethereum/go-ethereum v1.8.27
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace (
	github.com/dappledger/AnnChain => ../AnnChain
	github.com/dappledger/ann-go-sdk => ../ann-go-sdk
)
