package ui

import (
	"fmt"
	"strings"
)

const (
	mapWidth    = 76
	bitmapRows  = 40 // 2x vertical resolution for half-block rendering
	displayRows = 20
)

// landData defines land areas at 4° latitude resolution (40 rows).
// Each row contains pairs of (startCol, endCol) inclusive.
// Col 0 = 180°W, Col 75 ≈ 175°E. Each col ≈ 4.74°.
// Row 0 = 80°N, Row 39 = 80°S. Each row = 4°.
var landData = [bitmapRows][]int{
	// Row 0 (80-76°N): Arctic islands
	{15, 17, 50, 52, 55, 56},
	// Row 1 (76-72°N): Canadian Arctic, Greenland, Svalbard, N Russia
	{14, 22, 28, 32, 41, 41, 50, 52, 54, 75},
	// Row 2 (72-68°N): Alaska, Canadian Arctic, Greenland, Norway, Russia
	{3, 4, 12, 23, 27, 33, 41, 43, 46, 75},
	// Row 3 (68-64°N): Alaska, N Canada, Greenland, Iceland, Scandinavia, Russia
	{3, 5, 10, 24, 27, 33, 37, 38, 41, 75},
	// Row 4 (64-60°N): Alaska, Canada, S Greenland, Scandinavia, Russia
	{3, 5, 9, 24, 28, 32, 40, 75},
	// Row 5 (60-56°N): Alaska, Canada, UK, Scandinavia, Russia
	{3, 5, 9, 24, 38, 39, 40, 72},
	// Row 6 (56-52°N): Canada, Newfoundland, UK, Europe, Russia, Kamchatka
	{9, 25, 37, 70, 74, 75},
	// Row 7 (52-48°N): S Canada, UK, W/Central Europe, Russia, Kamchatka
	{9, 22, 38, 68, 74, 75},
	// Row 8 (48-44°N): N USA, France, Central Europe, Ukraine, Kazakhstan
	{10, 22, 38, 60, 62, 66},
	// Row 9 (44-40°N): USA, Iberia, Italy, Turkey, Central Asia, China, Hokkaido
	{10, 22, 37, 66, 68, 69},
	// Row 10 (40-36°N): USA, Portugal, Italy, Turkey, Iran, China, Korea, Japan
	{11, 21, 37, 38, 40, 42, 44, 54, 56, 68},
	// Row 11 (36-32°N): S USA, Morocco, Tunisia, Middle East, Pakistan, China, Japan
	{12, 21, 36, 38, 40, 42, 44, 68},
	// Row 12 (32-28°N): Mexico/S USA, N Africa, Iraq/Iran, India, S China
	{13, 21, 35, 50, 52, 65},
	// Row 13 (28-24°N): Mexico, Sahara, Arabia, India, S China, Taiwan
	{14, 20, 34, 49, 53, 65},
	// Row 14 (24-20°N): Mexico, Cuba, Sahara, Arabia, India, Myanmar/Vietnam
	{14, 19, 21, 21, 33, 48, 54, 61},
	// Row 15 (20-16°N): Mexico, Cuba, Sahel, Yemen, India, SE Asia, Philippines
	{14, 18, 20, 21, 33, 48, 54, 61, 64, 65},
	// Row 16 (16-12°N): C America, W Africa, Ethiopia, India, SE Asia, Philippines
	{17, 22, 33, 42, 44, 47, 55, 57, 59, 61, 64, 65},
	// Row 17 (12-8°N): C America, W Africa, Ethiopia, S India, SE Asia
	{17, 19, 34, 42, 44, 47, 55, 57, 59, 62, 64, 64},
	// Row 18 (8-4°N): Panama, Colombia/Venezuela, W/C Africa, Malay, Borneo
	{18, 24, 34, 47, 57, 57, 60, 61, 63, 65},
	// Row 19 (4-0°N): Colombia, Guianas, C Africa, Sumatra, Borneo
	{19, 22, 24, 25, 38, 47, 60, 65},
	// Row 20 (0-4°S): Amazon, C/E Africa, Sumatra, Borneo, Sulawesi
	{19, 27, 40, 48, 60, 65},
	// Row 21 (4-8°S): Brazil, C/E Africa, Java, Borneo, Papua
	{22, 28, 42, 48, 62, 65, 70, 71},
	// Row 22 (8-12°S): Brazil, E Africa, Timor, N Australia
	{23, 29, 44, 48, 65, 71},
	// Row 23 (12-16°S): Brazil, E Africa, N Australia
	{23, 29, 44, 48, 65, 72},
	// Row 24 (16-20°S): Brazil/Bolivia, Mozambique, Madagascar, Australia
	{23, 29, 45, 48, 49, 50, 64, 72},
	// Row 25 (20-24°S): Brazil, Namibia/Zimbabwe, Madagascar, Australia
	{23, 28, 42, 48, 49, 50, 64, 73},
	// Row 26 (24-28°S): Brazil/Paraguay, Namibia/S Africa, Australia
	{23, 28, 42, 48, 64, 73},
	// Row 27 (28-32°S): S Brazil/Uruguay, S Africa, Australia
	{24, 27, 43, 48, 64, 73},
	// Row 28 (32-36°S): Chile/Argentina, S Africa, SE Australia
	{22, 26, 43, 47, 68, 73},
	// Row 29 (36-40°S): Argentina, S Africa tip, SE Australia, NZ
	{22, 25, 44, 46, 69, 73, 75, 75},
	// Row 30 (40-44°S): Argentina, Tasmania, NZ
	{22, 25, 72, 72, 75, 75},
	// Row 31 (44-48°S): Chile/Argentina, NZ
	{22, 24, 75, 75},
	// Row 32 (48-52°S): Patagonia
	{22, 24},
	// Row 33 (52-56°S): S Patagonia
	{22, 23},
	// Row 34 (56-60°S): Tierra del Fuego
	{22, 23},
	// Row 35 (60-64°S): ocean
	{},
	// Row 36 (64-68°S): Antarctic Peninsula
	{22, 23},
	// Row 37 (68-72°S): Antarctic coast
	{20, 24, 46, 54, 62, 70},
	// Row 38 (72-76°S): Antarctica
	{16, 26, 38, 75},
	// Row 39 (76-80°S): Antarctica
	{14, 28, 36, 75},
}

func generateBitmap() [bitmapRows][mapWidth]bool {
	var grid [bitmapRows][mapWidth]bool
	for y := 0; y < bitmapRows; y++ {
		ranges := landData[y]
		for i := 0; i+1 < len(ranges); i += 2 {
			for x := ranges[i]; x <= ranges[i+1] && x < mapWidth; x++ {
				grid[y][x] = true
			}
		}
	}
	return grid
}

func RenderMap(issLat, issLon float64, userLat, userLon float64, showUser bool) string {
	grid := generateBitmap()

	// Convert to bitmap coordinates (40 rows, 76 cols)
	issX, issY := latLonToBitmap(issLat, issLon)
	var userX, userY int
	if showUser {
		userX, userY = latLonToBitmap(userLat, userLon)
	}

	var sb strings.Builder

	// Longitude header
	sb.WriteString("     180W      120W       60W        0        60E       120E      180E\n")

	for row := 0; row < displayRows; row++ {
		topY := row * 2
		botY := row * 2 + 1

		// Latitude label (center of the display row)
		lat := 80.0 - float64(row)*8.0 - 4.0
		if lat > 0 {
			sb.WriteString(fmt.Sprintf(" %2.0f°N", lat))
		} else if lat < 0 {
			sb.WriteString(fmt.Sprintf(" %2.0f°S", -lat))
		} else {
			sb.WriteString("  0° ")
		}

		for x := 0; x < mapWidth; x++ {
			issDisplayRow := issY / 2
			isISS := issDisplayRow == row && issX == x

			userDisplayRow := userY / 2
			isUser := showUser && userDisplayRow == row && userX == x

			if isISS {
				sb.WriteString(ISSMarkerStyle.Render("★"))
				continue
			}
			if isUser {
				sb.WriteString(UserMarkerStyle.Render("◉"))
				continue
			}

			top := grid[topY][x]
			bot := false
			if botY < bitmapRows {
				bot = grid[botY][x]
			}

			switch {
			case top && bot:
				sb.WriteString(LandStyle.Render("█"))
			case top:
				sb.WriteString(LandStyle.Render("▀"))
			case bot:
				sb.WriteString(LandStyle.Render("▄"))
			default:
				sb.WriteString(OceanStyle.Render("·"))
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func latLonToBitmap(lat, lon float64) (int, int) {
	x := int((lon + 180.0) / 360.0 * float64(mapWidth))
	y := int((80.0 - lat) / 160.0 * float64(bitmapRows))

	if x < 0 {
		x = 0
	}
	if x >= mapWidth {
		x = mapWidth - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= bitmapRows {
		y = bitmapRows - 1
	}

	return x, y
}
