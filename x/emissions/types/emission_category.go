package types

func GetAllCategories() []EmissionCategory {
	return []EmissionCategory{
		EmissionCategory_ObserverEmission,
		EmissionCategory_ValidatorEmission,
		EmissionCategory_TssSignerEmission,
	}
}
