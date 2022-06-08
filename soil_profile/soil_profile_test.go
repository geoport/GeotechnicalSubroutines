package soil_profile

import (
	np "github.com/geoport/numpy4go/vectors"
	"reflect"
	"testing"
)

var layer1 = Layer{Thickness: 1,
	SoilClass: "SC", DryUnitWeight: 1.8, SaturatedUnitWeight: 2}
var layer2 = Layer{Thickness: 1.4,
	SoilClass: "SP", DryUnitWeight: 1.7, SaturatedUnitWeight: 2.1}
var layer3 = Layer{Thickness: 3.4,
	SoilClass: "SM", DryUnitWeight: 1.9, SaturatedUnitWeight: 2.2}
var layer4 = Layer{
	SoilClass: "Diğer", LiquidLimit: 10, PlasticLimit: 4,
	SoilClassManuel: "SC", SoilDefinitionManuel: "Kil"}

var soilProfile = SoilProfile{
	Layers: []Layer{layer1, layer2, layer3},
	Gwt:    1,
}

func TestCalcPI(t *testing.T) {
	layer4.CalcPI()
	expected := 6.0
	if expected != layer4.PlasticityIndex {
		t.Errorf("Expected %v, got %v", expected, layer4.PlasticityIndex)
	}
}

func TestLayer_SelectSoilClass(t *testing.T) {
	layer4.SelectSoilClass()
	expected := "SC"
	if expected != layer4.
		SoilClassSelected {
		t.Errorf("Expected %v, got %v", expected, layer4.
			SoilClassSelected)
	}
}

func TestGetLayerDepths(t *testing.T) {
	outputLayerDepths := soilProfile.GetLayerDepths()
	expectedLayerDepths := []float64{1, 2.4, 5.8}

	if reflect.DeepEqual(outputLayerDepths, expectedLayerDepths) == false {
		t.Errorf("Expected %v, got %v", expectedLayerDepths, outputLayerDepths)
	}
}

func TestGetLayerCenters(t *testing.T) {
	outputLayerCenters := soilProfile.GetLayerCenters()
	expectedLayerCenters := []float64{0.5, 1.7, 4.1}

	if reflect.DeepEqual(outputLayerCenters, expectedLayerCenters) == false {
		t.Errorf("Expected %v, got %v", expectedLayerCenters, outputLayerCenters)
	}
}

func TestGetLayerIndex(t *testing.T) {
	testInputs := []float64{-1, 2.4, 7, 3}
	expectedOutputs := []int{0, 1, 2, 2}

	for i, inp := range testInputs {
		output := soilProfile.GetLayerIndex(inp)
		if output != expectedOutputs[i] {
			t.Errorf("Expected %v, got %v", expectedOutputs[i], output)
		}
	}
}

func TestGetPropFloat(t *testing.T) {
	expectedOutputs := []float64{1, 1.4, 3.4}
	output := soilProfile.GetPropFloat("Thickness")

	if reflect.DeepEqual(output, expectedOutputs) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs, output)
	}

}

func TestGetPropString(t *testing.T) {
	expectedOutputs := []string{"SC", "SP", "SM"}
	output := soilProfile.GetPropString("" +
		"SoilClass")

	if reflect.DeepEqual(output, expectedOutputs) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs, output)
	}

}

func TestCalcNormalStress(t *testing.T) {
	layers := []Layer{layer1, layer2, layer3}
	SP1 := SoilProfile{Layers: layers, Gwt: 1}
	SP2 := SoilProfile{Layers: layers, Gwt: 0.5}
	SP3 := SoilProfile{Layers: layers, Gwt: 10}

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

func TestEffectiveStress(t *testing.T) {
	layers := []Layer{layer1, layer2, layer3}
	SP := SoilProfile{Layers: layers, Gwt: 1}

	checkPoints := []float64{0.5, 1.5}

	expectedOutputs := []float64{0.9, 2.36}

	output := np.Apply(checkPoints, SP.CalcEffectiveStress)

	if reflect.DeepEqual(np.Round(output, 2), expectedOutputs) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs, output)
	}
}
