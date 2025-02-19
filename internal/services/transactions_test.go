package services

import (
	"bytes"
	"context"
	"testing"

	"github.com/NethermindEth/juno/internal/db"
	"github.com/NethermindEth/juno/internal/db/transaction"
	"google.golang.org/protobuf/proto"
)

var txs = []*transaction.Transaction{
	{
		Hash: decodeString("49eb3544a95587518b0d2a32b9e456cb05b32e0085ebc0bcecb8ef2e15dc3a2"),
		Tx: &transaction.Transaction_Invoke{Invoke: &transaction.InvokeFunction{
			ContractAddress:    decodeString("7e1b2de3dc9e3cf83278452786c23b384cf77a66c3073f94ab451ed0029b5af"),
			EntryPointSelector: decodeString("317eb442b72a9fae758d4fb26830ed0d9f31c8e7da4dbff4e8c59ea6a158e7f"),
			CallData: [][]byte{
				decodeString("1d24663eb96a2b9a568b88ff520d04779c03ae8ce3087aa55bea1a34f07c6f7"),
				decodeString("2"),
				decodeString("1ab006325ae5196978cfcaefd3d748e8583079ebb33b402537a4b4c174e16c6"),
				decodeString("77575e21a0326adb26b58aaa1ef7139ba8e3164ab54411dbf9a5809b8d6ea8"),
			},
			Signature: nil,
			MaxFee:    decodeString("0"),
		}},
	},
	{
		Hash: decodeString("50398c6ec05a07642e5bd52c656e1650f3b057361283ecbb19d4062199e4626"),
		Tx: &transaction.Transaction_Invoke{Invoke: &transaction.InvokeFunction{
			ContractAddress:    decodeString("3e875a858f9a0229e4a59cb72a4086d324b9b2148242694f2dd12d59d993b62"),
			EntryPointSelector: decodeString("27c3334165536f239cfd400ed956eabff55fc60de4fb56728b6a4f6b87db01c"),
			CallData: [][]byte{
				decodeString("4d6f00affbeb6239fe0eb3eb4afefddbaea71533c152f44a1cdd113c1fdeade"),
				decodeString("33ce93a3eececa5c9fc70da05f4aff3b00e1820b79587924d514bc76788991a"),
				decodeString("1"),
				decodeString("0"),
			},
			Signature: nil,
			MaxFee:    decodeString("0"),
		}},
	},
	{
		Hash: decodeString("1209ae3031dd69ef8ab4507dc4cc2c478d9a0414cb42225ce223670dee5cdcf"),
		Tx: &transaction.Transaction_Invoke{Invoke: &transaction.InvokeFunction{
			ContractAddress:    decodeString("764c36cfdc456e1f3565441938f958badcc0ce8f20b7ed5819af30ed18f245"),
			EntryPointSelector: decodeString("317eb442b72a9fae758d4fb26830ed0d9f31c8e7da4dbff4e8c59ea6a158e7f"),
			CallData: [][]byte{
				decodeString("7e24f78ee360727e9fdd55e2b847202057724bf1f7e5bc25ff78f7760b6895b"),
				decodeString("2"),
				decodeString("51378ba07a08230eab5af933c8e1bd905bc9436bf96ab5f173010eb022eb2a4"),
				decodeString("5f8a361ec261cb4b34d4481803903bb9b8e5c8768e24099aa85ad7f3e8f13b8"),
			},
			Signature: nil,
			MaxFee:    decodeString("0"),
		}},
	},
	{
		Hash: decodeString("e0a2e45a80bb827967e096bcf58874f6c01c191e0a0530624cba66a508ae75"),
		Tx: &transaction.Transaction_Deploy{Deploy: &transaction.Deploy{
			ContractAddressSalt: decodeString("546c86dc6e40a5e5492b782d8964e9a4274ff6ecb16d31eb09cee45a3564015"),
			ConstructorCallData: [][]byte{
				decodeString("06cf6c2f36d36b08e591e4489e92ca882bb67b9c39a3afccf011972a8de467f0"),
				decodeString("7ab344d88124307c07b56f6c59c12f4543e9c96398727854a322dea82c73240"),
			},
		}},
	},
	{
		Hash: decodeString("12c96ae3c050771689eb261c9bf78fac2580708c7f1f3d69a9647d8be59f1e1"),
		Tx: &transaction.Transaction_Deploy{Deploy: &transaction.Deploy{
			ContractAddressSalt: decodeString("12afa0f342ece0468ca9810f0ea59f9c7204af32d1b8b0d318c4f2fe1f384e"),
			ConstructorCallData: [][]byte{
				decodeString("cfc2e2866fd08bfb4ac73b70e0c136e326ae18fc797a2c090c8811c695577e"),
				decodeString("5f1dd5a5aef88e0498eeca4e7b2ea0fa7110608c11531278742f0b5499af4b3"),
			},
		}},
	},
}

func TestTransactionService_StoreTransaction(t *testing.T) {
	defer resetTransactionService()
	database := db.NewKeyValueDb(t.TempDir(), 0)
	TransactionService.Setup(database)
	err := TransactionService.Run()
	if err != nil {
		t.Errorf("error running the service: %s", err)
	}

	for _, tx := range txs {
		TransactionService.StoreTransaction(tx.Hash, tx)
	}
	TransactionService.Close(context.Background())
}

func TestManager_GetTransaction(t *testing.T) {
	defer resetTransactionService()
	database := db.NewKeyValueDb(t.TempDir(), 0)
	TransactionService.Setup(database)
	err := TransactionService.Run()
	if err != nil {
		t.Errorf("error running the service: %s", err)
	}
	// Insert all the transactions
	for _, tx := range txs {
		TransactionService.StoreTransaction(tx.Hash, tx)
	}
	// Get all the transactions and compare
	for _, tx := range txs {
		outTx := TransactionService.GetTransaction(tx.Hash)

		if !equalMessage(t, tx, outTx) {
			t.Errorf("transaction not equal after Put/Get operations")
		}
	}
	TransactionService.Close(context.Background())
}

func resetTransactionService() {
	TransactionService = transactionService{}
}

var receipts = []*transaction.TransactionReceipt{
	{
		TxHash:       decodeString("7932de7ec535bfd45e2951a35c06e13d22188cb7eb7b7cc43454ee63df78aff"),
		ActualFee:    decodeString("0"),
		Status:       0,
		StatusData:   "",
		MessagesSent: nil,
		L1OriginMessage: &transaction.MessageToL2{
			FromAddress: "0x659a00c33263d9254Fed382dE81349426C795BB6",
			Payload: [][]byte{
				decodeString("68a443797ed3eb691347e1d69e6480d1c3ad37acb0d6b1d17c311600002f3d6"),
				decodeString("2616da7c393d14000"),
				decodeString("0"),
				decodeString("b9d83d298d46c4dd73618f19a2a40084ce36476a"),
			},
		},
		Events: []*transaction.Event{
			{
				FromAddress: decodeString("da114221cb83fa859dbdb4c44beeaa0bb37c7537ad5ae66fe5e0efd20e6eb3"),
				Keys: [][]byte{
					decodeString("99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
				},
				Data: [][]byte{
					decodeString("0"),
					decodeString("68a443797ed3eb691347e1d69e6480d1c3ad37acb0d6b1d17c311600002f3d6"),
					decodeString("2616da7c393d14000"),
					decodeString("0"),
				},
			},
			{
				FromAddress: decodeString("1108cdbe5d82737b9057590adaf97d34e74b5452f0628161d237746b6fe69e"),
				Keys: [][]byte{
					decodeString("221e5a5008f7a28564f0eaa32cdeb0848d10657c449aed3e15d12150a7c2db3"),
				},
				Data: [][]byte{
					decodeString("68a443797ed3eb691347e1d69e6480d1c3ad37acb0d6b1d17c311600002f3d6"),
					decodeString("2616da7c393d14000"),
					decodeString("0"),
				},
			},
		},
	},
}

func TestManager_PutReceipt(t *testing.T) {
	defer resetTransactionService()
	database := db.NewKeyValueDb(t.TempDir(), 0)
	TransactionService.Setup(database)
	err := TransactionService.Run()
	if err != nil {
		t.Errorf("error running the service: %s", err)
	}
	for _, receipt := range receipts {
		TransactionService.StoreReceipt(receipt.TxHash, receipt)
	}
	TransactionService.Close(context.Background())
}

func TestManager_GetReceipt(t *testing.T) {
	defer resetTransactionService()
	database := db.NewKeyValueDb(t.TempDir(), 0)
	TransactionService.Setup(database)
	err := TransactionService.Run()
	if err != nil {
		t.Errorf("error running the service: %s", err)
	}
	for _, receipt := range receipts {
		TransactionService.StoreReceipt(receipt.TxHash, receipt)
	}
	for _, receipt := range receipts {
		outReceipt := TransactionService.GetReceipt(receipt.TxHash)

		if !equalMessage(t, receipt, outReceipt) {
			t.Errorf("receipt not equal after Put/Get operations")
		}
	}
	TransactionService.Close(context.Background())
}

func equalMessage(t *testing.T, a, b proto.Message) bool {
	aRaw, err := proto.Marshal(a)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}
	bRaw, err := proto.Marshal(b)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}
	return bytes.Compare(aRaw, bRaw) == 0
}
