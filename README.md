# 图片合成工具

一个用于合成图片、文本和图形元素的Go语言库，支持多种样式和布局功能。参考了java版本，https://gitee.com/dromara/image-combiner

## 功能特点

- 支持图片元素添加，可设置透明度、旋转和圆角
- 文本元素支持自动换行、字体样式和大小调整
- 矩形元素可自定义颜色和边框
- 多种缩放模式适应不同布局需求
- 完整的单元测试覆盖核心功能

## 安装方法

```bash
# 克隆仓库
git clone https://gitee.com/csn1024/image-combiner-go
cd 绘制图片
注意一定要安装字体，否则中文字会出不来，Alibaba-PuHuiTi-Medium.ttf
仓库里自带了这个字体，直接clone 即可
# 安装依赖
go mod download
```

## 使用示例

```go
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
	avatarImg.Alpha = 100

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
```

## 运行测试

```bash
# 运行所有测试
go test -v

# 查看测试覆盖率
go test -coverprofile=coverage.out
```

## 许可证

[Apache License 2.0](LICENSE)