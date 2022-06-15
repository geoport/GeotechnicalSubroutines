package data_structures

import (
	np "github.com/geoport/numpy4go/vectors"
	"reflect"
	"testing"
)

var soilProfile = SoilProfile{
	SoilClass:           []string{"SC", "SP", "SM"},
	DryUnitWeight:       []float64{1.8, 1.7, 1.9},
	SaturatedUnitWeight: []float64{2, 2.1, 2.2},
	Thickness:           []float64{1, 1.4, 3.4},
	Gwt:                 1,
}

var TestSPTData = SPTData{
	Ce:         0.75,
	Cb:         1,
	Cs:         1,
	Depth:      []float64{1.5, 3},
	N:          []int{11, 14},
	Correction: true,
}
var TestCPTData = CPTData{
	Depth:          []float64{1, 3, 5},
	PorePressure:   []float64{10, 20, 30},
	ConeResistance: []float64{11, 12, 13},
}
var TestMASWData = MASWData{
	Thickness: []float64{1, 2, 3},
	VS:        []float64{10, 20, 30},
}

func TestSoilProfile_GetLayerDepths(t *testing.T) {
	outputLayerDepths := soilProfile.GetLayerDepths()
	expectedLayerDepths := []float64{1, 2.4, 5.8}

	if reflect.DeepEqual(outputLayerDepths, expectedLayerDepths) == false {
		t.Errorf("Expected %v, got %v", expectedLayerDepths, outputLayerDepths)
	}
}

func TestSoilProfile_GetLayerCenters(t *testing.T) {
	outputLayerCenters := soilProfile.GetLayerCenters()
	expectedLayerCenters := []float64{0.5, 1.7, 4.1}

	if reflect.DeepEqual(outputLayerCenters, expectedLayerCenters) == false {
		t.Errorf("Expected %v, got %v", expectedLayerCenters, outputLayerCenters)
	}
}

func TestSoilProfile_GetLayerIndex(t *testing.T) {
	testInputs := []float64{-1, 2.4, 7, 3}
	expectedOutputs := []int{0, 1, 2, 2}

	for i, inp := range testInputs {
		output := soilProfile.GetLayerIndex(inp)
		if output != expectedOutputs[i] {
			t.Errorf("Expected %v, got %v", expectedOutputs[i], output)
		}
	}
}

func TestSoilProfile_CalcNormalStress(t *testing.T) {
	SP1 := soilProfile
	SP2 := soilProfile
	SP2.Gwt = 0.5
	SP3 := soilProfile
	SP3.Gwt = 10

	checkPoints := []float64{0, 1, 1.5, 3, 8}

	expectedOutputs1 := []float64{0, 1.8, 2.85, 6.06, 17.06}
	expectedOutputs2 := []float64{0, 1.9, 2.95, 6.16, 17.16}
	expectedOutputs3 := []float64{0, 1.8, 2.65, 5.32, 14.82}

	output1 := np.Apply(checkPoints, SP1.CalcNormalStress)
	output2 := np.Apply(checkPoints, SP2.CalcNormalStress)
	output3 := np.Apply(checkPoints, SP3.CalcNormalStress)

	if reflect.DeepEqual(np.Round(output1, 2), expectedOutputs1) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs1, output1)
	}
	if reflect.DeepEqual(np.Round(output2, 2), expectedOutputs2) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs2, output2)
	}
	if reflect.DeepEqual(np.Round(output3, 2), expectedOutputs3) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs3, output3)
	}
}

func TestSoilProfile_EffectiveStress(t *testing.T) {
	SP := soilProfile
	checkPoints := []float64{0.5, 1.5}

	expectedOutputs := []float64{0.9, 2.36}

	output := np.Apply(checkPoints, SP.CalcEffectiveStress)

	if reflect.DeepEqual(np.Round(output, 2), expectedOutputs) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs, output)
	}
}

func TestSoilProfile_GetLayerFields(t *testing.T) {
	expected := []string{"SoilClass",
		"SoilClassManuel",
		"SoilType",
		"SoilDefinition",
		"MaterialType",
		"Thickness",
		"DryUnitWeight",
		"SaturatedUnitWeight",
		"FineContent",
		"LiquidLimit",
		"PlasticLimit",
		"PlasticityIndex",
		"Cu",
		"Cohesion",
		"Phi",
		"WaterContent",
		"PoissonRatio",
		"ElasticModulus",
		"ShearModulus",
		"VoidRatio",
		"Cr",
		"Cc",
		"Gp",
		"Mv",
		"VS",
		"RQD",
		"IS50",
		"Kp",
		"DampingRatio"}

	output := soilProfile.GetLayerFields()
	if reflect.DeepEqual(output, expected) == false {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}

func TestSoilProfile_GetFieldProperties(t *testing.T) {
	expected := soilProfile.DryUnitWeight
	output := soilProfile.GetFieldProperties("DryUnitWeight").([]float64)
	if reflect.DeepEqual(output, expected) == false {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}

func TestSoilProfile_SetField(t *testing.T) {
	SP := soilProfile
	SP.SetField("FineContent", []float64{1.9, 2.1, 2.2})
	expected := []float64{1.9, 2.1, 2.2}
	output := SP.FineContent
	if reflect.DeepEqual(output, expected) == false {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}

func TestSoilProfile_CombineSPT(t *testing.T) {
	expectedThickness := []float64{1, 0.5, 0.9, 0.6}
	expectedSoilClass := []string{"SC", "SP", "SP", "SM"}
	expectedDryUnitWeight := []float64{1.8, 1.7, 1.7, 1.9}
	expectedN := []int{11, 11, 14, 14}
	outputSoilProfile := soilProfile.CombineSPT(TestSPTData)
	outputThickness := np.Round(outputSoilProfile.Thickness, 2)
	outputSoilClass := outputSoilProfile.SoilClass
	outputDryUnitWeight := np.Round(outputSoilProfile.DryUnitWeight, 2)
	outputN := outputSoilProfile.SPT
	if reflect.DeepEqual(outputThickness, expectedThickness) == false {
		t.Errorf("Expected %v, got %v", expectedThickness, outputThickness)
	}
	if reflect.DeepEqual(outputSoilClass, expectedSoilClass) == false {
		t.Errorf("Expected %v, got %v", expectedSoilClass, outputSoilClass)
	}
	if reflect.DeepEqual(outputDryUnitWeight, expectedDryUnitWeight) == false {
		t.Errorf("Expected %v, got %v", expectedDryUnitWeight, outputDryUnitWeight)
	}
	if reflect.DeepEqual(outputN, expectedN) == false {
		t.Errorf("Expected %v, got %v", expectedN, outputN)
	}
}

func TestSoilProfile_CombineCPT(t *testing.T) {
	expectedThickness := []float64{1, 1.4, 0.6, 2}
	expectedPorePressure := []float64{10, 20, 20, 30}
	outputSoilProfile := soilProfile.CombineCPT(TestCPTData)
	outputThickness := np.Round(outputSoilProfile.Thickness, 2)
	outputPorePressure := np.Round(outputSoilProfile.PorePressure, 2)
	if reflect.DeepEqual(outputThickness, expectedThickness) == false {
		t.Errorf("Expected %v, got %v", expectedThickness, outputThickness)
	}
	if reflect.DeepEqual(outputPorePressure, expectedPorePressure) == false {
		t.Errorf("Expected %v, got %v", expectedPorePressure, outputPorePressure)
	}
}

func TestSoilProfile_CombineVS(t *testing.T) {
	expectedThickness := []float64{1, 1.4, 0.6, 2.8, 0.2}
	expectedVS := []float64{10, 20, 20, 30, 30}
	outputSoilProfile := soilProfile.CombineVS(TestMASWData)
	outputThickness := np.Round(outputSoilProfile.Thickness, 2)
	outputVS := np.Round(outputSoilProfile.VS, 2)
	if reflect.DeepEqual(outputThickness, expectedThickness) == false {
		t.Errorf("Expected %v, got %v", expectedThickness, outputThickness)
	}
	if reflect.DeepEqual(outputVS, expectedVS) == false {
		t.Errorf("Expected %v, got %v", expectedVS, outputVS)
	}
}
