package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

// OutputFormat 输出图片格式枚举
type OutputFormat string

const (
	JPG OutputFormat = "jpg"
	PNG OutputFormat = "png"
)

// ZoomMode 缩放模式枚举
type ZoomMode string

const (
	Origin      ZoomMode = "origin"
	Width       ZoomMode = "width"
	Height      ZoomMode = "height"
	WidthHeight ZoomMode = "width_height"
)

// CombineElement 组合元素接口
type CombineElement interface {
	Draw(g *gg.Context, canvasWidth int)
}

// ImageElement 图片元素
type ImageElement struct {
	ImagePath   string      // 图片路径
	X, Y        int         // 位置坐标
	Width       int         // 宽度
	Height      int         // 高度
	Rotate      float64     // 旋转角度(度)
	Alpha       int         // 透明度(0-255)
	ZoomMode    ZoomMode    // 缩放模式
	RoundCorner int         // 圆角半径
	image       image.Image // 缓存的图片对象
}

// applyAlpha 为图片应用透明度
func applyAlpha(img image.Image, alpha int) image.Image {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)
	alphaRatio := float64(alpha) / 255.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := rgba.RGBAAt(x, y)
			pixel.A = uint8(float64(pixel.A) * alphaRatio)
			rgba.SetRGBA(x, y, pixel)
		}
	}
	return rgba
}

// TextElement 文本元素
type TextElement struct {
	Text          string      // 文本内容
	FontSize      float64     // 字体大小
	X, Y          int         // 文本位置坐标
	Color         color.Color // 文本颜色
	Rotate        float64     // 旋转角度(度)
	MaxLineWidth  int         // 最大行宽，超出则自动换行(像素)
	MaxLineCount  int         // 最大行数，超出部分将被截断
	LineHeight    float64     // 行高，默认1.5倍字体大小
	StrikeThrough bool        // 是否显示删除线
}

// RectangleElement 矩形元素
type RectangleElement struct {
	X, Y        int
	Width       int
	Height      int
	Color       color.Color
	RoundCorner int
}

// ImageCombiner 图片合成器
type ImageCombiner struct {
	width, height int
	context       *gg.Context
	elements      []CombineElement
	outputFormat  OutputFormat
	quality       float64
}

// NewImageCombiner 创建新的图片合成器
func NewImageCombiner(width, height int) *ImageCombiner {
	ctx := gg.NewContext(width, height)
	ctx.SetRGB(1, 1, 1)
	ctx.Clear()

	return &ImageCombiner{
		width:       width,
		height:      height,
		context:     ctx,
		outputFormat: PNG,
		quality:     1.0,
	}
}

// AddElement 添加元素到合成器
func (ic *ImageCombiner) AddElement(element CombineElement) {
	ic.elements = append(ic.elements, element)
}

// AddImageElement 添加图片元素
func (ic *ImageCombiner) AddImageElement(imagePath string, x, y int, zoomMode ZoomMode) (*ImageElement, error) {
	img, err := loadImage(imagePath)
	if err != nil {
		return nil, err
	}

	element := &ImageElement{
		image:    img,
		X:        x,
		Y:        y,
		ZoomMode: zoomMode,
		Alpha:    255,
	}

	ic.AddElement(element)
	return element, nil
}

// AddTextElement 添加文本元素
func (ic *ImageCombiner) AddTextElement(text string, fontSize float64, x, y int) *TextElement {
	element := &TextElement{
		Text:     text,
		FontSize: fontSize,
		X:        x,
		Y:        y,
		Color:    color.Black,
	}

	ic.AddElement(element)
	return element
}

// AddRectangleElement 添加矩形元素
func (ic *ImageCombiner) AddRectangleElement(x, y, width, height int) *RectangleElement {
	element := &RectangleElement{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Color:  color.Black,
	}

	ic.AddElement(element)
	return element
}

// Combine 执行图片合成
func (ic *ImageCombiner) Combine() (image.Image, error) {
	ctx := gg.NewContext(ic.width, ic.height)
	ctx.SetColor(color.White)
	ctx.Clear()

	for _, element := range ic.elements {
		element.Draw(ctx, ic.width)
	}

	return ctx.Image(), nil
}

// Save 将合成图片保存到文件
func (ic *ImageCombiner) Save(filePath string) error {
	img, err := ic.Combine()
	if err != nil {
		return err
	}

	switch ic.outputFormat {
	case JPG:
		return gg.SaveJPG(filePath, img, int(ic.quality*100))
	case PNG:
		return gg.SavePNG(filePath, img)
	default:
		return fmt.Errorf("unsupported output format: %s", ic.outputFormat)
	}
}

// ToBytes 将合成图片编码为[]byte返回
func (ic *ImageCombiner) ToBytes() ([]byte, error) {
	img, err := ic.Combine()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	switch ic.outputFormat {
	case JPG:
		options := jpeg.Options{Quality: int(ic.quality * 100)}
		if err := jpeg.Encode(&buf, img, &options); err != nil {
			return nil, err
		}
	case PNG:
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported output format: %s", ic.outputFormat)
	}
	return buf.Bytes(), nil
}

// loadImage 从路径加载图片
func loadImage(path string) (image.Image, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return decodeImage(resp.Body)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return decodeImage(file)
}

// decodeImage 解码图片数据
func decodeImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	return img, err
}

// Draw 实现CombineElement接口
func (ie *ImageElement) Draw(g *gg.Context, canvasWidth int) {
	// 实现图片绘制逻辑
	g.Push()
	defer g.Pop()

	// 获取原始图片尺寸
	origWidth := ie.image.Bounds().Dx()
	origHeight := ie.image.Bounds().Dy()

	// 根据ZoomMode计算缩放后的尺寸
	width, height := ie.Width, ie.Height
	switch ie.ZoomMode {
	case Origin:
		// 原始比例，不缩放
		width = origWidth
		height = origHeight
	case Width:
		// 指定宽度，高度按比例自动计算
		if origWidth > 0 {
			ratio := float64(origHeight) / float64(origWidth)
			height = int(float64(width) * ratio)
		}
	case Height:
		// 指定高度，宽度按比例自动计算
		if origHeight > 0 {
			ratio := float64(origWidth) / float64(origHeight)
			width = int(float64(height) * ratio)
		}
	case WidthHeight:
		// 指定高度和宽度，强制缩放
		width = ie.Width
		height = ie.Height
	}

	// 创建缩放后的图片
	// 直接使用resize进行缩放，获取image.Image对象
	scaledImg := resize.Resize(uint(width), uint(height), ie.image, resize.Lanczos3)

	// 处理圆角
	if ie.RoundCorner > 0 {
		// 创建圆角蒙版
		mask := gg.NewContext(width, height)

		// 计算最大可能的圆角半径（用于圆形效果）
		maxRadius := float64(width)
		if height < width {
			maxRadius = float64(height)
		}
		maxRadius /= 2

		// 如果设置的圆角半径大于等于最大半径，则使用最大半径形成圆形
		cornerRadius := float64(ie.RoundCorner)
		if cornerRadius >= maxRadius {
			cornerRadius = maxRadius
		}

		mask.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
		mask.Clip()
		mask.DrawImage(scaledImg, 0, 0)
		scaledImg = mask.Image()
	}

	// 应用透明度到图片
	modifiedImage := applyAlpha(scaledImg, ie.Alpha)

	// 处理旋转
	if ie.Rotate != 0 {
		g.Translate(float64(ie.X+width/2), float64(ie.Y+height/2))
		g.Rotate(gg.Radians(ie.Rotate))
		g.DrawImage(modifiedImage, -width/2, -height/2)
	} else {
		g.DrawImage(modifiedImage, ie.X, ie.Y)
	}
}

// GetWidth 计算文本元素的宽度，考虑自动换行后的最长行宽度
func (te *TextElement) GetWidth() float64 {
	// 创建足够大的上下文以确保文本测量准确性
	g := gg.NewContext(10000, 100)

	// 加载字体，与Draw方法保持一致

	// 优先加载系统字体确保测量一致性
	if err := g.LoadFontFace("Alibaba-PuHuiTi-Medium.ttf", te.FontSize); err != nil {
		if err := g.LoadFontFace("/Library/Fonts/Arial.ttf", te.FontSize); err != nil {
			if err := g.LoadFontFace("/System/Library/Fonts/PingFang.ttc", te.FontSize); err != nil {
				g.LoadFontFace("", te.FontSize)
			}
		}
	}

	// 初始化行集合
	var lines []string

	// 仅当设置了最大行宽时才进行换行处理
	if te.MaxLineWidth > 0 {
		runes := []rune(te.Text)
		if len(runes) == 0 {
			return 0
		}

		currentLine := []rune{}
		// 避免变量 shadow，重用之前声明的 lines 变量，移除重复声明
		// lines := []string{} 这行代码移除，因为之前已经声明过 lines 变量

		for _, r := range runes {
			testLine := append(currentLine, r)
			width, _ := g.MeasureString(string(testLine))

			if width > float64(te.MaxLineWidth) && len(currentLine) > 0 {
				lines = append(lines, string(currentLine))
				currentLine = []rune{r}
			} else {
				currentLine = testLine
			}
		}
		if len(currentLine) > 0 {
			lines = append(lines, string(currentLine))
		}

		// 应用最大行数限制
		if te.MaxLineCount > 0 && len(lines) > te.MaxLineCount {
			lines = lines[:te.MaxLineCount]
		}
	} else {
		// 不换行，整段文本作为一行
		lines = []string{te.Text}
	}

	if te.MaxLineCount > 0 && len(lines) > te.MaxLineCount {
		lines = lines[:te.MaxLineCount]
	}

	maxWidth := 0.0

	for _, line := range lines {
		width, _ := g.MeasureString(line)
		if width > maxWidth {
			maxWidth = width
		}

	}

	return maxWidth
}

// Draw 实现CombineElement接口，绘制文本元素并支持自动换行
func (te *TextElement) Draw(g *gg.Context, canvasWidth int) {
	g.Push()
	defer g.Pop()

	g.SetColor(te.Color)
	// 字体加载逻辑：尝试加载自定义字体，失败时降级使用系统字体
	// 优先加载系统字体确保测量一致性
	if err := g.LoadFontFace("Alibaba-PuHuiTi-Medium.ttf", te.FontSize); err != nil {
		if err := g.LoadFontFace("/Library/Fonts/Arial.ttf", te.FontSize); err != nil {
			if err := g.LoadFontFace("/System/Library/Fonts/PingFang.ttc", te.FontSize); err != nil {
				g.LoadFontFace("", te.FontSize) // 使用gg默认字体作为最后的备选
			}
		}
	}

	// 处理旋转文本
	if te.Rotate != 0 {
		g.Translate(float64(te.X), float64(te.Y))
		g.Rotate(gg.Radians(te.Rotate))
		g.DrawString(te.Text, 0, 0)
		// 旋转文本的删除线暂不支持
	} else {
		// 自动换行逻辑：仅当设置了最大行宽时启用
		if te.MaxLineWidth > 0 {
			// 将文本转换为rune切片处理中文
			runes := []rune(te.Text)
			if len(runes) == 0 {
				g.DrawString(te.Text, float64(te.X), float64(te.Y))

				// 绘制删除线
				if te.StrikeThrough {
					width := te.GetWidth()
					strikeY := float64(te.Y) - te.FontSize*0.4
					g.DrawLine(float64(te.X), strikeY, float64(te.X)+width, strikeY)
					g.Stroke()
				}
				return
			}

			// 计算行高：优先使用自定义行高，未设置时使用1.5倍字体大小
			lineHeight := te.LineHeight
			if lineHeight <= 0 {
				lineHeight = te.FontSize * 1.5
			}

			currentLine := []rune{}
			lines := []string{}

			// 按字符逐个添加，判断是否超出最大宽度
			for _, r := range runes {
				// 尝试添加当前字符
				testLine := append(currentLine, r)
				width, _ := g.MeasureString(string(testLine))

				// 如果超出最大宽度且当前行不为空，则换行
				if width > float64(te.MaxLineWidth) && len(currentLine) > 0 {
					lines = append(lines, string(currentLine))
					currentLine = []rune{r} // 新行从当前字符开始
				} else {
					currentLine = testLine
				}
			}
			// 添加最后一行
			if len(currentLine) > 0 {
				lines = append(lines, string(currentLine))
			}

			// 应用最大行数限制：截断超出部分
			if te.MaxLineCount > 0 && len(lines) > te.MaxLineCount {
				lines = lines[:te.MaxLineCount]
			}

			// 绘制所有文本行：默认左对齐，按行高偏移Y坐标
			for i, line := range lines {
				y := float64(te.Y) + float64(i)*lineHeight
				g.DrawString(line, float64(te.X), y)

				// 绘制删除线
				if te.StrikeThrough {
					width, _ := g.MeasureString(line)
					strikeY := y - te.FontSize*0.4 // 调整此值以垂直居中删除线
					g.SetLineWidth(1.0)
					g.DrawLine(float64(te.X), strikeY, float64(te.X)+width, strikeY)
					g.Stroke()
				}
			}
		} else {
			// 不启用自动换行：直接绘制完整文本
			g.DrawString(te.Text, float64(te.X), float64(te.Y))

			// 绘制删除线
			if te.StrikeThrough {
				width, _ := g.MeasureString(te.Text)
				strikeY := float64(te.Y) - te.FontSize*0.4
				g.SetLineWidth(2.0)
				g.DrawLine(float64(te.X), strikeY, float64(te.X)+width, strikeY)
				g.Stroke()
			}
		}
	}
}

// Draw 实现CombineElement接口
func (re *RectangleElement) Draw(g *gg.Context, canvasWidth int) {
	// 实现矩形绘制逻辑
	g.Push()
	defer g.Pop()

	g.SetColor(re.Color)
	if re.RoundCorner > 0 {
		g.DrawRoundedRectangle(float64(re.X), float64(re.Y), float64(re.Width), float64(re.Height), float64(re.RoundCorner))
	} else {
		g.DrawRectangle(float64(re.X), float64(re.Y), float64(re.Width), float64(re.Height))
	}
	g.Fill()
}
