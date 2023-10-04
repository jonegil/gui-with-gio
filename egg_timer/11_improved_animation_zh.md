---
layout: default
title: 附加 - 动画
nav_order: 3
parent: 煮蛋计时器
has_children: false
---

# 附加内容 - 改进动画

## 目标

本章节的目标是讨论与动画相关的稍微高级点的内容，包括：何时、如何使帧失效，这玩意具体是什么以及如何编写逻辑。

![The complete egg timer](egg_timer.gif)

## 大纲

附加章节的大纲如下：

- 首先，我们讨论什么是使帧失效
- 然后，我们来看两种不同的方法来使帧失效
- 最后，我们讨论一种替代模式，用于生成和控制动画

### 1. 什么是使帧失效？

Gio 只在产生 [FrameEvent](https://pkg.go.dev/gioui.org/io/system#FrameEvent) 时才会更新您看到的内容，例如按下按键、点击鼠标、小部件获得或失去焦点时，这是非常合理的，因为现代设备的刷新率普遍可以达到每秒120帧，且大部分情况下显示的内容与上一帧完全相同。

不过 `“大部分”` 并不是 `“所有”`

比如在进行动画绘制时，您希望它尽可能平滑地运行，为了实现这一点，我们需要让 Gio 持续重绘，并且是在不触发任何事件的情况下，所以我们需要明确告诉 Gio 需要这样做，这就是通过调用 `invalidate` 来完成的。

### 2. 两种使帧失效的方式

让我们一起看看：

- [op.InvalidateOp{}.Add(ops)](https://pkg.go.dev/gioui.org/op#InvalidateOp) 是最高效的，可以用于请求立即或将来重绘 `At time.Time`。
- [window.Invalidate](https://pkg.go.dev/gioui.org/app#Window.Invalidate) 的效率较低，用于外部触发的事件，但如果您想在布局代码之外使其失效，那么这也是正确的选项，比如在[第七章 - 进度条](07_progressbar_zh.html)中的那个全局、独立的进度脉冲发生器。
  (原文是 tick-generator，不太好翻译，这里全用 “脉冲” 来代替 "tick")

#### 示例一： op.InvalidateOp{}

为了展示 `op.InvaliateOp{}.Add(ops)`，我们将引用来自[架构文档](https://gioui.org/doc/architecture#animation)中写得很好的动画示例：

```go
// 来源：https://gioui.org/doc/architecture#animation

var startTime = time.Now()
var duration = 10 * time.Second

func drawProgressBar(ops *op.Ops, now time.Time) {
  // 基于当前时间计算要绘制的进度条的多少
  // 这是根据当前时间计算的
  elapsed := now.Sub(startTime)
  progress := elapsed.Seconds() / duration.Seconds()
  if progress < 1 {
    // 进度条尚未完成动画
    op.InvalidateOp{}.Add(ops)
  } else {
    progress = 1
  }

  defer op.Save(ops).Load()
  width := 200 * float32(progress)
  clip.Rect{Max: image.Pt(int(width), 20)}.Add(ops)
  paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
  paint.ColorOp{Color: color.NRGBA{G: 0x80, A: 0xFF}}.Add(ops)
  paint.PaintOp{}.Add(ops)
}
```

根据您的布局复杂性，这个示例可能会稍微有点占用资源，然而，在[优化渲染速度的更新](https://lists.sr.ht/~eliasnaur/gio/%3CCD3XWVXUTCG0.23LAQED4PF674%40themachine%3E)之后，系统负载会减小，如果这适用于您的应用程序，那么您可以放心使用它。

#### 示例二： w.Invalidate()

但是，如果您更喜欢使用全局脉冲 `progressIncrementer` 来控制帧速率，那么在[第七章 - 进度条](07_progressbar_zh.html)中有一个 `w.Invalidate()` 的示例，以及在 gioui.org 首页上的 [kitchen](https://github.com/gioui/gio-example/blob/main/kitchen/kitchen.go) 中也有一个示例。前者的代码：

```go
func main() {
  // 初始化通道以提供增加进度的脉冲 (这里原文是"Setup a separate channel to..." 怪怪的)
  progressIncrementer = make(chan float32)
  go func() {
    for {
      time.Sleep(time.Second / 25)
      progressIncrementer <- 0.004
    }
  }()
  // ... main的其余部分
}

func draw(w *app.Window) error {
  // ...
  for {
    select {
    // 处理FrameEvent并生成布局 ...

    // 在布局之外，监听增加器通道中的事件
    case p := <-progressIncrementer:
      if boiling && progress < 1 {
        progress += p
        w.Invalidate()
      }
    }
  }
```

#### 示例三： 用 op.InvalidateOp{} 代替 w.Invalidate() —— 会有什么效果

可以更好吗？为什么不使用 `IncrementOp{At time.Time}` 代替全局脉冲呢？我们需要将它移到 FrameEvent 内的 `Layout` 部分。以下是该刚性(Rigid)区域的代码：

```go
// 进度条
layout.Rigid(
  func(gtx C) D {
    bar := material.ProgressBar(th, progress)
    if boiling && progress < 1 {
      op.InvalidateOp{At: gtx.Now.Add(time.Second / 25)}.Add(&ops)
    }
    return bar.Layout(gtx)
  },
 )

```

查看动画中的代码。您可以在 [line 201](https://github.com/jonegil/gui-with-gio/blob/fc54ae4394fe92f79934e816bf54ac800e703daa/egg_timer/code/11_improved_animation/main.go#L201) 找到上述的 `op.InvalidateOp{}`，并在 [line 255](https://github.com/jonegil/gui-with-gio/blob/fc54ae4394fe92f79934e816bf54ac800e703daa/egg_timer/code/11_improved_animation/main.go#L255) 找到旧的`w.Invalidate()`。尝试更改其中一个，以查看哪个执行效果更好。

我尝试运行了三次 60 秒的倒计时，其中一个使用了每种使帧失效的方法，分别在我的 Macbook 和 Windows PC 上。在Mac上，一个是使用旧渲染器运行的，一个是通过 `GIORENDERER=forcecompute go run main.go` 运行的。

在没有使用 forcecompute 的情况下，在 2017 款 Macbook Air 上运行时，`op.InvalidateOp{}` 以大约 16-17% 的 CPU 占用运行 ，而 `w.Invalidate()` 则消耗了大约 18-19%，使用 forcecompute 之后，这些数据降低至大约 12% 和 15%，负载水平有点高，但差异并不是很大。不过，了解每种使帧失效技术的效果是值得的，在我的 Windows PC 上，负载都要小得许多，且两种技术之间没有明显的区别。

![Invalidate CPU load](11_invalidate_cpu_load.png)

### 一个通用的模式

我的朋友 Chris Waldon 提出了这种模式：

> _在我的 Gio 程序中，我发现了一种模式，我认为它很适合封装动画逻辑。它允许您在 API 后面隐藏帧失效的管理，并将进度作为时间函数进行计算。它完全去除了用于发生脉冲的 goroutine，尽管我有点担心它会使程序变得不那么酷。_

听起来很有趣，对吧？您可以在[他的仓库](https://github.com/whereswaldon/gui-with-gio/commit/83e43a39e75c5e6cb96985046a521ac553615d39)中找到源码，不过我们也在这里稍微聊一下：

#### 一个 animation 结构体

首先，用一个结构体取代状态变量`boiling`和`boilDuration`，该结构体储存我们何时开始以及应该持续多长时间：

```go
// 动画跟踪多帧间线性动画的进度
type animation struct {
  start    time.Time
  duration time.Duration
}

var anim animation
```

这个结构体允许我们创建方法来渲染下一帧和最终结束的动画，同时也是调用 `op.InvalidateOp{}.Add()` 的地方。

另外，注意它使用 `gtx.Now` 而不是 `time.Now`，最重要的目的是确保动画同步，也避免了来自 `time.Now` 的一些开销。

```go
// animate 在当前帧开始一个持续指定时间的动画
func (a *animation) animate(gtx layout.Context, duration time.Duration) {
  a.start = gtx.Now
  a.duration = duration
  op.InvalidateOp{}.Add(gtx.Ops)
}

// stop 立即结束动画
func (a *animation) stop() {
  a.duration = time.Duration(0)
}
```

最后一个方法用于报告进度：

- 动画是否仍在进行？
- 如果是，它目前进行了多少？

```go
// progress 返回动画当前是否正在运行以及它的完成百分比(如果在跑的话)
func (a animation) progress(gtx layout.Context) (animating bool, progress float32) {
  if gtx.Now.After(a.start.Add(a.duration)) {
    return false, 0
  }
  op.InvalidateOp{}.Add(gtx.Ops)
  return true, float32(gtx.Now.Sub(a.start)) / float32(a.duration)
}

```

#### 开始简化

有了这个结构，我们可以简化 `startButton.Clicked()` 的代码：

```go
case system.FrameEvent:
  gtx := layout.NewContext(&ops, e)
  boiling, progress := anim.progress(gtx)

  if startButton.Clicked() {
    // 开始(或停止)加热
    if boiling {
      anim.stop()
    } else {
      // 从输入框中读取
      inputString := boilDurationInput.Text()
      inputString = strings.TrimSpace(inputString)
      inputFloat, _ := strconv.ParseFloat(inputString, 32)
      anim.animate(gtx, time.Duration(inputFloat)*time.Second)
    }
  }

```

#### 整理未使用的代码

最后，由于我们删除了 `boilDuration` 变量，因此我们使用 `anim.duration.Seconds()` ：

```go
// 计算倒计时的文本
if boiling && progress < 1 {
  boilRemain := (1 - progress) * float32(anim.duration.Seconds())
```

#### 那之前的脉冲发生器呢？

与 `progressIncrementer` 通道相关的所有代码，包括变量、读写通道，都被删除了。

我们不会在这里更详细地讨论该模式，但我们知道有这种方法存在，且它还具有一些非常不错的功能，比如可以处理动画的状态。

## 注解

总结一切，我希望这为 Gio 中的动画提供了更多的信息，像往常一样，这取决于什么是最好的解决方案，展示一个令人印象深刻的演示？追求流畅的动画，最小化投入？分阶段动画，高端与低端用户硬件？要优雅还是要简洁，都取决于你的需求。

只需要权衡一下并了解一些有帮助的技巧。

---
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
