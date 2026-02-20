package token

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/operation/extras"
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	MintFactHint = hint.MustNewHint("mitum-token-mint-operation-fact-v0.0.1")
	MintHint     = hint.MustNewHint("mitum-token-mint-operation-v0.0.1")
)

type MintFact struct {
	TokenFact
	receiver base.Address
	amount   common.Big
}

func NewMintFact(
	token []byte,
	sender, contract base.Address,
	currency ctypes.CurrencyID,
	receiver base.Address,
	amount common.Big,
) MintFact {
	fact := MintFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(MintFactHint, token), sender, contract, currency,
		),
		receiver: receiver,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact MintFact) IsValid(b []byte) error {
	if err := fact.TokenFact.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := fact.receiver.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.contract.Equal(fact.receiver) {
		return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with contract address", fact.receiver)))
	}

	if !fact.amount.OverZero() {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("mint amount must be over zero, got %v", fact.amount)))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact MintFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MintFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.TokenFact.Bytes(),
		fact.receiver.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact MintFact) Receiver() base.Address {
	return fact.receiver
}

func (fact MintFact) Amount() common.Big {
	return fact.amount
}

func (fact MintFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	as = append(as, fact.TokenFact.Sender())
	as = append(as, fact.TokenFact.Contract())
	as = append(as, fact.receiver)

	return as, nil
}

func (fact MintFact) ActiveContractOwnerHandlerOnly() [][2]base.Address {
	return [][2]base.Address{{fact.contract, fact.sender}}
}

func (fact MintFact) DupKey() (map[ctypes.DuplicationKeyType][]string, error) {
	r := make(map[ctypes.DuplicationKeyType][]string)
	r[extras.DuplicationKeyTypeSender] = []string{fact.sender.String()}
	r[extras.DuplicationKeyTypeContractStatus] = []string{fact.contract.String()}

	return r, nil
}

type Mint struct {
	extras.ExtendedOperation
}

func NewMint(fact MintFact) Mint {
	return Mint{
		ExtendedOperation: extras.NewExtendedOperation(MintHint, fact),
	}
}
