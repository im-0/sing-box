package route

import (
	"context"
	"strings"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/log"
	dns "github.com/sagernet/sing-dns"
	F "github.com/sagernet/sing/common/format"
)

type RuleDomainStrategy struct {
	router         adapter.Router
	logger         log.ContextLogger
	domainStrategy dns.DomainStrategy
}

func NewRuleDomainStrategy(router adapter.Router, logger log.ContextLogger, domain_strategy dns.DomainStrategy) *RuleDomainStrategy {
	return &RuleDomainStrategy{
		router:         router,
		logger:         logger,
		domainStrategy: domain_strategy,
	}
}

func (r *RuleDomainStrategy) Resolve(ctx context.Context, metadata *adapter.InboundContext) bool {
	if metadata.InsideDomainStrategyRule {
		// Recursion detected: used within DNS rules
		r.logger.Error("domain_strategy is not supported within DNS rules")
		return false
	} else if metadata.AppliedDomainStrategy == r.domainStrategy {
		// Already resolved using the same strategy
		return true
	} else if metadata.Destination.IsFqdn() && r.domainStrategy != dns.DomainStrategyAsIS {
		var newMetadata adapter.InboundContext
		newMetadata = *metadata
		newMetadata.ResetRuleCache()
		newMetadata.InsideDomainStrategyRule = true

		addresses, err := r.router.Lookup(adapter.WithContext(ctx, &newMetadata), metadata.Destination.Fqdn, r.domainStrategy)
		if err == nil {
			metadata.DestinationAddresses = addresses
			r.logger.Debug("rule resolved ", metadata.Destination.Fqdn, " => [", strings.Join(F.MapToString(metadata.DestinationAddresses), " "), "]")
			metadata.AppliedDomainStrategy = r.domainStrategy
			return true
		} else {
			r.logger.Error("rule failed to resolve ", metadata.Destination.Fqdn, " addresses: ", err)
			metadata.AppliedDomainStrategy = dns.DomainStrategyAsIS
			return false
		}
	} else if metadata.Destination.IsIP() || len(metadata.DestinationAddresses) > 0 {
		// No FQDN, but we already have addresses => no need to resolve
		metadata.AppliedDomainStrategy = dns.DomainStrategyAsIS
		return true
	} else {
		r.logger.Error("rule failed to resolve addresses: no FQDN")
		metadata.AppliedDomainStrategy = dns.DomainStrategyAsIS
		return false
	}
}
