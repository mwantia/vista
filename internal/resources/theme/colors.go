package theme

import (
	"fmt"
	"math"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// AdjustBrightness adjusts the brightness of a color by the given percentage
// Positive percentage increases brightness, negative decreases it
// Handles both solid colors and adaptive colors
func AdjustBrightness(color lipgloss.TerminalColor, percent float64) lipgloss.TerminalColor {
	switch c := color.(type) {
	case lipgloss.Color:
		return lipgloss.Color(adjustHexBrightness(string(c), percent))
	case lipgloss.AdaptiveColor:
		return lipgloss.AdaptiveColor{
			Light: adjustHexBrightness(c.Light, percent),
			Dark:  adjustHexBrightness(c.Dark, percent),
		}
	default:
		return color
	}
}

// adjustHexBrightness adjusts the brightness of a hex color string
func adjustHexBrightness(hex string, percent float64) string {
	// Parse hex color
	r, g, b := parseHexColor(hex)
	if r == 0 && g == 0 && b == 0 && hex != "#000000" && hex != "#000" {
		// Failed to parse, return original
		return hex
	}

	// Convert to HSL
	h, s, l := rgbToHSL(r, g, b)

	// Adjust lightness
	l = math.Max(0, math.Min(1, l+percent/100))

	// Convert back to RGB
	r, g, b = hslToRGB(h, s, l)

	// Convert to hex
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// parseHexColor parses a hex color string to RGB values
func parseHexColor(hex string) (r, g, b uint8) {
	// Remove # if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	// Handle 3-digit hex colors
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}

	if len(hex) != 6 {
		return 0, 0, 0
	}

	// Parse each component
	if val, err := strconv.ParseUint(hex[0:2], 16, 8); err == nil {
		r = uint8(val)
	}
	if val, err := strconv.ParseUint(hex[2:4], 16, 8); err == nil {
		g = uint8(val)
	}
	if val, err := strconv.ParseUint(hex[4:6], 16, 8); err == nil {
		b = uint8(val)
	}

	return r, g, b
}

// rgbToHSL converts RGB to HSL color space
func rgbToHSL(r, g, b uint8) (h, s, l float64) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(math.Max(rf, gf), bf)
	min := math.Min(math.Min(rf, gf), bf)
	delta := max - min

	// Lightness
	l = (max + min) / 2

	if delta == 0 {
		// Grayscale
		return 0, 0, l
	}

	// Saturation
	if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2 - max - min)
	}

	// Hue
	switch max {
	case rf:
		h = (gf - bf) / delta
		if gf < bf {
			h += 6
		}
	case gf:
		h = (bf-rf)/delta + 2
	case bf:
		h = (rf-gf)/delta + 4
	}
	h /= 6

	return h, s, l
}

// hslToRGB converts HSL to RGB color space
func hslToRGB(h, s, l float64) (r, g, b uint8) {
	if s == 0 {
		// Grayscale
		val := uint8(l * 255)
		return val, val, val
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = uint8(hueToRGB(p, q, h+1.0/3.0) * 255)
	g = uint8(hueToRGB(p, q, h) * 255)
	b = uint8(hueToRGB(p, q, h-1.0/3.0) * 255)

	return r, g, b
}

// hueToRGB is a helper for HSL to RGB conversion
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}
