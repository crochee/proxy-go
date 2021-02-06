// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package dynamic

import "proxy-go/middlewares/balance"

type Config struct {
	Balancer balance.Node
}
