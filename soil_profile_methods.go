package GeotechnicalSubroutines

import (
	"fmt"
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

func (sp *SoilProfile) getLayerIndex(depth float64) float64 {
	layerDepths := sp.getLayerDepths()
	fmt.Println(layerDepths)

	return depth
}
