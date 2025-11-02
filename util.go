package main

import (
	"image"
	"image/color"
	"math"

	"github.com/kbinani/screenshot"
)

var myColors = []color.NRGBA{
	color.NRGBA{0xA6, 0xA9, 0xFF, 0xff},
	color.NRGBA{0xCC, 0x8C, 0xF5, 0xff},
	color.NRGBA{0xFF, 0xA6, 0xF9, 0xff},
	color.NRGBA{0xF1, 0xE6, 0xF7, 0xff},
	color.NRGBA{0xFF, 0xE6, 0x8C, 0xff},

	color.NRGBA{0x41, 0xbc, 0xa4, 0xff},
	color.NRGBA{0xA5, 0xED, 0x53, 0xff},
	// FFFD91
	color.NRGBA{0xFF, 0xFD, 0x91, 0xff},
	// FFD353
	color.NRGBA{0xFF, 0xD3, 0x53, 0xff},
	// FF7268
	color.NRGBA{0xFF, 0x72, 0x68, 0xff},
	// 7DD2D1
	color.NRGBA{0x7D, 0xD2, 0xD1, 0xff},
	// BAF06A
	color.NRGBA{0xBA, 0xF0, 0x6A, 0xff},
	// FF6A00
	color.NRGBA{0xFF, 0x6A, 0x00, 0xff},
}

func hsvToRgb(h, s, v float64) (r, g, b uint8) {
	h = math.Mod(h, 360)
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c
	var rp, gp, bp float64
	swi := int(h / 60)
	switch swi {
	case 0:
		rp, gp, bp = c, x, 0
	case 1:
		rp, gp, bp = x, c, 0
	case 2:
		rp, gp, bp = 0, c, x
	case 3:
		rp, gp, bp = 0, x, c
	case 4:
		rp, gp, bp = x, 0, c
	case 5:
		rp, gp, bp = c, 0, x
	default:
		rp, gp, bp = 0, 0, 0
	}
	r = uint8((rp + m) * 255.0)
	g = uint8((gp + m) * 255.0)
	b = uint8((bp + m) * 255.0)
	return
}

func captureFullScreenImage() (image.Image, int, int, error) {
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		return nil, 0, 0, err
	}
	r := img.Bounds()
	return img, r.Dx(), r.Dy(), nil
}

func heart(t float64) (x, y float64) {
	x = 16.0 * math.Pow(math.Sin(t), 3)
	y = 13.0*math.Cos(t) - 5.0*math.Cos(2*t) - 2.0*math.Cos(3*t) - math.Cos(4*t)
	return
}

// AdjustNRGBALightness 调整 NRGBA 颜色的亮度
// factor: 1.0 表示不变，>1.0 提高亮度，<1.0 降低亮度
// 注意：factor 应该在合理范围内（如 0.0-2.0），避免颜色溢出
func AdjustNRGBALightness(r, g, b float32, factor float32) (float32, float32, float32) {
	// 转换到 HSL
	h, s, l := rgbToHsl(r, g, b)

	// 调整亮度
	l = l * factor
	if l > 1.0 {
		l = 1.0
	}
	if l < 0.0 {
		l = 0.0
	}

	// 转换回 RGB
	return hslToRgb(h, s, l)
}

// RGB to HSL 转换
func rgbToHsl(r, g, b float32) (h, s, l float32) {
	max := r
	if g > max {
		max = g
	}
	if b > max {
		max = b
	}

	min := r
	if g < min {
		min = g
	}
	if b < min {
		min = b
	}

	l = (max + min) / 2.0

	if max == min {
		h = 0.0
		s = 0.0
	} else {
		delta := max - min
		if l < 0.5 {
			s = delta / (max + min)
		} else {
			s = delta / (2.0 - max - min)
		}

		if r == max {
			h = (g - b) / delta
		} else if g == max {
			h = 2.0 + (b-r)/delta
		} else {
			h = 4.0 + (r-g)/delta
		}

		h *= 60.0
		if h < 0 {
			h += 360.0
		}
	}

	return h, s, l
}

func getColorLightness(r, g, b float32) float32 {
	max := r
	if g > max {
		max = g
	}
	if b > max {
		max = b
	}

	min := r
	if g < min {
		min = g
	}
	if b < min {
		min = b
	}

	return (max + min) / 2.0
}

// HSL to RGB 转换
func hslToRgb(h, s, l float32) (r, g, b float32) {
	if s == 0.0 {
		r = l
		g = l
		b = l
		return
	}

	var q float32
	if l < 0.5 {
		q = l * (1.0 + s)
	} else {
		q = l + s - l*s
	}

	p := 2.0*l - q

	h = h / 360.0
	r = hueToRgb(p, q, h+1.0/3.0)
	g = hueToRgb(p, q, h)
	b = hueToRgb(p, q, h-1.0/3.0)

	return r, g, b
}

func hueToRgb(p, q, t float32) float32 {
	if t < 0 {
		t += 1.0
	}
	if t > 1 {
		t -= 1.0
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}
