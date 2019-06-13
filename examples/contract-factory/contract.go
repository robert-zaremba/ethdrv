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
	"bitbucket.org/sweetbridge/oracles/go-contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/ethdrv"
)

// Contract names
const (
	CtrSWC      = "SweetToken"
	CtrSWClogic = "SweetTokenLogic"
)

// Contracts is a list of available contracts
var availableContracts = []string{
	CtrSWC, CtrSWClogic,
}

// ContractFactory delivers methods to easily construct contracts
type ContractFactory interface {
	GetSWC() (*contracts.SweetToken, common.Address, errstack.E)
	GetSWClogic() (*contracts.SweetTokenLogic, common.Address, errstack.E)
	GetForwarderFactory() (*contracts.ForwarderFactory, common.Address, errstack.E)

	ethdrv.TxrFactory
}

type contractFactory struct {
	client    *ethclient.Client
	sf        ethdrv.SchemaFactory
	txrF      ethdrv.TxrFactory
	isTestRPC bool

	addrs map[string]common.Address
}

// NewContractFactory is a default contract provider based on truffle schema files.
func NewContractFactory(c *ethclient.Client, sf ethdrv.SchemaFactory, txrF ethdrv.TxrFactory, isTestRPC bool) ContractFactory {
	return contractFactory{c, sf, txrF, isTestRPC,
		map[string]common.Address{}}
}

// Txo implements TxrFactory interface
func (cf contractFactory) Txo() *bind.TransactOpts {
	txo := cf.txrF.Txo()
	// nonce, err := cf.client.PendingNonceAt(context.TODO(), txo.From)
	// logger.Info("nonce", "val", nonce)
	// if err != nil {
	// 	logger.Error("Can't get pending nonce", err)
	// 	txo.Nonce = big.NewInt(1)
	// 	return txo
	// }
	// if !cf.isTestRPC {
	// 	nonce++
	// }
	// txo.Nonce = big.NewInt(int64(nonce))

	return txo
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

func (cf contractFactory) GetSWC() (c *contracts.SweetToken, addr common.Address, err errstack.E) {
	addr, err = cf.mkContract(CtrSWC, func(addr common.Address) (err2 error) {
		c, err2 = contracts.NewSweetToken(addr, cf.client)
		return
	})
	return
}

func (cf contractFactory) mkContract(ctrName string, constructor func(common.Address) error) (common.Address, errstack.E) {
	addr, errE := cf.getSchemaAddres(ctrName)
	if errE != nil {
		return addr, errE
	}
	if err := constructor(addr); err != nil {
		return addr, errstack.WrapAsInfF(err, "Can't construct %s contract instance", ctrName)
	}
	return addr, nil
}
