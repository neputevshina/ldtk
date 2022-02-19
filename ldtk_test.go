package ldtk_test

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/neputevshina/ldtk"
)

type object = map[string]interface{}
type array = []interface{}

// TestEntityMapping checks for correctness of conversion from entity array to map.
func TestEntityMapping(t *testing.T) {
	f, err := os.Open("./2021.ldtk")
	if err != nil {
		t.Fatal(err)
	}
	jd := json.NewDecoder(f)
	proj1 := ldtk.Project{}
	err = jd.Decode(&proj1)
	if err != nil {
		t.Fatal("json is corrupted")
	}
	proj2 := object{}
	f.Seek(0, 0)
	_ = json.NewDecoder(f).Decode(&proj2)
	_ = proj2["defs"].(object)

	for _, v := range proj2["defs"].(object)["entities"].(array) {
		assert := func(c bool) {
			if !c {
				t.Fail()
			}
		}
		obj := v.(object)
		k := obj["identifier"].(string)
		ent := proj1.EntityDefs[k]

		assert(ent.Color == obj["color"].(string))
		assert(ent.Identifier == obj["identifier"].(string))
		assert(ent.Height == int(obj["height"].(float64)))
		assert(ent.Width == int(obj["width"].(float64)))
		assert(ent.PivotX == obj["pivotX"].(float64))
		assert(ent.PivotY == obj["pivotY"].(float64))
		if ent.TileID == nil {
			assert(obj["tileId"] == nil)
		} else {
			assert(*ent.TileID == int(obj["tileId"].(float64)))
		}
		if ent.TilesetID == nil {
			assert(obj["tilesetId"] == nil)
		} else {
			assert(*ent.TilesetID == int(obj["tilesetId"].(float64)))
		}
		assert(ent.UID == int(obj["uid"].(float64)))
	}
}

func TestColor(t *testing.T) {
	f, err := os.Open("./2021.ldtk")
	if err != nil {
		t.Fatal(err)
	}
	jd := json.NewDecoder(f)
	proj1 := ldtk.Project{}
	err = jd.Decode(&proj1)
	if err != nil {
		t.Fatal("json is corrupted")
	}
	proj2 := object{}
	f.Seek(0, 0)
	_ = json.NewDecoder(f).Decode(&proj2)
	_ = proj2["defs"].(object)

	conv := func(b byte) string {
		return strconv.FormatUint(uint64(b), 16)
	}
	c := proj1.BgColor
	c2 := "#" + conv(c.R) + conv(c.G) + conv(c.B)
	if strings.ToLower(proj2["bgColor"].(string)) != c2 {
		t.Error("expected", proj2["bgColor"].(string)+",", "got", c2)
	}
}
