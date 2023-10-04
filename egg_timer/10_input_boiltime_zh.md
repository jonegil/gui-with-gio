---
layout: default
title: 第十章 - 输入
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第十章 - 设置煮沸时间

## 目标

本章节的目标是添加一个输入框以设置煮沸时间。

![The complete egg timer](egg_timer.gif)

## 大纲

代码在几个方面进行了更改

1. 导入来自 Gio 的 [gioui.org/text](https://pkg.go.dev/gioui.org/text) 包，以及标准库中的字符串和数字操作
2. 添加第四个刚性(Rigid)区域来放一个 [widget.Editor() (编辑器小部件)](https://pkg.go.dev/gioui.org/widget#Editor)
3. 为按钮添加一些逻辑，让它看起来稍微更合理些

就这么多，看码：

## 代码

### 1. 新的包

```go
import (
  "fmt"
  "strconv"
  "strings"

  "gioui.org/text"
)
```

标准库相当分散，用于字符串和数字操作的有用功能来自这些包：

- `fmt` 将用于将浮点数转换为字符串
- `strconv` 将用于将字符串转换为浮点数
- `strings` 将用于修剪输入字符串中的空格

来自 Gio 的包：

- [gioui.org/text](https://pkg.go.dev/gioui.org/text) 提供了用于处理文本的支持类型，其中很多是用于字体支持和缓存，不过我们仅使用它来进行文本对齐。

### 2. 编辑器小部件

编辑器小部件就是用于输入蛋应该煮多长时间的输入框。

**一些用于输入框的变量**

与按钮一样，我们需要一个用于输入框本身的变量。因此，我们首先声明了一个 [widget.Editor](https://pkg.go.dev/gioui.org/widget#Editor) 变量。

我们还创建了一个变量用于保存输入字段中的实际数字值，并将其命名为 `boilDuration`，另外，要知道这些变量之间不存在魔法让它们直接同步起来，我们稍后将编写代码实现从输入字段中读取值，并将这些值存储在 `boilDuration` 中。慢慢来，一切都会在适当的时候进行。

有了这些，我们现在可以在 `draw()` 函数的顶部找到以下行：

```go
  // boilDurationInput 是一个用于输入煮沸时长的文本输入框
  var boilDurationInput widget.Editor

  // 鸡蛋是否正在加热？
  var boiling bool
  var boilDuration float32
```

### 3. 从输入框中读取

我们需要检查输入框内容的唯一时间点是用户点击开始按钮时，因此，我们将逻辑放在这一个 `if{}` 代码块内。

```go
if startButton.Clicked() {
  //...

  // 从输入框中读取
  inputString := boilDurationInput.Text()
  inputString = strings.TrimSpace(inputString)
  inputFloat, _ := strconv.ParseFloat(inputString, 32)
  boilDuration = float32(inputFloat)
  boilDuration = boilDuration / (1 - progress)
}
```

前几行挺直白的：

- `boilDurationInput.Text()` 返回输入框内的文本字符串
- `strings.TrimSpace()` 从字符串中删除开头结尾的空格字符（如果有的话）
- `strconv.ParseFloat()` 将文本转换为浮点数，注意第二个参数 `bitsize` 是 32，根据[标准库文档](https://pkg.go.dev/strconv#ParseFloat)：
  - > _ParseFloat将字符串 s 转换为具有由 bitSize 指定的精度的浮点数：32 表示 float32，64 表示 float64。当 bitSize=32 时，输出的结果仍然是 float64 类型，但可以在不改变其值的情况下转换为 float32。_
  - 啊哈，所以我们需要明确将其转换为 `float32`

最后，有一个将 `progress` 与 `boilDuration` 相关联的技巧。例如，如果一次煮沸完成了 20%，用户输入了新的时间 10 秒，我们可以合理地假设用户是在希望额外增加 10 秒，而不是继续倒数完成了 20% 之后剩下的 8 秒，此时我们可以通过除以 `(1-progress)` 将其放大到 12.5。

还有其他解决方案，比如重新调整进度条，但需要通过调整 `progress` 状态变量。为了简单起见，我们在这里跳过了这一步，但在你的程序中需要注意状态变量是否还与其他逻辑有关联。

### 4. 布局所有元素

现在，让我们向世界呈现我们刚写的新功能。在 flexbox 中，我们为输入框创建一个单独的刚性区域，由于它位于鸡蛋下方和进度条上方，它将成为这四个布局中的第二个：

```go
layout.Flex{
  // 垂直对齐，从上到下
  Axis: layout.Vertical,
  // 空白空间位于开头，即顶部
  Spacing: layout.SpaceStart,
}.Layout(gtx,
  // 1. 鸡蛋
  layout.Rigid(
    //...
  )
  // 2. 输入框
  layout.Rigid(
    // 在此处添加用于显示输入框的新代码
  )
  // 3. 进度条
  layout.Rigid(
    //...
  )
  // 4. 按钮
  layout.Rigid(
    //...
  )
)
```

### 5. 深入一下输入框

有了概述之后，现在让我们看看第二个刚性部分：

**自定义样式的输入框**

我们首先将 `boilDurationInput` 变量包装在 Material Design 主题中，我们借此机会添加一个[输入提示](https://pkg.go.dev/gioui.org/widget/material#EditorStyle)。

```go
// 输入框
layout.Rigid(
  func(gtx C) D {
    // 将编辑器包装在 Material Design 中
    ed := material.Editor(th, &boilDurationInput, "sec")
```

**定义属性**

此时，`boilDurationInput` 仍然只是一个空字段，因此我们将进行一些配置：

```go
    // 定义输入框的属性
    boilDurationInput.SingleLine = true
    boilDurationInput.Alignment = text.Middle
```

- `SingleLine` 强制输入框始终为一行高，否则，用户在按回车键时输入框会变大
- `Alignment` 使文本在框内居中

**倒计时**

因为查看剩余时间很有用，所以接下来我们在输入框内显示倒计时：

```go
if boiling && progress < 1 {
  boilRemain := (1 - progress) * boilDuration
  // 格式化为 1 位小数
  inputStr := fmt.Sprintf("%.1f", math.Round(float64(boilRemain)*10)/10)
  // 更新输入框中的文本
  boilDurationInput.SetText(inputStr)
}
```

当我们在煮沸过程中，我们在这里定义一个新的 `boilRemain`，它保存煮沸完成前的剩余时间，计算方式为 `(1-progress)`。

由于 [math.Round()](https://pkg.go.dev/math#Round) 不允许四舍五入到指定的小数位数，我们必须使用一个小技巧。

- 首先乘以 10。
- 然后将小数部分四舍五入为 0。
- 然后除以 10。
- 最后将结果转换为带有 1 位小数的文本。
  <br>有点紧凑，但应该很直观。

最后，再次使用 Gio，我们调用
[SetText](https://pkg.go.dev/gioui.org/widget#Editor.SetText)，用 `inputStr` 替换框中的文本。

**布局**

输入框做好了，现在我们开始进行布局：

```go
    // 定义边距 ...
    margins := layout.Inset{
      Top:    unit.Dp(0),
      Right:  unit.Dp(170),
      Bottom: unit.Dp(40),
      Left:   unit.Dp(170),
    }
    // ... 和边框 ...
    border := widget.Border{
      Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
      CornerRadius: unit.Dp(3),
      Width:        unit.Dp(2),
    }
    // ... 然后嵌套布局，一个包含在另一个内部
    return margins.Layout(gtx,
      func(gtx C) D {
        return border.Layout(gtx, ed.Layout)
      },
    )
  },
),
```

1. 定义边距，这里左边和右边设置得比较大，以使输入框差不多刚好够用于输入两三位数字。
1. 定义一个带有自定义颜色和圆角的边框。
1. 结合 1+2 并返回尺寸(Dimensions)。

有没有发现最终的布局现在就像一个俄罗斯套娃一样？边距包含边框，边框包含编辑器，它们都返回自己的 `layout.Dimensions`：

### 6. 进度条

最后，如果上一次倒计时已经完成，进度条将被重置，这能让我们连续煮好几个鸡蛋，很棒！

为了将其呈现给用户，我们扩展了绘制按钮的代码：

```go
func(gtx C) D {
  var text string
  if !boiling {
    text = "Start"
  }
  if boiling && progress < 1 {
    text = "Stop"
  }
  if boiling && progress >= 1 {
    text = "Finished"
  }
  btn := material.Button(th, &startButton, text)
  return btn.Layout(gtx)
},
```

- 如果不在倒计时状态，则显示 "Start"
- 如果正在倒计时但尚未完成，则显示 "Stop"
- 如果正在倒计时且已完成，则显示 "Finished"

可以在这里添加更多功能，例如，当倒计时结束时，尝试更改[背景颜色](https://pkg.go.dev/gioui.org/widget/material#ButtonStyle)吧！

## 最后的注释

就这么多了，感谢您的陪伴，我希望这篇教程能让您尝试自己动手开发 GUI 程序。

我们只初步涉足了 Gio，这个框架中还有比我们在这里展示的更多功能，不过我们都已经一起煮了那个蛋，这是不是意味着我们走了很远？

如果您喜欢这篇教程，请在 Github 上给它加一个星标，收到这些鼓励的象征会极大地激励我。

还有，~~如果~~当您着手编写自己的项目，不论大小，都请给我来封信，我很想听听您的想法。

此外，请确保在 Gio 的网站、时事通讯和社区上关注 Gio。

现在，终于到了早餐的时间，来猜猜我要吃什么。

---

等等 —— 不止这些，在完成了这十章之后，我感觉还有一些额外的功能可能会引起您的兴趣，准备好了吗？

[好好好！](11_improved_animation_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
