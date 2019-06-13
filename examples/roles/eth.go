/* Roles interface example
This package shows how to use a smart-contract interface with go interfaces.
For demonstration we define a simple RBAC interface.
*/
package roles

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/robert-zaremba/errstack"
)

// WithRoles provides interface for contracts which implement Roles contract.
type WithRoles interface {
	Owner(opts *bind.CallOpts) (common.Address, error)
	HasRole(opts *bind.CallOpts, roleName string) (bool, error)
	SenderHasRole(opts *bind.CallOpts, roleName string) (bool, error)
}

// SenderHasRole checks if contract sender has specified role
func SenderHasRole(roleName string, ctr WithRoles) (bool, errstack.E) {
	hasRole, txErr := ctr.SenderHasRole(nil, roleName)
	return hasRole, errstack.WrapAsIOf(txErr, "checking role in %T", ctr)
}

// SenderIsOwnerOrHasRole checks if `who` is owner or sender has role
func SenderIsOwnerOrHasRole(who common.Address, roleName string, ctr WithRoles) (bool, errstack.E) {
	owner, txErr := ctr.Owner(nil)
	if txErr != nil {
		return false, errstack.WrapAsIOf(txErr, "reading owner from %T", ctr)
	}
	if owner == who {
		return true, nil
	}
	return SenderHasRole(roleName, ctr)
}
