package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dgraph-io/badger/v2"
	jsondc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go/cmd/util/cmd/common"
	"github.com/onflow/flow-go/model/flow"
	"github.com/sanity-io/litter"
)

func main() {

	heightStr := os.Args[1]

	height, err := strconv.ParseUint(heightStr, 10, 64)
	if err != nil {
		panic(err)
	}
	dir := "badger"
	db, err := badger.Open(badger.DefaultOptions(dir).WithReadOnly(true))
	if err != nil {
		panic(err)
	}
	storage := common.InitStorages(db)

	block, err := storage.Blocks.ByHeight(height)
	if err != nil {
		panic(err)
	}

	transactions := []*flow.TransactionBody{}
	for _, guarantee := range block.Payload.Guarantees {
		collection, err := storage.Collections.ByID(guarantee.CollectionID)
		if err != nil {
			panic(err)
		}

		transactions = append(transactions, collection.Transactions...)
	}
	txResults := []*flow.TransactionResult{}
	for _, tx := range transactions {
		txr, err := storage.TransactionResults.ByBlockIDTransactionID(block.ID(), tx.ID())
		if err != nil {
			panic(err)
		}
		txResults = append(txResults, txr)

	}

	// get all events for a block
	blockEvents, err := storage.Events.ByBlockID(block.ID())
	if err != nil {
		panic(err)
	}

	for _, ev := range blockEvents {
		_, err := jsondc.Decode(ev.Payload)
		if err != nil {
			panic(err)
		}
	}

	litter.Dump(block)
	fmt.Println("transactions", len(transactions))
	fmt.Println("txResult", len(txResults))
	fmt.Println("events", len(blockEvents))

}
