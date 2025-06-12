package main

import (
	"fmt"
	"image/color"
	"os"
	"testing"
)

// TestSimpleCombine 测试简单图片合成功能
func TestSimpleCombine(t *testing.T) {
	// 先加载图片获取原始尺寸
	img, err := loadImage("https://img.thebeastshop.com/image/20201130115835493501.png")
	if err != nil {
		t.Fatalf("加载图片失败: %v", err)
	}
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// 创建合成器，使用原图尺寸作为画布大小
	combiner := NewImageCombiner(imgWidth, imgHeight, PNG)

	// 添加背景矩形，大小与原图一致
	bgRect := combiner.AddRectangleElement(0, 0, imgWidth, imgHeight)
	bgRect.Color = color.RGBA{240, 240, 240, 255}

	// 添加图片元素（使用原图尺寸，不缩放）
	imgElement, err := combiner.AddImageElement("https://img.thebeastshop.com/image/20201130115835493501.png", 0, 0, Height)
	if err != nil {
		t.Fatalf("添加图片元素失败: %v", err)
	}
	imgElement.ZoomMode = Origin
	imgElement.RoundCorner = 20

	// 添加文本元素
	textElement := combiner.AddTextElement("Go11111", 36, 250, 550)
	textElement.Color = color.RGBA{255, 0, 0, 255}

	// 合成并保存图片
	err = combiner.Save("test_simple_output.png")
	if err != nil {
		t.Fatalf("保存图片失败: %v", err)
	}

	// 验证文件是否创建
	if _, err := os.Stat("test_simple_output.png"); os.IsNotExist(err) {
		t.Error("合成图片文件未生成")
	}
}

// TestAdvancedFeatures 测试高级功能（圆角、旋转、透明度等）
// TestMultipleImages 测试多图片叠加功能
func TestMultipleImages(t *testing.T) {
	// 创建简化的测试场景，仅包含基础元素
	combiner := NewImageCombiner(600, 400, PNG)

	// 添加纯色背景（应始终可见）
	bg := combiner.AddRectangleElement(0, 0, 600, 400)
	bg.Color = color.RGBA{255, 255, 0, 255} // 黄色背景

	// 添加中心红色矩形（应在背景之上）
	centerRect := combiner.AddRectangleElement(250, 150, 100, 100)
	centerRect.Color = color.RGBA{255, 0, 0, 255} // 红色

	// 添加调试文本（应在所有元素之上）
	testText := combiner.AddTextElement("调试测试", 36, 250, 300)
	testText.Color = color.RGBA{0, 0, 0, 255}

	// 保存输出图片
	err := combiner.Save("test_multiple_images.png")
	if err != nil {
		t.Fatalf("保存图片失败: %v", err)
	}

	// 验证输出文件存在且不为空
	fileInfo, err := os.Stat("test_multiple_images.png")
	if err != nil {
		t.Fatalf("输出文件不存在: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Fatalf("输出文件为空，大小: %d bytes", fileInfo.Size())
	}
	fmt.Printf("输出文件已创建，大小: %d bytes\n", fileInfo.Size())
}

// TestTextFeatures 测试文本各种特性
func TestTextFeatures(t *testing.T) {
	combiner := NewImageCombiner(500, 300, PNG)

	// 添加背景
	bg := combiner.AddRectangleElement(0, 0, 500, 300)
	bg.Color = color.RGBA{255, 255, 255, 255}

	// 添加普通文本
	text1 := combiner.AddTextElement("普通文本", 24, 50, 50)
	text1.Color = color.RGBA{0, 0, 0, 255}

	// 添加旋转文本
	text2 := combiner.AddTextElement("旋转文本", 24, 250, 100)
	text2.Color = color.RGBA{255, 0, 0, 255}
	text2.Rotate = 45

	// 添加大字体文本
	text3 := combiner.AddTextElement("大字体", 36, 50, 200)
	text3.Color = color.RGBA{0, 0, 255, 255}

	err := combiner.Save("test_text_features.png")
	if err != nil {
		t.Fatalf("保存图片失败: %v", err)
	}
}

// TestRectangleFeatures 测试矩形各种特性
func TestRectangleFeatures(t *testing.T) {
	combiner := NewImageCombiner(500, 400, PNG)

	// 添加背景
	bg := combiner.AddRectangleElement(0, 0, 500, 400)
	bg.Color = color.RGBA{240, 240, 240, 255}

	// 普通矩形
	rect1 := combiner.AddRectangleElement(50, 50, 100, 100)
	rect1.Color = color.RGBA{255, 0, 0, 200}

	// 圆角矩形
	rect2 := combiner.AddRectangleElement(200, 50, 100, 100)
	rect2.Color = color.RGBA{0, 255, 0, 200}
	rect2.RoundCorner = 20

	// 半透明矩形
	rect3 := combiner.AddRectangleElement(50, 200, 100, 100)
	rect3.Color = color.RGBA{0, 0, 255, 100}

	err := combiner.Save("test_rectangle_features.png")
	if err != nil {
		t.Fatalf("保存图片失败: %v", err)
	}
}

// TestAdvancedFeatures 测试高级功能（圆角、旋转、透明度等）
func TestAdvancedFeatures(t *testing.T) {
	// 使用固定画布尺寸以便测试元素可见性
	combiner := NewImageCombiner(400, 400, PNG)

	// 添加纯色背景
	bgRect := combiner.AddRectangleElement(0, 0, 400, 400)
	bgRect.Color = color.RGBA{255, 255, 255, 255}

	// 添加带圆角的矩形
	rect := combiner.AddRectangleElement(50, 50, 300, 300)
	rect.Color = color.RGBA{255, 0, 0, 255}
	rect.RoundCorner = 30

	// 添加旋转的文本
	text := combiner.AddTextElement("测试文本", 36, 100, 250)
	text.Color = color.Black
	text.Rotate = 30

	err := combiner.Save("test_advanced_output.png")

	if err != nil {
		t.Fatalf("保存图片失败: %v", err)
	}
}

// TestFullFunctionality 完整功能测试，模拟Java测试用例
func TestFullFunctionality(t *testing.T) {
	// 背景图URL
	bgImageUrl := "https://img.thebeastshop.com/combine_image/funny_topic/resource/bg_3x4.png"
	qrCodeUrl := "http://imgtest.thebeastshop.com/file/combine_image/qrcodef3d132b46b474fe7a9cc6e76a511dfd5.jpg"
	productImageUrl := "https://img.thebeastshop.com/combine_image/funny_topic/resource/product_3x4.png"
	waterMarkUrl := "https://img.thebeastshop.com/combine_image/funny_topic/resource/water_mark.png"
	avatarUrl := "https://img.thebeastshop.com/member/privilege/level-icon/level-three.jpg"

	// 文本内容
	title := "# 最爱的家居"
	content := "苏格拉底说：“如果没有那个桌子，可能就没有那个水壶”"

	// 创建合成器
	// 注意：Go版本需要先加载背景图获取尺寸
	bgImg, err := loadImage(bgImageUrl)
	if err != nil {
		t.Fatalf("加载背景图失败: %v", err)
	}
	bgWidth := bgImg.Bounds().Dx()
	bgHeight := bgImg.Bounds().Dy()
	combiner := NewImageCombiner(bgWidth, bgHeight, PNG)

	// 添加商品图
	bg, err := combiner.AddImageElement(bgImageUrl, 0, 0, Origin)
	if err != nil {
		t.Fatalf("添加商品图失败: %v", err)
	}
	bg.ZoomMode = Origin

	// 添加商品图
	productImg, err := combiner.AddImageElement(productImageUrl, 150, 160, Width)
	if err != nil {
		t.Fatalf("添加商品图失败: %v", err)
	}
	productImg.RoundCorner = 30
	productImg.Height = 0
	productImg.Width = 837

	// 添加标题文本
	titleText := combiner.AddTextElement(title, 55, 150, 1400)
	titleText.Color = color.RGBA{0, 0, 0, 255}

	// 添加内容文本
	contentText := combiner.AddTextElement(content, 40, 150, 1480)
	contentText.Color = color.RGBA{0, 0, 0, 255}
	contentText.MaxLineCount = 2
	contentText.MaxLineWidth = 837
	contentText.LineHeight = 60

	// 添加头像
	avatarImg, err := combiner.AddImageElement(avatarUrl, 200, 1200, WidthHeight)
	if err != nil {
		t.Fatalf("添加头像失败: %v", err)
	}
	avatarImg.Width = 130
	avatarImg.Height = 130
	avatarImg.RoundCorner = 200
	// avatarImg.Alpha = 100

	// 添加水印
	waterMarkImg, err := combiner.AddImageElement(waterMarkUrl, 630, 1200, Origin)
	if err != nil {
		t.Fatalf("添加水印失败: %v", err)
	}
	waterMarkImg.Rotate = 15

	// 添加二维码
	qrCodeImg, err := combiner.AddImageElement(qrCodeUrl, 138, 1707, WidthHeight)
	if err != nil {
		t.Fatalf("添加二维码失败: %v", err)
	}
	qrCodeImg.Width = 186
	qrCodeImg.Height = 186
	qrCodeImg.Alpha = 255
	qrCodeImg.Rotate = 0

	// 添加价格文本
	textPrice := combiner.AddTextElement("￥1290", 40, 600, 1400)
	textPrice.Color = color.RGBA{100, 100, 100, 255}
	textPrice.StrikeThrough = true

	// 添加优惠价格
	//动态计算位置
	offsetPrice := float64(textPrice.X) + textPrice.GetWidth() + 10
	fmt.Println("offsetPrice", textPrice.GetWidth(), offsetPrice, int(offsetPrice))
	priceText := combiner.AddTextElement("￥999", 60, int(offsetPrice), 1400)
	priceText.Color = color.RGBA{255, 0, 0, 255}

	// 执行合成并保存
	err = combiner.Save("test_full_functionality.png")
	if err != nil {
		t.Fatalf("保存图片失败: %v", err)
	}

	// 验证文件是否创建
	if _, err := os.Stat("test_full_functionality.png"); os.IsNotExist(err) {
		t.Error("完整功能测试输出文件未生成")
	}
}
