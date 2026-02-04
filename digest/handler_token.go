package digest

import (
	"net/http"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/types"
)

var (
	HandlerPathToken        = `/token/{contract:(?i)` + ctypes.REStringAddressString + `}`
	HandlerPathTokenBalance = `/token/{contract:(?i)` + ctypes.REStringAddressString + `}/account/{address:(?i)` + ctypes.REStringAddressString + `}` // revive:disable-line:line-length-limit
)

func SetHandlers(hd *cdigest.Handlers) {
	get := 1000
	_ = hd.SetHandler(HandlerPathTokenBalance, HandleTokenBalance, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.SetHandler(HandlerPathToken, HandleToken, true, get, get).
		Methods(http.MethodOptions, "GET")
}

func HandleToken(hd *cdigest.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleTokenInGroup(hd, contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, hd.ExpireShortLived())
		}
	}
}

func handleTokenInGroup(hd *cdigest.Handlers, contract string) (interface{}, error) {
	switch design, err := Token(hd.Database(), contract); {
	case err != nil:
		return nil, err
	default:
		hal, err := buildTokenHal(hd, contract, *design)
		if err != nil {
			return nil, err
		}
		return hd.Encoder().Marshal(hal)
	}
}

func buildTokenHal(hd *cdigest.Handlers, contract string, design types.Design) (cdigest.Hal, error) {
	h, err := hd.CombineURL(HandlerPathToken, "contract", contract)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(design, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func HandleTokenBalance(hd *cdigest.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	account, err, status := cdigest.ParseRequest(w, r, "address")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleTokenBalanceInGroup(hd, contract, account)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, hd.ExpireShortLived())
		}
	}
}

func handleTokenBalanceInGroup(hd *cdigest.Handlers, contract, account string) (interface{}, error) {
	switch amount, err := TokenBalance(hd.Database(), contract, account); {
	case err != nil:
		return nil, err
	default:
		hal, err := buildTokenBalanceHal(hd, contract, account, amount)
		if err != nil {
			return nil, err
		}
		return hd.Encoder().Marshal(hal)
	}
}

func buildTokenBalanceHal(hd *cdigest.Handlers, contract, account string, amount *common.Big) (cdigest.Hal, error) {
	var hal cdigest.Hal

	if amount == nil {
		hal = cdigest.NewEmptyHal()
	} else {
		h, err := hd.CombineURL(HandlerPathTokenBalance, "contract", contract, "address", account)
		if err != nil {
			return nil, err
		}

		hal = cdigest.NewBaseHal(struct {
			Amount common.Big `json:"amount"`
		}{Amount: *amount}, cdigest.NewHalLink(h, nil))
	}

	return hal, nil
}
