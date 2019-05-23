package api

type API struct {
	Key      string
	Password string
	Account  string
	Meter    string
	HubID    string // for SmartPost

	FedExURL string
}
