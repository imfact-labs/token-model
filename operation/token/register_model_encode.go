package token

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/token-model/types"
)

func (fact *RegisterModelFact) unpack(enc encoder.Encoder,
	symbol, name, decimal, initialSupply string,
) error {
	fact.symbol = types.TokenSymbol(symbol)
	fact.name = name

	big, err := common.NewBigFromString(decimal)
	if err != nil {
		return err
	}
	fact.decimal = big

	big, err = common.NewBigFromString(initialSupply)
	if err != nil {
		return err
	}
	fact.initialSupply = big

	return nil
}
