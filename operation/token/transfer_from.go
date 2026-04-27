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
	TransferFromFactHint = hint.MustNewHint("mitum-token-transfer-from-operation-fact-v0.0.1")
	TransferFromHint     = hint.MustNewHint("mitum-token-transfer-from-operation-v0.0.1")
)

var MaxTransferFromItems = 100

type TransferFromFact struct {
	base.BaseFact
	sender   base.Address
	items    []TransferFromItem
	currency types.CurrencyID
}

func NewTransferFromFact(
	token []byte,
	sender base.Address,
	items []TransferFromItem,
	currency types.CurrencyID,
) TransferFromFact {
	fact := TransferFromFact{
		BaseFact: base.NewBaseFact(TransferFromFactHint, token),
		sender:   sender,
		items:    items,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact TransferFromFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items for TransferFromFact")))
	} else if l > int(MaxTransferFromItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxTransferFromItems)))
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

		key := item.contract.String() + "-" + item.target.String() + "-" + item.receiver.String()
		if _, found := founds[key]; found {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(
					errors.Errorf(
						"target account %v and receiver account %v in contract account %v",
						item.target, item.receiver, item.contract)))
		}

		founds[key] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact TransferFromFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFromFact) Bytes() []byte {
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

func (fact TransferFromFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact TransferFromFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferFromFact) Items() []TransferFromItem {
	return fact.items
}

func (fact TransferFromFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact TransferFromFact) Addresses() ([]base.Address, error) {
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

func (fact TransferFromFact) FeeBase() (types.CurrencyID, int, int, bool) {
	return fact.Currency(), len(fact.items), len(fact.Bytes()), extras.HasItem
}

func (fact TransferFromFact) FeePayer() base.Address {
	return fact.sender
}

func (fact TransferFromFact) FactUser() base.Address {
	return fact.sender
}

func (fact TransferFromFact) Signer() base.Address {
	return fact.sender
}

func (fact TransferFromFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

func (fact TransferFromFact) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)
	dupSet := make(map[string]struct{}, len(fact.items))
	for _, item := range fact.items {
		key := fmt.Sprintf("%s:%s", item.Contract().String(), item.Target().String())
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

type TransferFrom struct {
	extras.ExtendedOperation
}

func (op TransferFrom) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)

	if err := extras.AddOperationFeePayerDupKeys(r, op); err != nil {
		return nil, err
	}

	return r, nil
}

func NewTransferFrom(fact TransferFromFact) TransferFrom {
	return TransferFrom{
		ExtendedOperation: extras.NewExtendedOperation(TransferFromHint, fact),
	}
}
