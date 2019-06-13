// Copyright (c) 2017 Robert Zaremba
// Copyright (c) 2017 Sweetbridge Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ethereum

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robert-zaremba/errstack"
)

// ContractFactory delivers methods to easily construct contracts
type ContractFactory interface {
	TxrFactory
}

type contractFactory struct {
	client    *ethclient.Client
	sf        SchemaFactory
	txrF      TxrFactory
	isTestRPC bool

	addrs map[string]common.Address
}

// NewContractFactory is a default contract provider based on truffle schema files.
func NewContractFactory(c *ethclient.Client, sf SchemaFactory, txrF TxrFactory, isTestRPC bool) ContractFactory {
	return contractFactory{c, sf, txrF, isTestRPC,
		map[string]common.Address{}}
}

// Txo implements TxrFactory interface
func (cf contractFactory) Txo() *bind.TransactOpts {
	return cf.txrF.Txo()
}

// Addr returns signer address
func (cf contractFactory) Addr() common.Address {
	return cf.txrF.Addr()
}

func (cf contractFactory) getSchemaAddres(contractName string) (common.Address, errstack.E) {
	if addr, ok := cf.addrs[contractName]; ok {
		return addr, nil
	}
	_, addr, err := cf.sf.ReadGetAddress(contractName)
	if err != nil {
		return addr, err
	}
	cf.addrs[contractName] = addr
	return addr, nil
}

func (cf contractFactory) mkContract(ctrName string, constructor func(common.Address) error) (common.Address, errstack.E) {
	addr, errE := cf.getSchemaAddres(ctrName)
	if errE != nil {
		return addr, errE
	}
	if err := constructor(addr); err != nil {
		return addr, errstack.WrapAsIOf(err, "Can't create new %q contract instance", ctrName)
	}
	return addr, nil
}
