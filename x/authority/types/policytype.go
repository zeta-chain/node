package types

import "fmt"

// Validate PolicyType validates the policy type.
// Valid policy types are groupEmergency, groupOperational, and groupAdmin.
func (p PolicyType) Validate() error {
	if p != PolicyType_groupEmergency && p != PolicyType_groupAdmin &&
		p != PolicyType_groupOperational {
		return fmt.Errorf("invalid policy type: %s", p)
	}
	return nil
}
