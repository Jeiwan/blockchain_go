package main

import (
	"flag"
	"fmt"
	"log"

	"os"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
	fmt.Println()
	fmt.Println("Exploring cmds:")
	fmt.Println("  generatePrivKey - generate KeyPair for exploring")
	fmt.Println("  getPubKey -privKey PRIKEY - generate PubKey from privateKey")
	fmt.Println("  getAddress -pubKey PUBKEY - convert pubKey to address")
	fmt.Println("  getPubKeyHash -address Address - get pubKeyHash of an address")
	fmt.Println("  validateAddress -addr Address - validate an address")
	fmt.Println("  getBlock -hash BlockHash - get a block with BlockHash")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	generatePrivKeyCmd := flag.NewFlagSet("generatePrivKey", flag.ExitOnError)
	getPubKeyCmd := flag.NewFlagSet("getPubKey", flag.ExitOnError)
	getAddressCmd := flag.NewFlagSet("getAddress", flag.ExitOnError)
	getPubKeyHashCmd := flag.NewFlagSet("getPubKeyHash", flag.ExitOnError)
	validateAddrCmd := flag.NewFlagSet("validateAddress", flag.ExitOnError)
	getBlockCmd := flag.NewFlagSet("getBlock", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
	privateKey := getPubKeyCmd.String("privKey", "", "generate PubKey based on this")
	pubKey := getAddressCmd.String("pubKey", "", "the key where address generated")
	pubKeyAddress := getPubKeyHashCmd.String("address", "", "the pub address")
	address := validateAddrCmd.String("addr", "", "the public address")
	blockHash := getBlockCmd.String("hash", "", "the block hash")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "validateAddress":
		err := validateAddrCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "generatePrivKey":
		err := generatePrivKeyCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getPubKey":
		err := getPubKeyCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getPubKeyHash":
		err := getPubKeyHashCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getAddress":
		err := getAddressCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBlock":
		err := getBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress, nodeID)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress, nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses(nodeID)
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)
	}

	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeMiner)
	}

	if generatePrivKeyCmd.Parsed() {
		cli.generatePrivKey()
	}

	if getPubKeyCmd.Parsed() {
		if *privateKey == "" {
			getPubKeyCmd.Usage()
			os.Exit(1)
		}
		cli.getPubKey(*privateKey)
	}

	if getAddressCmd.Parsed() {
		if *pubKey == "" {
			getAddressCmd.Usage()
			os.Exit(1)
		}

		cli.getAddress(*pubKey)
	}

	if getPubKeyHashCmd.Parsed() {
		if *pubKeyAddress == "" {
			getPubKeyHashCmd.Usage()
			os.Exit(1)
		}

		cli.getPubKeyHash(*pubKeyAddress)
	}

	if validateAddrCmd.Parsed() {
		if *address == "" {
			validateAddrCmd.Usage()
			os.Exit(1)
		}

		cli.validateAddr(*address)
	}

	if getBlockCmd.Parsed() {
		if *blockHash == "" {
			getBlockCmd.Usage()
			os.Exit(1)
		}

		cli.printBlock(*blockHash, nodeID)
	}

}
