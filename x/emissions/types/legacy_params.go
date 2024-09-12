package types

import "gopkg.in/yaml.v2"

func (m *LegacyParams) String() string {
	out, err := yaml.Marshal(m)
	if err != nil {
		return ""
	}
	return string(out)
}
