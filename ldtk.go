// Package ldtk implements LDTk's (https://ldtk.io/json) JSON schema in user-friendly format.
//
// Usage:
//	package main
//
//	import (
//		"json"
//		_ `embed`
//
//		"github.com/neputevshina/ldtk"
//	)
//
//	//go:embed ./levels.ldtk
//	var _project []byte
//
//	var Project ldtk.Project
//
//	func init() {
//		json.Unmarshal(_project, &Project)
//	}
//
//	func main() {
//		// your game
//	}
//
package ldtk

import (
	"encoding/json"
	"image/color"
)

// TODO: copy-paste comments from original schema from https://ldtk.io/json

type Project struct {
	WorldLayout     string
	WorldGridWidth  int
	WorldGridHeight int
	BgColor         color.RGBA
	ExternalLevels  bool
	Layers          []LayerDef
	Entities        map[string]EntityDef
	Tilesets        []TilesetDef
	Enums           []interface{} // TODO: implement when it's time to
	ExternalEnums   []interface{} // and this
	LevelFields     []interface{} // this one too
	Levels          map[int]Level
}

func (p *Project) UnmarshalJSON(data []byte) error {
	// this is done once in loading time, we can afford it
	// also we can probably ignore convention from json.Unmarshaler description
	var project struct {
		WorldLayout     string `json:"worldLayout"`
		WorldGridWidth  int    `json:"worldGridWidth"`
		WorldGridHeight int    `json:"worldGridHeight"`
		BgColor         string `json:"bgColor"`
		ExternalLevels  bool   `json:"externalLevels"`
		Defs            struct {
			Layers        []LayerDef    `json:"layers"`
			Entities      []EntityDef   `json:"entities"`
			Tilesets      []TilesetDef  `json:"tilesets"`
			Enums         []interface{} `json:"enums"`
			ExternalEnums []interface{} `json:"externalEnums"`
			LevelFields   []interface{} `json:"levelFields"`
		} `json:"defs"`
		Levels []Level `json:"levels"`
	}

	err := json.Unmarshal(data, &project)
	if err != nil {
		return err
	}
	p.Entities = make(map[string]EntityDef, len(project.Defs.Entities))
	for _, v := range project.Defs.Entities {
		p.Entities[v.Identifier] = v
	}

	p.Layers = project.Defs.Layers
	p.Tilesets = project.Defs.Tilesets
	p.Enums = project.Defs.Enums
	p.ExternalEnums = project.Defs.ExternalEnums
	p.LevelFields = project.Defs.LevelFields

	p.WorldLayout = project.WorldLayout
	p.WorldGridWidth = project.WorldGridWidth
	p.WorldGridHeight = project.WorldGridHeight
	p.BgColor = hex(project.BgColor)
	p.ExternalLevels = project.ExternalLevels

	p.Levels = make(map[int]Level, len(project.Levels))
	for _, v := range project.Levels {
		p.Levels[v.UID] = v
	}

	return nil
}

type IntGridValue struct {
	Value      int         `json:"value"`
	Identifier interface{} `json:"identifier"`
	Color      string      `json:"color"`
}

type LayerDef struct {
	Type                  string         `json:"__type"`
	Identifier            string         `json:"identifier"`
	UID                   int            `json:"uid"`
	GridSize              int            `json:"gridSize"`
	DisplayOpacity        int            `json:"displayOpacity"`
	PxOffsetX             int            `json:"pxOffsetX"`
	PxOffsetY             int            `json:"pxOffsetY"`
	IntGridValues         []IntGridValue `json:"intGridValues"`
	AutoTilesetDefUID     *int           `json:"autoTilesetDefUid"`
	AutoSourceLayerDefUID *int           `json:"autoSourceLayerDefUid"`
	TilesetDefUID         int            `json:"tilesetDefUid"`
}

type EntityDef struct {
	Identifier string  `json:"identifier"`
	UID        int     `json:"uid"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Color      string  `json:"color"` // TODO: this
	TilesetID  *int    `json:"tilesetId"`
	TileID     *int    `json:"tileId"`
	PivotX     float64 `json:"pivotX"`
	PivotY     float64 `json:"pivotY"`
}

type TilesetDef struct {
	Width      int    `json:"__cWid"`     // Grid-based height.
	Height     int    `json:"__cHei"`     // Grid-based width.
	Identifier string `json:"identifier"` // Unique string identifier.
	UID        int    `json:"uid"`        // Unique int identifier.

	// Path to the source file, relative to the current project JSON file.
	RelPath string `json:"relPath"`

	PixelWidth   int `json:"pxWid"`        // Image width in pixels.
	PixelHeight  int `json:"pxHei"`        // Image height in pixels.
	TileGridSize int `json:"tileGridSize"` // TODO: explain
	Spacing      int `json:"spacing"`      // Unique String identifier.
	Padding      int `json:"padding"`      // Distance in pixels from image borders.

	// Optional Enum definition UID used for this tileset metadata.
	TagsSourceEnumUID *int `json:"tagsSourceEnumUid"`

	// Tileset tags using Enum values specified by tagsSourceEnumId.
	// This array contains 1 element per Enum value, which contains
	// an array of all Tile IDs that are tagged with it.
	EnumTags []struct {
		EnumValueId string `json:"enumValueId"`
		TileIds     []int  `json:"tileIds"`
	} `json:"enumTags"`

	// An array of custom tile metadata.
	CustomData []struct {
		Data   string `json:"data"`
		TileId []int  `json:"tileId"`
	} `json:"customData"`
}

type Tile struct {
	Px  [2]int `json:"px"`
	Src [2]int `json:"src"`
	F   int    `json:"f"`
	T   int    `json:"t"`
}

type LayerInstance struct {
	Identifier         string        `json:"__identifier"`
	Type               string        `json:"__type"`
	CWid               int           `json:"__cWid"`
	CHei               int           `json:"__cHei"`
	GridSize           int           `json:"__gridSize"`
	Opacity            int           `json:"__opacity"`
	PxTotalOffsetX     int           `json:"__pxTotalOffsetX"`
	PxTotalOffsetY     int           `json:"__pxTotalOffsetY"`
	TilesetDefUID      int           `json:"__tilesetDefUid"`
	TilesetRelPath     string        `json:"__tilesetRelPath"`
	LevelID            int           `json:"levelId"`
	LayerDefUID        int           `json:"layerDefUid"`
	PxOffsetX          int           `json:"pxOffsetX"`
	PxOffsetY          int           `json:"pxOffsetY"`
	Visible            bool          `json:"visible"`
	IntGridCsv         []int         `json:"intGridCsv"`
	AutoLayerTiles     []interface{} `json:"autoLayerTiles"`
	OverrideTilesetUID int           `json:"overrideTilesetUid"`
	GridTiles          []Tile        `json:"gridTiles"`
	EntityInstances    []interface{} `json:"entityInstances"`
}

type Level struct {
	Identifier      string          `json:"identifier"`
	UID             int             `json:"uid"`
	WorldX          int             `json:"worldX"`
	WorldY          int             `json:"worldY"`
	PxWid           int             `json:"pxWid"`
	PxHei           int             `json:"pxHei"`
	BgColor         string          `json:"__bgColor"`
	BgPos           interface{}     `json:"__bgPos"`
	BgRelPath       *string         `json:"bgRelPath"`       // don't
	ExternalRelPath *string         `json:"externalRelPath"` // don't
	FieldInstances  []interface{}   `json:"fieldInstances"`  // TODO
	LayerInstances  []LayerInstance `json:"layerInstances"`
	Neighbours      []struct {
		UID       int    `json:"levelUid"`
		Direction string `json:"dir"`
	} `json:"__neighbours"`
}

func hex(s string) color.RGBA {
	if len(s) != 7 || s[0] != '#' {
		panic("incorrect color string '" + s + "'")
	}
	c := color.RGBA{}
	s = s[1:]
	nib := func(i int) uint8 {
		if s[i] >= 'A' && s[i] <= 'F' {
			return s[i] - 'A' + 10
		}
		return s[i] - '0'
	}
	c.R = nib(0)<<4 + nib(1)
	c.G = nib(2)<<4 + nib(3)
	c.B = nib(4)<<4 + nib(5)
	c.A = 0xff
	return c
}
