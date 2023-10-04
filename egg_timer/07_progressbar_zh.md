---
layout: default
title: 第七章 - 进度条
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第七章 - 实现一个进度条

## 目标

本章节的目标是添加一个进度条。

![Progressbar](07_progressbar.gif)

## 大纲

自从开始编写这个系列以来，我一直期待着这一章。此次将涵盖相当多的内容并引入多个新的想法：

- 尝试一个新的小部件：`material.Progressbar`
- 开始使用状态变量来控制行为
- 使用两种并发技术：一种用于创建和共享一个稳定的脉冲用于增长进度条，另一种用于在独立的通信操作中进行 _选择(select{})_

让我们依次看看这些内容。

## 功能一：进度条

显然，进度条是一根显示进度的条形图，但是显示什么进度？以及如何控制它？它应该以多快的速度增长，是否可以暂停，甚至可以倒退？从[文档](https://pkg.go.dev/gioui.org/widget/material?utm_source=gopls#ProgressBar)中我们可以看到 `ProgressBar(th *Theme, progress float32)` 接收 0 到 1 之间的十进制数作为进度。

### 代码

我们在 `main()` 外部声明进度变量，以便只设置一次并且可以在整个程序中访问：

```go
// 根级别，不在 main() 内
var progress float32
```

为了布局进度条，我们依然使用强大的 Flexbox 并通过刚性(Rigid)区域添加它：

```go
// 在 case System.FrameEvent: 内部
layout.Flex{
  // ...
}.Layout(gtx,
  layout.Rigid(
    func(gtx C) D {
      bar := material.ProgressBar(th, progress)  // 这里使用了progress变量
      return bar.Layout(gtx)
    },
  ),

```

重点来了，小部件本身没有状态，状态在程序的其他部分中维护，小部件只知道如何显示我们发送给它的进度，任何增加、暂停、倒退或重置的逻辑都在小部件外部控制。

## 功能一：状态变量

我们刚刚提到了 `progress`，一个包含状态的变量，另外一个有用的状态是开始按钮是否已被点击，在我们的程序中就是用于表示鸡蛋是否开始加热（或者说是否开始倒计时）。

### 代码

```go
// 鸡蛋是否正在加热？
var boiling bool
```
我们希望在单击开始按钮时切换该布尔值，因此，我们侦听来自 GUI 的 `system.FrameEvent` 并检查 `startButton.Clicked()` 是否为 true。

```go
case system.FrameEvent:
  gtx := layout.NewContext(&ops, e)
  // 让我们尝试一下 flexbox 布局概念
  if startButton.Clicked() {
    boiling = !boiling
  }
```

再次声明，按钮唯一的任务就只是告诉我们它刚刚被点击了，除此之外的任何操作都要由程序的其余部分来执行。

比如按下前后按钮上的文本应该是什么。我们在调用 `material.Button()` 函数之前通过首先检查 `boiling` 的状态来决定。

```go
// ... 之前用于创建按钮的相同函数
func(gtx C) D {
  var text string
  if !boiling {
    text = "Start"
  } else {
    text = "Stop"
  }
  btn := material.Button(th, &startButton, text)
  return btn.Layout(gtx)
},
```

## 功能一：稳定的脉冲

一个良好的进度条必须平稳而准确地增长。为了实现这一点，我们首先创建一个单独的 goroutine，该协程提供稳定的脉冲，稍后，在我们侦听事件的地方，我们会捕获这些脉冲并增长进度条。

### 代码

首先是脉冲发生器：

```go
// 定义进度变量，一个通道和一个变量
var progressIncrementer chan float32
var progress float32

func main() {
  // 初始化通道以提供增加进度的脉冲 (这里原文是"Setup a separate channel to..." 怪怪的)
  progressIncrementer = make(chan float32)
  go func() {
    for {
      time.Sleep(time.Second / 25)
      progressIncrementer <- 0.004
    }
  }()

  // ...
```

`progressIncrementer` 是我们发送值的[通道](https://tour.golang.org/concurrency/2)，在这种情况下，值的类型为 `float32`。

同样，这是在匿名函数中实现的，调用时开始执行，这意味着这个循环会在整个程序的执行过程中一直运行，每过 1/25 秒，一个 0.004 会被传入通道。

稍后我们通过 `draw(w *app.window)` 中的以下代码从通道中接收值：

```go
  // .. 在 draw() 里面
  for {
    select {
      // 侦听窗口中的事件
      case e := <-w.Events():
        // ...

      // 侦听脉冲发生器通道中的事件
      case p := <-progressIncrementer:
        if boiling && progress < 1 {
        progress += p
        w.Invalidate()
      }
    }
  }

```

在之前的章节中，我们使用 `for e := range w.Events()` 对事件进行了遍历，在这里，我们使用了一个带有 [select](https://tour.golang.org/concurrency/5) 的循环，这是 Go 的并发特性，`select` 在等待一个可以执行的 `case` 语句时会保持阻塞。

- 事件可以是来自窗口的，此时使用 `e := <- w.Events()` 来获取它
- 如果事件来自进度脉冲，则使用 `p := <- progressIncrementer` 来获取它

如果控制变量 `boiling` 为 true 且 `progress` 小于 1，我们就将 `p` 增加到 `progress` 变量中。由于 `p` 为 0.004 且进度每秒增加 25 次，所以需要在 10 秒之后 `progress` 才能达到 1。可以随意调整这两个参数，以找到你喜欢的速度和平滑度。

最后，我们通过调用 `w.Invalidate()` 强制窗口绘制，这个方法的作用是告诉 Gio 旧的绘图现在已经 `“无效”`，必须进行新的绘制。如果没有这样的通知，Gio 只有在通过点击鼠标、按下按钮或其他事件强制绘制时才会更新。每帧都执行无效化的代价会有些高昂，也有一些替代方法存在，不过这是一个略微高级的话题，所以现在让我们先把它放一边，在 [附加内容 - 改进动画](11_improved_animation_zh.md) 中再回头来看。

使用这样的通道，我们实现了：

1. 精确的时间控制，可以精确地控制执行时间
2. 一致的时间控制，在高低端硬件上获得几乎无差的结果
3. 并发的时间控制，让程序的其余部分继续运行

虽然所有这些都是有道理的，但第二点值得额外注意，如果重新编译应用程序时去掉 `time.Sleep(time.Second / 25)`，您的计算机将全力以赴地尽快渲染，这可能会消耗大量的 CPU 资源，从而还可能会耗尽电池，它还确保设备上的时间控制一致，所有设备都以相同的速度进行脉冲。例如，包含来自 3 台不同计算机的 pprof 的文件包含了一个 1/25 秒的延迟，确保了相同的最终结果。请查看一下。

**更新**

7月28日，[Elias Naur 发布了一个更新](https://lists.sr.ht/~eliasnaur/gio/%3CCD3XWVXUTCG0.23LAQED4PF674%40themachine%3E)，优化了渲染速度：

> _gpu: [compute] 缓存并重用前一帧的绘制操作：这一变化实现了一种自动分层方案，使得只有帧中发生变更的部分需要计算，没有这种优化，CPU fallback 就不会可行。_

这里也有更详细的解释：[7月社区会议](https://www.youtube.com/watch?v=HC4Cg78l-9U)

## 注解

通过组合所有这些概念与技术，我们现在有了一个易于控制的有状态程序，用户界面告诉我们发生了什么事情，程序的其余部分使用这一信息来处理业务。我们不得不使用一些技巧，包括 `channel` 和 `select`。现在我们已经将这些技巧加入到知识库中，我们将能够在下一章中添加一些自定义图形。

---

[下一章](08_egg_as_circle_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }