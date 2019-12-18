package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = `
	createChain --address ADDRESS "create a blockchain" 
	send --from FROM --to TO --amount AMOUNT "send coin from FROM to TO"
	getBalance --address Address "get balance of the address"
	printChain	"print all blocks" 
`

const CreateChainCmdString = "createChain"
const GetBalanceCmdString = "getBalance"
const SendCmdString = "send"
const PrintChainCmdString = "printChain"

type CLI struct {
}

func (cli *CLI) printUsage() {
	fmt.Println("Invalid input!")
	fmt.Println(usage)
	os.Exit(1)
}

func (cli *CLI) parameterCheck() {
	if len(os.Args) < 2 {
		cli.printUsage()
	}

}

func (cli *CLI) Run() {
	cli.parameterCheck()

	createChainCmd := flag.NewFlagSet(CreateChainCmdString, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(GetBalanceCmdString, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(SendCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)

	createChainCmdPara := createChainCmd.String("address", "", "address info!")
	getBalanceCmdPara := getBalanceCmd.String("address", "", "address info!")
	sendFromCmdPara := sendCmd.String("from", "", "from address info!")
	sendToCmdPara := sendCmd.String("to", "", "to address info!")
	sendAmountCmdPara := sendCmd.Float64("amount", 0, "amount info!")

	switch os.Args[1] {
	case CreateChainCmdString:
		err := createChainCmd.Parse(os.Args[2:])
		CheckErr("Run1()", err)

		if createChainCmd.Parsed() {
			if *createChainCmdPara == "" {
				fmt.Println("createChain address should not be empty!")
				cli.printUsage()
			}

			cli.CreateChain(*createChainCmdPara)
		}

	case SendCmdString:
		err := sendCmd.Parse(os.Args[2:])
		CheckErr("Run2()", err)

		if sendCmd.Parsed() {
			if *sendFromCmdPara == "" || *sendToCmdPara == "" || *sendAmountCmdPara <= 0 {
				fmt.Println("send cmd parameters invalid!")
				cli.printUsage()
			}

			cli.Send(*sendFromCmdPara, *sendToCmdPara, *sendAmountCmdPara)
		}

	case GetBalanceCmdString:
		err := getBalanceCmd.Parse(os.Args[2:])
		CheckErr("Run3()", err)

		if getBalanceCmd.Parsed() {
			if *getBalanceCmdPara == "" {
				fmt.Println("address should not be empty!")
				cli.printUsage()
			}

			cli.GetBalance(*getBalanceCmdPara)
		}

	case PrintChainCmdString:
		err := printChainCmd.Parse(os.Args[2:])
		CheckErr("Run4()", err)
		if printChainCmd.Parsed() {
			cli.PrintChain()
		}
	default:
		cli.printUsage()
	}
}
