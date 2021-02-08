// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: l30002214
// Create: 2021/2/8

// Package dynamic
package dynamic

type Switch struct {
	ServiceName string      `json:"service_name" binding:"required"`
	Add         bool        `json:"add"`
	Node        BalanceNode `json:"node" binding:"required"`
}
