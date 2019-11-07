package client

import (
	"fmt"
	"strconv"
	"strings"

	cmn "github.com/dappledger/AnnChain/cmd/client/commons"
	"github.com/dappledger/AnnChain/eth/accounts/abi"
	"github.com/dappledger/AnnChain/eth/common"
)

func StringToAddr(addrStr string) common.Address {
	if strings.Index(addrStr, "0x") == 0 {
		addrStr = addrStr[2:]
	}

	addrBytes := common.Hex2Bytes(addrStr)
	addr := common.BytesToAddress(addrBytes)
	return addr
}

func parseData(methodName string, abiDef *abi.ABI, params []interface{}) (string, error) {
	args, err := parseArgs(methodName, abiDef, params)
	if err != nil {
		return "", err
	}
	data, err := abiDef.Pack(methodName, args...)
	if err != nil {
		return "", err
	}

	var hexData string
	for _, b := range data {
		hexDataP := strconv.FormatInt(int64(b), 16)
		if len(hexDataP) == 1 {
			hexDataP = "0" + hexDataP
		}
		hexData += hexDataP
	}
	return hexData, nil
}

func parseArgs(methodName string, abiDef *abi.ABI, params []interface{}) ([]interface{}, error) {
	var method abi.Method
	if methodName == "" {
		method = abiDef.Constructor
	} else {
		var ok bool
		method, ok = abiDef.Methods[methodName]
		if !ok {
			return nil, fmt.Errorf("no such method")
		}
	}

	if params == nil {
		params = []interface{}{}
	}
	if len(params) != len(method.Inputs) {
		return nil, fmt.Errorf("unmatched params %x:%d", params, len(method.Inputs))
	}

	var args []interface{}
	for i := range params {
		a, err := cmn.ParseArg(method.Inputs[i], params[i])
		if err != nil {
			fmt.Println(fmt.Sprintf("fail to parse args %v into %s: %v ", params[i], method.Inputs[i].Name, err))
			return nil, err
		}
		args = append(args, a)
	}
	return args, nil
}

func CreateContractData(abiJson abi.ABI, byteCode string, params []interface{}) ([]byte, error) {
	initParam, err := parseData("", &abiJson, params)
	if err != nil {
		return nil, err
	}

	if initParam != "" && strings.Index(initParam, "0x") == 0 {
		initParam = initParam[2:]
	}

	data := common.Hex2Bytes(byteCode + initParam)

	return data, nil
}

func CallContractData(abiJson abi.ABI, conMethod string, inputData []interface{}) ([]byte, error) {
	args, err := parseArgs(conMethod, &abiJson, inputData)
	if err != nil {
		return nil, err
	}

	callData, err := abiJson.Pack(conMethod, args...)
	if err != nil {
		return nil, err
	}

	return callData, nil
}
