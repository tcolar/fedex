package api

type API struct {
	Key      string `json:"key"`
	Password string `json:"password"`
	Account  string `json:"account"`
	Meter    string `json:"meter"`
	HubID    string `json:"hubID"` // for SmartPost

	FedExURL string `json:"fedexURL"`
}
