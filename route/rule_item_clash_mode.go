package route

import (
	"context"
	"strings"

	"github.com/sagernet/sing-box/adapter"
)

var _ RuleItem = (*ClashModeItem)(nil)

type ClashModeItem struct {
	router adapter.Router
	mode   string
}

func NewClashModeItem(router adapter.Router, mode string) *ClashModeItem {
	return &ClashModeItem{
		router: router,
		mode:   mode,
	}
}

func (r *ClashModeItem) Match(ctx context.Context, metadata *adapter.InboundContext) bool {
	clashServer := r.router.ClashServer()
	if clashServer == nil {
		return false
	}
	return strings.EqualFold(clashServer.Mode(), r.mode)
}

func (r *ClashModeItem) String() string {
	return "clash_mode=" + r.mode
}
