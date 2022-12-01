---
layout: default
title: Teleprompter
nav_order: 3
has_children: true
has_toc: false
---

# Teleprompter - Animation and interaction

## Goals

This project continues where the [egg timer](../egg_timer/) leaves off. The timer was a good start and gave us the foundation to build an app. But we're not done. Especially we should look closer at how how to deal with user input, both keyboard and mouse.

![Mr_Gorbachev_tear_down_this_wall](teleprompter_Mr_Gorbachev.gif)

To do that we'll build what's known as a [teleprompter](https://en.wikipedia.org/wiki/Teleprompter). A teleprompter is a device that displays and scrolls text. TV studios buy these for thousands of dollars. But we, armed with Go and Gio, scoff at such largesse. Why buy when you can build, huh? 

Our app should

1. Read text from a `.txt` file. The speaker should be able to display personal files with ease.
1. An **auto scroll** that is easy to start, stop, pause, speed up and slow down. 
1. Some **manual scroll** as well, so the user can roll up or down a speech at will. Thus we must read input from f.ex. a mousepad or scroll-wheel.
1. Easily adjust **font size** and **text width**, this should work on screens both big and small.
1. A **focus bar** would be nice, helping the user easily read the text. It should be easy to position. 
1. Gesticulation friendly, i.e. all controls can be done with one hand. 


---

**Please note:**
We're not doing politics here. When building this, it was simply more interesting with famous speeches. No matter your origin or partisanship, [Ronald Reagan's 1987 Berlin Wall Speech](https://www.nytimes.com/video/world/europe/100000006815729/reagan-berlin-wall.html) played a part in ending the cold war. Both sides, East and West, deserve respect and praise for seeking and finding a mostly peaceful solution. We must learn from our past, repeat the best, and fight the worst. The world needs this. 


On that note, let's (sc)roll!

---

[Let's start](01_setup.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }
