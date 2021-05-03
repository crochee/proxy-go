// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package middleware

import (
	"fmt"
	"net/http"
)

type Handler interface {
	Name() HandlerName
	http.Handler
}

type HandlerName string

var (
	LoadBalancer HandlerName = "LoadBalancer"
	Logger       HandlerName = "Logger"
	RateLimiter  HandlerName = "RateLimiter"
	Recovery     HandlerName = "Recovery"
	Switcher     HandlerName = "Switcher"
	Cross        HandlerName = "Cross"
)

type Action string

var (
	Add    Action = "Add"
	Delete Action = "Delete"
	Update Action = "Update"
	Get    Action = "Get"
	List   Action = "List"
)

func CompleteAction(name HandlerName, action Action) string {
	return fmt.Sprintf("%s:%s", name, action)
}
