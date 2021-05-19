package types

// Capability is a type of access gating mechanism. If present on the user
// account access is granted, otherwise not.
type Capability string

const (
	// CapabilityModifyCI is required for modifying CI properties such as adding or removing a repo.
	CapabilityModifyCI Capability = "modify:ci"
	// CapabilityModifyUser allows you to modify users; including caps.
	CapabilityModifyUser Capability = "modify:user"
	// CapabilitySubmit allows manual submissions
	CapabilitySubmit Capability = "submit"
	// CapabilityCancel allows cancels
	CapabilityCancel Capability = "cancel"
)

// AllCapabilities comprises the superuser account's list of capabilities.
var AllCapabilities = []Capability{CapabilityModifyCI, CapabilityModifyUser, CapabilitySubmit, CapabilityCancel}
