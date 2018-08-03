/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/zhigui/zigledger/core/chaincode/shim"
	pb "github.com/zhigui/zigledger/protos/peer"
)

const (
	//func name
	GetBalance string = "getBalance"
	GetAccount string = "getAccount"
	Transfer   string = "transfer"
	Sender     string = "sender"
	CalcFee    string = "calcFee"
)

// User chaincode for token operations
// After a token issued, users can use this chaincode to make query or transfer operations.
type tokenChaincode struct {
}

// Init func
func (t *tokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("token user chaincode Init.")
	return shim.Success([]byte("Init success."))
}

// Invoke func
func (t *tokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("token user chaincode Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case GetBalance:
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2.")
		}
		return t.getBalance(stub, args)

	case GetAccount:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		return t.getAccount(stub, args)

	case Transfer:
		if len(args) != 3 {
			return shim.Error("Incorrect number of arguments. Expecting 3")
		}
		return t.transfer(stub, args)

	case Sender:
		sender, err := stub.GetSender()
		if err != nil {
			return shim.Error("Get sender failed.")
		}
		return shim.Success([]byte(sender))

	case CalcFee:
		if len(args) != 1 {
			return shim.Error("CalcFee failed.")
		}
		return t.calcFee(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"getBalance\", \"getAccount\", \"transfer\", \"sender\", \"calcFee\".")
}

// getBalance
// Get the balance of a specific token type in an account
func (t *tokenChaincode) getBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string           // Address
	var BalanceType string // Token type
	var err error

	A = strings.ToLower(args[0])
	BalanceType = args[1]
	// Get the state from the ledger
	account, err := stub.GetAccount(A)
	if err != nil {
		return shim.Error("account does not exists")
	}

	if account == nil || account.Balance[BalanceType] == nil {
		return shim.Error("Nil amount for " + A)
	}
	result := make(map[string]string)
	result[BalanceType] = account.Balance[BalanceType].String()
	balanceJson, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		return shim.Error(jsonErr.Error())
	}
	return shim.Success([]byte(balanceJson))
}

// getAccount
// Get the balances of all token types in an account
func (t *tokenChaincode) getAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Address
	var err error

	A = strings.ToLower(args[0])
	// Get the state from the ledger
	account, err := stub.GetAccount(A)
	if err != nil {
		return shim.Error("account does not exists")
	}

	if account == nil {
		return shim.Error("Nil amount for " + A)
	}
	result := make(map[string]string)
	for key, value := range account.Balance {
		result[key] = value.String()
	}
	balanceJson, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		return shim.Error(jsonErr.Error())
	}
	return shim.Success([]byte(balanceJson))
}

// transfer
// Send tokens to the specified address
func (t *tokenChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var B string           // To address
	var BalanceType string // Token type
	var err error

	B = strings.ToLower(args[0])
	BalanceType = args[1]

	// Amount
	amount := big.NewInt(0)
	_, good := amount.SetString(args[2], 10)
	if !good {
		return shim.Error("Expecting integer value for amount")
	}
	err = stub.Transfer(B, BalanceType, amount)
	if err != nil {
		return shim.Error("transfer error" + err.Error())
	}
	return shim.Success(nil)
}

func (t *tokenChaincode) calcFee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fee, err := stub.CalcFee(string(args[0]))
	if err != nil {
		return shim.Error("Query fee failed.")
	}
	res := map[string]interface{}{
		"fee": fee,
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resJson)
}

func main() {
	err := shim.Start(new(tokenChaincode))
	if err != nil {
		fmt.Printf("Error starting tokenChaincode: %s", err)
	}
}

