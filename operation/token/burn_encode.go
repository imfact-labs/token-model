package token

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
)

func (fact *BurnFact) unpack(enc encoder.Encoder,
	ta, am string,
) error {
	switch a, err := base.DecodeAddress(ta, enc); {
	case err != nil:
		return err
	default:
		fact.target = a
	}

	big, err := common.NewBigFromString(am)
	if err != nil {
		return err
	}
	fact.amount = big

	return nil
}
