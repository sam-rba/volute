package compressor

import (
	"github.com/BurntSushi/toml"
	"io/fs"
	fp "path/filepath"

	"github.com/sam-anthony/volute/mass"
	"github.com/sam-anthony/volute/util"
)

const root = "compressor/res/"

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
	// MaxPR is the pressure ratio at MaxY.
	MaxPR float32
}

// [manufacturer][series][model]
var compressors = make(map[string]map[string]map[string]Compressor)

func init() {
	// Walk root, looking for .toml files describing a compressor.
	// Parse these toml files, create a Compressor and add it to compressors.
	err := fp.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if fp.Ext(path) != ".toml" {
			return nil
		}

		path = path[len(root):]

		// manufacturer/series, model
		manSer, mod := fp.Split(path)
		manSer = fp.Clean(manSer)         // Clean trailing slash
		mod = mod[:len(mod)-len(".toml")] // Trim .toml extension
		// manufacturer, series
		man, ser := fp.Split(manSer)
		man = fp.Clean(man) // Clean trailing slash

		var exists bool
		_, exists = compressors[man]
		if !exists {
			compressors[man] = make(map[string]map[string]Compressor)
		}
		_, exists = compressors[man][ser]
		if !exists {
			compressors[man][ser] = make(map[string]Compressor)
		}

		tomlFile := fp.Join(root, path)

		var c Compressor
		_, err = toml.DecodeFile(tomlFile, &c)
		if err != nil {
			return err
		}

		// Replace .toml with .jpg
		imageFile := tomlFile[:len(tomlFile)-len(".toml")] + ".jpg"
		c.FileName = imageFile

		// Must parse MaxFlow seperately because the MassFlowRateUnit
		// is stored as a string and must be converted with
		// FlowRateUnitFromString().
		flow := struct {
			FlowVal  float32
			FlowUnit string
		}{}
		_, err = toml.DecodeFile(tomlFile, &flow)
		if err != nil {
			return err
		}
		u, err := mass.FlowRateUnitFromString(flow.FlowUnit)
		if err != nil {
			return err
		}
		c.MaxFlow = mass.NewFlowRate(flow.FlowVal, u)

		compressors[man][ser][mod] = c

		return nil
	})
	util.Check(err)
}

func Compressors() map[string]map[string]map[string]Compressor {
	return compressors
}
