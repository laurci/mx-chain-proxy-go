package facade_test

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/kyber"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/facade/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/assert"
)

func TestNewElrondProxyFacade_NilAccountProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		nil,
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.FaucetProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilAccountProcessor, err)
}

func TestNewElrondProxyFacade_NilTransactionProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		nil,
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.FaucetProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilTransactionProcessor, err)
}

func TestNewElrondProxyFacade_NilGetValuesProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		nil,
		&mock.HeartbeatProcessorStub{},
		&mock.FaucetProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilSCQueryService, err)
}

func TestNewElrondProxyFacade_NilHeartbeatProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		nil,
		&mock.FaucetProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilHeartbeatProcessor, err)
}

func TestNewElrondProxyFacade_NilFaucetProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		nil,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilFaucetProcessor, err)
}

func TestNewElrondProxyFacade_ShouldWork(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.FaucetProcessorStub{},
	)

	assert.NotNil(t, epf)
	assert.Nil(t, err)
}

func TestElrondProxyFacade_GetAccount(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string) (account *data.Account, e error) {
				wasCalled = true
				return &data.Account{}, nil
			},
		},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.FaucetProcessorStub{},
	)

	_, _ = epf.GetAccount("")

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SendTransaction(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SendTransactionCalled: func(tx *data.Transaction) (int, string, error) {
				wasCalled = true

				return 0, "", nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.FaucetProcessorStub{},
	)

	_, _, _ = epf.SendTransaction(&data.Transaction{})

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SendUserFunds(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string) (*data.Account, error) {
				return &data.Account{
					Nonce: uint64(0),
				}, nil
			},
		},
		&mock.TransactionProcessorStub{
			SendTransactionCalled: func(tx *data.Transaction) (int, string, error) {
				wasCalled = true
				return 0, "", nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.FaucetProcessorStub{
			SenderDetailsFromPemCalled: func(receiver string) (crypto.PrivateKey, string, error) {
				return getPrivKey(), "rcvr", nil
			},
			GenerateTxForSendUserFundsCalled: func(senderSk crypto.PrivateKey, senderPk string, senderNonce uint64, receiver string, value *big.Int) (*data.Transaction, error) {
				return &data.Transaction{}, nil
			},
		},
	)

	_ = epf.SendUserFunds("", big.NewInt(0))

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_GetDataValue(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{
			ExecuteQueryCalled: func(query *process.SCQuery) (*vmcommon.VMOutput, error) {
				wasCalled = true
				return &vmcommon.VMOutput{}, nil
			},
		},
		&mock.HeartbeatProcessorStub{},
		&mock.FaucetProcessorStub{},
	)

	_, _ = epf.ExecuteQuery(nil)

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_GetHeartbeatData(t *testing.T) {
	t.Parallel()

	expectedResults := &data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				ReceivedShardID: 0,
				ComputedShardID: 1,
			},
		},
	}
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{
			GetHeartbeatDataCalled: func() (*data.HeartbeatResponse, error) {
				return expectedResults, nil
			},
		},
		&mock.FaucetProcessorStub{},
	)

	actualResult, _ := epf.GetHeartbeatData()

	assert.Equal(t, expectedResults, actualResult)
}

func getPrivKey() crypto.PrivateKey {
	keyGen := signing.NewKeyGenerator(kyber.NewBlakeSHA256Ed25519())
	sk, _ := keyGen.GeneratePair()

	return sk
}
