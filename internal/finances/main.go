package finances

import (
	"github.com/gorilla/mux"
)

type Holding struct {
	ARSAmount int64
	MutualFundsValuation int64
	CEDEARsValuation int64
	MEPUSDAmount int64
	MEPUSDQuote float64
}

func (h *Holding) Valuation() int64 {
	return h.ARSAmount +
		h.MutualFundsValuation +
		h.CEDEARsValuation +
		int64(float64(h.MEPUSDAmount) * h.MEPUSDQuote)
}

func AttachTo(r *mux.Router) {
	// TODO: Implement
}
