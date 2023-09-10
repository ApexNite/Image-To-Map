/*
Copyright © 2023 ApexNite <c8gkxkui@duck.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/klauspost/compress/zlib"

	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"
	"github.com/sqweek/dialog"
)

var ColorsMap = map[color.RGBA]string{
	{51, 112, 204, 255}:  "deep_ocean",
	{64, 132, 226, 255}:  "close_ocean",
	{85, 174, 240, 255}:  "shallow_waters",
	{247, 232, 152, 255}: "sand",
	{226, 147, 75, 255}:  "soil_low",
	{182, 111, 58, 255}:  "soil_high",
	{246, 45, 20, 255}:   "lava0",
	{255, 103, 0, 255}:   "lava1",
	{255, 172, 0, 255}:   "lava2",
	{255, 222, 0, 255}:   "lava3",
	{91, 94, 92, 255}:    "hills",
	{65, 69, 69, 255}:    "mountains",
	{126, 175, 70, 255}:  "soil_low:grass_low",
	{95, 131, 60, 255}:   "soil_high:grass_high",
	{240, 177, 33, 255}:  "soil_low:savanna_low",
	{207, 147, 27, 255}:  "soil_high:savanna_high",
	{140, 220, 106, 255}: "soil_low:enchanted_low",
	{118, 177, 83, 255}:  "soil_high:enchanted_high",
	{103, 118, 66, 255}:  "soil_low:mushroom_low",
	{85, 99, 56, 255}:    "soil_high:mushroom_high",
	{111, 85, 108, 255}:  "soil_low:corruption_low",
	{83, 63, 81, 255}:    "soil_high:corruption_high",
	{156, 54, 38, 255}:   "soil_low:infernal_low",
	{104, 55, 45, 255}:   "soil_high:infernal_high",
	{70, 160, 82, 255}:   "soil_low:jungle_low",
	{31, 112, 32, 255}:   "soil_high:jungle_high",
	{77, 72, 62, 255}:    "soil_low:swamp_low",
	{69, 62, 52, 255}:    "soil_high:swamp_high",
	{132, 147, 113, 255}: "soil_low:wasteland_low",
	{108, 119, 89, 255}:  "soil_high:wasteland_high",
	{232, 199, 110, 255}: "soil_low:desert_low",
	{225, 186, 90, 255}:  "soil_high:desert_high",
	{255, 150, 176, 255}: "soil_low:candy_low",
	{95, 214, 203, 255}:  "soil_high:candy_high",
	{255, 150, 176, 255}: "soil_low:crystal_low",
	{251, 135, 164, 255}: "soil_high:crystal_high",
	{209, 231, 113, 255}: "soil_low:lemon_low",
	{138, 207, 85, 255}:  "soil_high:lemon_high",
	{153, 188, 219, 255}: "soil_low:permafrost_low",
	{109, 0, 205, 255}:   "soil_low:water_bomb",
	{238, 81, 131, 255}:  "soil_low:tumor_low",
	{254, 24, 100, 255}:  "soil_high:tumor_high",
	{69, 200, 66, 255}:   "soil_low:biomass_low",
	{65, 168, 64, 255}:   "soil_high:biomass_high",
	{143, 147, 57, 255}:  "soil_low:pumpkin_low",
	{105, 108, 2, 255}:   "soil_high:pumpkin_high",
	{158, 166, 163, 255}: "soil_low:cybertile_low",
	{133, 136, 134, 255}: "soil_high:cybertile_high",
	{193, 153, 124, 255}: "soil_low:road",
	{131, 76, 76, 255}:   "soil_low:fuse",
	{168, 102, 58, 255}:  "soil_low:field",
	{163, 0, 0, 255}:     "soil_low:tnt",
	{180, 61, 204, 255}:  "soil_low:fireworks",
	{127, 0, 0, 255}:     "soil_low:tnt_timed",
	{153, 0, 0, 255}:     "soil_low:landmine",
	{186, 213, 211, 255}: "soil_low:frozen_low",
	{211, 228, 227, 255}: "soil_high:frozen_high",
	{175, 245, 241, 255}: "soil_low:snow_sand",
	{167, 214, 244, 255}: "soil_low:ice",
	{226, 237, 236, 255}: "soil_low:snow_hills",
	{252, 253, 253, 255}: "soil_low:snow_block",
}

var TilesMap = map[string]color.RGBA{
	"deep_ocean":      {51, 112, 204, 255},
	"close_ocean":     {64, 132, 226, 255},
	"shallow_waters":  {85, 174, 240, 255},
	"sand":            {247, 232, 152, 255},
	"soil_low":        {226, 147, 75, 255},
	"soil_high":       {182, 111, 58, 255},
	"lava0":           {246, 45, 20, 255},
	"lava1":           {255, 103, 0, 255},
	"lava2":           {255, 172, 0, 255},
	"lava3":           {255, 222, 0, 255},
	"hills":           {91, 94, 92, 255},
	"mountains":       {65, 69, 69, 255},
	"grass_low":       {126, 175, 70, 255},
	"grass_high":      {95, 131, 60, 255},
	"savanna_low":     {240, 177, 33, 255},
	"savanna_high":    {207, 147, 27, 255},
	"enchanted_low":   {140, 220, 106, 255},
	"enchanted_high":  {118, 177, 83, 255},
	"mushroom_low":    {103, 118, 66, 255},
	"mushroom_high":   {85, 99, 56, 255},
	"corruption_low":  {111, 85, 108, 255},
	"corruption_high": {83, 63, 81, 255},
	"infernal_low":    {156, 54, 38, 255},
	"infernal_high":   {104, 55, 45, 255},
	"jungle_low":      {70, 160, 82, 255},
	"jungle_high":     {31, 112, 32, 255},
	"swamp_low":       {77, 72, 62, 255},
	"swamp_high":      {69, 62, 52, 255},
	"wasteland_low":   {132, 147, 113, 255},
	"wasteland_high":  {108, 119, 89, 255},
	"desert_low":      {232, 199, 110, 255},
	"desert_high":     {225, 186, 90, 255},
	"candy_low":       {255, 150, 176, 255},
	"candy_high":      {95, 214, 203, 255},
	"crystal_low":     {255, 150, 176, 255},
	"crystal_high":    {251, 135, 164, 255},
	"lemon_low":       {209, 231, 113, 255},
	"lemon_high":      {138, 207, 85, 255},
	"permafrost_low":  {153, 188, 219, 255},
	"permafrost_high": {180, 207, 229, 255},
	"water_bomb":      {109, 0, 205, 255},
	"tumor_low":       {238, 81, 131, 255},
	"tumor_high":      {254, 24, 100, 255},
	"biomass_low":     {69, 200, 66, 255},
	"biomass_high":    {65, 168, 64, 255},
	"pumpkin_low":     {143, 147, 57, 255},
	"pumpkin_high":    {105, 108, 2, 255},
	"cybertile_low":   {158, 166, 163, 255},
	"cybertile_high":  {133, 136, 134, 255},
	"road":            {193, 153, 124, 255},
	"fuse":            {131, 76, 76, 255},
	"field":           {168, 102, 58, 255},
	"tnt":             {163, 0, 0, 255},
	"fireworks":       {180, 61, 204, 255},
	"tnt_timed":       {127, 0, 0, 255},
	"landmine":        {153, 0, 0, 255},
	"frozen_low":      {186, 213, 211, 255},
	"frozen_high":     {211, 228, 227, 255},
	"snow_sand":       {175, 245, 241, 255},
	"ice":             {167, 214, 244, 255},
	"snow_hills":      {226, 237, 236, 255},
	"snow_block":      {252, 253, 253, 255},
}

type Configuration struct {
	Algorithm string
	Included  []string
	Strength  float32
}

type MapData struct {
	SaveVersion   int       `json:"saveVersion"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	MapStats      MapStats  `json:"mapStats"`
	WorldLaws     WorldLaws `json:"worldLaws"`
	TileMap       []string  `json:"tileMap"`
	TileArray     [][]int   `json:"tileArray"`
	TileAmounts   [][]int   `json:"tileAmounts"`
	Fire          []string  `json:"fire"`
	ConwayEater   []string  `json:"conwayEater"`
	ConwayCreator []string  `json:"conwayCreator"`
	FrozenTiles   []string  `json:"frozen_tiles"`
	Tiles         []string  `json:"tiles"`
	Cities        []string  `json:"cities"`
	ActorsData    []string  `json:"actors_data"`
	Buildings     []string  `json:"buildings"`
	Kingdoms      []string  `json:"kingdoms"`
	Clans         []string  `json:"clans"`
	Alliances     []string  `json:"alliances"`
	Wars          []string  `json:"wars"`
	Plots         []string  `json:"plots"`
	Relations     []string  `json:"relations"`
	Cultures      []string  `json:"cultures"`
}

type MapStats struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	WorldTime    int    `json:"worldTime"`
	EraID        string `json:"era_id"`
	EraNextID    string `json:"era_next_id"`
	EraMonthNext int    `json:"era_month_next"`
	Deaths       int    `json:"deaths"`
	DeathsOther  int    `json:"deaths_other"`
	IDUnit       int    `json:"id_unit"`
	IDBuilding   int    `json:"id_building"`
}

type WorldLaws struct {
	List []List `json:"list"`
}

type List struct {
	Name    string `json:"name"`
	BoolVal *bool  `json:"boolVal,omitempty"`
}

func main() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		log.Fatal(err)
	}

	palette := []color.Color{}
	imagePath, _ := dialog.File().Filter("Images (*.jpg;*.jpeg;*.png)", "jpg", "jpeg", "png").Load()

	if len(configuration.Included) > 0 {
		for _, v := range configuration.Included {
			if findColor(palette, TilesMap[v]) == -1 {
				palette = append(palette, TilesMap[v])
			}
		}
	}

	img, err := getImageFromFilePath(imagePath)

	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	width, height := (bounds.Max.X/64)*64, (bounds.Max.Y/64)*64
	img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	img = ditherImage(img, palette, configuration.Algorithm, configuration.Strength)
	out, err := os.Create("./preview.png")

	if err != nil {
		log.Fatal(err)
	}

	png.Encode(out, img)
	defer out.Close()

	tileArray, tileAmounts, tileMap := readImage(img)

	f := new(bool)
	*f = false

	mapData := MapData{
		SaveVersion: 13,
		Width:       width / 64,
		Height:      height / 64,
		MapStats: MapStats{
			Name:         "BigBot's Inauspicious Kingdom",
			EraID:        "age_hope",
			EraMonthNext: 3000,
		},
		WorldLaws: WorldLaws{
			[]List{
				{Name: "world_law_diplomacy"},
				{Name: "world_law_peaceful_monsters", BoolVal: f},
				{Name: "world_law_hunger"},
				{Name: "world_law_vegetation_random_seeds"},
				{Name: "world_law_vegetation_seeds"},
				{Name: "world_law_grow_minerals"},
				{Name: "world_law_grow_grass"},
				{Name: "world_law_biome_overgrowth"},
				{Name: "world_law_kingdom_expansion"},
				{Name: "world_law_old_age"},
				{Name: "world_law_animals_spawn"},
				{Name: "world_law_animals_babies"},
				{Name: "world_law_rebellions"},
				{Name: "world_law_border_stealing"},
				{Name: "world_law_erosion"},
				{Name: "world_law_forever_lava", BoolVal: f},
				{Name: "world_law_disasters_nature"},
				{Name: "world_law_disasters_other"},
				{Name: "world_law_angry_civilians", BoolVal: f},
				{Name: "world_law_civ_babies"},
				{Name: "world_law_forever_tumor_creep", BoolVal: f},
				{Name: "world_law_civ_army"},
				{Name: "world_law_civ_limit_population_100", BoolVal: f},
				{Name: "age_hope"},
				{Name: "age_sun"},
				{Name: "age_dark"},
				{Name: "age_tears"},
				{Name: "age_moon"},
				{Name: "age_chaos"},
				{Name: "age_wonders"},
				{Name: "age_ice"},
				{Name: "age_ash"},
				{Name: "age_despair"},
			},
		},
		TileMap:       tileMap,
		TileArray:     tileArray,
		TileAmounts:   tileAmounts,
		Fire:          []string{},
		ConwayEater:   []string{},
		ConwayCreator: []string{},
		FrozenTiles:   []string{},
		Tiles:         []string{},
		Cities:        []string{},
		ActorsData:    []string{},
		Buildings:     []string{},
		Kingdoms:      []string{},
		Clans:         []string{},
		Alliances:     []string{},
		Wars:          []string{},
		Plots:         []string{},
		Relations:     []string{},
		Cultures:      []string{},
	}

	jsonData, err := json.Marshal(mapData)

	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer
	writer, err := zlib.NewWriterLevel(&buffer, zlib.BestCompression)

	if err != nil {
		log.Fatal(err)
	}

	writer.Write(jsonData)
	writer.Close()
	os.WriteFile("./map.wbox", []byte(buffer.String()), 0644)
}

func ditherImage(img image.Image, palette []color.Color, algorithm string, strength float32) image.Image {
	d := dither.NewDitherer(palette)

	switch algorithm {
	case "FloydSteinberg":
		d.Matrix = dither.ErrorDiffusionStrength(dither.FloydSteinberg, strength)
	case "FalseFloydSteinberg":
		d.Matrix = dither.ErrorDiffusionStrength(dither.FalseFloydSteinberg, strength)
	case "Atkinson":
		d.Matrix = dither.ErrorDiffusionStrength(dither.Atkinson, strength)
	case "Stucki":
		d.Matrix = dither.ErrorDiffusionStrength(dither.Stucki, strength)
	case "Burkes":
		d.Matrix = dither.ErrorDiffusionStrength(dither.Burkes, strength)
	case "JarvisJudiceNinke":
		d.Matrix = dither.ErrorDiffusionStrength(dither.JarvisJudiceNinke, strength)
	case "Simple2D":
		d.Matrix = dither.ErrorDiffusionStrength(dither.Simple2D, strength)
	case "StevenPigeon":
		d.Matrix = dither.ErrorDiffusionStrength(dither.StevenPigeon, strength)
	case "Sierra":
		d.Matrix = dither.ErrorDiffusionStrength(dither.Sierra, strength)
	case "Sierra2":
		d.Matrix = dither.ErrorDiffusionStrength(dither.Sierra2, strength)
	case "SierraLite":
		d.Matrix = dither.ErrorDiffusionStrength(dither.SierraLite, strength)
	default:
		log.Println("Unknown algorithm, switching to default")
		d.Matrix = dither.ErrorDiffusionStrength(dither.SierraLite, strength)
	}

	return d.Dither(img)
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func contains(strSlice []string, str string) bool {
	for _, v := range strSlice {
		if v == str {
			return true
		}
	}

	return false
}

func findColor(colorSlice []color.Color, color color.Color) int {
	for i, v := range colorSlice {
		if v == color {
			return i
		}
	}

	return -1
}

func readImage(img image.Image) ([][]int, [][]int, []string) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	rgba := image.NewRGBA(bounds)
	tileMap := []string{}
	tiles := map[string]int{}
	tileArray := make([][]int, height)
	tileAmounts := make([][]int, height)
	tileCount := 0

	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	for y := 0; y < height; y++ {
		var previousNum int
		oppositeIndex := height - 1 - y
		tilesAmount := 0

		for x := 0; x < width; x++ {
			index := (y*width + x) * 4
			pix := rgba.Pix[index : index+4]
			tile := ColorsMap[color.RGBA{pix[0], pix[1], pix[2], pix[3]}]

			if !contains(tileMap, tile) {
				tileMap = append(tileMap, tile)
				tiles[tile] = tileCount
				tileCount++
			}

			num := tiles[tile]

			if x == 0 || previousNum == num {
				tilesAmount++
			} else {
				tileArray[oppositeIndex] = append(tileArray[oppositeIndex], previousNum)
				tileAmounts[oppositeIndex] = append(tileAmounts[oppositeIndex], tilesAmount)
				tilesAmount = 1
			}

			if x == width-1 {
				tileArray[oppositeIndex] = append(tileArray[oppositeIndex], num)
				tileAmounts[oppositeIndex] = append(tileAmounts[oppositeIndex], tilesAmount)
			}

			previousNum = num
		}
	}

	return tileArray, tileAmounts, tileMap
}
