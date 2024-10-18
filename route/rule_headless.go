package route

import (
	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	dns "github.com/sagernet/sing-dns"
	E "github.com/sagernet/sing/common/exceptions"
)

func NewHeadlessRule(router adapter.Router, logger log.ContextLogger, options option.HeadlessRule) (adapter.HeadlessRule, error) {
	switch options.Type {
	case "", C.RuleTypeDefault:
		if !options.DefaultOptions.IsValid() {
			return nil, E.New("missing conditions")
		}
		return NewDefaultHeadlessRule(router, logger, options.DefaultOptions)
	case C.RuleTypeLogical:
		if !options.LogicalOptions.IsValid() {
			return nil, E.New("missing conditions")
		}
		return NewLogicalHeadlessRule(router, logger, options.LogicalOptions)
	default:
		return nil, E.New("unknown rule type: ", options.Type)
	}
}

var _ adapter.HeadlessRule = (*DefaultHeadlessRule)(nil)

type DefaultHeadlessRule struct {
	abstractDefaultRule
}

func NewDefaultHeadlessRule(router adapter.Router, logger log.ContextLogger, options option.DefaultHeadlessRule) (*DefaultHeadlessRule, error) {
	rule := &DefaultHeadlessRule{
		abstractDefaultRule{
			invert: options.Invert,
		},
	}
	if dns.DomainStrategy(options.DomainStrategy) != dns.DomainStrategyAsIS {
		if router != nil {
			rule.domainStrategy = NewRuleDomainStrategy(router, logger, dns.DomainStrategy(options.DomainStrategy))
		}
	}
	if len(options.Network) > 0 {
		item := NewNetworkItem(options.Network)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.Domain) > 0 || len(options.DomainSuffix) > 0 {
		item := NewDomainItem(options.Domain, options.DomainSuffix)
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	} else if options.DomainMatcher != nil {
		item := NewRawDomainItem(options.DomainMatcher)
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.DomainKeyword) > 0 {
		item := NewDomainKeywordItem(options.DomainKeyword)
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.DomainRegex) > 0 {
		item, err := NewDomainRegexItem(options.DomainRegex)
		if err != nil {
			return nil, E.Cause(err, "domain_regex")
		}
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.SourceIPCIDR) > 0 {
		item, err := NewIPCIDRItem(true, options.SourceIPCIDR)
		if err != nil {
			return nil, E.Cause(err, "source_ip_cidr")
		}
		rule.sourceAddressItems = append(rule.sourceAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	} else if options.SourceIPSet != nil {
		item := NewRawIPCIDRItem(true, options.SourceIPSet)
		rule.sourceAddressItems = append(rule.sourceAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.IPCIDR) > 0 {
		item, err := NewIPCIDRItem(false, options.IPCIDR)
		if err != nil {
			return nil, E.Cause(err, "ipcidr")
		}
		rule.destinationIPCIDRItems = append(rule.destinationIPCIDRItems, item)
		rule.allItems = append(rule.allItems, item)
	} else if options.IPSet != nil {
		item := NewRawIPCIDRItem(false, options.IPSet)
		rule.destinationIPCIDRItems = append(rule.destinationIPCIDRItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.SourcePort) > 0 {
		item := NewPortItem(true, options.SourcePort)
		rule.sourcePortItems = append(rule.sourcePortItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.SourcePortRange) > 0 {
		item, err := NewPortRangeItem(true, options.SourcePortRange)
		if err != nil {
			return nil, E.Cause(err, "source_port_range")
		}
		rule.sourcePortItems = append(rule.sourcePortItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.Port) > 0 {
		item := NewPortItem(false, options.Port)
		rule.destinationPortItems = append(rule.destinationPortItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.PortRange) > 0 {
		item, err := NewPortRangeItem(false, options.PortRange)
		if err != nil {
			return nil, E.Cause(err, "port_range")
		}
		rule.destinationPortItems = append(rule.destinationPortItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.ProcessName) > 0 {
		item := NewProcessItem(options.ProcessName)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.ProcessPath) > 0 {
		item := NewProcessPathItem(options.ProcessPath)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.ProcessPathRegex) > 0 {
		item, err := NewProcessPathRegexItem(options.ProcessPathRegex)
		if err != nil {
			return nil, E.Cause(err, "process_path_regex")
		}
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.PackageName) > 0 {
		item := NewPackageNameItem(options.PackageName)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.WIFISSID) > 0 {
		if router != nil {
			item := NewWIFISSIDItem(router, options.WIFISSID)
			rule.items = append(rule.items, item)
			rule.allItems = append(rule.allItems, item)
		}
	}
	if len(options.WIFIBSSID) > 0 {
		if router != nil {
			item := NewWIFIBSSIDItem(router, options.WIFIBSSID)
			rule.items = append(rule.items, item)
			rule.allItems = append(rule.allItems, item)
		}
	}
	if len(options.AdGuardDomain) > 0 {
		item := NewAdGuardDomainItem(options.AdGuardDomain)
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	} else if options.AdGuardDomainMatcher != nil {
		item := NewRawAdGuardDomainItem(options.AdGuardDomainMatcher)
		rule.destinationAddressItems = append(rule.destinationAddressItems, item)
		rule.allItems = append(rule.allItems, item)
	}
	return rule, nil
}

var _ adapter.HeadlessRule = (*LogicalHeadlessRule)(nil)

type LogicalHeadlessRule struct {
	abstractLogicalRule
}

func NewLogicalHeadlessRule(router adapter.Router, logger log.ContextLogger, options option.LogicalHeadlessRule) (*LogicalHeadlessRule, error) {
	r := &LogicalHeadlessRule{
		abstractLogicalRule{
			rules:  make([]adapter.HeadlessRule, len(options.Rules)),
			invert: options.Invert,
		},
	}
	switch options.Mode {
	case C.LogicalTypeAnd:
		r.mode = C.LogicalTypeAnd
	case C.LogicalTypeOr:
		r.mode = C.LogicalTypeOr
	default:
		return nil, E.New("unknown logical mode: ", options.Mode)
	}
	for i, subRule := range options.Rules {
		rule, err := NewHeadlessRule(router, logger, subRule)
		if err != nil {
			return nil, E.Cause(err, "sub rule[", i, "]")
		}
		r.rules[i] = rule
	}
	return r, nil
}
