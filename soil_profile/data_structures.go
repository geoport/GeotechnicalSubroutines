package soil_profile

// Layer is a struct that contains the properties of single layer of a soil profile
type Layer struct {
	soilClass              string
	soilClassManuel        string
	soilClassSelected      string
	soilType               string
	soilDefinition         string
	soilDefinitionManuel   string
	soilDefinitionSelected string
	materialType           string
	thickness              float64
	dryUnitWeight          float64
	saturatedUnitWeight    float64
	fineContent            float64
	liquidLimit            float64
	plasticLimit           float64
	plasticityIndex        float64
	Cu                     float64
	cohesion               float64
	phi                    float64
	waterContent           float64
	poissonRatio           float64
	elasticModulus         float64
	shearModulus           float64
	voidRatio              float64
	Cr                     float64
	Cc                     float64
	Gp                     float64
	mv                     float64
	VS                     float64
	RQD                    float64
	IS50                   float64
	kp                     float64
	dampingRatio           float64
}

// SoilProfile is a struct that contains the properties of a soil profile
type SoilProfile struct {
	Layers       []Layer
	gwt          float64
	checkGwt     bool
	densityUnit  string
	pressureUnit string
}
