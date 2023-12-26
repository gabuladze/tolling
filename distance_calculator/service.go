package main

import (
	"math"

	"github.com/gabuladze/tolling/types"
)

type DistanceCalculator interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type DistanceCalculatorService struct {
	lastCoords map[int][2]float64 // obuID -> [lat, long]
}

func NewDistanceCalculatorService() DistanceCalculator {
	return &DistanceCalculatorService{
		lastCoords: map[int][2]float64{},
	}
}

func (dcs *DistanceCalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	_, ok := dcs.lastCoords[data.OBUID]
	if !ok {
		dcs.lastCoords[data.OBUID] = [2]float64{data.Lat, data.Long}
		return 0.0, nil
	}

	distance := calculateDistance(dcs.lastCoords[data.OBUID][0], dcs.lastCoords[data.OBUID][1], data.Lat, data.Long)
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow((x2-x1), 2) + math.Pow((y2-y1), 2))
}
