package types

func ParseStringToEmissionCategory(category string) EmissionCategory {
	c := EmissionCategory_value[category]
	return EmissionCategory(c)
}
