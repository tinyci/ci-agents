package types

// OAuthToken contains a token and any metadata we want to assign to it.
type OAuthToken struct {
	Token    string   `json:"token"`
	Scopes   []string `json:"scopes"`
	Username string   `json:"username"`
}

// Can returns true if the scope is present in the scopes list.
func (oat *OAuthToken) Can(scope string) bool {
	var found bool

	for _, s := range oat.Scopes {
		if s == scope {
			found = true
			break
		}
	}

	return found
}
