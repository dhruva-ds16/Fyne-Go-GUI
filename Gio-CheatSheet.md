This code snippet starts up the application loop that talks to the OS and starts the windows logic in a separate goroutine

```
func main() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
```

This code snippet creates a theme which contains all the fonts and colour settings

```
func run(w *app.Window) error {
	th := material.NewTheme()
```

Events:

Communication with the OS happens through events. DestroyEvent means that the user pressed the close button and FrameEvent means that the program should handle the input and render a new frame

```
for {
	switch e := w.NextEvent().(type) {
	case app.DestroyEvent:
		return e.Err
	case app.FrameEvent:
```

Text Drawing:

To draw text, the application needs to first define a new graphics context and then define the properties of the text. Then, the application should pass the operations to the GPU

```
// This graphics context is used for managing the rendering state.
gtx := app.NewContext(&ops, e)

// Define an large label with an appropriate text:
title := material.H1(th, "Hello, Gio")

// Change the color of the label.
maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
title.Color = maroon

// Change the position of the label.
title.Alignment = text.Middle

// Draw the label to the graphics context.
title.Layout(gtx)

// Pass the drawing operations to the GPU.
e.Frame(gtx.Ops)
```

---

Splitting the window:

```
type SplitVisual struct{}

func (s SplitVisual) Layout(gtx layout.Context, left, right layout.Widget) layout.Dimensions {
	leftsize := gtx.Constraints.Min.X / 2
	rightsize := gtx.Constraints.Min.X - leftsize

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftsize, gtx.Constraints.Max.Y))
		left(gtx)
	}

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
		trans := op.Offset(image.Pt(leftsize, 0)).Push(gtx.Ops)
		right(gtx)
		trans.Pop()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
```

Usage:

```
func exampleSplitVisual(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return SplitVisual{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Left", red)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Right", blue)
	})
}

func FillWithLabel(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions {
	ColorBox(gtx, gtx.Constraints.Max, backgroundColor)
	return layout.Center.Layout(gtx, material.H3(th, text).Layout)
}
```

Split Ratio

```
type SplitRatio struct {
	// Ratio keeps the current layout.
	// 0 is center, -1 completely to the left, 1 completely to the right.
	Ratio float32
}

func (s SplitRatio) Layout(gtx layout.Context, left, right layout.Widget) layout.Dimensions {
	proportion := (s.Ratio + 1) / 2
	leftsize := int(proportion * float32(gtx.Constraints.Max.X))

	rightoffset := leftsize
	rightsize := gtx.Constraints.Max.X - rightoffset

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftsize, gtx.Constraints.Max.Y))
		left(gtx)
	}

	{
		trans := op.Offset(image.Pt(rightoffset, 0)).Push(gtx.Ops)
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
		right(gtx)
		trans.Pop()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
```

Usage

```
func exampleSplitRatio(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return SplitRatio{Ratio: -0.3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Left", red)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Right", blue)
	})
}
```