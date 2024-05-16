package userparam

type RegisterRequest struct {
	PublicKey string `json:"public_key"`
}

type RegisterResponse struct {
	EncryptedCode string `json:"encrypted_code"`
}
