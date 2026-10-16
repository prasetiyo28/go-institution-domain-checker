package checker

// Institution mencerminkan skema data institusi pendidikan
type Institution struct {
	Name          string   `json:"name"`
	Country       string   `json:"country,omitempty"`
	Domains       []string `json:"domains"`
	WebPages      []string `json:"web_pages,omitempty"`
	AlphaTwoCode  string   `json:"alpha_two_code,omitempty"`
	StateProvince string   `json:"state-province,omitempty"`
}

// CheckResult berisi hasil pengecekan domain
type CheckResult struct {
	Domain      string       `json:"domain"`
	Extension   string       `json:"extension"`
	IsAvailable bool         `json:"is_available"`
	Institution *Institution `json:"institution,omitempty"`
}
