package compressor

import (
	"time"

	"github.com/sam-anthony/volute/mass"
	"github.com/sam-anthony/volute/util"
)

type Compressor struct {
	FileName string
	// MinX is the distance of the y-axis from left of image in pixels.
	MinX int
	// MinY is the distance of the x-axis from the top of the image in
	//pixels.
	MinY int
	// MaxX is the distance of the end of the graph from the left of the
	// image in pixels.
	MaxX int
	// MaxY is the distance of the top of the graph from the top of the
	// image in pixels.
	MaxY int
	// MaxFlow is the mass flow rate at MaxX.
	MaxFlow mass.FlowRate
	// MaxPressureRatio is the pressure ratio at MaxY.
	MaxPressureRatio float32
}

func GarrettG25660() Compressor {
	maxFlow, err := mass.NewFlowRate(
		mass.Mass{70, mass.Pound},
		time.Minute,
		mass.PoundsPerMinute,
	)
	util.Check(err)
	return Compressor{
		"compressor/res/GarrettG25660.jpg",
		204,
		1885,
		1665,
		25,
		maxFlow,
		4.0,
	}
}
