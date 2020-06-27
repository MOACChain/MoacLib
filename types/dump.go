// Copyright 2017  The MOAC Foundation
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"strconv"

	"github.com/MOACChain/MoacLib/common"
	pb "github.com/MOACChain/MoacLib/proto"
)

type GetContractInfoReq struct {
	SubChainAddr common.Address
	Request      []*pb.StorageRequest
}

func ScreeningStorage(storage map[string]string, request []*pb.StorageRequest) map[string]string {
	resp := make(map[string]string)
	for _, val := range request {
		key := common.Bytes2Hex(val.Storagekey)
		position := common.Bytes2Hex(val.Position)
		structformat := val.Structformat
		switch val.Reqtype {
		case 0:
			for k, value := range storage {
				resp[k] = value
			}
		case 1:
			if len(position) == 0 {
				var num int64
				strlen := storage[key]
				if len(strlen) > 2 {
					num, _ = strconv.ParseInt(strlen[2:], 16, 64)
				} else {
					num, _ = strconv.ParseInt(strlen, 16, 64)
				}
				resp[key] = storage[key]
				keys := common.KeytoKey(key)
				for i := int64(0); i < num; i++ {
					if len(structformat) != 0 {
						key0 := keys
						for j := 0; j < len(structformat); j++ {
							if structformat[j] == '1' {
								resp[key0] = storage[key0]
							} else if structformat[j] == '2' {
								resp[key0] = storage[key0]
								var num0 int64
								strlen0 := storage[key0]
								if len(strlen0) > 2 {
									num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
								} else {
									num0, _ = strconv.ParseInt(strlen0, 16, 64)
								}
								key1 := common.KeytoKey(key0)
								for k := int64(0); k < num0; k++ {
									resp[key1] = storage[key1]
									key1 = common.IncreaseHexByOne(key1)
								}
							} else if structformat[j] == '3' {
								nlen := len(storage[key0])
								if nlen == 66 {
									resp[key0] = storage[key0]
								} else if nlen == 2 {
									resp[key0] = storage[key0]
									key1 := common.KeytoKey(key0)
									resp[key1] = storage[key1]
									key1 = common.IncreaseHexByOne(key1)
									resp[key1] = storage[key1]
								} else if nlen > 2 && nlen < 66 {
									resp[key0] = storage[key0]
									if nlen < 7 {
										num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
										key1 := common.KeytoKey(key0)
										for i := num - 1; i > 0; {
											resp[key1] = storage[key1]
											key1 = common.IncreaseHexByOne(key1)
											i = i - 64
										}
									}
								}
							}
							key0 = common.IncreaseHexByOne(key0)
						}
					} else {
						// resp[keys] = storage[keys]
						key0 := keys
						nlen := len(storage[key0])
						if nlen == 66 {
							resp[key0] = storage[key0]
						} else if nlen == 2 {
							resp[key0] = storage[key0]
							key1 := common.KeytoKey(key0)
							resp[key1] = storage[key1]
							key1 = common.IncreaseHexByOne(key1)
							resp[key1] = storage[key1]
						} else if nlen > 2 && nlen < 66 {
							resp[key0] = storage[key0]
							if nlen < 7 {
								num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
								key1 := common.KeytoKey(key0)
								for i := num - 1; i > 0; {
									if storage[key1] != "" {
										resp[key1] = storage[key1]
									}
									key1 = common.IncreaseHexByOne(key1)
									i = i - 64
								}
							}
						}
					}
					keys = common.IncreaseHexByOne(keys)
				}

			} else {
				num, _ := strconv.ParseInt(position, 16, 64)
				keys := common.KeytoKey(key)
				keys = common.IncreaseHexByNum(num, keys)
				if len(structformat) != 0 {
					key0 := keys
					for j := 0; j < len(structformat); j++ {
						if structformat[j] == '1' {
							resp[key0] = storage[key0]
						} else if structformat[j] == '2' {
							resp[key0] = storage[key0]
							var num0 int64
							strlen0 := storage[key0]
							if len(strlen0) > 2 {
								num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
							} else {
								num0, _ = strconv.ParseInt(strlen0, 16, 64)
							}
							key1 := common.KeytoKey(key0)
							for k := int64(0); k < num0; k++ {
								resp[key1] = storage[key1]
								key1 = common.IncreaseHexByOne(key1)
							}
						} else if structformat[j] == '3' {
							nlen := len(storage[key0])
							if nlen == 66 {
								resp[key0] = storage[key0]
							} else if nlen == 2 {
								resp[key0] = storage[key0]
								key1 := common.KeytoKey(key0)
								resp[key1] = storage[key1]
								key1 = common.IncreaseHexByOne(key1)
								resp[key1] = storage[key1]
							} else if nlen > 2 && nlen < 66 {
								resp[key0] = storage[key0]
								if nlen < 7 {
									num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
									key1 := common.KeytoKey(key0)
									for i := num - 1; i > 0; {
										resp[key1] = storage[key1]
										key1 = common.IncreaseHexByOne(key1)
										i = i - 64
									}
								}
							}
						}
						key0 = common.IncreaseHexByOne(key0)
					}
				} else {
					// resp[keys] = storage[keys]
					key0 := keys
					nlen := len(storage[key0])
					if nlen == 66 {
						resp[key0] = storage[key0]
					} else if nlen == 2 {
						resp[key0] = storage[key0]
						key1 := common.KeytoKey(key0)
						resp[key1] = storage[key1]
						key1 = common.IncreaseHexByOne(key1)
						resp[key1] = storage[key1]
					} else if nlen > 2 && nlen < 66 {
						resp[key0] = storage[key0]
						if nlen < 7 {
							num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
							key1 := common.KeytoKey(key0)
							for i := num - 1; i > 0; {
								if storage[key1] != "" {
									resp[key1] = storage[key1]
								}
								key1 = common.IncreaseHexByOne(key1)
								i = i - 64
							}
						}
					}
				}
			}
		case 2:
			keys := common.KeytoKey(position + key)
			if len(structformat) != 0 {
				key0 := keys
				for j := 0; j < len(structformat); j++ {
					if structformat[j] == '1' {
						resp[key0] = storage[key0]
					} else if structformat[j] == '2' {
						resp[key0] = storage[key0]
						var num0 int64
						strlen0 := storage[key0]
						if len(strlen0) > 2 {
							num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
						} else {
							num0, _ = strconv.ParseInt(strlen0, 16, 64)
						}

						key1 := common.KeytoKey(key0)
						for k := int64(0); k < num0; k++ {
							resp[key1] = storage[key1]
							key1 = common.IncreaseHexByOne(key1)
						}
					} else if structformat[j] == '3' {
						nlen := len(storage[key0])
						if nlen == 66 {
							resp[key0] = storage[key0]
						} else if nlen == 2 {
							resp[key0] = storage[key0]
							key1 := common.KeytoKey(key0)
							resp[key1] = storage[key1]
							key1 = common.IncreaseHexByOne(key1)
							resp[key1] = storage[key1]
						} else if nlen > 2 && nlen < 66 {
							resp[key0] = storage[key0]
							if nlen < 7 {
								num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
								key1 := common.KeytoKey(key0)
								for i := num - 1; i > 0; {
									resp[key1] = storage[key1]
									key1 = common.IncreaseHexByOne(key1)
									i = i - 64
								}
							}
						}
					}
					key0 = common.IncreaseHexByOne(key0)
				}
			} else {
				// resp[keys] = storage[keys]
				key0 := keys
				nlen := len(storage[key0])
				if nlen == 66 {
					resp[key0] = storage[key0]
				} else if nlen == 2 {
					resp[key0] = storage[key0]
					key1 := common.KeytoKey(key0)
					resp[key1] = storage[key1]
					key1 = common.IncreaseHexByOne(key1)
					resp[key1] = storage[key1]
				} else if nlen > 2 && nlen < 66 {
					resp[key0] = storage[key0]
					if nlen < 7 {
						num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
						key1 := common.KeytoKey(key0)
						for i := num - 1; i > 0; {
							if storage[key1] != "" {
								resp[key1] = storage[key1]
							}
							key1 = common.IncreaseHexByOne(key1)
							i = i - 64
						}
					}
				}
			}
		case 3:
			key0 := key
			for j := 0; j < len(structformat); j++ {
				if structformat[j] == '1' {
					resp[key0] = storage[key0]
				} else if structformat[j] == '2' {
					resp[key0] = storage[key0]
					var num0 int64
					strlen0 := storage[key0]
					if len(strlen0) > 2 {
						num0, _ = strconv.ParseInt(strlen0[2:], 16, 64)
					} else {
						num0, _ = strconv.ParseInt(strlen0, 16, 64)
					}
					key1 := common.KeytoKey(key0)
					for k := int64(0); k < num0; k++ {
						resp[key1] = storage[key1]
						key1 = common.IncreaseHexByOne(key1)
					}
				} else if structformat[j] == '3' {
					nlen := len(storage[key0])
					if nlen == 66 {
						resp[key0] = storage[key0]
					} else if nlen == 2 {
						resp[key0] = storage[key0]
						key1 := common.KeytoKey(key0)
						resp[key1] = storage[key1]
						key1 = common.IncreaseHexByOne(key1)
						resp[key1] = storage[key1]
					} else if nlen > 2 && nlen < 66 {
						resp[key0] = storage[key0]
						if nlen < 7 {
							num, _ := strconv.ParseInt(storage[key0][2:], 16, 64)
							key1 := common.KeytoKey(key0)
							for i := num - 1; i > 0; {
								resp[key1] = storage[key1]
								key1 = common.IncreaseHexByOne(key1)
								i = i - 64
							}
						}
					}
				}
				key0 = common.IncreaseHexByOne(key0)
			}
		case 4:
			resp[key] = storage[key]
		case 5:
			nlen := len(storage[key])
			if nlen == 66 {
				resp[key] = storage[key]
			} else if nlen == 2 {
				resp[key] = storage[key]
				key0 := common.KeytoKey(key)
				resp[key0] = storage[key0]
				key0 = common.IncreaseHexByOne(key0)
				resp[key0] = storage[key0]
			} else if nlen > 2 && nlen < 66 {
				resp[key] = storage[key]
				if nlen < 7 {
					num, _ := strconv.ParseInt(storage[key][2:], 16, 64)
					key0 := common.KeytoKey(key)
					for i := num - 1; i > 0; {
						resp[key0] = storage[key0]
						key0 = common.IncreaseHexByOne(key0)
						i = i - 64
					}
				}
			}
		default:

		}
	}
	return resp
}
