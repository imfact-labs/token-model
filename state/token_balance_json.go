package state

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/token-model/utils"
)

type TokenBalanceStateValueJSONMarshaler struct {
	hint.BaseHinter
	Amount common.Big `json:"amount"`
}

func (s TokenBalanceStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TokenBalanceStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Amount:     s.Amount,
	})
}

type TokenBalanceStateValueJSONUnmarshaler struct {
	Amount string `json:"amount"`
}

func (s *TokenBalanceStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u TokenBalanceStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	big, err := common.NewBigFromString(u.Amount)
	if err != nil {
		return e.Wrap(err)
	}
	s.Amount = big

	return nil
}
