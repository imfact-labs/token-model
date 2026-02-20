package token

import (
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
)

func (fact *TokenFact) unpack(enc encoder.Encoder,
	sa, ca, cid string,
) error {
	fact.currency = types.CurrencyID(cid)

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return err
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		fact.contract = a
	}

	return nil
}
