package types

import (
	"github.com/imfact-labs/currency-model/utils/bsonenc"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/token-model/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (p Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":        p.Hint().String(),
			"total_supply": p.totalSupply,
			"approve_list": p.approveList,
		},
	)
}

type PolicyBSONUnmarshaler struct {
	Hint        string   `bson:"_hint"`
	TotalSupply string   `bson:"total_supply"`
	ApproveList bson.Raw `bson:"approve_list"`
}

func (p *Policy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*p))

	var u PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return p.unpack(enc, ht, u.TotalSupply, u.ApproveList)
}
