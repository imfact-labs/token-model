package types

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/token-model/utils"
)

func (p *Policy) unpack(enc encoder.Encoder, ht hint.Hint, ts string, bap []byte) error {
	e := util.StringError(utils.ErrStringUnPack(*p))

	p.BaseHinter = hint.NewBaseHinter(ht)

	big, err := common.NewBigFromString(ts)
	if err != nil {
		return e.Wrap(err)
	}
	p.totalSupply = big

	hap, err := enc.DecodeSlice(bap)
	if err != nil {
		return e.Wrap(err)
	}

	al := make([]ApproveBox, len(hap))
	for i, h := range hap {
		ap, ok := h.(ApproveBox)
		if !ok {
			return e.Wrap(util.ErrInvalid.Errorf("expected %T, not %T", ApproveBox{}, h))
		}

		al[i] = ap
	}
	p.approveList = al

	return nil
}
