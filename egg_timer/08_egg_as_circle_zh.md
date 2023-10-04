---
layout: default
title: 第八章 - 画圆圈
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第八章 - 画一个圆圈

## 目标

本章节的目标是绘制自定义图形，它类似于一个鸡蛋（先当它是）。

![Egg as circle](08_egg_as_circle.gif)

## 概述

本章介绍了自定义图形的代码，程序中的圆形由 Gio 进行绘制，而不是直接呈现静态图片，为了实现这一点，我们结合使用了以下内容：

- **op/clip** 以定义我们可以绘制的区域
- **op/paint** 以填充该区域
- 一些参数来设置颜色

## 新包

有一些新的包被导入：

- [image](https://pkg.go.dev/image) & [image/color](https://pkg.go.dev/image/color) - Go 的标准 2D 图像库。

  - Nigel Tao 发在 Go blog 上的[这篇文章](https://blog.golang.org/image)写得很棒。

- [f32](https://pkg.go.dev/gioui.org/f32) - 因为 Go 的图像库基于 `int`，而 Gio 的某些功能使用 `float32`，所以 f32 重新实现了两个主要类型的浮点版本， `Points(点)` 和 `Rectangles(矩形)`。

- [op/clip](https://pkg.go.dev/gioui.org/op/clip) - 用于定义要绘制的区域，在这个区域之外的绘制将被忽略。

- [op/paint](https://pkg.go.dev/gioui.org/op/paint) - 包含用于填充形状颜色的绘制操作。

## 点(Points) 和 矩形(Rectangles)

点和矩形被广泛使用，因此值得引用上文出现的 Nigel 的博客。

点是坐标，而矩形由点定义：

```go
type Point struct {
  X, Y float32
}

type Rectangle struct {
  Min, Max Point
}
```

一个点是 XY 坐标对，坐标轴向右和向下为正（原点 = 左上角）。它既不是像素也不是网格方格，一个点没有固有的宽度、高度或颜色，在下面的图示中使用了一个彩色方块来呈现。

![Point](08_image_package_point.png)

```go
p := image.Point{2, 1}
```

一个矩形包含了最小点和最大点，其中 Min.X <= X < Max.X，Min.Y <= Y < Max.Y。它没有固有的颜色，在下面的图示中用细线对其描边，并放大了它们的最小和最大点。

![Rectangle](08_image_package_rectangle.png)

```go
r := image.Rect(2, 1, 5, 5)
```

为了方便起见，image.Rect(x0, y0, x1, y1) 等同于 `Rectangle{Point{x0, y0}, Point{x1, y1}}`，但更容易输入。它还交换了最小和最大值，以确保它们的格式一致。

就这么多，看看码：

## 代码

```go
layout.Rigid(
  func(gtx C) D {
    circle := clip.Ellipse{
      // 固定 x 坐标，此时尝试一下调整窗口大小
      Min: image.Pt(80, 0),
      Max: image.Pt(320, 240),
      // 动态计算 x 坐标，再试试调整窗口大小
      //Min: image.Pt(gtx.Constraints.Max.X/2-120, 0),
      //Max: image.Pt(gtx.Constraints.Max.X/2+120, 240),
    }.Op(gtx.Ops)
    color := color.NRGBA{R: 200, A: 255}
    paint.FillShape(gtx.Ops, color, circle)
    d := image.Point{Y: 400}
    return layout.Dimensions{Size: d}
  },
),
```

## 注解

首先，我们使用 `clip.Ellipse{}` 定义了一个圆形。它通过在一个矩形中内切一个 `椭圆(Ellipse)` 实现，矩形的尺寸由左上角和右下角分别指定，分别是 `Min` 和 `Max`。

在上面的代码中，圆形是写死的，如果尝试调整窗口大小，你会发现这可能不是你想要的。要进行调整，只需注释掉硬编码的坐标，并取消注释接下来的两行，这将引入动态定位，你可以尝试调整这些尺寸，熟悉一下当调整窗口大小或改变实现椭圆的矩形时，圆形是如何上下移动的。

**注意：** 如果你在使用高 DPI 显示器，并且正好以 125% 的缩放比例来显示，则硬编码坐标的另一个问题将浮现。Gio 很好地实现了 `Dp` _设备独立像素_ **D**evice independent **p**ixels，它确保 1 个 Dp 在不同的显示器和分辨率上具有相同的视觉大小，但在这里进行硬编码时，这种动态性被覆盖了，125% 的分辨率缩放会将 `400 Dp` 宽的窗口（比如在 `app.NewWindow` 中定义好）转换为 `500 像素` 的上下文。你可以通过检查 `gtx.Constraints` 并通过 `gtx.Dp()` 将像素转换为 Dp 来查看这一点。

`color.NRGBA` 定义了圆形中要填充的颜色，需要注意，透明(Alpha)通道默认为 0 (全透明)，所以我们这里直接将它设置为 255。

`paint.FillShape` 用 `color` 填充了这个形状。

最后，我们就像一个小部件一样返回 `Dimensions`，高度为 400。

---

[下一章](09_egg_as_egg_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
