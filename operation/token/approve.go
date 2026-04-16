package token

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/operation/extras"
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	ApproveFactHint = hint.MustNewHint("mitum-token-approve-operation-fact-v0.0.1")
	ApproveHint     = hint.MustNewHint("mitum-token-approve-operation-v0.0.1")
)

var MaxApproveItems = 100

type ApproveFact struct {
	base.BaseFact
	sender   base.Address
	items    []ApproveItem
	currency types.CurrencyID
}

func NewApproveFact(
	token []byte,
	sender base.Address,
	items []ApproveItem,
	currency types.CurrencyID,
) ApproveFact {
	fact := ApproveFact{
		BaseFact: base.NewBaseFact(ApproveFactHint, token),
		sender:   sender,
		items:    items,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact ApproveFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items for ApproveFact")))
	} else if l > int(MaxTransferItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxApproveItems)))
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

		if _, found := founds[item.contract.String()+"-"+item.approved.String()]; found {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(errors.Errorf("contract account %v", item.contract)))
		}

		founds[item.Contract().String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact ApproveFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveFact) Bytes() []byte {
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

func (fact ApproveFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact ApproveFact) Sender() base.Address {
	return fact.sender
}

func (fact ApproveFact) Items() []ApproveItem {
	return fact.items
}

func (fact ApproveFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact ApproveFact) Addresses() ([]base.Address, error) {
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

func (fact ApproveFact) FeeBase() (types.CurrencyID, int, int, bool) {
	return fact.Currency(), len(fact.items), len(fact.Bytes()), extras.HasItem
}

func (fact ApproveFact) FeePayer() base.Address {
	return fact.sender
}

func (fact ApproveFact) FactUser() base.Address {
	return fact.sender
}

func (fact ApproveFact) Signer() base.Address {
	return fact.sender
}

func (fact ApproveFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

func (fact ApproveFact) DupKey() (map[types.DuplicationKeyType][]string, error) {
	r := make(map[types.DuplicationKeyType][]string)
	r[extras.DuplicationKeyTypeSender] = []string{fact.sender.String()}
	dupSet := make(map[string]struct{}, len(fact.items))
	for _, item := range fact.items {
		_, found := dupSet[item.contract.String()]
		if !found {
			r[extras.DuplicationKeyTypeContractStatus] = append(
				r[extras.DuplicationKeyTypeContractStatus],
				item.contract.String(),
			)
			dupSet[item.contract.String()] = struct{}{}
		}
	}

	return r, nil
}

type Approve struct {
	extras.ExtendedOperation
}

func NewApprove(fact ApproveFact) Approve {
	return Approve{
		ExtendedOperation: extras.NewExtendedOperation(ApproveHint, fact),
	}
}
