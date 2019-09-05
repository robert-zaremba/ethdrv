// Here should be a go-abigen generated code

package cf

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// SweetTokenABI is the input ABI used to generate the binding from.
const SweetTokenABI = "generated ABI"

// SweetToken is a mock of auto generated Go binding around an Ethereum contract.
type SweetToken struct{}

// NewSweetToken creates a new instance of SweetToken, bound to a specific deployed contract.
func NewSweetToken(address common.Address, backend bind.ContractBackend) (*SweetToken, error) {

	return &SweetToken{}, nil
}
