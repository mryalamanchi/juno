// Package starknet contains all the functions related to Starknet State and Synchronization
// with Layer 2
package starknet

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/NethermindEth/juno/internal/config"
	"github.com/NethermindEth/juno/internal/log"
	base "github.com/NethermindEth/juno/pkg/common"
	"github.com/NethermindEth/juno/pkg/db"
	"github.com/NethermindEth/juno/pkg/feeder"
	"github.com/NethermindEth/juno/pkg/felt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
)

const latestBlockSynced = "latestBlockSynced"
const blockOfStarknetDeploymentContractMainnet = 13627000
const blockOfStarknetDeploymentContractGoerli = 5853000
const MaxChunk = 10000

// Synchronizer represents the base struct for Ethereum Synchronization
type Synchronizer struct {
	ethereumClient         *ethclient.Client
	feederGatewayClient    *feeder.Client
	db                     *db.Databaser
	MemoryPageHash         base.Dictionary
	GpsVerifier            base.Dictionary
	latestMemoryPageBlock  int64
	latestGpsVerifierBlock int64
	facts                  []string
	lock                   sync.RWMutex
}

// NewSynchronizer creates a new Synchronizer
func NewSynchronizer(db *db.Databaser) *Synchronizer {
	client, err := ethclient.Dial(config.Runtime.Ethereum.Node)
	if err != nil {
		log.Default.With("Error", err).Fatal("Unable to connect to Ethereum Client")
	}
	fClient := feeder.NewClient(config.Runtime.Starknet.FeederGateway, "/feeder_gateway", nil)
	return &Synchronizer{
		ethereumClient:      client,
		feederGatewayClient: fClient,
		db:                  db,
		MemoryPageHash:      base.Dictionary{},
		GpsVerifier:         base.Dictionary{},
		facts:               make([]string, 0),
	}
}

func (s *Synchronizer) initialBlockForStarknetContract() int64 {
	id, err := s.ethereumClient.ChainID(context.Background())
	if err != nil {
		return 0
	}
	if id.Int64() == 1 {
		return blockOfStarknetDeploymentContractMainnet
	}
	return blockOfStarknetDeploymentContractGoerli
}

func (s *Synchronizer) latestBlockQueried() (int64, error) {
	get, err := (*s.db).Get([]byte(latestBlockSynced))
	if err != nil {
		return 0, err
	}
	if get == nil {
		return 0, nil
	}
	var ret uint64
	buf := bytes.NewBuffer(get)
	err = binary.Read(buf, binary.BigEndian, &ret)
	if err != nil {
		return 0, err
	}
	return int64(ret), nil
}

func (s *Synchronizer) updateLatestBlockQueried(block int64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(block))
	err := (*s.db).Put([]byte(latestBlockSynced), b)
	if err != nil {
		log.Default.With("Block", block, "Key", latestBlockSynced).
			Info("Couldn't store the latest synced block")
		return err
	}
	return nil
}

type contractsStruct struct {
	contract  abi.ABI
	eventName string
	topic     string
	address   common.Address
}

type eventStruct struct {
	address         common.Address
	event           map[string]interface{}
	transactionHash common.Hash
}

func (s *Synchronizer) loadEvents(contracts map[common.Address]contractsStruct, eventChan chan eventStruct) error {
	addresses := make([]common.Address, 0)

	topics := make([][]common.Hash, 0)
	for k, v := range contracts {
		addresses = append(addresses, k)
		topics = append(topics, []common.Hash{common.BytesToHash([]byte(v.contract.Events[v.eventName].Sig))})
		log.Default.With("Topic to search", topics[len(topics)-1][0].Hex(), "topics from etherscan", v.topic,
			"Signature", v.contract.Events[v.eventName].Sig, "Bytes2Hex", common.Bytes2Hex([]byte(v.contract.Events[v.eventName].Sig))).
			Info("Checking topics")
	}
	latestBlockNumber, err := s.ethereumClient.BlockNumber(context.Background())
	if err != nil {
		log.Default.With("Error", err).Error("Couldn't get the latest block")
		return err
	}

	initialBlock := s.initialBlockForStarknetContract()
	increment := uint64(MaxChunk)
	i := uint64(initialBlock)
	for i < latestBlockNumber {
		log.Default.With("From Block", i, "To Block", i+increment).Info("Fetching logs....")
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(i)),
			ToBlock:   big.NewInt(int64(i + increment)),
			Addresses: addresses,
			Topics: [][]common.Hash{{
				common.HexToHash("0x73b132cb33951232d83dc0f1f81c2d10f9a2598f057404ed02756716092097bb"),
				common.HexToHash("0xb8b9c39aeba1cfd98c38dfeebe11c2f7e02b334cbe9f05f22b442a5d9c1ea0c5"),
				common.HexToHash("0x9866f8ddfe70bb512b2f2b28b49d4017c43f7ba775f1a20c61c13eea8cdac111"),
			}},
		}

		starknetLogs, err := s.ethereumClient.FilterLogs(context.Background(), query)
		if err != nil {
			log.Default.With("Error", err, "Initial block", i, "End block", i+increment, "Addresses", addresses).
				Info("Couldn't get logs")
			break
		}
		log.Default.With("Count", len(starknetLogs)).Info("Logs fetched")
		for _, vLog := range starknetLogs {
			log.Default.With("Log Fetched", contracts[vLog.Address].eventName, "BlockHash", vLog.BlockHash.Hex(), "BlockNumber", vLog.BlockNumber,
				"TxHash", vLog.TxHash.Hex()).Info("Event Fetched")
			event := map[string]interface{}{}

			err = contracts[vLog.Address].contract.UnpackIntoMap(event, contracts[vLog.Address].eventName, vLog.Data)
			if err != nil {
				log.Default.With("Error", err).Info("Couldn't get LogStateTransitionFact from event")
				continue
			}
			eventChan <- eventStruct{
				event:           event,
				address:         contracts[vLog.Address].address,
				transactionHash: vLog.TxHash,
			}
		}
		i += increment
	}
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(latestBlockNumber)),
		Addresses: addresses,
	}
	hLog := make(chan types.Log)
	sub, err := s.ethereumClient.SubscribeFilterLogs(context.Background(), query, hLog)
	if err != nil {
		log.Default.Info("Couldn't subscribe for incomming blocks")
		return err
	}
	for {
		select {
		case err := <-sub.Err():
			log.Default.With("Error", err).Info("Error getting the latest logs")
		case vLog := <-hLog:
			log.Default.With("Log Fetched", contracts[vLog.Address].eventName, "BlockHash", vLog.BlockHash.Hex(),
				"BlockNumber", vLog.BlockNumber, "TxHash", vLog.TxHash.Hex()).
				Info("Event Fetched")
			event := map[string]interface{}{}
			err = contracts[vLog.Address].contract.UnpackIntoMap(event, contracts[vLog.Address].eventName, vLog.Data)
			if err != nil {
				log.Default.With("Error", err).Info("Couldn't get event from log")
				continue
			}
			eventChan <- eventStruct{
				event:           event,
				address:         contracts[vLog.Address].address,
				transactionHash: vLog.TxHash,
			}
		}
	}
}

func (s *Synchronizer) latestBlockOnChain() (uint64, error) {
	number, err := s.ethereumClient.BlockNumber(context.Background())
	if err != nil {
		log.Default.With("Error", err).Error("Couldn't get the latest block")
		return 0, err
	}
	return number, nil
}

func (s *Synchronizer) FetchStarknetState() error {
	log.Default.Info("Starting to update state")

	contractAddresses, err := s.feederGatewayClient.GetContractAddresses()
	if err != nil {
		log.Default.With("Error", err).Info("Couldn't get Contract Address from Feeder Gateway")
		return err
	}
	event := make(chan eventStruct)

	contracts := make(map[common.Address]contractsStruct)

	// Add Starknet contract
	starknetAddress := common.HexToAddress(contractAddresses.Starknet)
	starknetContract, err := loadContract(config.Runtime.Starknet.ContractAbiPathConfig.StarknetAbiPath)
	if err != nil {
		return err
	}
	contracts[starknetAddress] = contractsStruct{
		contract:  starknetContract,
		topic:     "0x9866f8ddfe70bb512b2f2b28b49d4017c43f7ba775f1a20c61c13eea8cdac111",
		eventName: "LogStateTransitionFact",
	}

	// Add Gps Statement Verifier contract
	gpsStatementVerifierAddress := common.HexToAddress("0xa739B175325cCA7b71fcB51C3032935Ef7Ac338F")
	gpsStatementVerifierContract, err := loadContract(config.Runtime.Starknet.ContractAbiPathConfig.GpsVerifierAbiPath)
	if err != nil {
		return err
	}
	contracts[gpsStatementVerifierAddress] = contractsStruct{
		contract:  gpsStatementVerifierContract,
		topic:     "0x73b132cb33951232d83dc0f1f81c2d10f9a2598f057404ed02756716092097bb",
		eventName: "LogMemoryPagesHashes",
	}
	// Add Memory Page Fact Registry contract
	memoryPageFactRegistryAddress := common.HexToAddress(config.Runtime.Starknet.MemoryPageFactRegistryContract)
	memoryContract, err := loadContract(config.Runtime.Starknet.ContractAbiPathConfig.MemoryPageAbiPath)
	if err != nil {
		return err
	}
	contracts[memoryPageFactRegistryAddress] = contractsStruct{
		contract:  memoryContract,
		topic:     "0xb8b9c39aeba1cfd98c38dfeebe11c2f7e02b334cbe9f05f22b442a5d9c1ea0c5",
		eventName: "LogMemoryPageFactContinuous",
	}
	go func() {

		err = s.loadEvents(contracts, event)
		if err != nil {
			log.Default.With("Error", err).Info("Couldn't get Fact from Starknet Contract Events")
			close(event)
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 5)

		for {
			select {
			case <-ticker.C:
				if len(s.facts) == 0 {
					continue
				}
				if s.GpsVerifier.Exist(s.facts[0]) {
					s.lock.Lock()
					go s.processMemoryPages(s.facts[0])
					s.facts = s.facts[1:]
					s.lock.Unlock()
					return
				}
			}
		}
	}()

	for {
		select {
		case l, ok := <-event:
			if !ok {
				return fmt.Errorf("couldn't read event from logs")
			}
			// Process GpsStatementVerifier contract
			factHash, ok := l.event["factHash"]
			pagesHashes, ok1 := l.event["pagesHashes"]

			if ok && ok1 {
				b := make([]byte, 0)
				for _, v := range factHash.([32]byte) {
					b = append(b, v)
				}

				s.GpsVerifier.Add(common.BytesToHash(b).Hex(), pagesHashes.([][32]byte))
			}
			// Process MemoryPageFactRegistry contract
			if memoryHash, ok := l.event["memoryHash"]; ok {

				key := common.BytesToHash(memoryHash.(*big.Int).Bytes()).Hex()
				value := l.transactionHash
				s.MemoryPageHash.Add(key, value)
			}
			if fact, ok := l.event["stateTransitionFact"]; ok {

				b := make([]byte, 0)
				for _, v := range fact.([32]byte) {
					b = append(b, v)
				}

				s.lock.Lock()
				s.facts = append(s.facts, common.BytesToHash(b).Hex())
				s.lock.Unlock()

			}

		}
	}
}

func (s *Synchronizer) FetchStarknetFact(starknetAddress common.Address, fact chan factChan) error {
	log.Default.Info("Starting to update state")
	latestBlockNumber, err := s.ethereumClient.BlockNumber(context.Background())
	if err != nil {
		log.Default.With("Error", err).Error("Couldn't get the latest block")
		return err
	}
	starknetAbi, err := loadContract(config.Runtime.Starknet.ContractAbiPathConfig.StarknetAbiPath)
	if err != nil {
		log.Default.With("Error", err).Error("Couldn't get the contracts from the ABI")
		return err
	}

	initialBlock := s.initialBlockForStarknetContract()
	increment := uint64(MaxChunk)
	i := uint64(initialBlock)
	for i < latestBlockNumber {
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(i)),
			ToBlock:   big.NewInt(int64(i + increment)),
			Addresses: []common.Address{
				starknetAddress,
			},
		}
		starknetLogs, err := s.ethereumClient.FilterLogs(context.Background(), query)
		if err != nil {
			log.Default.With("Error", err, "Initial block", i, "End block", i+increment).
				Info("Couldn't get logs from starknet contract")
			continue
		}
		for _, vLog := range starknetLogs {
			p, err := vLog.MarshalJSON()
			if err != nil {
				log.Default.With("Error", err).Error("Couldn't unmarshal starknet log")
				return err
			}

			log.Default.With("BlockHash", vLog.BlockHash.Hex(), "BlockNumber", vLog.BlockNumber,
				"TxHash", vLog.TxHash.Hex(), "Data", string(p)).Info("Starknet Event Fetched")
			event := map[string]interface{}{}
			err = starknetAbi.UnpackIntoMap(event, "LogStateTransitionFact", vLog.Data)
			if err != nil {
				log.Default.With("Error", err).Info("Couldn't get LogStateTransitionFact from event")
				continue
			}
			str := fmt.Sprintf("%v", event["stateTransitionFact"].([32]byte))
			log.Default.With("Value", str).Info("Got Log State Transaction Fact")
			fact <- factChan{
				block: int64(vLog.BlockNumber),
				fact:  event["stateTransitionFact"].([32]byte),
			}
		}
		i += increment
	}
	return nil
}

type factChan struct {
	block int64
	fact  [32]byte
}

// UpdateState keeps updated the Starknet State in a process
func (s *Synchronizer) UpdateState() error {
	log.Default.Info("Starting to update state")
	//if config.Runtime.Starknet.FastSync {
	//	s.fastSync()
	//	return nil
	//}

	err := s.FetchStarknetState()
	if err != nil {
		return err
	}

	contractAddresses, err := s.feederGatewayClient.GetContractAddresses()
	if err != nil {
		log.Default.With("Error", err).Info("Couldn't get Contract Address from Feeder Gateway")
		return err
	}
	fact := make(chan factChan)
	go func() {

		err = s.FetchStarknetFact(common.HexToAddress(contractAddresses.Starknet), fact)
		if err != nil {
			log.Default.With("Error", err).Info("Couldn't get Fact from Starknet Contract Events")
			close(fact)
		}
	}()

	for {
		select {
		case l, ok := <-fact:
			if !ok {
				return fmt.Errorf("couldn't read fact from starknet")
			}
			log.Default.With("Fact", common.BytesToHash(l.fact[:]).String(), "Block Number", l.block).
				Info("Getting Fact from Starknet Contract")
			memoryPages := make(chan [][]byte)
			s.memoryPagesFromFact(l, memoryPages)
		}
	}
}

func loadContract(abiPath string) (abi.ABI, error) {
	log.Default.With("Contract", abiPath).Info("Loading contract")
	b, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return abi.ABI{}, err
	}
	contractAbi, err := abi.JSON(strings.NewReader(string(b)))
	if err != nil {
		return abi.ABI{}, err
	}
	return contractAbi, nil
}

// Close closes the client for the Layer 1 Ethereum node
func (s *Synchronizer) Close(ctx context.Context) {
	// notest
	log.Default.Info("Closing Layer 1 Synchronizer")
	select {
	case <-ctx.Done():
		s.ethereumClient.Close()
	default:
	}
}

func (s *Synchronizer) memoryPagesFromFact(l factChan, pages chan [][]byte) {
	_, err := loadContract(config.Runtime.Starknet.ContractAbiPathConfig.MemoryPageAbiPath)
	if err != nil {
		return
	}
}

func (s *Synchronizer) GpsVerifierEvents(l factChan, pages chan [][]byte) {
	_, err := loadContract(config.Runtime.Starknet.ContractAbiPathConfig.MemoryPageAbiPath)
	if err != nil {
		log.Default.With("Error", err).Info("Couldn't load contract MemoryPages")
		return
	}

}

func (s *Synchronizer) fastSync() {
	latestBlockQueried, err := s.latestBlockQueried()
	if err != nil {
		log.Default.With("Error", err).Info("Couldn't get latest Block queried")
		return
	}
	blockIterator := int(latestBlockQueried)
	lastBlockHash := ""
	for {
		newValueForIterator, newBlockHash := s.updateStateForOneBlock(blockIterator, lastBlockHash)
		if newBlockHash == lastBlockHash {
			break
		}
		lastBlockHash = newBlockHash
		blockIterator = newValueForIterator
	}

	ticker := time.NewTicker(time.Minute * 2)
	for {
		select {
		case <-ticker.C:
			newValueForIterator, newBlockHash := s.updateStateForOneBlock(blockIterator, lastBlockHash)
			if newBlockHash == lastBlockHash {
				break
			}
			lastBlockHash = newBlockHash
			blockIterator = newValueForIterator
		}
	}
}

func (s *Synchronizer) updateStateForOneBlock(blockIterator int, lastBlockHash string) (int, string) {
	log.Default.With("Number", blockIterator).Info("Updating StarkNet State")
	update, err := s.feederGatewayClient.GetStateUpdate("", strconv.Itoa(blockIterator))
	if err != nil {
		log.Default.With("Error", err).Info("Couldn't get state update")
		return blockIterator, lastBlockHash
	}
	if lastBlockHash == update.BlockHash {
		return blockIterator, lastBlockHash
	}
	go func() {
		err = s.updateState(update)
		if err != nil {
			log.Default.With("Error", err).Info("Couldn't update state")
			return
		}
		err = s.updateLatestBlockQueried(int64(blockIterator))
		if err != nil {
			log.Default.With("Error", err).Info("Couldn't save latest block queried")
			return
		}
	}()
	return blockIterator + 1, update.BlockHash
}

func (s *Synchronizer) updateState(update feeder.StateUpdateResponse) error {
	log.Default.With("Block Hash", update.BlockHash, "New Root", update.NewRoot, "Old Root", update.OldRoot).
		Info("Updating state")
	s.updateAbiAndCode(update)
	return nil
}

func (s *Synchronizer) processMemoryPages(fact string) {
	pages := make([][]byte, 0)

	// Get memory pages hashes using fact
	var memoryPages [][32]byte
	memoryPages = (s.GpsVerifier.Get(fact)).([][32]byte)

	// iterate over each memory page
	for _, v := range memoryPages {
		h := make([]byte, 0)

		for _, s := range v {
			h = append(h, s)
		}
		// Get transactionsHash based on the memory page
		hash := common.BytesToHash(h)
		transactionHash := s.MemoryPageHash.Get(hash.Hex())
		//	transaction_str = self.memory_page_transactions_map[
		//		int.from_bytes(memory_page_hash, "big")
		//]
		log.Default.With("Hash", hash.Hex()).Info("Getting transaction...")
		txn, _, err := s.ethereumClient.TransactionByHash(context.Background(), transactionHash.(common.Hash))
		if err != nil {
			log.Default.With("Error", err, "Transaction Hash", v).
				Error("Couldn't retrieve transactions")
			return
		}
		// Get the inputs of the transaction from Layer 1
		// Append to the memory pages
		pages = append(pages, txn.Data())
	}
	// pages should contain all txn information
}

type stateToSave struct {
	address       felt.Felt
	contractState struct {
		code    string
		storage []feeder.KV
	}
}

func (s *Synchronizer) updateAbiAndCode(update feeder.StateUpdateResponse) {
	for _, v := range update.StateDiff.DeployedContracts {
		code, err := s.feederGatewayClient.GetCode(v.Address, update.BlockHash, "")
		if err != nil {
			return
		}
		log.Default.With("Contract Address", v.Address, "Block Hash", update.BlockHash, "Code", code).
			Info("Got code and ABI")
		// TODO: Store code and ABI, where to store it? How to store it?

		var address felt.Felt

		err = address.UnmarshalJSON([]byte(v.Address))
		if err != nil {
			log.Default.With("Error", err, "Address", v.Address).Info("Couldn't get felt from address")
			return
		}

		// TODO: Save state to trie
		_ = stateToSave{
			address: address,
			contractState: struct {
				code    string
				storage []feeder.KV
			}{
				code[0], // TODO set how the code is retrieved
				update.StateDiff.StorageDiffs[v.Address],
			},
		}

	}
}

func (s *Synchronizer) updateBlocksAndTransactions(update feeder.StateUpdateResponse) {
	block, err := s.feederGatewayClient.GetBlock(update.BlockHash, "")
	if err != nil {
		return
	}
	log.Default.With("Block Hash", update.BlockHash).
		Info("Got block")
	// TODO: Store block, where to store it? How to store it?

	for _, bTxn := range block.Transactions {
		transactionInfo, err := s.feederGatewayClient.GetTransaction(bTxn.TransactionHash, "")
		if err != nil {
			return
		}
		log.Default.With("Transaction Hash", transactionInfo.Transaction.TransactionHash).
			Info("Got transactions of block")
		// TODO: Store transactions, where to store it? How to store it?

	}

}
