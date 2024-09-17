package models

type IPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Org      string `json:"org"`
	Loc      string `json:"loc"`
}
