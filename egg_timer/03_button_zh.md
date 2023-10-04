---
layout: default
title: 第三章 - 按钮
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第三章 - 实现一个按钮

## 目标

本章节的目标是添加一个按钮，能按是当然的，但它还会有优雅的悬停、点击动画。

![A button](03_button.gif)

## 大纲

本节将介绍许多新组件，我们不会深入探讨，重点是程序的整体结构。不要过分注重细节，专注于大局就好。

我们首先看一下导入的新包，有点多，所以这里得花点时间，之后，我们将看一下 `operations` 和 `widgets` 如何结合在一起创建按钮。

在最后，我们会涉及到 [Material Design](https://material.io/)，这是一种在 Gio 中也可用的成熟的用户界面框架。

为了使事情井井有条，让我们先讨论包的导入，然后再讨论主函数。

## 新包导入

### 代码

```go
import (
  "os"

  "gioui.org/app"
  "gioui.org/io/system"
  "gioui.org/layout"
  "gioui.org/op"
  "gioui.org/unit"
  "gioui.org/widget"
  "gioui.org/widget/material"
)
```

### 注解

`os`、`app` 和 `unit` 我们之前已经了解了，其余的：

- [io/system](https://pkg.go.dev/gioui.org/io/system) - 提供从窗口传来的高级事件，最主要的是 `system.FrameEvent`，它请求一个新的帧，新的帧通过一系列操作来定义，比如要显示什么、如何处理输入。

- [layout](https://pkg.go.dev/gioui.org/layout) - 定义布局的有用部分，例如 _尺寸(dimensions)_、_约束(constraints)_ 和 _方向(directions)_。此外，它还包括广泛用于Web和用户界面开发的布局概念，即 [Flexbox](https://pkg.go.dev/gioui.org/layout#Flex)，这是非常常见的 Web 和用户界面开发技术，在众多的介绍中，我推荐从 [Mozilla](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox) 那里了解其他有关的信息。

- [op](https://pkg.go.dev/gioui.org/op) - 在 Gio 中，_操作(oprations/ops)_ 是核心，它们用于更新界面，有用于绘制、处理输入、更改窗口属性、缩放、旋转等操作。有趣的是，还有一种称为 [macros](https://pkg.go.dev/gioui.org/op#MacroOp) 的东西，可以记录要稍后执行的操作。总的来说，操作列表是一个 _可变堆栈_，你可以控制其流程。

- [widget](https://pkg.go.dev/gioui.org/widget) - 提供了 UI 组件的底层功能，如状态跟踪和事件处理：鼠标是否悬停在按钮上？它被点击了吗，如果是的话点击了多少次？

- [widget/material](https://pkg.go.dev/gioui.org/widget/material) - `widget` 提供了功能，而 `widget/material` 定义了一个主题，重点来了，界面实际上分为两部分：

  1. 实际的小部件，具有状态
  1. 小部件的绘制，完全无状态

  这是有意为之的，目的是提高小部件的可重用性和灵活性，我们稍后会使用到这一点。

  默认主题看起来很好，这也是我们将要使用的，但通过设置属性（颜色、文本大小、字体属性等），也可以很轻松地调整它。

  - 另外，Gio 在一个专用的仓库中扩展了基本功能，称为 [gio-x](https://pkg.go.dev/gioui.org/x)，其中有开发中的[更多 Material Design 组件](https://pkg.go.dev/gioui.org/x/component)，包括导航栏和工具提示。

## 主函数

包导完了，上码：（变长了点，不过还是很简单）

### 代码

```go
func main() {
  go func() {
    // 创建新窗口
    w := app.NewWindow(
      app.Title("Egg timer"),
      app.Size(unit.Dp(400), unit.Dp(600)),
    )

    // ops 是用户界面的操作
    var ops op.Ops

    // startButton 是一个可点击的小组件
    var startButton widget.Clickable

    // th 定义了 Material Design 样式
    th := material.NewTheme()

    // 在窗口中监听事件
    for e := range w.Events() {

      // 判断事件的类型
      switch e := e.(type) {

      // 当应用程序需要重新渲染时发送的事件
      case system.FrameEvent:
        gtx := layout.NewContext(&ops, e)
        btn := material.Button(th, &startButton, "Start")
        btn.Layout(gtx)
        e.Frame(gtx.Ops)
      }
    }
    os.Exit(0)
  }()
  app.Main()
}
```

### 注解

1. 从头开始，我们首先定义了主函数并调用了匿名的 goroutine。

2. 然后，我们定义了窗口 `w`。

3. 设置了三个新变量：

- `ops` 是用户界面的操作。

- `startButton` 就是我们的按钮，一个可点击的小部件。

- `th` 是 Material Design 主题。
  - 更新于 Gio 0.2 (2023.07)：
    - Gio 现在默认使用系统字体，少了一件需要记住的事情。如果你想用别的字体，请阅读[2023年7月时事通讯](https://gioui.org/news/2023-07)。
    - 另外，你可能想了解一下 Go 自家专用的高质量 TrueType 字体？看看这篇[有趣的博客](https://blog.golang.org/go-fonts)，而且一定要访问作者 [Bigelow & Holmes](https://bigelowandholmes.typepad.com) 的网站。

4. `for e := range w.Events() ` 循环更有意思：
  - `w.Events()` 获取了传递事件的 _通道_，我们只需一直监听这个通道。

  - 然后...这个 `e := e.(type)` 是什么鬼？实际上这是一个很巧妙的语法，称为[类型切换(type switch)](https://tour.golang.org/methods/16)，它能让我们根据正在处理的事件的 `类型(type)` 采取不同的操作。

  - 在我们的情况下，我们只关心事件是否是 `system.FrameEvent` 类型，如果是的话：

    - 我们定义了一个新的 _图形上下文(graphical context/gtx)_，它接收指向 `ops` 以及 `事件` 的指针。

    - `btn` 被声明为实际用于渲染的按钮，具有主题 `th` 和指向 `startButton` 小部件的指针，还有要显示的文本（文本纯粹是显示在按钮上的内容，并不是按钮中的有状态小部件 `startButton` 的一部分）。

    - 需要注意，按钮 `btn` 需要在图形上下文 `gtx` 上进行 _布局_，这很关键，布局由按钮自己进行，这实现了在调整窗口大小时，按钮会自动重新布局，无论画布的大小或形状如何。

    - 此时，我们直接得到了所有的鼠标悬停和点击动画，它们全包含在主题中。

    - 最后，我们通过实际将操作 `ops` 从上下文 `gtx` 发送到 FrameEvent `e` 来完成渲染呈现。

5. 最后我们需要调用 `app.Main()`，别忘了。

真是一个长篇大论，如果你还在继续阅读，感谢你的陪伴，最后我们可以用三行总结整个章节：

```go
  gtx := layout.NewContext(&ops, e)
  b := material.Button(th, &startButton, "Start")
  b.Layout(gtx)
```

如果你可以很轻松地阅读这些，那就完全没问题了。

---

[下一章](04_button_low_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
