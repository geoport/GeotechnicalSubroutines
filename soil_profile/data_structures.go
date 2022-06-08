package soil_profile

// Layer is a struct that contains the properties of single layer of a soil profile
type Layer struct {
	SoilClass              string
	SoilClassManuel        string
	SoilClassSelected      string
	SoilType               string
	SoilDefinition         string
	SoilDefinitionManuel   string
	SoilDefinitionSelected string
	MaterialType           string
	Thickness              float64
	DryUnitWeight          float64
	SaturatedUnitWeight    float64
	FineContent            float64
	LiquidLimit            float64
	PlasticLimit           float64
	PlasticityIndex        float64
	Cu                     float64
	Cohesion               float64
	Phi                    float64
	WaterContent           float64
	PoissonRatio           float64
	ElasticModulus         float64
	ShearModulus           float64
	VoidRatio              float64
	Cr                     float64
	Cc                     float64
	Gp                     float64
	mv                     float64
	VS                     float64
	RQD                    float64
	IS50                   float64
	Lp                     float64
	DampingRatio           float64
}

// SoilProfile is a struct that contains the properties of a soil profile
type SoilProfile struct {
	Layers       []Layer
	Gwt          float64
	CheckGwt     bool
	DensityUnit  string
	PressureUnit string
}
