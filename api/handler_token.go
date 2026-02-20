package api

import (
	"github.com/imfact-labs/token-model/digest"
	"net/http"

	apic "github.com/imfact-labs/currency-model/api"
	"github.com/imfact-labs/currency-model/common"
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/token-model/types"
)

var (
	HandlerPathToken        = `/token/{contract:(?i)` + ctypes.REStringAddressString + `}`
	HandlerPathTokenBalance = `/token/{contract:(?i)` + ctypes.REStringAddressString + `}/account/{address:(?i)` + ctypes.REStringAddressString + `}` // revive:disable-line:line-length-limit
)

func SetHandlers(hd *apic.Handlers) {
	get := 1000
	_ = hd.SetHandler(HandlerPathTokenBalance, HandleTokenBalance, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.SetHandler(HandlerPathToken, HandleToken, true, get, get).
		Methods(http.MethodOptions, "GET")
}

func HandleToken(hd *apic.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := apic.CacheKeyPath(r)
	if err := apic.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := apic.ParseRequest(w, r, "contract")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleTokenInGroup(hd, contract)
	}); err != nil {
		apic.HTTP2HandleError(w, err)
	} else {
		apic.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)
		if !shared {
			apic.HTTP2WriteCache(w, cachekey, hd.ExpireShortLived())
		}
	}
}

func handleTokenInGroup(hd *apic.Handlers, contract string) (interface{}, error) {
	switch design, err := digest.Token(hd.Database(), contract); {
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

func buildTokenHal(hd *apic.Handlers, contract string, design types.Design) (apic.Hal, error) {
	h, err := hd.CombineURL(HandlerPathToken, "contract", contract)
	if err != nil {
		return nil, err
	}

	hal := apic.NewBaseHal(design, apic.NewHalLink(h, nil))

	return hal, nil
}

func HandleTokenBalance(hd *apic.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := apic.CacheKeyPath(r)
	if err := apic.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := apic.ParseRequest(w, r, "contract")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	account, err, status := apic.ParseRequest(w, r, "address")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handleTokenBalanceInGroup(hd, contract, account)
	}); err != nil {
		apic.HTTP2HandleError(w, err)
	} else {
		apic.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)
		if !shared {
			apic.HTTP2WriteCache(w, cachekey, hd.ExpireShortLived())
		}
	}
}

func handleTokenBalanceInGroup(hd *apic.Handlers, contract, account string) (interface{}, error) {
	switch amount, err := digest.TokenBalance(hd.Database(), contract, account); {
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

func buildTokenBalanceHal(hd *apic.Handlers, contract, account string, amount *common.Big) (apic.Hal, error) {
	var hal apic.Hal

	if amount == nil {
		hal = apic.NewEmptyHal()
	} else {
		h, err := hd.CombineURL(HandlerPathTokenBalance, "contract", contract, "address", account)
		if err != nil {
			return nil, err
		}

		hal = apic.NewBaseHal(struct {
			Amount common.Big `json:"amount"`
		}{Amount: *amount}, apic.NewHalLink(h, nil))
	}

	return hal, nil
}
