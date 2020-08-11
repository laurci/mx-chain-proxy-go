package data

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
)

// Transaction represents the structure that maps and validates user input for publishing a new transaction
type Transaction struct {
	// This field is used to tag transactions for send-multiple route
	Index     int    `json:"-"`
	Nonce     uint64 `form:"nonce" json:"nonce"`
	Value     string `form:"value" json:"value"`
	Receiver  string `form:"receiver" json:"receiver"`
	Sender    string `form:"sender" json:"sender"`
	GasPrice  uint64 `form:"gasPrice" json:"gasPrice,omitempty"`
	GasLimit  uint64 `form:"gasLimit" json:"gasLimit,omitempty"`
	Data      []byte `form:"data" json:"data,omitempty"`
	Signature string `form:"signature" json:"signature,omitempty"`
	ChainID   string `form:"chainID" json:"chainID"`
	Version   uint32 `form:"version" json:"version"`
}

// FullTransaction is a transaction featuring all data saved in the full history
type FullTransaction struct {
	Type       string                 `json:"type"`
	Hash       string                 `json:"hash,omitempty"`
	Nonce      uint64                 `json:"nonce,omitempty"`
	Round      uint64                 `json:"round,omitempty"`
	Epoch      uint32                 `json:"epoch,omitempty"`
	Value      string                 `json:"value,omitempty"`
	Receiver   string                 `json:"receiver,omitempty"`
	Sender     string                 `json:"sender,omitempty"`
	GasPrice   uint64                 `json:"gasPrice,omitempty"`
	GasLimit   uint64                 `json:"gasLimit,omitempty"`
	Data       []byte                 `json:"data,omitempty"`
	Code       string                 `json:"code,omitempty"`
	Signature  string                 `json:"signature,omitempty"`
	SndShard   uint32                 `json:"sndShardID,omitempty"`
	RcvShard   uint32                 `json:"rcvShardID,omitempty"`
	BlockNonce uint64                 `json:"blockNonce,omitempty"`
	MBHash     string                 `json:"miniblockHash,omitempty"`
	BlockHash  string                 `json:"blockHash,omitempty"`
	Status     core.TransactionStatus `json:"status,omitempty"`
}

// GetTransactionResponseData follows the format of the data field of get transaction response
type GetTransactionResponseData struct {
	Transaction transaction.ApiTransactionResult `json:"transaction"`
}

// GetTransactionResponse defines a response from the node holding the transaction sent from the chain
type GetTransactionResponse struct {
	Data  GetTransactionResponseData `json:"data"`
	Error string                     `json:"error"`
	Code  string                     `json:"code"`
}

// transactionWrapper is a wrapper over a normal transaction in order to implement the interface needed in elrond-go
// for computing gas cost for a transaction
type transactionWrapper struct {
	transaction     *Transaction
	pubKeyConverter core.PubkeyConverter
}

// NewTransactionWrapper returns a new instance of transactionWrapper
func NewTransactionWrapper(transaction *Transaction, pubKeyConverter core.PubkeyConverter) (*transactionWrapper, error) {
	if transaction == nil {
		return nil, ErrNilTransaction
	}
	if pubKeyConverter == nil {
		return nil, ErrNilPubKeyConverter
	}

	return &transactionWrapper{
		transaction:     transaction,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// GetValue will return the value of the transaction
func (tw *transactionWrapper) GetValue() *big.Int {
	valueBigInt, ok := big.NewInt(0).SetString(tw.transaction.Value, 10)
	if !ok {
		return big.NewInt(0)
	}

	return valueBigInt
}

// GetRcvAddr will return the receiver address in byte slice format
func (tw *transactionWrapper) GetRcvAddr() []byte {
	rcvrBytes, _ := tw.pubKeyConverter.Decode(tw.transaction.Receiver)
	return rcvrBytes
}

// GetGasLimit will return the gas limit of the transaction
func (tw *transactionWrapper) GetGasLimit() uint64 {
	return tw.transaction.GasLimit
}

// GetGasPrice will return the gas price of the transaction
func (tw *transactionWrapper) GetGasPrice() uint64 {
	return tw.transaction.GasPrice
}

// GetData will return the data of the tx
func (tw *transactionWrapper) GetData() []byte {
	return tw.transaction.Data
}

// TransactionResponseData represents the format of the data field of a transaction response
type TransactionResponseData struct {
	TxHash string `json:"txHash"`
}

// ResponseTransaction defines a response tx holding the resulting hash
type ResponseTransaction struct {
	Data  TransactionResponseData `json:"data"`
	Error string                  `json:"error"`
	Code  string                  `json:"code"`
}

// MultipleTransactionsResponseData holds the data which is returned when sending a bulk of transactions
type MultipleTransactionsResponseData struct {
	NumOfTxs  uint64         `json:"txsSent"`
	TxsHashes map[int]string `json:"txsHashes"`
}

// ResponseMultipleTransactions defines a response from the node holding the number of transactions sent to the chain
type ResponseMultipleTransactions struct {
	Data  MultipleTransactionsResponseData `json:"data"`
	Error string                           `json:"error"`
	Code  string                           `json:"code"`
}

// TxCostResponseData follows the format of the data field of a transaction cost request
type TxCostResponseData struct {
	TxCost uint64 `json:"txGasUnits"`
}

// ResponseTxCost defines a response from the node holding the transaction cost
type ResponseTxCost struct {
	Data  TxCostResponseData `json:"data"`
	Error string             `json:"error"`
	Code  string             `json:"code"`
}

// ResponseTxStatus defines a response from the node holding the transaction status
type ResponseTxStatus struct {
	Status string `json:"status"`
}

// FundsRequest represents the data structure needed as input for sending funds from a node to an address
type FundsRequest struct {
	Receiver string   `form:"receiver" json:"receiver"`
	Value    *big.Int `form:"value" json:"value,omitempty"`
	TxCount  int      `form:"txCount" json:"txCount,omitempty"`
}

// ResponseFunds defines the response structure for the node's generate-and-send-multiple endpoint
type ResponseFunds struct {
	Message string `json:"message"`
}
