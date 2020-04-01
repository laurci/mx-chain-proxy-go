package data

// Account defines the data structure for an account
type Account struct {
	Address  string `json:"address"`
	Nonce    uint64 `json:"nonce"`
	Balance  string `json:"balance"`
	Code     string `json:"code"`
	CodeHash []byte `json:"codeHash"`
	RootHash []byte `json:"rootHash"`
}

// ValidatorApiResponse represents the data which is fetched from each validator for returning it in API call
type ValidatorApiResponse struct {
	TempRating               float32 `json:"tempRating"`
	NumLeaderSuccess         uint32  `json:"numLeaderSuccess"`
	NumLeaderFailure         uint32  `json:"numLeaderFailure"`
	NumValidatorSuccess      uint32  `json:"numValidatorSuccess"`
	NumValidatorFailure      uint32  `json:"numValidatorFailure"`
	Rating                   float32 `json:"rating"`
	TotalNumLeaderSuccess    uint32  `json:"totalNumLeaderSuccess"`
	TotalNumLeaderFailure    uint32  `json:"totalNumLeaderFailure"`
	TotalNumValidatorSuccess uint32  `json:"totalNumValidatorSuccess"`
	TotalNumValidatorFailure uint32  `json:"totalNumValidatorFailure"`
}

// ResponseAccount defines a wrapped account that the node respond with
type ResponseAccount struct {
	AccountData Account `json:"account"`
}
