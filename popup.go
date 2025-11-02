package main

import (
	"bytes"
	"container/list"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"github.com/anthonynsimon/bild/blur"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var messages = []string{
	"愿你每天都有好心情",
	"新的一天，加油！",
	"烟花灿烂，如你一般！",
	"心怀浪漫，活成闪光",
	"世界因你而亮丽",
	"我爱你，永远爱你！",
	"每天都很想你！",
	"我在上海很想你！",
	"我想和你结婚！",
	"你是我的唯一！",
	"每天要开心哦！",
	"记得喝水！",
	"想我就和我打电话！",
	"又是想你的一天！",
	"你是我的宝贝！",
	"要记得想我哟！",
	"我爱你的全部！",
	"小笨蛋，记得想我！",
	"心里满满的都是你",
	"我想你有一百个滇池那么多！",
	"嫁给我，让我照顾你！",
	"人生有你才有意义！",
}

const (
	popupImageWidth  = 350 * 2
	popupImageHeight = 310 * 2
	popupWidth       = 270 * 2
	popupHeight      = 160 * 2
	popupTitleH      = 27 * 2
)

type float32Color struct {
	R, G, B, A float32
}

// Popup 表示一个弹窗消息
type Popup struct {
	X, Y       float64
	Zoom       float64
	Color      float32Color
	TextColor  float32Color
	Alpha      float32
	Life       float64
	MaxAge     float64
	popupImage *ebiten.Image
	textImage  *ebiten.Image
}

func (p *Popup) Draw(screen *ebiten.Image) {
	// 背景卡片
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(p.Zoom, p.Zoom)
	op.GeoM.Translate(p.X-popupImageWidth/2*p.Zoom, p.Y-popupImageHeight/2*p.Zoom)
	op.ColorScale.Scale(p.Color.R, p.Color.G, p.Color.B, 1.0)
	op.ColorScale.ScaleAlpha(p.Alpha)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(p.popupImage, op)

	// 绘制文字（居中）
	tPoint := p.textImage.Bounds().Size()
	tw, th := float64(tPoint.X)*p.Zoom, float64(tPoint.Y)*p.Zoom
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(p.Zoom, p.Zoom)
	opts.GeoM.Translate(p.X-tw/2, p.Y+(popupTitleH*p.Zoom-th)/2)
	opts.ColorScale.Scale(p.TextColor.R, p.TextColor.G, p.TextColor.B, 1.0)
	opts.ColorScale.ScaleAlpha(p.Alpha)
	opts.Filter = ebiten.FilterLinear
	screen.DrawImage(p.textImage, opts)
}

// PopupManager 管理所有弹窗
type PopupManager struct {
	heartPopups []*Popup
	randomPops  list.List
	face        *text.GoTextFace
	heartTime   float64
	heartTick   int
	popupImage  *ebiten.Image
	textImages  []*ebiten.Image
}

// 初始化所有颜色窗口
func (pm *PopupManager) InitPopup() {
	popupX := (popupImageWidth - popupWidth) / 2
	popupY := (popupImageHeight - popupHeight) / 2
	cornerR := 10 * 2 // 圆角半径
	titleColor := color.NRGBA{180, 180, 180, 255}
	whiteColor := color.NRGBA{255, 255, 255, 255}

	closeR := 8 * 2
	closeCX := popupX + popupWidth - 12*2 - closeR
	closeCY := popupY + popupTitleH/2

	bg := image.NewRGBA(image.Rect(0, 0, popupImageWidth, popupImageHeight))

	// 绘制阴影
	drawShadow(bg, popupX, popupY, popupX+popupWidth, popupY+popupHeight, cornerR, 4*2, 6*2, 8*2)

	// 高斯模糊
	bg = blur.Gaussian(bg, 40.0)

	// 窗口主体（圆角矩形，白色）
	fillRoundedRect(bg, popupX, popupY, popupX+popupWidth, popupY+popupHeight, cornerR, whiteColor)

	// 标题栏（顶部圆角需要与窗口圆角对齐）
	fillRoundedRect(bg, popupX, popupY, popupX+popupWidth, popupY+popupTitleH, cornerR, titleColor)
	// 为了在标题栏下方保持直角（避免上半部分的圆角再次填充），我们需要在标题栏底部画一条与背景同色的矩形覆盖小段
	// 但这里因为标题栏是位于窗口顶端，用一个普通矩形清除底部外溢部分（简单处理）
	draw.Draw(bg, image.Rect(popupX, popupY+cornerR, popupX+popupWidth, popupY+popupTitleH), &image.Uniform{titleColor}, image.Point{}, draw.Src)
	// 在标题栏上绘制关闭按钮（右上角）
	drawCloseButton(bg, closeCX, closeCY, closeR)

	img := ebiten.NewImageFromImage(bg)

	// 标题文字
	faceSource, _ := text.NewGoTextFaceSource(bytes.NewReader(notoSansSC))
	face := &text.GoTextFace{Source: faceSource, Size: 16 * 2}
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(popupX)+6*2, float64(popupY)+8*2)
	opts.ColorScale.SetR(1.0)
	opts.ColorScale.SetG(1.0)
	opts.ColorScale.SetB(1.0)
	text.Draw(img, "WH Love SJH", face, opts)

	pm.popupImage = img
}

// 初始化所有文字
func (pm *PopupManager) InitTexts() {
	faceSource, _ := text.NewGoTextFaceSource(bytes.NewReader(notoSansSC))
	pm.face = &text.GoTextFace{Source: faceSource, Size: 20 * 2}

	pm.textImages = make([]*ebiten.Image, len(messages))
	for i, msg := range messages {
		width, height := text.Measure(msg, pm.face, 20*2)

		img := ebiten.NewImage(int(width), int(height))
		// 绘制文字
		text.Draw(img, msg, pm.face, nil)

		pm.textImages[i] = img
	}
}

// 获取随机文字图片
func (pm *PopupManager) GetRandomTextImage() *ebiten.Image {
	return pm.textImages[rand.Intn(len(pm.textImages))]
}

// 新建弹窗
func (pm *PopupManager) NewPopup(x, y, MaxAge float64, color color.NRGBA, alpha float32, zoom float64) *Popup {
	p := &Popup{
		X:          x,
		Y:          y,
		Zoom:       zoom,
		Life:       0,
		MaxAge:     MaxAge,
		popupImage: pm.popupImage,
		textImage:  pm.GetRandomTextImage(),
		Alpha:      alpha,
	}
	p.Color.R = float32(color.R) / 255.0
	p.Color.G = float32(color.G) / 255.0
	p.Color.B = float32(color.B) / 255.0

	// 判断颜色亮度，并设置文字颜色
	light := getColorLightness(p.Color.R, p.Color.G, p.Color.B)
	if light <= 0.5 {
		p.TextColor.R, p.TextColor.G, p.TextColor.B = AdjustNRGBALightness(p.Color.R, p.Color.G, p.Color.B, 2.4)
	} else {
		p.TextColor.R, p.TextColor.G, p.TextColor.B = AdjustNRGBALightness(p.Color.R, p.Color.G, p.Color.B, 0.3)
	}

	return p
}
func (pm *PopupManager) NewHeart(x, y float64) {
	pm.heartPopups = append(pm.heartPopups, pm.NewPopup(x, y, 9999, myColors[2], 0.9, 0.5))
}

// 新建随机弹窗
func (pm *PopupManager) NewRandom(x, y float64, shuttle bool) {
	if shuttle {
		pm.randomPops.PushFront(pm.NewPopup(x, y, 6+rand.Float64()*16, myColors[rand.Intn(len(myColors))], 0.9, 0.25))
	} else {
		pm.randomPops.PushBack(pm.NewPopup(x, y, 4+rand.Float64()*10, myColors[rand.Intn(len(myColors))], 0.9, 0.25))
	}
}

// 更新弹窗状态
func (pm *PopupManager) Update(dt float64) {
	windowWidth2 := windowWidth / 2
	windowHeight2 := windowHeight / 2

	element := pm.randomPops.Front()
	for element != nil {
		p := element.Value.(*Popup)
		p.Life += dt
		progress := float32(p.Life / p.MaxAge)
		if progress > 0.7 {
			p.Alpha = (1.0 - (progress-0.7)/0.3) * 0.9 // 淡出
		}
		if p.Life >= p.MaxAge {
			next := element.Next()
			pm.randomPops.Remove(element)
			element = next
		} else {
			element = element.Next()
		}

		if pm.heartTime >= 6.0 {
			p.X = (p.X-windowWidth2)*1.0045 + windowWidth2
			p.Y = (p.Y-windowHeight2)*1.0045 + windowHeight2
			p.Zoom *= 1.002
		}
	}

	s := 35 * windowHeight / 1440
	if pm.heartTime < math.Pi*2 {
		// 生成心形
		if pm.heartTick%9 == 0 {
			x, y := heart(pm.heartTime)
			x = windowWidth/2 + x*s
			y = windowHeight/2 - y*s
			pm.NewHeart(x, y)
			pm.heartTime += dt * 3
		}
		pm.heartTick += 1
	} else {
		if rand.Intn(100) < 40 {
			if pm.heartTime >= 6.0 {
				x := rand.Float64()*windowWidth/6 + windowWidth*5/12
				y := rand.Float64()*windowHeight/6 + windowHeight*5/12
				pm.NewRandom(x, y, true)
			} else {
				x := rand.Float64() * windowWidth
				y := rand.Float64() * windowHeight
				pm.NewRandom(x, y, false)
			}
		}

		if len(pm.heartPopups) > 0 {
			// 爱心放大
			windowWidth2 := windowWidth / 2
			windowHeight2 := windowHeight / 2
			for _, p := range pm.heartPopups {
				p.X = (p.X-windowWidth2)*1.005 + windowWidth2
				p.Y = (p.Y-windowHeight2)*1.005 + windowHeight2
				p.Zoom *= 1.002
			}
			if pm.heartTime > 12.3 {
				pm.heartPopups = nil
			}

			pm.heartTime += dt
		}
	}
}

// 绘制所有弹窗
func (pm *PopupManager) Draw(screen *ebiten.Image) {
	for element := pm.randomPops.Front(); element != nil; element = element.Next() {
		p := element.Value.(*Popup)
		p.Draw(screen)
	}
	for _, p := range pm.heartPopups {
		p.Draw(screen)
	}
}
