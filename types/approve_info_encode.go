package types

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/token-model/utils"
)

func (a *ApproveInfo) unpack(enc encoder.Encoder, ht hint.Hint, ac, am string) error {
	e := util.StringError(utils.ErrStringUnPack(*a))

	a.BaseHinter = hint.NewBaseHinter(ht)

	switch ad, err := base.DecodeAddress(ac, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		a.account = ad
	}

	amount, err := common.NewBigFromString(am)
	if err != nil {
		return e.Wrap(err)
	}
	a.amount = amount

	return nil
}
