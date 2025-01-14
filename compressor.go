package main

import (
	"github.com/BurntSushi/toml"
	"io/fs"
	fp "path/filepath"
	"strings"
)

const root = "compressor_maps/"

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
	MaxFlow MassFlowRate
	// MaxPR is the pressure ratio at MaxY.
	MaxPR float32
}

// [manufacturer][series][model]
var Compressors = make(map[string]map[string]map[string]Compressor)

func init() {
	// Walk root, looking for .toml files describing a compressor.
	// Parse these toml files, create a Compressor and add it to Compressors.
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

		if _, ok := Compressors[man]; !ok { // Manufacturer does NOT exist
			Compressors[man] = make(map[string]map[string]Compressor)
		}
		if _, ok := Compressors[man][ser]; !ok { // Series does NOT exist
			Compressors[man][ser] = make(map[string]Compressor)
		}

		tomlFile := fp.Join(root, path)
		var c Compressor
		if _, err = toml.DecodeFile(tomlFile, &c); err != nil {
			return err
		}
		c.FileName = strings.TrimSuffix(tomlFile, ".toml") + ".jpg"
		c.MaxFlow, err = readMaxFlow(tomlFile)
		if err != nil {
			return err
		}
		Compressors[man][ser][mod] = c
		return nil
	})
	Check(err)
}

func readMaxFlow(tomlFile string) (MassFlowRate, error) {
	flow := struct {
		FlowVal  float32
		FlowUnit string
	}{}
	if _, err := toml.DecodeFile(tomlFile, &flow); err != nil {
		return -1, err
	}
	unit, err := ParseMassFlowRateUnit(flow.FlowUnit)
	if err != nil {
		return -1, err
	}
	return MassFlowRate(flow.FlowVal) * unit, nil
}
