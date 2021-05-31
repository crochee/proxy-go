package dynamic

type Config struct {
	Retry          *Retry          `json:"retry,omitempty" yaml:"retry,omitempty"`
	AccessLog      *LogInfo        `json:"access_log,omitempty" yaml:"access_log,omitempty"`
	Trace          *TraceInfo      `json:"trace,omitempty" yaml:"trace,omitempty"`
	Balance        *BalanceConfig  `json:"balance,omitempty" yaml:"balance,omitempty"`
	RateLimit      *RateLimit      `json:"rate_limit,omitempty" yaml:"rate_limit,omitempty"`
	CrossDomain    bool            `json:"cross_domain,omitempty" yaml:"cross_domain,omitempty"`
	CircuitBreaker *CircuitBreaker `json:"circuit_breaker,omitempty" yaml:"circuit_breaker,omitempty"`
	Recovery       bool            `json:"recovery,omitempty" yaml:"recovery,omitempty"`
}
