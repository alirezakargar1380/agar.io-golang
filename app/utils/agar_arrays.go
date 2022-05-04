package agar_arrays

import (
	"github.com/alirezakargar1380/agar.io-golang/app/trigonometric_circle"
)

type Agars struct {
	Agars []trigonometric_circle.AgarDe
}

func (agars *Agars) GETAgarIndexWithId(Id int) int {
	index := 0
	for i := 0; i < len(agars.Agars); i++ {
		if agars.Agars[i].Id == Id {
			index = i
		}
	}

	return index
}
