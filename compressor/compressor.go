package compressor

import (
	"time"

	"github.com/sam-anthony/volute/mass"
	"github.com/sam-anthony/volute/util"
)

type Compressor struct {
	Name     string
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

var compressors = make(map[string]map[string]map[string]Compressor)

func init() {
	compressors["Garrett"] = make(map[string]map[string]Compressor)
	compressors["Garrett"]["G"] = make(map[string]Compressor)
	compressors["Garrett"]["G"]["25-660"] = garrettG25660()

	compressors["BorgWarner"] = make(map[string]map[string]Compressor)
	compressors["BorgWarner"]["K"] = make(map[string]Compressor)
	compressors["BorgWarner"]["K"]["03"] = borgwarnerK03()
}

func Compressors() map[string]map[string]map[string]Compressor {
	return compressors
}

func garrettG25660() Compressor {
	maxFlow, err := mass.NewFlowRate(
		mass.Mass{70, mass.Pound},
		time.Minute,
		mass.PoundsPerMinute,
	)
	util.Check(err)
	return Compressor{
		"Garrett G25-660",
		"compressor/res/garrett/g/25-660.jpg",
		204,
		1885,
		1665,
		25,
		maxFlow,
		4.0,
	}
}

func borgwarnerK03() Compressor {
	maxFlow, err := mass.NewFlowRate(
		mass.Mass{0.13, mass.Kilogram},
		time.Second,
		mass.KilogramsPerSecond,
	)
	util.Check(err)
	return Compressor{
		"BorgWarner K03",
		"compressor/res/borgwarner/k/03.jpg",
		30,
		714,
		876,
		4,
		maxFlow,
		2.8,
	}
}
