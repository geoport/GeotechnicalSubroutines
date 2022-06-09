package data_structures

type RequestData struct {
	SoilProfile  SoilProfile    `json:"soil_profile"`
	BuildingData BuildingData   `json:"building_data"`
	SeismicData  SeismicData    `json:"eq_information"`
	FieldData    FieldTestsData `json:"investigation_information"`
}

type FieldTestsData struct {
	SPT SPTData           `json:"SPTLog"`
	CPT CPTData           `json:"CPTLog"`
	VS  MASWData          `json:"VSLog"`
	PS  PressureMeterData `json:"PSLog"`
}

// SoilProfile is a struct that contains the properties of a soil profile
type SoilProfile struct {
	SoilClass              []string  `json:"soil_class"`
	SoilClassManuel        []string  `json:"soil_class_manuel"`
	SoilClassSelected      []string  `json:"soil_class_selected"`
	SoilType               []string  `json:"soil_type"`
	SoilDefinition         []string  `json:"soil_definition"`
	SoilDefinitionManuel   []string  `json:"soil_definition_manuel"`
	SoilDefinitionSelected []string  `json:"soil_definition_selected"`
	MaterialType           []string  `json:"material_type"`
	Thickness              []float64 `json:"thickness"`
	DryUnitWeight          []float64 `json:"dry_unit_weight"`
	SaturatedUnitWeight    []float64 `json:"saturated_unit_weight"`
	FineContent            []float64 `json:"fine_content"`
	LiquidLimit            []float64 `json:"liquid_limit"`
	PlasticLimit           []float64 `json:"plastic_limit"`
	PlasticityIndex        []float64 `json:"plasticity_index"`
	Cu                     []float64 `json:"Cu"`
	Cohesion               []float64 `json:"cohesion"`
	Phi                    []float64 `json:"phi"`
	WaterContent           []float64 `json:"water_content"`
	PoissonRatio           []float64 `json:"poisson_ratio"`
	ElasticModulus         []float64 `json:"elastic_modulus"`
	ShearModulus           []float64 `json:"shear_modulus"`
	VoidRatio              []float64 `json:"void_ratio"`
	Cr                     []float64 `json:"Cr"`
	Cc                     []float64 `json:"Cc"`
	Gp                     []float64 `json:"Gp"`
	Mv                     []float64 `json:"mv"`
	VS                     []float64 `json:"VS"`
	RQD                    []float64 `json:"RQD"`
	IS50                   []float64 `json:"IS50"`
	Kp                     []float64 `json:"Kp"`
	DampingRatio           []float64 `json:"damping"`
	Gwt                    float64   `json:"gwt"`
	CheckGwt               bool      `json:"check_gwt"`
	DensityUnit            string    `json:"density_unit"`
	PressureUnit           string    `json:"pressure_unit"`
}

// BuildingData is a struct that contains the properties of a soil profile
type BuildingData struct {
	FoundationType      string  `json:"Foundation_Type"`
	Df                  float64 `json:"Df"`
	FoundationBaseAngle float64 `json:"Base_Angle"`
	B                   float64 `json:"B"`
	L                   float64 `json:"L"`
	SlopeAngle          float64 `json:"Slope_Angle"`
	Vx                  float64 `json:"HB"`
	Vy                  float64 `json:"HL"`
	FrictionCoefficient float64 `json:"FSS"`
	Q                   float64 `json:"Q"`
}

// SeismicData is a struct that contains the properties of earthquake data
type SeismicData struct {
	Mw  float64 `json:"Mw"`
	PGA float64 `json:"PGA"`
}

// SPTData is a struct that contains the properties of SPT log
type SPTData struct {
	Ce         float64   `json:"Ce"`
	Cb         float64   `json:"Cb"`
	Cs         float64   `json:"Cs"`
	Correction bool      `json:"Correction"`
	Depth      []float64 `json:"Depth"`
	N          []float64 `json:"N"`
}

// CPTData is a struct that contains the properties of CPT log
type CPTData struct {
	Depth          []float64 `json:"Depth"`
	ConeResistance []float64 `json:"Cone_Resistance"`
	PorePressure   []float64 `json:"Pore_Pressure"`
}

// MASWData is a struct that contains the properties of MASW log
type MASWData struct {
	Thickness []float64 `json:"Thickness"`
	VS        []float64 `json:"VS"`
	VP        []float64 `json:"VP"`
}

// PressureMeterData is a struct that contains the properties of a pressure meter
type PressureMeterData struct {
	Depth       []float64 `json:"Depth"`
	Pressure    []float64 `json:"PL"`
	NetPressure []float64 `json:"PL_net"`
}
