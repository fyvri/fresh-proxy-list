package entity

type Proxy struct {
	Category  string  `json:"category" yaml:"category"`
	Proxy     string  `json:"proxy" yaml:"proxy"`
	IP        string  `json:"ip"  yaml:"ip"`
	Port      string  `json:"port" yaml:"port"`
	TimeTaken float64 `json:"time_taken" yaml:"time_taken"`
	CheckedAt string  `json:"checked_at" yaml:"checked_at"`
}

type AdvancedProxy struct {
	Proxy      string   `json:"proxy" yaml:"proxy"`
	IP         string   `json:"ip" yaml:"ip"`
	Port       string   `json:"port" yaml:"port"`
	TimeTaken  float64  `json:"time_taken" yaml:"time_taken"`
	CheckedAt  string   `json:"checked_at" yaml:"checked_at"`
	Categories []string `json:"categories" yaml:"categories"`
}
