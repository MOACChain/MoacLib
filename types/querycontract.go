// Copyright 2017 The MOAC-core Authors

package types

import (
	"github.com/innowells/moac-lib/common"
)

type QueryContract struct {
	Block           uint           `json:"queryInBlock" gencodec:"required"`
	ContractAddress common.Address `json:"contractAddress"`
}
