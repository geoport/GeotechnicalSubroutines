package GeotechnicalSubroutines

import (
	"reflect"
	"testing"
)

func TestGetLayerDepths(t *testing.T) {
	layer1 := Layer{thickness: 1}
	layer2 := Layer{thickness: 1.4}
	layer3 := Layer{thickness: 3.4}

	soilProfile := SoilProfile{
		Layers: []Layer{layer1, layer2, layer3},
	}

	outputLayerDepths := soilProfile.getLayerDepths()
	expectedLayerDepths := []float64{1, 2.4, 5.8}

	if reflect.DeepEqual(outputLayerDepths, expectedLayerDepths) == false {
		t.Errorf("Expected %v, got %v", expectedLayerDepths, outputLayerDepths)
	}
}

func TestGetLayerCenters(t *testing.T) {
	layer1 := Layer{thickness: 1}
	layer2 := Layer{thickness: 1.4}
	layer3 := Layer{thickness: 3.4}

	soilProfile := SoilProfile{
		Layers: []Layer{layer1, layer2, layer3},
	}

	outputLayerCenters := soilProfile.getLayerCenters()
	expectedLayerCenters := []float64{0.5, 1.7, 4.1}

	if reflect.DeepEqual(outputLayerCenters, expectedLayerCenters) == false {
		t.Errorf("Expected %v, got %v", expectedLayerCenters, outputLayerCenters)
	}
}
