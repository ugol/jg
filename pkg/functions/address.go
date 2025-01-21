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
	"fmt"
	"math"
	"os"

	"github.com/cnkei/gospline"
	"github.com/jrnd-io/jr/pkg/ctx"
)

const (
	earthRadius     = 6371000 // in meters
	degreesPerMeter = 1.0 / earthRadius * 180.0 / math.Pi
)

// BuildingNumber generates a random building number of max n digits
func BuildingNumber(n int) string {
	building := make([]byte, Random.Intn(n)+1)
	for i := range building {
		building[i] = digits[Random.Intn(len(digits))]
	}
	return string(building)
}

// Capital returns a random Capital
func Capital() string {
	return Word("capital")
}

// CapitalAt returns Capital at given index
func CapitalAt(index int) string {
	return WordAt("capital", index)
}

// Cardinal return a random cardinal direction, in long or short form
func Cardinal(short bool) string {
	if short {
		directions := []string{"N", "S", "E", "O", "NE", "NO", "SE", "SO"}
		return directions[Random.Intn(len(directions))]
	}

	directions := []string{"North", "South", "East", "Ovest", "North-East", "North-Ovest", "South-East", "South-Ovest"}
	return directions[Random.Intn(len(directions))]
}

// City returns a random City
func City() string {
	c := Word("city")
	ctx.JrContext.Ctx["_city"] = c
	ctx.JrContext.CityIndex = ctx.JrContext.LastIndex
	return c
}

// CityAt returns City at given index
func CityAt(index int) string {
	return WordAt("city", index)
}

// Country returns the ISO 3166 Country selected with locale
func Country() string {
	countryIndex := ctx.JrContext.CountryIndex
	if countryIndex == -1 {
		return Word("country")
	}

	return WordAt("country", countryIndex)
}

// CountryRandom returns a random ISO 3166 Country
func CountryRandom() string {
	return Word("country")
}

// CountryAt returns an ISO 3166 Country at a given index
func CountryAt(index int) string {
	return WordAt("country", index)
}

// Latitude returns a random latitude between -90 and 90
func Latitude() string {
	latitude := -90 + Random.Float64()*(180)
	return fmt.Sprintf("%.4f", latitude)
}

// Longitude returns a random longitude between -180 and 180
func Longitude() string {
	longitude := -180 + Random.Float64()*(360)
	return fmt.Sprintf("%.4f", longitude)
}

// NearbyGPS returns a random latitude longitude within a given radius in meters
func NearbyGPS(latitude float64, longitude float64, radius int) string {
	radiusInMeters := float64(radius)

	// Generate a random angle in radians
	randomAngle := Random.Float64() * 2 * math.Pi

	// Calculate the distance from the center point
	distanceInMeters := Random.Float64() * radiusInMeters

	// Convert the distance to degrees
	distanceInDegrees := distanceInMeters * degreesPerMeter

	// Calculate the new latitude and longitude
	newLatitude := latitude + (distanceInDegrees * math.Cos(randomAngle))
	newLongitude := longitude + (distanceInDegrees * math.Sin(randomAngle))

	return fmt.Sprintf("%.4f %.4f", newLatitude, newLongitude)
}

// NearbyGPSOnPolyline generates a random latitude and longitude within a specified radius (in meters)
// along a given polyline path. Upon reaching the end of the path, the function reverses direction and continues
// generating points along the same path, allowing for continuous point generation in a back-and-forth pattern.
func NearbyGPSOnPolyline(radius int) string {
	ctx.JrContext.CtxGeoJsonLock.Lock()
	defer ctx.JrContext.CtxGeoJsonLock.Unlock()

	// Ensure the path is available and has enough points
	if len(ctx.JrContext.CtxGeoJson) < 2 {
		println("Path must contain at least two points.")
		os.Exit(1)
	}

	// Get the current point on the path
	currentPoint := ctx.JrContext.CtxGeoJson[0]
	currentLat, currentLon := currentPoint[0], currentPoint[1]

	// Update last known point if there are recent saved coordinates
	if len(ctx.JrContext.CtxLastPointLat) == 1 {
		currentLat = ctx.JrContext.CtxLastPointLat[len(ctx.JrContext.CtxLastPointLat)-1]
	}
	if len(ctx.JrContext.CtxLastPointLon) == 1 {
		currentLon = ctx.JrContext.CtxLastPointLon[len(ctx.JrContext.CtxLastPointLon)-1]
	}

	// Convert radius to float for calculations
	radiusInMeters := float64(radius)

	// Find the next point on the polyline at the specified distance
	nextPoint, nuovoIndex, newDirection := findNextPoint(ctx.JrContext.CtxGeoJson, []float64{currentLat, currentLon}, ctx.JrContext.CtxForward, ctx.JrContext.CtxIndex, radiusInMeters)
	ctx.JrContext.CtxForward = newDirection
	ctx.JrContext.CtxIndex = nuovoIndex

	// Update the context with the new valid point, maintaining a maximum ctx of 10 points
	ctx.JrContext.CtxLastPointLat = append(ctx.JrContext.CtxLastPointLat, nextPoint[0])
	ctx.JrContext.CtxLastPointLon = append(ctx.JrContext.CtxLastPointLon, nextPoint[1])
	// Keep only the last point in the ctx
	if len(ctx.JrContext.CtxLastPointLat) > 1 {
		ctx.JrContext.CtxLastPointLat = ctx.JrContext.CtxLastPointLat[1:]
	}
	if len(ctx.JrContext.CtxLastPointLon) > 1 {
		ctx.JrContext.CtxLastPointLon = ctx.JrContext.CtxLastPointLon[1:]
	}

	// Return the coordinates of the valid point
	return fmt.Sprintf("%.12f %.12f", nextPoint[0], nextPoint[1])
}

// distanceEuclidean calculates the Euclidean distance between two points (p1 and p2) represented as []float64.
// It returns the straight-line distance between the points in meters.
func distanceEuclidean(p1, p2 []float64) float64 {
	dx := p2[0] - p1[0]             // Difference in the x-coordinates
	dy := p2[1] - p1[1]             // Difference in the y-coordinates
	return math.Sqrt(dx*dx + dy*dy) // Apply the Euclidean distance formula
}

// distanceHaversine calculates the Haversine distance between two points (p1 and p2) on the Earth's surface,
// given their latitude and longitude coordinates in degrees. It returns the distance in meters.
func distanceHaversine(p1, p2 []float64) float64 {
	lat1 := p1[0] * math.Pi / 180               // Convert latitude of point 1 from degrees to radians
	lat2 := p2[0] * math.Pi / 180               // Convert latitude of point 2 from degrees to radians
	deltaLat := (p2[0] - p1[0]) * math.Pi / 180 // Difference in latitudes (in radians)
	deltaLng := (p2[1] - p1[1]) * math.Pi / 180 // Difference in longitudes (in radians)

	// Haversine formula for calculating the distance between two points on a sphere
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)) // Angular distance in radians

	// Earth’s radius in meters (mean radius)
	return earthRadius * c // Final distance in meters
}

// findNextPoint calculates the point on a polyline after traveling a specified distance.
// If the end of the polyline is reached, it reverses the direction and continues the search.
// Returns the destination point, the last index, and the final direction.
func findNextPoint(polyline [][]float64, startingPoint []float64, direction bool, currentIndex int, distanceToTravel float64) ([]float64, int, bool) {
	currentPoint := startingPoint // Initialize the current point as the starting point
	currentDirection := direction // Set the current direction (true for forward, false for backward)

	for distanceToTravel > 0 { // Continue until the remaining distance is exhausted
		// Determine the next index based on the current direction
		nextIndex := currentIndex
		if currentDirection {
			nextIndex++ // Move forward if the direction is true
		} else {
			nextIndex-- // Move backward if the direction is false
		}

		// Check if it’s necessary to reverse the direction at the end of the polyline
		if nextIndex >= len(polyline) { // If we reach the end of the polyline
			currentDirection = false      // Reverse the direction
			nextIndex = len(polyline) - 2 // Move to the second-to-last point
		} else if nextIndex < 0 { // If we reach the beginning of the polyline
			currentDirection = true // Reverse the direction
			nextIndex = 1           // Move to the second point
		}

		// Calculate the distance to the next point on the polyline
		nextPoint := polyline[nextIndex]
		segmentDistance := 0.0
		if distanceToTravel > 1000.00 { // Use Haversine for larger distances
			segmentDistance = distanceHaversine(currentPoint, nextPoint)
		} else { // Use Euclidean for smaller distances
			segmentDistance = distanceHaversine(currentPoint, nextPoint)
		}

		// If the remaining distance is within this segment, calculate the destination point
		if distanceToTravel <= segmentDistance {
			// Calculate the point along the segment where the remaining distance ends
			ratio := distanceToTravel / segmentDistance
			destinationPoint := []float64{
				currentPoint[0] + (nextPoint[0]-currentPoint[0])*ratio,
				currentPoint[1] + (nextPoint[1]-currentPoint[1])*ratio,
			}
			return destinationPoint, currentIndex, currentDirection
		}

		// Update the remaining distance and move to the next point
		distanceToTravel -= segmentDistance
		currentPoint = nextPoint
		currentIndex = nextIndex
	}

	// Return the final point if the entire distance is traversed
	return currentPoint, currentIndex, currentDirection
}

// NearbyGPSIntoPolygon generates a random latitude and longitude within a specified radius (in meters)
// from an initial point and checks if the generated point falls within the boundaries of a polygon
// defined in a GeoJSON file. If successful, it returns the coordinates as a formatted string.
func NearbyGPSIntoPolygon(latitude float64, longitude float64, radius int) string {
	// Lock the GeoJSON context to ensure thread safety
	ctx.JrContext.CtxGeoJsonLock.Lock()
	defer ctx.JrContext.CtxGeoJsonLock.Unlock()

	// Default starting point: either use the provided coordinates or the last known point if available from ctx
	lastLat := latitude
	lastLon := longitude

	// Update last known point if there are recent saved coordinates
	if len(ctx.JrContext.CtxLastPointLat) == 1 {
		lastLat = ctx.JrContext.CtxLastPointLat[len(ctx.JrContext.CtxLastPointLat)-1]
	}
	if len(ctx.JrContext.CtxLastPointLon) == 1 {
		lastLon = ctx.JrContext.CtxLastPointLon[len(ctx.JrContext.CtxLastPointLon)-1]
	}

	// Predict the next point if there is enough data for interpolation
	if len(ctx.JrContext.CtxLastPointLat) >= 2 && len(ctx.JrContext.CtxLastPointLon) >= 2 {
		lastLat, lastLon = predictNextPoint(ctx.JrContext.CtxLastPointLat, ctx.JrContext.CtxLastPointLon)
	}

	// Ensure that the GeoJSON polygon has enough vertices (at least 3) to form a valid shape
	if len(ctx.JrContext.CtxGeoJson) < 3 {
		return fmt.Sprintf("%.12f %.12f", lastLat, lastLon)
	}

	// Convert radius to float for calculations
	radiusInMeters := float64(radius)

	attempts := 0

	// Loop until a valid point within the polygon is found
	for {
		if attempts > 10 {
			// Slightly expand the search radius to ensure coverage
			radiusInMeters *= 1.1
		}
		// Generate a random angle and distance within the specified radius
		randomAngle := Random.Float64() * 2 * math.Pi
		distanceInMeters := Random.Float64() * radiusInMeters

		// Convert the distance from meters to degrees (assuming small distances for simplicity)
		distanceInDegrees := distanceInMeters * degreesPerMeter

		// Calculate new latitude and longitude based on the random angle and distance
		newLatitude := lastLat + (distanceInDegrees * math.Cos(randomAngle))
		newLongitude := lastLon + (distanceInDegrees * math.Sin(randomAngle))

		// Check if the generated point lies within the specified polygon
		if isPointInPolygon([]float64{newLatitude, newLongitude}, ctx.JrContext.CtxGeoJson) {
			// Update the context with the new valid point, maintaining a maximum ctx of 10 points
			ctx.JrContext.CtxLastPointLat = append(ctx.JrContext.CtxLastPointLat, newLatitude)
			ctx.JrContext.CtxLastPointLon = append(ctx.JrContext.CtxLastPointLon, newLongitude)

			// Keep the last 10 points in the ctx
			if len(ctx.JrContext.CtxLastPointLat) > 10 {
				ctx.JrContext.CtxLastPointLat = ctx.JrContext.CtxLastPointLat[1:]
			}
			if len(ctx.JrContext.CtxLastPointLon) > 10 {
				ctx.JrContext.CtxLastPointLon = ctx.JrContext.CtxLastPointLon[1:]
			}
			// Return the coordinates of the valid point
			return fmt.Sprintf("%.12f %.12f", newLatitude, newLongitude)
			// return fmt.Sprintf("%.12f %.12f %d %.12f", newLatitude, newLongitude, attempts, radiusInMeters)
		}
		attempts++
		// Retry if the generated point is not within the polygon boundaries
	}
}

// NearbyGPSIntoPolygonWithoutStart
func NearbyGPSIntoPolygonWithoutStart(radius int) string {
	latitude, longitude := selectRandomPoint(ctx.JrContext.CtxGeoJson)
	return NearbyGPSIntoPolygon(latitude, longitude, radius)
}

// isPointInPolygon checks if a given point lies within a specified polygon.
func isPointInPolygon(point []float64, vertices [][]float64) bool {
	x, y := point[1], point[0]
	n := len(vertices)
	// A polygon must have at least 3 vertices
	if n < 3 {
		return false
	}
	intersections := 0
	for i := 0; i < n; i++ {
		x1, y1 := vertices[i][0], vertices[i][1]
		x2, y2 := vertices[(i+1)%n][0], vertices[(i+1)%n][1]
		if (y1 > y) != (y2 > y) {
			xInt := (y-y1)*(x2-x1)/(y2-y1) + x1
			if x < xInt {
				intersections++
			}
		}
	}
	return intersections%2 == 1
}

// predictNextPoint predicts the next latitude and longitude using cubic spline interpolation
func predictNextPoint(latitudes, longitudes []float64) (float64, float64) {
	if len(longitudes) != len(latitudes) {
		println("Need at least two points and matching latitude/longitude arrays: latitudes: ", len(latitudes), " e longitudes:", len(longitudes))
		os.Exit(1)
	}

	if len(latitudes) < 2 && len(longitudes) < 2 {
		fmt.Println("Need at least two points !")
		return 0, 0
	}

	// Create X values based on indices (0, 1, 2, ...) for even spacing assumption
	x := make([]float64, len(latitudes))
	for i := 0; i < len(latitudes); i++ {
		x[i] = float64(i)
	}

	// Create splines for latitude and longitude
	latSpline := gospline.NewCubicSpline(x, latitudes)
	lonSpline := gospline.NewCubicSpline(x, longitudes)

	// Predict the next index position (next point in the sequence)
	nextX := float64(len(latitudes))

	// Interpolate to get the predicted latitude and longitude for nextX
	nextLatitude := latSpline.At(nextX)
	nextLongitude := lonSpline.At(nextX)

	return nextLatitude, nextLongitude
}

// selectRandomPoint selects a random point within the polygon defined by the given coordinates.
func selectRandomPoint(coords [][]float64) (float64, float64) {
	if len(coords) == 0 {
		return 0, 0 // Return zero values if no coordinates are provided
	}

	// Calculate the bounding box of the coordinates
	minX, minY, maxX, maxY := boundingBox(coords)

	var randX, randY float64

	// Loop until a valid point within the polygon is found
	for {
		// Generate a random point within the bounding box
		x := Random.Float64()*(maxX-minX) + minX
		y := Random.Float64()*(maxY-minY) + minY

		// Check if the generated point is within the polygon
		if isPointInPolygon([]float64{y, x}, coords) {
			randX = x
			randY = y
			break // Exit the loop once a valid point is found
		}
	}

	return randY, randX // Return as (latitude, longitude)
}

// boundingBox calculates the minimum and maximum coordinates (bounding box) of the provided vertices.
// It returns the coordinates of the bottom-left (minX, minY) and top-right (maxX, maxY) corners
// of the bounding box that encompasses all the given vertices.
func boundingBox(vertices [][]float64) (minX, minY, maxX, maxY float64) {
	if len(vertices) == 0 {
		return 0, 0, 0, 0 // Return zero values if no vertices are provided
	}

	// Initialize min and max values based on the first vertex
	minX, minY = vertices[0][0], vertices[0][1]
	maxX, maxY = minX, minY

	// Iterate through vertices to find the bounding box
	for _, vertex := range vertices {
		x, y := vertex[0], vertex[1]
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}
	return
}

// State returns a random State
func State() string {
	s := Word("state")
	ctx.JrContext.Ctx["_state"] = s
	ctx.JrContext.CountryIndex = ctx.JrContext.LastIndex
	return s
}

// StateAt returns State at given index
func StateAt(index int) string {
	return WordAt("state", index)
}

// StateShort returns a random short State
func StateShort() string {
	return Word("state_short")
}

// StateShortAt returns short State at given index
func StateShortAt(index int) string {
	return WordAt("state_short", index)
}

// Street returns a random street
func Street() string {
	return Word("street")
}

// Zip returns a random Zip code
func Zip() string {
	cityIndex := ctx.JrContext.CityIndex

	if cityIndex == -1 {
		z := Word("zip")
		zip, _ := Regex(z)
		return zip
	}

	return ZipAt(cityIndex)
}

// ZipAt returns Zip code at given index
func ZipAt(index int) string {
	z := WordAt("zip", index)
	zip, _ := Regex(z)
	return zip
}
