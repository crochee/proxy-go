// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package dynamic

type Config struct {
	Balancer BalanceNode
	Limit    RateLimit
}
