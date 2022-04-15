package GeotechnicalSubroutines

import (
	np "github.com/geoport/numpy4go/vectors"
	"reflect"
	"testing"
)

var layer1 = Layer{thickness: 1, soilClass: "SC", dryUnitWeight: 1.8, saturatedUnitWeight: 2}
var layer2 = Layer{thickness: 1.4, soilClass: "SP", dryUnitWeight: 1.7, saturatedUnitWeight: 2.1}
var layer3 = Layer{thickness: 3.4, soilClass: "SM", dryUnitWeight: 1.9, saturatedUnitWeight: 2.2}

var soilProfile = SoilProfile{
	Layers: []Layer{layer1, layer2, layer3},
	gwt:    1,
}

func TestGetLayerDepths(t *testing.T) {
	outputLayerDepths := soilProfile.getLayerDepths()
	expectedLayerDepths := []float64{1, 2.4, 5.8}

	if reflect.DeepEqual(outputLayerDepths, expectedLayerDepths) == false {
		t.Errorf("Expected %v, got %v", expectedLayerDepths, outputLayerDepths)
	}
}

func TestGetLayerCenters(t *testing.T) {
	outputLayerCenters := soilProfile.getLayerCenters()
	expectedLayerCenters := []float64{0.5, 1.7, 4.1}

	if reflect.DeepEqual(outputLayerCenters, expectedLayerCenters) == false {
		t.Errorf("Expected %v, got %v", expectedLayerCenters, outputLayerCenters)
	}
}

func TestGetLayerIndex(t *testing.T) {
	testInputs := []float64{-1, 2.4, 7, 3}
	expectedOutputs := []int{0, 1, 2, 2}

	for i, inp := range testInputs {
		output := soilProfile.getLayerIndex(inp)
		if output != expectedOutputs[i] {
			t.Errorf("Expected %v, got %v", expectedOutputs[i], output)
		}
	}
}

func TestGetPropFloat(t *testing.T) {
	expectedOutputs := []float64{1, 1.4, 3.4}
	output := soilProfile.getPropFloat("thickness")

	if reflect.DeepEqual(output, expectedOutputs) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs, output)
	}

}

func TestGetPropString(t *testing.T) {
	expectedOutputs := []string{"SC", "SP", "SM"}
	output := soilProfile.getPropString("soilClass")

	if reflect.DeepEqual(output, expectedOutputs) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs, output)
	}

}

func TestCalcNormalStress(t *testing.T) {
	layers := []Layer{layer1, layer2, layer3}
	SP1 := SoilProfile{Layers: layers, gwt: 1}
	SP2 := SoilProfile{Layers: layers, gwt: 0.5}
	SP3 := SoilProfile{Layers: layers, gwt: 10}

	checkPoints := []float64{0, 1, 1.5, 3, 8}

	expectedOutputs1 := []float64{0, 1.8, 2.85, 6.06, 17.06}
	expectedOutputs2 := []float64{0, 1.9, 2.95, 6.16, 17.16}
	expectedOutputs3 := []float64{0, 1.8, 2.65, 5.32, 14.82}

	output1 := np.Apply(checkPoints, SP1.calcNormalStress)
	output2 := np.Apply(checkPoints, SP2.calcNormalStress)
	output3 := np.Apply(checkPoints, SP3.calcNormalStress)

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
	SP := SoilProfile{Layers: layers, gwt: 1}

	checkPoints := []float64{0.5, 1.5}

	expectedOutputs := []float64{0.9, 2.36}

	output := np.Apply(checkPoints, SP.calcEffectiveStress)

	if reflect.DeepEqual(np.Round(output, 2), expectedOutputs) == false {
		t.Errorf("Expected %v, got %v", expectedOutputs, output)
	}
}
