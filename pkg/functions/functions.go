// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package functions

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/jrnd-io/jr/pkg/constants"
	"github.com/jrnd-io/jr/pkg/ctx"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FunctionsMap() template.FuncMap {
	return fmap
}

var Random = rand.New(rand.NewSource(0))
var data = map[string][]string{}
var fmap = map[string]interface{}{

	// text utilities
	"atoi":                     Atoi,
	"itoa":                     strconv.Itoa,
	"concat":                   func(a string, b string) string { return a + b },
	"counter":                  Counter,
	"first":                    func(s string) string { return s[:1] },
	"firstword":                func(s string) string { return strings.Split(s, " ")[0] },
	"from":                     Word,
	"from_at":                  WordAt,
	"from_shuffle":             WordShuffle,
	"from_n":                   WordShuffleN,
	"join":                     strings.Join,
	"len":                      Len,
	"lower":                    strings.ToLower,
	"lorem":                    Lorem,
	"markov":                   Nonsense,
	"random":                   func(s []string) string { return s[Random.Intn(len(s))] },
	"randoms":                  func(s string) string { a := strings.Split(s, "|"); return a[Random.Intn(len(a))] },
	"random_index":             RandomIndex,
	"random_string":            RandomString,
	"random_string_vocabulary": RandomStringVocabulary,
	"regex":                    Regex,
	"repeat":                   strings.Repeat,
	"replaceall":               strings.ReplaceAll,
	"sentence":                 Sentence,
	"sentence_prefix":          SentencePrefix,
	"squeeze":                  func(s string) string { return strings.ReplaceAll(s, " ", "") },
	"squeezechars":             func(s, c string) string { return strings.ReplaceAll(s, c, "") },
	"split":                    strings.Split,
	"substr":                   func(start, length int, s string) string { return s[start:length] },
	"trim":                     strings.TrimSpace,
	"trimchars":                strings.Trim,
	"title":                    cases.Title(language.English).String,
	"upper":                    strings.ToUpper,

	// math utilities
	"add":          func(a, b int) int { return a + b },
	"div":          func(a, b int) int { return a / b },
	"format_float": func(f string, v float32) string { return fmt.Sprintf(f, v) },
	"integer":      func(min, max int) int { return min + Random.Intn(max-min) },
	"integer64":    func(min, max int64) int64 { return min + Random.Int63n(max-min) },
	"floating":     func(min, max float32) float32 { return min + Random.Float32()*(max-min) },
	"sub":          func(a, b int) int { return a - b },
	"max":          math.Max,
	"min":          math.Min,
	"minint":       Minint,
	"maxint":       Maxint,
	"mod":          func(a, b int) int { return a % b },
	"mul":          func(a, b int) int { return a * b },

	// networking and time utilities
	"http_method":       HttpMethod,
	"ip":                Ip,
	"ipv6":              Ipv6,
	"ip_known_protocol": IpKnownProtocol,
	"ip_known_port":     IpKnownPort,
	"mac":               Mac,
	"password":          Password,
	"useragent":         UserAgent,

	// people related utilities
	"cf":             CodiceFiscale,
	"company":        Company,
	"email":          Email,
	"email_provider": EmailProvider,
	"email_work":     WorkEmail,
	"gender":         Gender,
	"middlename":     Middlename,
	"name":           Name,
	"name_m":         NameM,
	"name_f":         NameF,
	"ssn":            Ssn,
	"surname":        Surname,
	"user":           User,
	"username":       Username,

	// address
	"building":                              BuildingNumber,
	"cardinal":                              Cardinal,
	"capital":                               Capital,
	"capital_at":                            CapitalAt,
	"city":                                  City,
	"city_at":                               CityAt,
	"country":                               Country,
	"country_random":                        CountryRandom,
	"country_at":                            CountryAt,
	"latitude":                              Latitude,
	"longitude":                             Longitude,
	"nearby_gps":                            NearbyGPS,
	"nearby_gps_into_polygon":               NearbyGPSIntoPolygon,
	"nearby_gps_into_polygon_without_start": NearbyGPSIntoPolygonWithoutStart,
	"nearby_gps_on_polyline":                NearbyGPSOnPolyline,
	"state":                                 State,
	"state_at":                              StateAt,
	"state_short":                           StateShort,
	"state_short_at":                        StateShortAt,
	"street":                                Street,
	"zip":                                   Zip,
	"zip_at":                                ZipAt,

	// finance
	"account":      Account,
	"amount":       Amount,
	"bitcoin":      Bitcoin,
	"card":         CreditCard,
	"cardCVV":      CreditCardCVV,
	"cusip":        Cusip,
	"ethereum":     Ethereum,
	"isin":         Isin,
	"sedol":        Sedol,
	"stock_symbol": StockSymbol,
	"swift":        Swift,
	"valor":        Valor,
	"wkn":          Wkn,

	// time and dates
	"birthdate":        BirthDate,
	"date_between":     DateBetween,
	"dates_between":    DatesBetween,
	"future":           Future,
	"past":             Past,
	"recent":           Recent,
	"just_passed":      Justpassed,
	"format_timestamp": FormatTimestamp,
	"now":              Now,
	"now_sub":          Nowsub,
	"now_add":          Nowadd,
	"soon":             Soon,
	"unix_time_stamp":  UnixTimeStamp,

	// phone
	"country_code":    CountryCode,
	"country_code_at": CountryCodeAt,
	"imei":            Imei,
	"phone":           Phone,
	"phone_at":        PhoneAt,
	"mobile_phone":    MobilePhone,
	"mobile_phone_at": MobilePhoneAt,

	// generic utilities
	"array":    func(count int) []int { return make([]int, count) },
	"bool":     RandomBool,
	"image":    Image,
	"image_of": ImageOf,
	"index_of": IndexOf,
	"key":      func(name string, n int) string { return fmt.Sprintf("%s%d", name, Random.Intn(n)) },
	"seed":     Seed,
	"uuid":     UniqueId,
	"yesorno":  YesOrNo,
	"inject":   Inject,

	// context utilities
	"add_v_to_list":            AddValueToList,
	"random_v_from_list":       RandomValueFromList,
	"random_n_v_from_list":     RandomNValuesFromList,
	"get_v_from_list_at_index": GetValueFromListAtIndex,
	"get_v":                    GetV,
	"set_v":                    SetV,
	"fromcsv":                  FromCsv,
}

func Atoi(s string) int {
	if len(s) == 0 {
		return 0
	}

	i, _ := strconv.Atoi(s)
	return i
}

// Seed sets seeds and can be used in a template
func Seed(rndSeed int64) string {
	SetSeed(rndSeed)
	return ""
}

// SetSeed sets seeds for all random JR objects
func SetSeed(rndSeed int64) {
	Random.Seed(rndSeed)
	uuid.SetRand(Random)
}

// AddValueToList adds value v to Context list l
func AddValueToList(l string, v string) string {
	ctx.JrContext.CtxListLock.Lock()
	defer ctx.JrContext.CtxListLock.Unlock()
	ctx.JrContext.CtxList[l] = append(ctx.JrContext.CtxList[l], v)
	return ""
}

// GetV gets value s from Context
func GetV(s string) string {
	ctx.JrContext.CtxLock.RLock()
	defer ctx.JrContext.CtxLock.RUnlock()
	return ctx.JrContext.Ctx[s]
}

// SetV adds value v to Context
func SetV(s string, v string) string {
	ctx.JrContext.CtxLock.Lock()
	defer ctx.JrContext.CtxLock.Unlock()
	ctx.JrContext.Ctx[s] = v
	return ""
}

// IndexOf returns the index of the s string in a file
func IndexOf(s string, name string) int {
	_, err := Cache(name)
	if err != nil {
		return -1
	}
	words := data[name]
	index := sort.Search(len(words), func(i int) bool { return strings.ToLower(words[i]) >= strings.ToLower(s) })

	if index < len(words) && words[index] == s {
		return index
	}

	return -1
}

// Len returns number of words (lines) in a word file
func Len(name string) string {
	_, err := Cache(name)
	if err != nil {
		return ""
	}
	l := len(data[name])
	return strconv.Itoa(l)
}

// RandomIndex returns a random index in a word file
func RandomIndex(name string) string {
	_, err := Cache(name)
	if err != nil {
		return ""
	}
	words := data[name]
	ctx.JrContext.LastIndex = Random.Intn(len(words))
	return strconv.Itoa(ctx.JrContext.LastIndex)
}

// RandomValueFromList returns a random value from Context list l
func RandomValueFromList(s string) string {
	ctx.JrContext.CtxListLock.RLock()
	defer ctx.JrContext.CtxListLock.RUnlock()
	list := ctx.JrContext.CtxList[s]
	l := len(list)
	if l != 0 {
		return list[Random.Intn(l)]
	}

	return ""
}

// GetValuesFromList returns a value from Context list l at index
func GetValueFromListAtIndex(s string, index int) string {

	ctx.JrContext.CtxListLock.RLock()
	defer ctx.JrContext.CtxListLock.RUnlock()
	list := ctx.JrContext.CtxList[s]
	l := len(list)
	if l != 0 && index < l {
		return list[index]
	}

	return ""
}

// RandomNValuesFromList returns a random value from Context list l
func RandomNValuesFromList(s string, n int) []string {
	ctx.JrContext.CtxListLock.RLock()
	defer ctx.JrContext.CtxListLock.RUnlock()
	list := ctx.JrContext.CtxList[s]
	l := len(list)
	if l != 0 {
		ints := findNDifferentInts(n, l)
		results := make([]string, len(ints))
		for i := range ints {
			results[i] = list[i]
		}
		return results
	}

	return []string{""}
}

// Word returns a random string from a list of strings in a file.
func Word(name string) string {
	_, err := Cache(name)
	if err != nil {
		return ""
	}
	words := data[name]
	ctx.JrContext.LastIndex = Random.Intn(len(words))
	return words[ctx.JrContext.LastIndex]
}

// WordAt returns a string at a given position in a list of strings in a file.
func WordAt(name string, index int) string {
	_, err := Cache(name)
	if err != nil {
		return ""
	}
	words := data[name]
	return words[index]
}

// WordShuffle returns a shuffled list of strings in a file.
func WordShuffle(name string) []string {
	_, err := Cache(name)
	if err != nil {
		return []string{""}
	}
	words := data[name]
	return WordShuffleN(name, len(words))
}

// wordShuffleN return a subset of n elements in a list of string in a file.
func WordShuffleN(name string, n int) []string {
	_, err := Cache(name)
	if err != nil {
		return []string{""}
	}
	words := data[name]
	Random.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	number := Minint(n, len(words))
	return words[:number]
}

// Minint returns the minimum between two ints
func Minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Maxint returns the minimum between two ints
func Maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Cache is used to internally Cache data from word files
func Cache(name string) (bool, error) {

	templateDir := fmt.Sprintf("%s/%s", constants.JR_SYSTEM_DIR, "templates")

	v := data[name]
	if v != nil {
		return false, nil
	}

	locale := strings.ToLower(ctx.JrContext.Locale)
	filename := fmt.Sprintf("%s/data/%s/%s", os.ExpandEnv(templateDir), locale, name)
	if locale != "us" && !(fileExists(filename)) {
		filename = fmt.Sprintf("%s/data/%s/%s", os.ExpandEnv(templateDir), "us", name)
	}
	data[name] = initialize(filename)
	if len(data[name]) == 0 {
		return false, fmt.Errorf("no words found in %s", filename)
	}

	return true, nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}

	return false
}

// Helper function to generate n different integers from 0 to length
func findNDifferentInts(n, max int) []int {

	n = Minint(n, max)
	ints := make([]int, n)

	// Generate n different random indices of maximum length
	for i := 0; i < n; {
		index := Random.Intn(max)
		if !contains(ints, index) {
			ints[i] = index
			i++
		}
	}

	return ints
}

// Helper function to check if an int is in a slice of ints
func contains(values []int, value int) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func initialize(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error in closing file")
		}
	}(file)

	var words []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words
}

func InitCSV(csvpath string) {
	// Loads the csv file in the context
	if len(csvpath) == 0 {
		return
	}

	var csvHeaders = make(map[int]string)
	csvValues := make(map[int]map[string]string)

	if _, err := os.Stat(csvpath); err != nil {
		println("File does not exist: ", csvpath)
		os.Exit(1)
	}

	file, err := os.Open(csvpath)

	if err != nil {
		println("Error opening file:", csvpath, "error:", err)
		os.Exit(1)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		println("Error reading CSV file:", csvpath, "error:", err)
		os.Exit(1)
	}

	for row := 0; row < len(lines); row++ {
		var aRow = lines[row]
		for col := 0; col < len(aRow); col++ {
			var value = aRow[col]
			// print(" ROW -> ",row)
			if row == 0 {
				// print(" H: ", value)
				csvHeaders[col] = strings.Trim(value, " ")
			} else {

				val, exists := csvValues[row-1]
				if exists {
					val[csvHeaders[col]] = strings.Trim(value, " ")
					csvValues[row-1] = val
				} else {
					var localmap = make(map[string]string)
					localmap[csvHeaders[col]] = strings.Trim(value, " ")
					csvValues[row-1] = localmap
				}
				//	print(" V: ", value)
			}
		}
		// println()
	}

	ctx.JrContext.CtxCSV = csvValues
}

func InitGeoJson(geojsonpath string) {
	if len(geojsonpath) == 0 {
		return
	}
	// Geometry represents a GeoJSON geometry object
	type Geometry struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}

	// Feature represents a GeoJSON feature object
	type Feature struct {
		Type       string                 `json:"type"`
		Geometry   Geometry               `json:"geometry"`
		Properties map[string]interface{} `json:"properties"`
	}

	// FeatureCollection represents a GeoJSON feature collection
	type FeatureCollection struct {
		Type     string    `json:"type"`
		Features []Feature `json:"features"`
	}

	var geometry *Geometry

	// Read the GeoJSON file
	data, err := os.ReadFile(geojsonpath)
	if err != nil {
		println("Error decoding GeoJson file:", geojsonpath, "error:", err)
		os.Exit(1)
	}

	// Check the type of GeoJSON object (FeatureCollection or single Feature)
	var geoJSONType struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &geoJSONType); err != nil {
		println("invalid GeoJSON format: %w", err)
		os.Exit(1)
	}

	switch geoJSONType.Type {
	case "FeatureCollection":
		// Parse as FeatureCollection
		var featureCollection FeatureCollection
		if err := json.Unmarshal(data, &featureCollection); err != nil {
			println("error parsing FeatureCollection: %w", err)
			os.Exit(1)
		}
		if len(featureCollection.Features) > 0 {
			geometry = &featureCollection.Features[0].Geometry
		}
	case "Feature":
		// Parse as single Feature
		var feature Feature
		if err := json.Unmarshal(data, &feature); err != nil {
			println("error parsing Feature: %w", err)
			os.Exit(1)
		}
		geometry = &feature.Geometry
	default:
		println("unsupported GeoJSON type: %w", geoJSONType.Type)
		os.Exit(1)
	}

	if geometry == nil {
		println("no geometry found in GeoJSON")
		os.Exit(1)
	}

	var coordinates [][]float64
	switch geometry.Type {
	case "Polygon":
		var points [][][]float64
		json.Unmarshal(geometry.Coordinates, &points)
		coordinates = points[0]
	case "LineString":
		var points [][]float64
		json.Unmarshal(geometry.Coordinates, &points)
		coordinates = points
	default:
		println("unsupported GeoJSON type: %w", geoJSONType.Type)
		os.Exit(1)
	}

	ctx.JrContext.CtxGeoJson = coordinates
}
