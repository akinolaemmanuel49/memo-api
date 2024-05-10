package response

type AuthResponse struct {
	Tokens `json:",omitempty"`
	User   `json:"profile,omitempty"`
}
