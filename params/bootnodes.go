// Copyright 2015 The go-ethereum Authors
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

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main MOAC network.
var MainnetBootnodes = []string{
	// Go Bootnodes
	"enode://3d8ba7cef2dcc8e25bd508e6a42d7b36a957afe8dc5cd9603ffd0857aae68b2768a06de6054816edec407a4579b15b5c7494e9c7e8323b015d6e9d518576e9eb@18.233.50.84:30333",
	"enode://6f29307c8715502c11791704759d88910e091f5d4a6daf509fc9d0d8653967fd2cda76c94731e2db6654ff889a684a55f9a6bc6f0b3144282eb2c3e7cbccb4b1@18.211.187.153:30333",
	"enode://1cf0eb6052ee692e85cf5447c9b1ef95c57f23e2efe1cd8cd5a01b44feadc3395f3aa02e0b1130dffcfdac3225b3ce8d773bfd148efadc3846353337a45c99a2@54.71.244.228:30333",
	"enode://bebab9587e939fab7e21922864cdf6668193e1fcfe5f69ef128216c96e5e08dae1a8d9b99e0cb158bd040c9ac249c3569d43675f0e637fb2cf9e2df9e39738d7@52.34.175.72:30333",
	"enode://0333b61a2d6fb25fb592b68fdf64b324ab0b27ecad815cf47b459dc83f56d191c81da1eb6d015ca14c87b822e1f2446112b9d15790209ef2aa8f22afd3f978c4@124.70.166.158:30333",
	"enode://b413ccff047e38415b4d6fc20ad79461956943d62cdfcfa9d210d8c522419cc4db1b229c7da5ad15e92a4228363fe54d712c215ef854e6992f4292c07c66744e@117.78.9.104:30333",
	"enode://5e59505a2aac6c24df868d3d5fc7047b2f50512773409c0a2c38eb628fdd58aa4826909109b284f3e0e8d2727badb285700cae334ed6d889ef0cb24a47df358e@121.37.217.111:30333",
	"enode://9675508dfc11685d762f0b95f9f939af872baf2d2f289e6b142768030dd62a8f413fc7baa6bd2744a4ec7bf8d001ed59dc82e47e46702092097dcc7a7f08beee@139.198.122.215:30333",
	"enode://68c48e272a78e288c9e86ca2c272cf51ed9dbbe305e8f257cf347b4282dc04bc404ca25178025661216b4f15d610b1d536ab3c8610ca4357251ba85451a889c4@18.130.240.247:30333",
	"enode://624c26083a9be2ac0624bc5e4d88633aca89f62bbca1e7929a05bfc9f06a42202168b12ed5303f95e26296f8092d07b61071eedaccfe7027ceb0e7330589bc5c@13.211.142.153:30333",
	//Canada
	"enode://4c8fac5ed1f8e5676d8f6b722fdb69b289dd58fe5844f0ae036b9d257bc6c4ebb94641cd590f8d7ed282c7690217eaf17898712b974ce23cfccbe167edc1faa0@35.182.237.1:30333",
	//Singapore
	"enode://ab485ba14e3955cedeff78ef760bc34c450f80396ea5f2c9a2669e2aa7d16dec014a6772a5eacd3726f8ee83cc0877ce2c0cf2a9a95b6e3e239df15bcbe52d58@47.88.237.205:30333",
	//Tokyo, Japan
	"enode://3303754238a2a2f07aadf7683224768c562822517e270b38fad7cf6df9c1af6d94a42ce361bbe965a38f2a2b1dfc2873518a8c20a5a004e8aa105d3bc3dd4a8b@47.74.9.125:30333",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes
var TestnetBootnodes = []string{
	"enode://271c55ef39be9208e6ad75c935061412b39e51dd97a8e4dbba7d358e91132fd7c79ee687228edea3fd9c833b6ce9c365365aa526999956914d7ac81d00576e76@18.217.180.94:30333",
	"enode://089554d6929600b9c70bbd6e1c12594697d0aec43127b9b29c6eb96faf06884fd284f56c3de64155d65e540b59a43f4fd07d8802b4b5e95b2922531e6096c2d5@52.15.143.41:30333",
	"enode://4489532a497654b67f5e6417f19950193208ec3ac06b0f13babdb315bf091cc650850a9a1f8d01b33d80dfddd0a0e58c4b5b11a399f7d1cd86eb9d1f6269fe60@18.188.171.176:30333",
	"enode://0eb6b47f65c2225981d6511dfa6472547cec15c854eb87c73af1db44a943353d32f5366d0dbb16573271ed56397cd8ee5b73402816507a367d522b3d30cfe369@124.70.63.159:30333",
	"enode://1c744f4ed3cfe94f27f0fa80287fe54b4755b9c3f0a526f16516da6312af352613276e4365e158ee8fc4aeeee90dc851a95149333bb6c75038b4668612aaea62@124.70.63.132:30333",
	"enode://e88a288b19e55a79cdeabd9c5326602dbb2ba5a5e29eca66bdddf20d598cf7126fce9080b53399ef02e200891cd590b6ff9c23a8c15789f04790d1e69054bd9c@39.108.79.40:30333",
	"enode://074d95c52cd573a7e1525eb8ef303f245401aeed1d7e2c603661ad9962ecdd3b590aaf9eacc5f0532c6afe97216da86a570d7b48778347856f3c431e78078197@139.198.18.165:30333",
	"enode://832bdd83b2f5c89b68a2ac733d7e31d26bbf904d5a5316e857802c0055b1eed3bd5daa2894bc9ac570ae0bb0c761d05c61621db23ffcd3ac2556ce3dc0503086@139.198.177.202:30333",
}

// SubnetBootnodes are the enode URLs of the subnet P2P bootstrape nodes
var SubnetBootnodes = []string{}
