package soil_profile

import (
	np "github.com/geoport/numpy4go/vectors"
	"reflect"
)

//CalcPI calculates the Plasticity index of a layer
func (L *Layer) CalcPI() {
	L.PlasticityIndex = L.LiquidLimit - L.PlasticLimit
}

//SelectSoilClass assigns selected soil class and soil definition to the Layer
func (L *Layer) SelectSoilClass() {
	if L.SoilClass == "Diğer" {
		L.SoilClassSelected = L.SoilClassManuel
		L.SoilDefinitionSelected = L.SoilDefinitionManuel
	} else {
		L.SoilClassSelected = L.SoilClass
		L.SoilDefinitionSelected = L.SoilDefinition
	}
}

//GetLayerDepths returns the level of bottom of each layer in the soil profile
func (sp *SoilProfile) GetLayerDepths() []float64 {
	var depths []float64

	for i, layer := range sp.Layers {
		if i == 0 {
			depths = append(depths, layer.Thickness)
		} else {
			depths = append(depths, depths[i-1]+layer.Thickness)
		}
	}

	return depths
}

//GetLayerCenters returns the center level of each layer in the soil profile
func (sp *SoilProfile) GetLayerCenters() []float64 {
	var centers []float64
	var center float64
	depths := sp.GetLayerDepths()

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

// GetLayerIndex returns the index of the layer that contains the given depth
func (sp *SoilProfile) GetLayerIndex(depth float64) int {
	layerDepths := sp.GetLayerDepths()
	if len(layerDepths) == 1 || depth <= layerDepths[0] {
		return 0
	} else if depth >= layerDepths[len(layerDepths)-1] {
		return len(layerDepths) - 1
	} else {
		diff := np.SumWith(layerDepths, -depth)
		for _, i := range np.Arange(1, float64(len(layerDepths)), 1) {
			prevDiff := diff[int(i-1)]
			currDiff := diff[int(i)]
			if currDiff == 0 {
				return int(i)
			}
			if prevDiff < 0 && currDiff > 0 {
				return int(i)
			}
		}
	}
	return 0
}

// GetPropFloat returns the property values of all the layers as float
func (sp *SoilProfile) GetPropFloat(prop string) []float64 {
	var props []float64
	for _, layer := range sp.Layers {
		r := reflect.ValueOf(layer)
		f := reflect.Indirect(r).FieldByName(prop)
		props = append(props, f.Float())
	}
	return props
}

// GetPropString returns the property values of all the layers as string
func (sp *SoilProfile) GetPropString(prop string) []string {
	var props []string
	for _, layer := range sp.Layers {
		r := reflect.ValueOf(layer)
		f := reflect.Indirect(r).FieldByName(prop)
		props = append(props, f.String())
	}
	return props
}

// CalcNormalStress returns the normal stress at the given depth
func (sp *SoilProfile) CalcNormalStress(depth float64) float64 {
	Stresses := []float64{0}
	gammaDry := sp.GetPropFloat("DryUnitWeight")
	gammaSaturated := sp.GetPropFloat("SaturatedUnitWeight")
	layerDepths := sp.GetLayerDepths()
	layerIndex := sp.GetLayerIndex(depth)

	var H1, H0, H float64
	for i := range np.Arange(0, float64(layerIndex+1), 1) {
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

		if sp.Gwt >= H1 {
			H = H1
		} else if H0 >= sp.Gwt {
			H = H0
		} else {
			H = sp.Gwt
		}

		stress := (H-H0)*gammaDry[i] + gammaSaturated[i]*(H1-H)
		Stresses = append(Stresses, stress+Stresses[i])
	}
	return Stresses[len(Stresses)-1]
}

// calcEffectiveStress returns the effective stress at the given depth
func (sp *SoilProfile) CalcEffectiveStress(depth float64) float64 {
	normalStress := sp.CalcNormalStress(depth)
	if sp.Gwt >= depth {
		return normalStress
	} else {
		return normalStress - (depth-sp.Gwt)*0.981
	}
}
