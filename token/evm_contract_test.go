package evm

import (
	"fmt"
	"github.com/dappledger/AnnChain/eth/common"
	"github.com/dappledger/AnnChain/eth/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

const senderPirvKey = "48deaa73f328f38d5fcb29d076b2b639c8491f97d245fc22e95a86366687903a"

var (
	senderNonce uint64
	conAddr     common.Address
)

func TestCreateCon(t *testing.T) {
	err := CreateConAddr(senderPirvKey)
	assert.Nil(t, err)
}

func TestCallCon(t *testing.T) {
	var msgData []byte

	senderAddrBytes, err := GetAddrBytes(senderPirvKey)
	assert.Nil(t, err)
	res, err := ContractCall(senderAddrBytes, conAddr, msgData)
	assert.Nil(t, err)
	fmt.Println("=======res", common.Bytes2Hex(res))
}

func CreateConAddr(privKeyStr string) error {
	senderAddrBytes, err := GetAddrBytes(privKeyStr)

	contractAddr := crypto.CreateAddress(senderAddrBytes, senderNonce)

	fmt.Println("=======contractAddr", contractAddr.String())

	return err
}

func GetAddrBytes(privKeyStr string) (addr common.Address, err error) {
	privBytes := common.Hex2Bytes(privKeyStr)

	privkey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return common.Address{}, err
	}
	addr = crypto.PubkeyToAddress(privkey.PublicKey)

	return addr, nil
}

func ContractCall(senderAddr, conAddr common.Address, msgData []byte) ([]byte, error) {
	to := AccountRef(conAddr)
	caller := AccountRef(senderAddr)

	// precompiles := PrecompiledContractsByzantium

	contract := NewContract(caller, to, big.NewInt(0), 0)

	var (
		code []byte
		evm  EVM
	)
	codeHash := crypto.Keccak256Hash(code)

	contract.SetCallCode(&conAddr, codeHash, code)

	ret, err := runIn(&evm, contract, msgData, false)

	return ret, err
}

const EVMGasLimit uint64 = 100000000

var (
	eInterpreters []Interpreter
	eInterpreter  Interpreter
	evmConfig     = Config{EVMGasLimit: EVMGasLimit}
)

func runIn(evm *EVM, contract *Contract, input []byte, readOnly bool) ([]byte, error) {
	precompiles := PrecompiledContractsByzantium

	if p := precompiles[*contract.CodeAddr]; p != nil {
		return p.Run(input)
	}

	NewInterpreter(evm)

	for _, interpreter := range eInterpreters {
		if interpreter.CanRun(contract.Code) {
			if eInterpreter != interpreter {
				// Ensure that the interpreter pointer is set back
				// to its current value upon return.
				defer func(i Interpreter) {
					eInterpreter = i
				}(eInterpreter)
				eInterpreter = interpreter
			}
			return interpreter.Run(contract, input, readOnly)
		}
	}

	return nil, fmt.Errorf("ErrNoCompatibleInterpreter")
}

func NewInterpreter(evm *EVM) {
	eInterpreters = append(eInterpreters, NewEVMInterpreter(evm, evmConfig))
}

func Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {
	if len(contract.Code) == 0 {
		return nil, nil
	}

	var (
		op    OpCode        // current opcode
		mem   = NewMemory() // bound memory
		stack = newstack()  // local stack
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		pc = uint64(0) // program counter
		//cost uint64
		// copies used by tracer
		//pcCopy  uint64 // needed for the deferred Tracer
		//gasCopy uint64 // for Tracer to log gas remaining before execution
		//logged  bool   // deferred Tracer should ignore already logged steps
	)

	contract.Input = input

	op = contract.GetOp(pc)

	var in *EVMInterpreter
	operation := in.cfg.JumpTable[op]
	res, err := operation.execute(&pc, in, contract, mem, stack)

	if operation.returns {
		in.returnData = res
	}

	switch {
	case err != nil:
		return nil, err
	case operation.reverts:
		return res, fmt.Errorf("errExecutionReverted")
	case operation.halts:
		return res, nil
	case !operation.jumps:
		pc++
	}

	return nil, nil
}
