package data_structures

import (
	np "github.com/geoport/numpy4go/vectors"
	"reflect"
	"sort"
)

//GetLayerFields returns the fields of a soil layer
func (sp *SoilProfile) GetLayerFields() []string {
	var fields []string
	nonLayerFields := []string{"CheckGwt", "DensityUnit", "PressureUnit", "Gwt", "SPT", "ConeResistance", "PorePressure"}
	val := reflect.ValueOf(sp).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i).Name
		if !np.Contains(nonLayerFields, field) {
			fields = append(fields, field)
		}
	}
	return fields
}

//GetFieldProperties returns the values of the given field for each layer in the soil profile
func (sp SoilProfile) GetFieldProperties(field string) any {
	return reflect.ValueOf(sp).FieldByName(field).Interface()
}

//SetField sets the value of the given field for each layer in the soil profile
func (sp *SoilProfile) SetField(fieldName string, value interface{}) {
	v := reflect.ValueOf(sp).Elem()

	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		fieldNames[typeField.Name] = i
	}
	fieldNum := fieldNames[fieldName]

	fieldVal := v.Field(fieldNum)
	fieldVal.Set(reflect.ValueOf(value))

}

//GetLayerDepths returns the level of bottom of each layer in the soil profile
func (sp *SoilProfile) GetLayerDepths() []float64 {
	return np.Cumsum(sp.Thickness)
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
	if layerDepths[len(layerDepths)-1] < depth {
		return len(layerDepths) - 1
	}
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

// CalcNormalStress returns the normal stress at the given depth
func (sp *SoilProfile) CalcNormalStress(depth float64) float64 {
	Stresses := []float64{0}
	gammaDry := sp.DryUnitWeight
	gammaSaturated := sp.SaturatedUnitWeight
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

// IsCohesive returns true if the soil class of the layer at given depth is cohesive
func (sp *SoilProfile) IsCohesive(depth float64) bool {
	cohesiveSoils := []string{"SW-SC", "SP-SC", "SC", "SC-SM", "CL", "CL-ML", "CH,OH"}
	layerIndex := sp.GetLayerIndex(depth)
	soilClass := sp.SoilClass[layerIndex]
	for _, cohesiveSoil := range cohesiveSoils {
		if soilClass == cohesiveSoil {
			return true
		}
	}
	return false
}

// CombineSPT SPT log with soil profile
func (sp *SoilProfile) CombineSPT(sptLog SPTData) SoilProfile {
	sptDepth := sptLog.Depth
	N := sptLog.N
	layerDepths := sp.GetLayerDepths()
	var combinedDepths []float64
	combinedDepths = append(combinedDepths, layerDepths...)
	combinedDepths = np.Unique(append(combinedDepths, sptDepth...))
	sort.Sort(sort.Float64Slice(combinedDepths))
	layerFields := sp.GetLayerFields()
	newSoilProfile := SoilProfile{}
	maxSPTDepth := sptDepth[len(sptDepth)-1]
	sptProfile := SoilProfile{Thickness: sptDepth}
	for _, field := range layerFields {
		var oldFieldValuesString []string
		var oldFieldValuesFloat []float64
		var newFieldValuesString []string
		var newFieldValuesFloat []float64
		var newN []int
		isStringField := np.Contains([]string{"SoilClass", "SoilType", "SoilDefinition", "MaterialType"}, field)
		if isStringField {
			oldFieldValuesString = sp.GetFieldProperties(field).([]string)
		} else {
			oldFieldValuesFloat = sp.GetFieldProperties(field).([]float64)
		}

		for i := range combinedDepths {
			var depthPrev float64
			if i == 0 {
				depthPrev = 0
			} else {
				depthPrev = combinedDepths[i-1]
			}
			depthCurrent := combinedDepths[i]
			if depthCurrent <= maxSPTDepth {
				if field == "Thickness" {
					newFieldValuesFloat = append(newFieldValuesFloat, depthCurrent-depthPrev)
				} else {
					layerIndex := sp.GetLayerIndex(depthCurrent)

					if isStringField {
						if len(oldFieldValuesString) > layerIndex {
							newFieldValuesString = append(newFieldValuesString, oldFieldValuesString[layerIndex])
						}
					} else {
						if len(oldFieldValuesFloat) > layerIndex {
							newFieldValuesFloat = append(newFieldValuesFloat, oldFieldValuesFloat[layerIndex])
						}
					}
				}
				sptIndex := sptProfile.GetLayerIndex(depthCurrent)
				newN = append(newN, N[sptIndex])
			}
		}

		if isStringField {
			newSoilProfile.SetField(field, newFieldValuesString)
		} else {
			newSoilProfile.SetField(field, newFieldValuesFloat)
		}
		newSoilProfile.SetField("SPT", newN)
	}
	return newSoilProfile
}

// CombineCPT CPT log with soil profile
func (sp *SoilProfile) CombineCPT(cptLog CPTData) SoilProfile {
	cptDepth := cptLog.Depth
	coneResistance := cptLog.ConeResistance
	layerDepths := sp.GetLayerDepths()
	var combinedDepths []float64
	combinedDepths = append(combinedDepths, layerDepths...)
	combinedDepths = np.Unique(append(combinedDepths, cptDepth...))
	sort.Sort(sort.Float64Slice(combinedDepths))
	layerFields := sp.GetLayerFields()
	newSoilProfile := SoilProfile{}
	maxCPTDepth := cptDepth[len(cptDepth)-1]

	cptProfile := SoilProfile{Thickness: cptDepth}
	for _, field := range layerFields {
		var oldFieldValuesString []string
		var oldFieldValuesFloat []float64
		var newFieldValuesString []string
		var newFieldValuesFloat []float64
		var newConeResistance []float64
		isStringField := np.Contains([]string{"SoilClass", "SoilType", "SoilDefinition", "MaterialType"}, field)
		if isStringField {
			oldFieldValuesString = sp.GetFieldProperties(field).([]string)
		} else {
			oldFieldValuesFloat = sp.GetFieldProperties(field).([]float64)
		}

		for i := range combinedDepths {
			var depthPrev float64
			if i == 0 {
				depthPrev = 0
			} else {
				depthPrev = combinedDepths[i-1]
			}
			depthCurrent := combinedDepths[i]
			if depthCurrent <= maxCPTDepth {
				if field == "Thickness" {
					newFieldValuesFloat = append(newFieldValuesFloat, depthCurrent-depthPrev)
				} else {
					layerIndex := sp.GetLayerIndex(depthCurrent)

					if isStringField {
						if len(oldFieldValuesString) > layerIndex {
							newFieldValuesString = append(newFieldValuesString, oldFieldValuesString[layerIndex])
						}
					} else {
						if len(oldFieldValuesFloat) > layerIndex {
							newFieldValuesFloat = append(newFieldValuesFloat, oldFieldValuesFloat[layerIndex])
						}
					}
				}
				cptIndex := cptProfile.GetLayerIndex(depthCurrent)
				newConeResistance = append(newConeResistance, coneResistance[cptIndex])
			}
		}

		if isStringField {
			newSoilProfile.SetField(field, newFieldValuesString)
		} else {
			newSoilProfile.SetField(field, newFieldValuesFloat)
		}
		newSoilProfile.SetField("ConeResistance", newConeResistance)
	}
	return newSoilProfile
}

// CombineVS VS log with soil profile
func (sp *SoilProfile) CombineVS(vsLog MASWData) SoilProfile {
	vsDepth := np.Cumsum(vsLog.Thickness)
	VS := vsLog.VS
	layerDepths := sp.GetLayerDepths()
	var combinedDepths []float64
	combinedDepths = append(combinedDepths, layerDepths...)
	combinedDepths = np.Unique(append(combinedDepths, vsDepth...))
	sort.Sort(sort.Float64Slice(combinedDepths))
	layerFields := sp.GetLayerFields()
	newSoilProfile := SoilProfile{}
	maxCPTDepth := vsDepth[len(vsDepth)-1]

	vsProfile := SoilProfile{Thickness: vsDepth}
	for _, field := range layerFields {
		var oldFieldValuesString []string
		var oldFieldValuesFloat []float64
		var newFieldValuesString []string
		var newFieldValuesFloat []float64
		var newVS []float64

		isStringField := np.Contains([]string{"SoilClass", "SoilType", "SoilDefinition", "MaterialType"}, field)

		if isStringField {
			oldFieldValuesString = sp.GetFieldProperties(field).([]string)
		} else {
			oldFieldValuesFloat = sp.GetFieldProperties(field).([]float64)
		}

		for i := range combinedDepths {
			var depthPrev float64
			if i == 0 {
				depthPrev = 0
			} else {
				depthPrev = combinedDepths[i-1]
			}
			depthCurrent := combinedDepths[i]
			if depthCurrent <= maxCPTDepth {
				if field == "Thickness" {
					newFieldValuesFloat = append(newFieldValuesFloat, depthCurrent-depthPrev)
				} else {
					layerIndex := sp.GetLayerIndex(depthCurrent)

					if isStringField {
						if len(oldFieldValuesString) > layerIndex {
							newFieldValuesString = append(newFieldValuesString, oldFieldValuesString[layerIndex])
						}
					} else {
						if len(oldFieldValuesFloat) > layerIndex {
							newFieldValuesFloat = append(newFieldValuesFloat, oldFieldValuesFloat[layerIndex])
						}
					}
				}
				vsIndex := vsProfile.GetLayerIndex(depthCurrent)
				newVS = append(newVS, VS[vsIndex])
			}
		}

		if isStringField {
			newSoilProfile.SetField(field, newFieldValuesString)
		} else {
			newSoilProfile.SetField(field, newFieldValuesFloat)
		}
		newSoilProfile.SetField("VS", newVS)
	}
	return newSoilProfile
}
