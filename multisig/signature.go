package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	privStr     = "3b95f02dee796480b8761640371e3843af71ca574788ac266cd92df3ef71b375"
	conAddrStr  = "0x63a9804bfe5c8eed6b1925bb6894f77dca990d5a"
	fromAddrStr = "0xca35b7d915458ef540ade6068dfe2f44e8fa733c"
	toAddrStr   = "0x14723a09acff6d2a60dcdf7aa4aff308fddc160c"
)

func main() {
	// admin's private key
	privkey, err := crypto.ToECDSA(common.Hex2Bytes(privStr))
	if err != nil {
		log.Fatal(err)
	}

	// contract's address
	conAddr := common.HexToAddress(conAddrStr).Bytes()
	// transfer's address
	from := common.HexToAddress(fromAddrStr).Bytes()
	// recipient's address
	to := common.HexToAddress(toAddrStr).Bytes()
	// the amount of tokens to be spent
	value := Int64ToBytes(100)

	operationHash := crypto.Keccak256(conAddr, from, to, value)
	fmt.Printf("keccak256 operationHash: %x\n", operationHash)
	sig, err := crypto.Sign(operationHash, privkey)
	if err != nil {
		log.Fatal("Sign failed: ", err)
	}

	fmt.Printf("signature: 0x%x", sig)
}

func Int64ToBytes(i int) []byte {
	var buf = make([]byte, 32)
	PutUint256(buf, uint64(i))
	return buf
}

func PutUint256(b []byte, v uint64) {
	_ = b[31] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 248)
	b[1] = byte(v >> 240)
	b[2] = byte(v >> 232)
	b[3] = byte(v >> 224)
	b[4] = byte(v >> 216)
	b[5] = byte(v >> 208)
	b[6] = byte(v >> 200)
	b[7] = byte(v >> 192)
	b[8] = byte(v >> 184)
	b[9] = byte(v >> 176)
	b[10] = byte(v >> 168)
	b[11] = byte(v >> 160)
	b[12] = byte(v >> 152)
	b[13] = byte(v >> 144)
	b[14] = byte(v >> 136)
	b[15] = byte(v >> 128)
	b[16] = byte(v >> 120)
	b[17] = byte(v >> 112)
	b[18] = byte(v >> 104)
	b[19] = byte(v >> 96)
	b[20] = byte(v >> 88)
	b[21] = byte(v >> 80)
	b[22] = byte(v >> 72)
	b[23] = byte(v >> 64)
	b[24] = byte(v >> 56)
	b[25] = byte(v >> 48)
	b[26] = byte(v >> 40)
	b[27] = byte(v >> 32)
	b[28] = byte(v >> 24)
	b[29] = byte(v >> 16)
	b[30] = byte(v >> 8)
	b[31] = byte(v)
}
