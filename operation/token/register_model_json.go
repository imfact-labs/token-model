package token

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/operation/extras"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/token-model/types"
)

type RegisterModelFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Symbol        types.TokenSymbol `json:"symbol"`
	Name          string            `json:"name"`
	Decimal       common.Big        `json:"decimal"`
	InitialSupply common.Big        `json:"initial_supply"`
}

func (fact RegisterModelFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterModelFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Symbol:                 fact.symbol,
		Name:                   fact.name,
		Decimal:                fact.decimal,
		InitialSupply:          fact.initialSupply,
	})
}

type RegisterModelFactJSONUnMarshaler struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Decimal       string `json:"decimal"`
	InitialSupply string `json:"initial_supply"`
}

func (fact *RegisterModelFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	var uf RegisterModelFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	if err := fact.unpack(enc, uf.Symbol, uf.Name, uf.Decimal, uf.InitialSupply); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op RegisterModel) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *RegisterModel) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
