package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/bitly/go-simplejson"
	"gopkg.in/urfave/cli.v1"

	cmn "github.com/dappledger/AnnChain/cmd/client/commands"
	"github.com/dappledger/AnnChain/eth/accounts/abi"
	"github.com/dappledger/AnnChain/eth/common"
	"github.com/dappledger/AnnChain/eth/crypto"
	"github.com/dappledger/AnnChain/eth/params"
	. "github.com/dappledger/ann-contracts/evm/client"
	. "github.com/dappledger/ann-contracts/evm/vm"
)

const EVMGasLimit uint64 = 100000000

var (
	evmConfig = Config{EVMGasLimit: EVMGasLimit}
	evm       = NewEVM(Context{}, nil, params.MainnetChainConfig, evmConfig)
)

var (
	EVMCommands = cli.Command{
		Name:     "contract",
		Usage:    "operations for evm contract",
		Subcommands: []cli.Command{
			{
				Name:   "create",
				Usage:  "create a new contract",
				Action: CreateContract,
				Flags: []cli.Flag{
					toolFlags.callf,
					toolFlags.abif,
					toolFlags.nonce,

				},
			}, {
				Name:   "call",
				Usage:  "execute a new contract",
				Action: CallContract,
				Flags: []cli.Flag{
					toolFlags.callf,
					toolFlags.abif,
				},
			},
		},
	}
)

func CreateContract(ctx *cli.Context) error {
	nonce := ctx.Uint64("nonce")

	_, addr, _, _, err := createContract(ctx)
	if err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	contractAddr := crypto.CreateAddress(addr, nonce)
	fmt.Println("contract address:", contractAddr.String())

	return nil
}

func createContract(ctx *cli.Context) (ret []byte, addr common.Address, jsonParams *simplejson.Json, abiJson abi.ABI, err error) {
	jsonParams, err = getParamsJSON(ctx)
	if err != nil {
		return nil, common.Address{}, nil, abi.ABI{}, cli.NewExitError(err.Error(), 127)
	}

	addr = StringToAddr(jsonParams.Get("address").MustString())
	sender := AccountRef(addr)

	abiDefinition, err := fileData(ctx.String("abif"))
	if err != nil {
		return nil, common.Address{}, nil, abi.ABI{}, cli.NewExitError(err.Error(), 127)
	}

	abiJson, err = abi.JSON(strings.NewReader(string(abiDefinition)))
	if err != nil {
		return nil, common.Address{}, nil, abi.ABI{}, cli.NewExitError(err.Error(), 127)
	}

	byteCode := jsonParams.Get("bytecode").MustString()
	if len(byteCode) == 0 {
		return nil, common.Address{}, nil, abi.ABI{}, cli.NewExitError("please give me the contract's bytecode", 127)
	}

	initParam := jsonParams.Get("initParam").MustArray()
	data, err := CreateContractData(abiJson, byteCode, initParam)

	ret, _, _, err = evm.Create(sender, data, 0, big.NewInt(0))
	if err != nil {
		return nil, common.Address{}, nil, abi.ABI{}, cli.NewExitError(err.Error(), 127)
	}

	return ret, addr, jsonParams, abiJson, nil
}

func CallContract(ctx *cli.Context) error {
	ret, addr, jsonParams, abiJson, err := createContract(ctx)
	if err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	method := jsonParams.Get("method").MustString()
	inputParams := jsonParams.Get("params").MustArray()

	input, err := CallContractData(abiJson, method, inputParams)
	if err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	sender := AccountRef(addr)
	res, _, err := evm.Call(sender, common.Address{}, input, ret, 0, big.NewInt(0))
	if err != nil {
		return cli.NewExitError(err.Error(), 127)
	}

	parseResult, err := cmn.UnpackResult(method, abiJson, string(res))
	responseJSON, err := json.Marshal(parseResult)
	fmt.Println("result:", string(responseJSON))

	return nil
}

func getParamsJSON(ctx *cli.Context) (*simplejson.Json, error) {
	var calljson string
	if ctx.String("callf") != "" {
		dat, err := fileData(ctx.String("callf"))
		if err != nil {
			return nil, err
		}
		calljson = string(dat)
		return simplejson.NewJson([]byte(calljson))
	} else {
		return nil, fmt.Errorf("callf is nil")
	}
}

func fileData(str string) ([]byte, error) {
	_, err := os.Stat(str)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(str)
}
