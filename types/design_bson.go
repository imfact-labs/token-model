package types

import (
	"github.com/imfact-labs/currency-model/utils/bsonenc"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/token-model/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (d Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":   d.Hint().String(),
			"symbol":  d.symbol,
			"name":    d.name,
			"decimal": d.decimal,
			"policy":  d.policy,
		},
	)
}

type DesignBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Symbol  string   `bson:"symbol"`
	Name    string   `bson:"name"`
	Decimal string   `bson:"decimal"`
	Policy  bson.Raw `bson:"policy"`
}

func (d *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*d))

	var u DesignBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return d.unpack(enc, ht, u.Symbol, u.Name, u.Decimal, u.Policy)
}
