package domain

type Config struct {
	Groups []Group `json:"groups" yaml:"groups"`
}

type Group struct {
	Name     string    `json:"name" yaml:"name"`
	Archives []Archive `json:"archives" yaml:"archives"`
	Skip     bool      `json:"skip,omitempty" yaml:"skip,omitempty"`
}

type Archive struct {
	Name     string   `json:"name" yaml:"name"`
	Type     string   `json:"type" yaml:"type"`
	Method   *Method  `json:"method,omitempty" yaml:"method,omitempty"`
	Files    []string `json:"files" yaml:"files"`
	Exclude  []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
	Rexclude []string `json:"rexclude,omitempty" yaml:"rexclude,omitempty"`
	Copy     []string `json:"copy,omitempty" yaml:"copy,omitempty"`
}

type Method struct {
	Level *int `json:"level,omitempty" yaml:"level,omitempty"`
}
