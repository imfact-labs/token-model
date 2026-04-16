package token

import (
	"fmt"

	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/operation/extras"
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/valuehash"
	"github.com/imfact-labs/token-model/operation/processor"
	"github.com/pkg/errors"
)

var (
	TransferFactHint = hint.MustNewHint("mitum-token-transfer-operation-fact-v0.0.1")
	TransferHint     = hint.MustNewHint("mitum-token-transfer-operation-v0.0.1")
)

var MaxTransferItems = 100

type TransferFact struct {
	base.BaseFact
	sender   base.Address
	items    []TransferItem
	currency types.CurrencyID
}

func NewTransferFact(
	token []byte,
	sender base.Address,
	items []TransferItem,
	currency types.CurrencyID,
) TransferFact {
	fact := TransferFact{
		BaseFact: base.NewBaseFact(TransferFactHint, token),
		sender:   sender,
		items:    items,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact TransferFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items for TransferFact")))
	} else if l > int(MaxTransferItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxTransferItems)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseFact,
		fact.sender,
		fact.currency,
	); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		if fact.sender.Equal(item.receiver) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with receiver", fact.sender)))
		}

		if item.contract.Equal(item.receiver) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with contract account", item.receiver)))
		}

		if !item.amount.OverZero() {
			return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("transfer amount must be over zero, got %v", item.amount)))
		}

		key := item.contract.String() + "-" + item.receiver.String()
		if _, found := founds[key]; found {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(
					errors.Errorf(
						"receiver account %v in contract account %v", item.receiver, item.contract)))
		}

		founds[key] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact TransferFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.currency.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact TransferFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact TransferFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferFact) Items() []TransferItem {
	return fact.items
}

func (fact TransferFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact TransferFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	for i := range fact.items {
		if ads, err := fact.items[i].Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, ads...)
		}
	}

	as = append(as, fact.Sender())

	return as, nil
}

func (fact TransferFact) FeeBase() (types.CurrencyID, int, int, bool) {
	return fact.Currency(), len(fact.items), len(fact.Bytes()), extras.HasItem
}

func (fact TransferFact) FeePayer() base.Address {
	return fact.sender
}

func (fact TransferFact) FactUser() base.Address {
	return fact.sender
}

func (fact TransferFact) Signer() base.Address {
	return fact.sender
}

func (fact TransferFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

func (fact TransferFact) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)
	r[extras.DuplicationKeyTypeSender] = []string{fact.sender.String()}
	dupSet := make(map[string]struct{}, len(fact.items))
	for _, item := range fact.items {
		key := fmt.Sprintf("%s:%s", item.Contract().String(), fact.sender.String())
		_, found := dupSet[key]
		if !found {
			r[processor.DuplicationTypeTokenSender] = append(
				r[processor.DuplicationTypeTokenSender],
				key,
			)
			dupSet[key] = struct{}{}
		}
	}

	return r, nil
}

type Transfer struct {
	extras.ExtendedOperation
}

func NewTransfer(fact TransferFact) Transfer {
	return Transfer{
		ExtendedOperation: extras.NewExtendedOperation(TransferHint, fact),
	}
}
