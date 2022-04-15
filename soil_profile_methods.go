package GeotechnicalSubroutines

import (
	np "github.com/geoport/numpy4go/vectors"
	"reflect"
)

//getLayerDepths returns the level of bottom of each layer in the soil profile
func (sp *SoilProfile) getLayerDepths() []float64 {
	var depths []float64

	for i, layer := range sp.Layers {
		if i == 0 {
			depths = append(depths, layer.thickness)
		} else {
			depths = append(depths, depths[i-1]+layer.thickness)
		}
	}

	return depths
}

//getLayerCenters returns the center level of each layer in the soil profile
func (sp *SoilProfile) getLayerCenters() []float64 {
	var centers []float64
	var center float64
	depths := sp.getLayerDepths()

	for i, depth := range depths {
		if i == 0 {
			center = depth / 2
		} else {
			center = depths[i-1] + (depth-depths[i-1])/2
		}
		centers = append(centers, center)
	}

	return centers
}

// getLayerIndex returns the index of the layer that contains the given depth
func (sp *SoilProfile) getLayerIndex(depth float64) int {
	layerDepths := sp.getLayerDepths()
	if len(layerDepths) == 1 || depth <= layerDepths[0] {
		return 0
	} else if depth >= layerDepths[len(layerDepths)-1] {
		return len(layerDepths) - 1
	} else {
		diff := np.IncrementBy(layerDepths, -depth)
		for _, i := range np.Arange(1, len(layerDepths), 1) {
			prevDiff := diff[i-1]
			currDiff := diff[i]
			if currDiff == 0 {
				return i
			}
			if prevDiff < 0 && currDiff > 0 {
				return i
			}
		}
	}
	return 0
}

// getPropFloat returns the property values of all the layers as float
func (sp *SoilProfile) getPropFloat(prop string) []float64 {
	var props []float64
	for _, layer := range sp.Layers {
		r := reflect.ValueOf(layer)
		f := reflect.Indirect(r).FieldByName(prop)
		props = append(props, f.Float())
	}
	return props
}

// getPropString returns the property values of all the layers as string
func (sp *SoilProfile) getPropString(prop string) []string {
	var props []string
	for _, layer := range sp.Layers {
		r := reflect.ValueOf(layer)
		f := reflect.Indirect(r).FieldByName(prop)
		props = append(props, f.String())
	}
	return props
}

// calcNormalStress returns the normal stress at the given depth
func (sp *SoilProfile) calcNormalStress(depth float64) float64 {
	Stresses := []float64{0}
	gammaDry := sp.getPropFloat("dryUnitWeight")
	gammaSaturated := sp.getPropFloat("saturatedUnitWeight")
	layerDepths := sp.getLayerDepths()
	layerIndex := sp.getLayerIndex(depth)

	var H1, H0, H float64
	for _, i := range np.Arange(0, layerIndex+1, 1) {
		if i == layerIndex {
			H1 = depth
		} else {
			H1 = layerDepths[i]
		}

		if i == 0 {
			H0 = 0
		} else {
			H0 = layerDepths[i-1]
		}

		if sp.gwt >= H1 {
			H = H1
		} else if H0 >= sp.gwt {
			H = H0
		} else {
			H = sp.gwt
		}

		stress := (H-H0)*gammaDry[i] + gammaSaturated[i]*(H1-H)
		Stresses = append(Stresses, stress+Stresses[i])
	}
	return Stresses[len(Stresses)-1]
}

// calcEffectiveStress returns the effective stress at the given depth
func (sp *SoilProfile) calcEffectiveStress(depth float64) float64 {
	normalStress := sp.calcNormalStress(depth)
	if sp.gwt >= depth {
		return normalStress
	} else {
		return normalStress - (depth-sp.gwt)*0.981
	}
}
