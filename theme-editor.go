package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"git.sr.ht/~whereswaldon/sprig/icons"
	sprigTheme "git.sr.ht/~whereswaldon/sprig/widget/theme"

	"git.sr.ht/~whereswaldon/colorpicker"
)

type ThemeEditorView struct {
	manager ViewManager

	PrimaryDefault colorpicker.State
	PrimaryDark    colorpicker.State
	PrimaryLight   colorpicker.State

	SecondaryDefault colorpicker.State
	SecondaryDark    colorpicker.State
	SecondaryLight   colorpicker.State

	BackgroundDefault colorpicker.State
	BackgroundDark    colorpicker.State
	BackgroundLight   colorpicker.State

	TextColor         colorpicker.State
	HintColor         colorpicker.State
	InvertedTextColor colorpicker.State

	ColorsList layout.List
	listElems  []colorListElement

	AncestorMux    colorpicker.MuxState
	DescendantMux  colorpicker.MuxState
	SelectedMux    colorpicker.MuxState
	SiblingMux     colorpicker.MuxState
	NonselectedMux colorpicker.MuxState

	MuxList      layout.List
	muxListElems []muxListElement

	*sprigTheme.Theme
	widgetTheme *material.Theme
}

type colorListElement struct {
	*colorpicker.State
	Label        string
	TargetColors []*color.RGBA
}

type muxListElement struct {
	*colorpicker.MuxState
	Label       string
	TargetColor **color.RGBA
}

var _ View = &ThemeEditorView{}

func NewThemeEditorView(theme *sprigTheme.Theme) View {
	c := &ThemeEditorView{
		Theme:       theme,
		widgetTheme: material.NewTheme(gofont.Collection()),
	}
	c.PrimaryDefault.SetColor(c.Theme.Primary.Default)
	c.PrimaryDark.SetColor(c.Theme.Primary.Dark)
	c.PrimaryLight.SetColor(c.Theme.Primary.Light)
	c.SecondaryDefault.SetColor(c.Theme.Secondary.Default)
	c.SecondaryDark.SetColor(c.Theme.Secondary.Dark)
	c.SecondaryLight.SetColor(c.Theme.Secondary.Light)
	c.BackgroundDefault.SetColor(c.Theme.Background.Default)
	c.BackgroundDark.SetColor(c.Theme.Background.Dark)
	c.BackgroundLight.SetColor(c.Theme.Background.Light)

	c.ColorsList.Axis = layout.Vertical
	c.listElems = []colorListElement{
		{
			Label: "Primary",
			TargetColors: []*color.RGBA{
				&c.Theme.Primary.Default,
				&c.Theme.Theme.Color.Primary,
			},
			State: &c.PrimaryDefault,
		},
		{
			Label: "Primary Light",
			TargetColors: []*color.RGBA{
				&c.Theme.Primary.Light,
			},
			State: &c.PrimaryLight,
		},
		{
			Label: "Primary Dark",
			TargetColors: []*color.RGBA{
				&c.Theme.Primary.Dark,
			},
			State: &c.PrimaryDark,
		},
		{
			Label: "Secondary",
			TargetColors: []*color.RGBA{
				&c.Theme.Secondary.Default,
			},
			State: &c.SecondaryDefault,
		},
		{
			Label: "Secondary Light",
			TargetColors: []*color.RGBA{
				&c.Theme.Secondary.Light,
			},
			State: &c.SecondaryLight,
		},
		{
			Label: "Secondary Dark",
			TargetColors: []*color.RGBA{
				&c.Theme.Secondary.Dark,
			},
			State: &c.SecondaryDark,
		},
		{
			Label: "Background",
			TargetColors: []*color.RGBA{
				&c.Theme.Background.Default,
			},
			State: &c.BackgroundDefault,
		},
		{
			Label: "Background Light",
			TargetColors: []*color.RGBA{
				&c.Theme.Background.Light,
			},
			State: &c.BackgroundLight,
		},
		{
			Label: "Background Dark",
			TargetColors: []*color.RGBA{
				&c.Theme.Background.Dark,
			},
			State: &c.BackgroundDark,
		},
		{
			Label: "Text",
			TargetColors: []*color.RGBA{
				&c.Theme.Theme.Color.Text,
			},
			State: &c.TextColor,
		},
		{
			Label: "Hint",
			TargetColors: []*color.RGBA{
				&c.Theme.Theme.Color.Hint,
			},
			State: &c.HintColor,
		},
		{
			Label: "Inverted Text",
			TargetColors: []*color.RGBA{
				&c.Theme.Theme.Color.InvText,
			},
			State: &c.InvertedTextColor,
		},
	}

	muxOptions := []colorpicker.MuxOption{}
	for _, elem := range c.listElems {
		if len(elem.TargetColors) < 1 || elem.TargetColors[0] == nil {
			continue
		}
		elem.SetColor(*elem.TargetColors[0])
		muxOptions = append(muxOptions, colorpicker.MuxOption{
			Label: elem.Label,
			Value: elem.TargetColors[0],
		})
	}
	c.muxListElems = []muxListElement{
		{
			Label:       "Ancestors",
			MuxState:    &c.AncestorMux,
			TargetColor: &c.Theme.Ancestors,
		},
		{
			Label:       "Descendants",
			MuxState:    &c.DescendantMux,
			TargetColor: &c.Theme.Descendants,
		},
		{
			Label:       "Selected",
			MuxState:    &c.SelectedMux,
			TargetColor: &c.Theme.Selected,
		},
		{
			Label:       "Siblings",
			MuxState:    &c.SiblingMux,
			TargetColor: &c.Theme.Siblings,
		},
		{
			Label:       "Unselected",
			MuxState:    &c.NonselectedMux,
			TargetColor: &c.Theme.Unselected,
		},
	}
	for _, mux := range c.muxListElems {
		*mux.MuxState = colorpicker.NewMuxState(muxOptions...)
	}

	return c
}

func (c *ThemeEditorView) BecomeVisible() {
}

func (c *ThemeEditorView) NavItem() *materials.NavItem {
	return &materials.NavItem{
		Name: "Theme",
		Icon: icons.CancelReplyIcon,
	}
}

func (c *ThemeEditorView) AppBarData() (bool, string, []materials.AppBarAction, []materials.OverflowAction) {
	return true, "Theme", []materials.AppBarAction{}, []materials.OverflowAction{}
}

func (c *ThemeEditorView) HandleClipboard(contents string) {
}

func (c *ThemeEditorView) Update(gtx layout.Context) {
	for _, elem := range c.listElems {
		if elem.Changed() {
			for _, target := range elem.TargetColors {
				*target = elem.Color()
			}
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}
	for _, elem := range c.muxListElems {
		if elem.Changed() {
			*elem.TargetColor = elem.Color()
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}
}

func (c *ThemeEditorView) Layout(gtx layout.Context) layout.Dimensions {
	return c.layoutPickers(gtx)
}

func (c *ThemeEditorView) layoutPickers(gtx layout.Context) layout.Dimensions {
	return c.ColorsList.Layout(gtx, len(c.listElems)+1, func(gtx C, index int) D {
		if index == len(c.listElems) {
			return c.layoutMuxes(gtx)
		}
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return sprigTheme.DrawRect(gtx, color.RGBA{A: 255}, f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}, 0)
			}),
			layout.Stacked(func(gtx C) D {
				return layout.UniformInset(unit.Dp(3)).Layout(gtx, func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Expanded(func(gtx C) D {
							return sprigTheme.DrawRect(gtx, color.RGBA{R: 255, G: 255, B: 255, A: 255}, f32.Point{
								X: float32(gtx.Constraints.Min.X),
								Y: float32(gtx.Constraints.Min.Y),
							}, 0)
						}),
						layout.Stacked(func(gtx C) D {
							elem := c.listElems[index]
							return colorpicker.Picker(c.widgetTheme, elem.State, elem.Label).Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (c *ThemeEditorView) layoutMuxes(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return sprigTheme.DrawRect(gtx, color.RGBA{R: 255, G: 255, B: 255, A: 255}, f32.Point{
				X: float32(gtx.Constraints.Min.X),
				Y: float32(gtx.Constraints.Min.Y),
			}, 0)
		}),
		layout.Stacked(func(gtx C) D {
			return c.MuxList.Layout(gtx, len(c.muxListElems), func(gtx C, index int) D {
				element := c.muxListElems[index]
				return colorpicker.Mux(c.widgetTheme, element.MuxState, element.Label).Layout(gtx)
			})
		}),
	)
}

func (c *ThemeEditorView) SetManager(mgr ViewManager) {
	c.manager = mgr
}
