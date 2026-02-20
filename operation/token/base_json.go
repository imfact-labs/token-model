package token

import (
	"github.com/imfact-labs/currency-model/common"
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
)

type TokenFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender   base.Address      `json:"sender"`
	Contract base.Address      `json:"contract"`
	Currency ctypes.CurrencyID `json:"currency"`
}

func (fact TokenFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TokenFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Currency:              fact.currency,
	})
}

type TokenFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Contract string `json:"contract"`
	Currency string `json:"currency"`
}

func (fact *TokenFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf TokenFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Sender, uf.Contract, uf.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (fact TokenFact) JSONMarshaler() TokenFactJSONMarshaler {
	return TokenFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Currency:              fact.currency,
	}
}
