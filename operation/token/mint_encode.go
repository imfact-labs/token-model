package token

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
)

func (fact *MintFact) unpack(enc encoder.Encoder,
	ra, am string,
) error {
	switch a, err := base.DecodeAddress(ra, enc); {
	case err != nil:
		return err
	default:
		fact.receiver = a
	}

	big, err := common.NewBigFromString(am)
	if err != nil {
		return err
	}
	fact.amount = big

	return nil
}
