---
layout: default
title: 第六章 - 边距
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第六章 - 带边距的下置按钮

## 目标

本章节的目标是在按钮的四周添加空白空间。

![Button with margin](06_button_low_margin.png)

## 大纲

在上一节重构完代码之后，这次我们只关注发生变化的行。再次声明，下面的这一切都发生在 `layout.Flex` 内。

## 代码 - 总体结构

为了突出结构，去掉一些细节可能会有用。
这里实际上只有三行关键代码：

1. 使用 `layout.Inset` 定义边距
2. 布局这些边距
3. 在这些边距内部定义并布局按钮

```go
layout.Flex{
  // ...
}.Layout(gtx,
  layout.Rigid(
    func(gtx C) D {
      // 首先使用 layout.Inset 定义按钮周围的边距...
      margins := layout.Inset{
        // ...
      }

      // ... 然后我们布局这些边距...
      margins.Layout(

        // ... 最后，在这些边距内部，我们定义并布局按钮
        func(gtx C) D {
          btn := material.Button(th, &startButton, "Start")
          return btn.Layout(gtx)
        },

      )

      }

  )
)

```

## 注解

上面这些就像一个带有按钮的甜甜圈。—— _有些比喻是有用的_，还记得吧？

![Button inside inset](06_button_inside_inset.jpeg)

边距是使用 [layout.Inset{}](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Inset) 创建的，它是一个结构体，用于定义小部件周围的空间：

```go
margins := layout.Inset{
  Top:    unit.Dp(25),
  Bottom: unit.Dp(25),
  Right:  unit.Dp(35),
  Left:   unit.Dp(35),
}
```

在这里，边距以 _设备独立像素_ **D**evice independent **p**ixels ([unit.Dp](https://pkg.go.dev/gioui.org/unit?utm_source=gopls#Dp)) 的形式给出，如果你想要四边等距，还可以直接用 `UniformInset()`。

## 代码 - 细节

最后，这里是整个 `case system.FrameEvent:` 的代码：

```go
case system.FrameEvent:
  gtx := layout.NewContext(&ops, e)
  // 让我们尝试一下 flexbox 布局概念
  layout.Flex{
    // 垂直对齐，从上到下
    Axis: layout.Vertical,
    // 空白空间位于开头，即顶部
    Spacing: layout.SpaceStart,
  }.Layout(gtx,
    layout.Rigid(
      func(gtx C) D {
        // 首先使用 layout.Inset 定义按钮周围的边距...
        margins := layout.Inset{
          Top:    unit.Dp(25),
          Bottom: unit.Dp(25),
          Right:  unit.Dp(35),
          Left:   unit.Dp(35),
        }
        // ... 然后我们布局这些边距...
        return margins.Layout(gtx,
          // ... 最后，在这些边距内部，我们定义并布局按钮
          func(gtx C) D {
            btn := material.Button(th, &startButton, "Start")
            return btn.Layout(gtx)
          },
        )
      },
    ),
  )
  e.Frame(gtx.Ops)

```

---

[下一章](07_progressbar_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }