---
layout: default
title: 第四章 - 布局
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第四章 - 将按钮下置

## 目标

显然，按钮不应该填满整个屏幕，让我们将按钮移到底部，为此，我们开始使用名为 [Flexbox](https://pkg.go.dev/gioui.org/layout#Flex) 的布局概念。

![A low button with a spacer below](04_button_low.png)

## 大纲

上一章讲解了程序的整体结构，现在，我们将深入研究 `system.FrameEvent` 并开始使用 Flexbox，如果您对此还不熟悉，请首先学习一下，例如[来自 Mozilla 的这篇文章](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox)。

## 整体结构

我们不会在这里重复所有代码，而是重点关注 `system.FrameEvent`：

### 代码

首先，我们移除掉细节以更好地查看结构：

```go
case system.FrameEvent:

  layout.Flex{
  // ...
  }.Layout( // ...
    // 我们插入了两个刚性元素：
    // 首先是一个按钮 ...
    layout.Rigid(),
    // ... 然后是一个空白间隔
    layout.Rigid(),
  }
```

### 注解

来看看这段代码的结构：

1. 首先，我们通过结构体 `layout.Flex{}` 定义了一个 `Flexbox`。
2. 然后，我们通过 `Layout(gtx, ...)` 向其发送一个要布局的 _子元素列表_（即gtx），图形上下文 _gtx_ 包含了子元素必须遵循的约束，并且可以包含任意数量的子元素。

我们列出的子元素都是由 `layout.Rigid()` 创建的：

1. 第一个是按钮的占位符，
2. 另外一个是用于包含按钮下方的空白间隔的占位符。

那么 [Rigid(刚性)](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Rigid) 是什么呢？简单来说，它的任务是填充其获得的空间。Rigid(刚性) 子元素首先被布局，而 [Flexed(弹性)](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Flexed) 子元素共享剩余空间，除此之外的子元素的位置按照它们的定义顺序决定。

#### 约束(Constraint) 和 尺寸(Dimensions)

在这一点上，让我们回顾一下将所有这些内容联系在一起的概念，即 **约束(C)** 和 **尺寸(D)**。

- [约束](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Constraints) 是小部件的最小和最大尺寸。小部件可以有多小或多大。
- [尺寸](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Dimensions) 是小部件的实际尺寸。它实际上有多小或多大。

父容器设置 **约束**，子元素用 **尺寸** 响应。父容器创建一个小部件并调用 `Layout()`，小部件用自己的尺寸响应，有效地进行自我布局。就与现实世界一样，不是所有的孩子都表现良好，正如任何孩子都可以证明父母的一些约束可能不太公平 —— 因此需要一些微妙的协商。在大部分情况下，**约束** 和 **尺寸** 就这样将所有内容联系在一起。

正如上面所示，布局操作是递归的，一个子元素可以拥有自己的子元素，布局(layout)本身可以包含布局(layout)，这一切都可以继续下去，从而实现用简单的组件构建复杂的结构，同时一切都可以追溯到根本。

## 详细代码

前面的那些内容层次有点高，现在是时候深入了解了，让我们详细看看 `system.FrameEvent`：

### 代码

```go
case system.FrameEvent:
  gtx := layout.NewContext(&ops, e)
  // 让我们尝试一下 flexbox 布局概念：
  layout.Flex{
    // 垂直对齐，从上到下
    Axis: layout.Vertical,
    // 空白空间位于开头，即顶部
    Spacing: layout.SpaceStart,
  }.Layout(gtx,
    // 我们插入了两个刚性元素：
    // 首先是一个按钮 ...
    layout.Rigid(
      func(gtx layout.Context) layout.Dimensions {
        btn := material.Button(th, &startButton, "Start")
        return btn.Layout(gtx)
      },
    ),
    // ... 然后是一个空白间隔
    layout.Rigid(
      // 空白间隔的高度为25dp
      layout.Spacer{Height: unit.Dp(25)}.Layout,
    ),
  )
  e.Frame(gtx.Ops)
```

### 注解

在 `layout.Flex{}` 内部，我们定义了两个属性：

1. Axis: 垂直对齐表示元素将放置在彼此上下
2. Spacing: 剩余空间将位于开头

![Button above spacer](04_button_above_spacer.jpg)

由于顺序很重要，您可以将其视为小部件从屏幕底部弹出，按钮先到达，然后空白间隔从底部出现并将按钮向上推。—— _所有的比喻都是错误的，但有些是有用的_。

现在让我们来看看对 `layout.Rigid()` 的两次调用：

- Rigid 接受一个 [小部件(Widget)](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Widget)。
- 一个小部件实际上就是返回它自己 **尺寸(Dimensions)** 的东西。
- 如何实现并不重要。下面是两种非常不同的方式：
  1. 在第一个 Rigid 中，我们传入了一个 `func()`，该函数通过 `btn.Layout()` 返回 **尺寸**。
  2. 在第二个 Rigid 中，我们创建了一个 `Spacer{}` 结构，调用其 `Layout` 方法，从而获得 **尺寸**。
- 从父容器的角度来看，这并不重要，只要子元素返回了 **尺寸** 就可以了。

这样就完成了小部件的布局，但这些小部件究竟是什么呢？

- 正如名称所示，`material.Button` 是基于 Material Design 的 [按钮](https://pkg.go.dev/gioui.org/widget/material?utm_source=gopls#Button)，就像我们在上一章详细介绍的那样。

- [Spacer](https://pkg.go.dev/gioui.org/layout#Spacer) 添加了空白间隔，这里通过 _Height_ 字段来定义。由于我们将整体布局定义为垂直，多余的空间位于顶部，所以这个空白间隔和其上面的按钮会出现在底部。就这样，我们创建了一些空间，使按钮从屏幕底部升起了一点。

~~(好工作)~~ **干得好**
我们已经走了很远的路程，学到了很多东西，感谢您一直的坚持，现在让我们继续深入研究代码库。

---

[下一章](05_button_low_refactored_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
