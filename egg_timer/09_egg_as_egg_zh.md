---
layout: default
title: 第九章 - 画鸡蛋
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第九章 - 画一个鸡蛋

## 目标

本章节的目标是绘制一个真正的蛋。

![An actual egg](09_egg_as_egg.gif)

## 大纲

在这里，我们利用基本的 Gio 功能来绘制一个像鸡蛋一样的图形。

## 代码

所有新代码都在之前画圆圈的 `layout.Rigid` 里面。

```go
layout.Rigid(
  func(gtx C) D {
    // 绘制蛋形的路径
    var eggPath clip.Path
    op.Offset(image.Pt(gtx.Dp(200), gtx.Dp(150))).Add(gtx.Ops)
    eggPath.Begin(gtx.Ops)
    // 从0度旋转到360度
    for deg := 0.0; deg <= 360; deg++ {

      // 在这个精彩的网站上找到了蛋形函数 (不开玩笑)，非常感激！
      // https://observablehq.com/@toja/egg-curve
      // 将角度转换为弧度
      rad := deg / 360 * 2 * math.Pi
      // 三角函数给出X和Y方向的距离
      cosT := math.Cos(rad)
      sinT := math.Sin(rad)
      // 用于定义蛋形的常数
      a := 110.0
      b := 150.0
      d := 20.0
      // XY坐标
      x := a * cosT
      y := -(math.Sqrt(b*b-d*d*cosT*cosT) + d*sinT) * sinT
      // 最终得到轮廓上的点
      p := f32.Pt(float32(x), float32(y))
      // 连线到这个点
      eggPath.LineTo(p)
    }
    // 结束路径绘制
    eggPath.Close()

    // 获取实际的剪辑区域
    eggArea := clip.Outline{Path: eggPath.End()}.Op()

    // 对形状内部填充颜色
    // color := color.NRGBA{R: 255, G: 239, B: 174, A: 255}
    color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress)), B: uint8(174 * (1 - progress)), A: 255}
    paint.FillShape(gtx.Ops, color, eggArea)

    d := image.Point{Y: 375}
    return layout.Dimensions{Size: d}
  },
),

```

主要思想是定义一个蛋形的 `clip.Path`，我们通过绘制一条线来实现，之后填充内部，忽略外部的任何绘制。

## 注解

首先，我们定义了新的路径：`var eggPath clip.Path`。

然后通过 `op.Offset()` 将点向右移动了 200dp、向下移动了 150dp，与之前一样，是从这个小部件内部的左上角开始计算的。有一点需要注意，我们不使用像素作为单位，而是将其转换为 Dp（设备独立像素），以确保用户在不同设备和分辨率上的体验是较为一致的。

现在我们位于蛋的中心，这是路径的起点，`eggPath.Begin()`。

从这里，我们旋转了整整 360 度，并继续绘制蛋的轮廓。我们使用了一个 [Hügelschäffer蛋](https://mathcurve.com/courbes2d.gb/oeuf/oeuf.shtml) 的方程，如 Torben Jansen [出色的交互式博客](https://observablehq.com/@toja/egg-curve)上所示，这个方程接受一个角度，计算出蛋壳到蛋中心的距离，从 0° 到 360°，它计算出了一堆点作为轮廓。数学真有意思！

![Hügelschäffer 蛋](09_torben_jansen.gif)

关于 Gio，最重要的是 for 循环中的最后一行的 `eggPath.LineTo(p)`，在这里，程序已经计算出了蛋形函数 360 度循环中的下一个点 `p`，我们使用 `eggPath.LineTo` 将画笔移到这个特定的坐标点来为两点连线。

完成了 for 循环后，蛋形轮廓就画完了，我们通过显式调用 `eggPath.Close()` 来最终完成它，它会关闭路径。

路径完成后，我们想要获取路径内部的区域，`clip.Outline{}` 创建的操作表示这个区域。

现在我们该为蛋上色了，颜色可以是静态的，但如果蛋的颜色从冷色渐变到暖色岂不是很酷？我是这么觉得的。记住进度是一个从 0 到 1 的变量，现在这个状态变量可以逐渐地改变蛋的颜色了，`* (1 - progress)` 只是另一种方式来表示 `“逐渐去除绿色和蓝色”`，当倒计时完成时，两者都是 0，我们只剩下红色。醒目！

最后，我们返回 `layout.Dimensions` ，这个小部件的高度。

---

[下一章](10_input_boiltime_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }