// window_draw.go
package main

import (
	"image"
	"image/color"
	"math"
)

// helper: clamp
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// 返回距离 (x,y) 到圆心(cx,cy) 的欧氏距离
func dist(x, y, cx, cy int) float64 {
	dx := float64(x - cx)
	dy := float64(y - cy)
	return math.Hypot(dx, dy)
}

// 在 rgba 上绘制一个圆角矩形（填充）
// x0,y0 - 左上， x1,y1 - 右下， r - 圆角半径， col - 填充色
func fillRoundedRect(img *image.RGBA, x0, y0, x1, y1, r int, col color.NRGBA) {
	b := img.Bounds()
	for y := clamp(y0, b.Min.Y, b.Max.Y); y < clamp(y1, b.Min.Y, b.Max.Y); y++ {
		for x := clamp(x0, b.Min.X, b.Max.X); x < clamp(x1, b.Min.X, b.Max.X); x++ {
			// 四角判断
			in := false
			// 中心区域直接填充
			if x >= x0+r && x < x1-r {
				in = true
			} else if y >= y0+r && y < y1-r {
				in = true
			} else {
				// 角上用圆形判断
				var cx, cy int
				if x < x0+r {
					cx = x0 + r
				} else {
					cx = x1 - r - 1
				}
				if y < y0+r {
					cy = y0 + r
				} else {
					cy = y1 - r - 1
				}
				if dist(x, y, cx, cy) <= float64(r) {
					in = true
				}
			}
			if in {
				offset := img.PixOffset(x, y)
				img.Pix[offset+0] = col.R
				img.Pix[offset+1] = col.G
				img.Pix[offset+2] = col.B
				img.Pix[offset+3] = col.A
			}
		}
	}
}

// 在 rgba 上绘制带 alpha 混合的像素（源覆盖到目标，按 src.A 混合）
func blendPixel(img *image.RGBA, x, y int, src color.RGBA) {
	if !(image.Pt(x, y).In(img.Bounds())) {
		return
	}
	off := img.PixOffset(x, y)
	dr := img.Pix[off+0]
	dg := img.Pix[off+1]
	db := img.Pix[off+2]
	da := img.Pix[off+3]

	sa := float64(src.A) / 255.0
	da_f := float64(da) / 255.0

	// premultiplied blending simple approximation:
	outA := sa + da_f*(1-sa)
	if outA == 0 {
		img.Pix[off+0], img.Pix[off+1], img.Pix[off+2], img.Pix[off+3] = 0, 0, 0, 0
		return
	}
	// convert to linear channel floats
	sr := float64(src.R) / 255.0
	sg := float64(src.G) / 255.0
	sb := float64(src.B) / 255.0

	dr_f := float64(dr) / 255.0
	dg_f := float64(dg) / 255.0
	db_f := float64(db) / 255.0

	outR := (sr*sa + dr_f*da_f*(1-sa)) / outA
	outG := (sg*sa + dg_f*da_f*(1-sa)) / outA
	outB := (sb*sa + db_f*da_f*(1-sa)) / outA

	img.Pix[off+0] = uint8(clamp(int(math.Round(outR*255.0)), 0, 255))
	img.Pix[off+1] = uint8(clamp(int(math.Round(outG*255.0)), 0, 255))
	img.Pix[off+2] = uint8(clamp(int(math.Round(outB*255.0)), 0, 255))
	img.Pix[off+3] = uint8(clamp(int(math.Round(outA*255.0)), 0, 255))
}

// 画一个圆角矩形的阴影：通过多次绘制稍微增大的半透明圆角矩形来模拟模糊阴影
func drawShadow(img *image.RGBA, x0, y0, x1, y1, r int, offsetX, offsetY int, maxBlur int) {
	// maxBlur 越大阴影越柔和（但绘制更慢）
	// 从外到内绘制，外部更透明
	for i := maxBlur; i >= 1; i-- {
		// alpha 随 i 递增（靠近窗口的部分更深）
		// alpha := uint8(maxBlur + (maxBlur-i)*1) // 可调
		alpha := uint8((maxBlur - i) * 4)
		expand := i * 4 // 每一层扩大一像素
		fillRoundedRect(img, x0+offsetX-expand, y0+offsetY-expand, x1+offsetX+expand, y1+offsetY+expand, r+expand, color.NRGBA{0, 0, 0, alpha})
	}
}

// 绘制 "X" 关闭按钮（一个圆内的白色 X）
func drawCloseButton(img *image.RGBA, cx, cy, radius int) {
	// 先画圆背景
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			if dist(x, y, cx, cy) <= float64(radius) {
				blendPixel(img, x, y, color.RGBA{200, 60, 60, 255})
			}
		}
	}
	// 画 X（用两条斜线）
	thickness := int(math.Max(1, float64(radius)/8.0))
	radius_2 := int(float32(radius) / 2.5)
	for dy := -thickness; dy <= thickness; dy++ {
		for t := -radius_2; t <= radius_2; t++ {
			// main diagonal
			xx := cx + t
			yy := cy + t + dy
			blendPixel(img, xx, yy, color.RGBA{255, 255, 255, 220})
			// other diagonal
			xx2 := cx + t
			yy2 := cy - t + dy
			blendPixel(img, xx2, yy2, color.RGBA{255, 255, 255, 220})
		}
	}
}
