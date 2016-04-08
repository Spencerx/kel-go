package kel

// ResourceGroup
// Site
// Instance
// Service
// EnvironmentVariable
// TLSKeyPair
// Build
// Deployment

// ResourceGroup represents a Kel resource group.
type ResourceGroup struct {
	Name     string `json:"name,omitempty"`
	Personal bool   `json:"personal"`
	Created  string `json:"created,omitempty"`

	srv *ResourceGroupService
}

// ResourceGroupUser represents a Kel resource group user.
type ResourceGroupUser struct {
}

// Site represents a Kel site.
type Site struct {
	Name string `json:"name,omitempty"`

	ResourceGroup *ResourceGroup `json:"-"`
}

// Instance represents a Kel instance.
type Instance struct {
	Label string `json:"label,omitempty"`

	Site *Site `json:"-"`
}

// Service represents a Kel service.
type Service struct {
	Name string `json:"name,omitempty"`

	Instance *Instance `json:"-"`
}
