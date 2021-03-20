package ctk

import (
	"fmt"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
)

// CDK type-tag for Box objects
const TypeBox cdk.CTypeTag = "ctk-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeBox, func() interface{} { return MakeBox() })
}

// Box Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Box
//	        +- ButtonBox
//	        +- VBox
//	        +- HBox
type Box interface {
	Container
	Buildable
	Orientable

	Init() (already bool)
	PackStart(child Widget, expand, fill bool, padding int)
	PackEnd(child Widget, expand, fill bool, padding int)
	Remove(w Widget)
	GetOrientation() (orientation cdk.Orientation)
	SetOrientation(orientation cdk.Orientation)
	GetHomogeneous() (value bool)
	SetHomogeneous(homogeneous bool)
	GetSpacing() (value int)
	SetSpacing(spacing int)
	ReorderChild(child Widget, position int)
	QueryChildPacking(child Widget) (expand bool, fill bool, padding int, packType PackType)
	SetChildPacking(child Widget, expand bool, fill bool, padding int, packType PackType)
	Build(builder Builder, element *CBuilderElement) error
	ShowAll()
	Add(child Widget)
	GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool)
	GetSizeRequest() (width, height int)
	Resize() cdk.EventFlag
	Draw(canvas cdk.Canvas) cdk.EventFlag
}

// The CBox structure implements the Box interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Box objects
type CBox struct {
	CContainer

	canvases   map[int]cdk.Canvas
	showHandle string
	hideHandle string
}

type cBoxChild struct {
	widget   Widget
	expand   bool
	fill     bool
	padding  int
	packType PackType
}

func MakeBox() (box *CBox) {
	box = NewBox(cdk.ORIENTATION_HORIZONTAL, false, 0)
	return
}

func NewBox(orientation cdk.Orientation, homogeneous bool, spacing int) *CBox {
	b := new(CBox)
	b.Init()
	b.SetOrientation(orientation)
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

// Box object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Box instance
func (b *CBox) Init() (already bool) {
	if b.InitTypeItem(TypeBox, b) {
		return true
	}
	b.CContainer.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE)
	b.SetFlags(APP_PAINTABLE)
	b.canvases = make(map[int]cdk.Canvas)
	b.showHandle = fmt.Sprintf("%v.show", b.ObjectName())
	b.hideHandle = fmt.Sprintf("%v.hide", b.ObjectName())
	_ = b.InstallBuildableProperty(PropertyDebugChildren, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyOrientation, cdk.StructProperty, true, cdk.ORIENTATION_HORIZONTAL)
	_ = b.InstallBuildableProperty(PropertyHomogeneous, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertySpacing, cdk.IntProperty, true, 0)
	_ = b.InstallChildProperty(PropertyBoxChildPackType, cdk.StructProperty, true, PackStart)
	_ = b.InstallChildProperty(PropertyBoxChildExpand, cdk.BoolProperty, true, false)
	_ = b.InstallChildProperty(PropertyBoxChildFill, cdk.BoolProperty, true, true)
	_ = b.InstallChildProperty(PropertyBoxChildPadding, cdk.IntProperty, true, 0)
	return false
}

// PackStart
//
// Adds child to box, packed with reference to the start of box. The child is
// packed after any other child packed with reference to the start of box.
//
// Parameters
//
// child     the Widget to be added to box
// expand    TRUE if the new child is to be given extra space allocated to box.
//           The extra space will be divided evenly between all children of box
//           that use this option
// fill      TRUE if space given to child by the expand option is actually
//           allocated to child, rather than just padding it. This parameter has
//           no effect if expand is set to FALSE. A child is always allocated
//           the full height of an HBox and the full width of a VBox. This
//           option affects the other dimension
// padding   extra space in pixels to put between this child and its neighbors,
//           over and above the global amount specified by spacing property.
//           If child is a widget at one of the reference ends of box , then
//           padding pixels are also put between child and the reference edge of
//           box
func (b *CBox) PackStart(child Widget, expand, fill bool, padding int) {
	b.LogDebug("PackStart(%v,%v,%v,%v)", child, expand, fill, padding)
	if f := b.Emit(SignalPackStart, b, child, expand, fill, padding); f == cdk.EVENT_PASS {
		child.Connect(SignalShow, b.showHandle, func([]interface{}, ...interface{}) cdk.EventFlag {
			child.LogDebug("signal show, resize: %v", b.ObjectName())
			child.SetFlags(VISIBLE)
			b.Resize()
			return cdk.EVENT_STOP
		})
		child.Connect(SignalHide, b.hideHandle, func([]interface{}, ...interface{}) cdk.EventFlag {
			child.LogDebug("signal hide, resize: %v", b.ObjectName())
			child.UnsetFlags(VISIBLE)
			b.Resize()
			return cdk.EVENT_STOP
		})
		child.SetParent(b)
		child.SetWindow(b.GetWindow())
		b.CContainer.AddWithProperties(child,
			PropertyBoxChildPackType, PackStart,
			PropertyBoxChildExpand, expand,
			PropertyBoxChildFill, fill,
			PropertyBoxChildPadding, padding,
		)
		b.canvases[child.ObjectID()] = cdk.NewCanvas(cdk.Point2I{}, cdk.Rectangle{}, child.GetThemeRequest().Content.Normal)
		b.Resize()
	}
}

// PackEnd
//
// Adds child to box, packed with reference to the end of box. The child is
// packed after (away from end of) any other child packed with reference to the
// end of box.
//
// Parameters
//
// child     the Widget to be added to box
// expand    TRUE if the new child is to be given extra space allocated to box.
//           The extra space will be divided evenly between all children of box
//           that use this option
// fill      TRUE if space given to child by the expand option is actually
//           allocated to child , rather than just padding it. This parameter
//           has no effect if expand is set to FALSE. A child is always
//           allocated the full height of an HBox and the full width of a VBox.
//           This option affects the other dimension
// padding   extra space in pixels to put between this child and its neighbors,
//           over and above the global amount specified by spacing property.
//           If child is a widget at one of the reference ends of box, then
//           padding pixels are also put between child and the reference edge of
//           box
func (b *CBox) PackEnd(child Widget, expand, fill bool, padding int) {
	b.LogDebug("PackEnd(%v,%v,%v,%v)", child, expand, fill, padding)
	if f := b.Emit(SignalPackEnd, b, child, expand, fill, padding); f == cdk.EVENT_PASS {
		child.Connect(SignalShow, b.showHandle, func([]interface{}, ...interface{}) cdk.EventFlag {
			child.LogDebug("signal show, resize: %v", b.ObjectName())
			child.SetFlags(VISIBLE)
			b.Resize()
			return cdk.EVENT_STOP
		})
		child.Connect(SignalHide, b.hideHandle, func([]interface{}, ...interface{}) cdk.EventFlag {
			child.LogDebug("signal hide, resize: %v", b.ObjectName())
			child.UnsetFlags(VISIBLE)
			b.Resize()
			return cdk.EVENT_STOP
		})
		child.SetParent(b)
		child.SetWindow(b.GetWindow())
		b.CContainer.AddWithProperties(child,
			PropertyBoxChildPackType, PackEnd,
			PropertyBoxChildExpand, expand,
			PropertyBoxChildFill, fill,
			PropertyBoxChildPadding, padding,
		)
		b.canvases[child.ObjectID()] = cdk.NewCanvas(cdk.Point2I{}, cdk.Rectangle{}, child.GetThemeRequest().Content.Normal)
		b.Resize()
	}
}

func (b *CBox) Remove(w Widget) {
	if err := b.Disconnect(SignalShow, b.showHandle); err != nil {
		b.LogErr(err)
	}
	if err := b.Disconnect(SignalHide, b.hideHandle); err != nil {
		b.LogErr(err)
	}
	b.CContainer.Remove(w)
}

// Returns the orientation of the Box
func (b *CBox) GetOrientation() (orientation cdk.Orientation) {
	var ok bool
	if v, err := b.GetStructProperty(PropertyOrientation); err != nil {
		b.LogErr(err)
	} else if orientation, ok = v.(cdk.Orientation); !ok && v != nil {
		b.LogError("invalid value stored in %v: %v (%T)", PropertyOrientation, v, v)
	}
	return
}

// Sets the orientation of the Box
func (b *CBox) SetOrientation(orientation cdk.Orientation) {
	if err := b.SetStructProperty(PropertyOrientation, orientation); err != nil {
		b.LogErr(err)
	}
	b.Resize()
}

// Returns whether the box is homogeneous (all children are the same size).
// See SetHomogeneous.
// Returns:
// 	TRUE if the box is homogeneous.
func (b *CBox) GetHomogeneous() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyHomogeneous); err != nil {
		b.LogErr(err)
	}
	return
}

// Sets the “homogeneous” property of box , controlling whether or not
// all children of box are given equal space in the box.
// Parameters:
// 	homogeneous	a boolean value, TRUE to create equal allotments,
// FALSE for variable allotments
func (b *CBox) SetHomogeneous(homogeneous bool) {
	if err := b.SetBoolProperty(PropertyHomogeneous, homogeneous); err != nil {
		b.LogErr(err)
	}
}

// Gets the value set by SetSpacing.
// Returns:
// 	spacing between children
func (b *CBox) GetSpacing() (value int) {
	var err error
	if value, err = b.GetIntProperty(PropertySpacing); err != nil {
		b.LogErr(err)
	}
	return
}

// Sets the “spacing” property of box , which is the number of pixels to
// place between children of box .
// Parameters:
// 	spacing	the number of pixels to put between children
func (b *CBox) SetSpacing(spacing int) {
	if err := b.SetIntProperty(PropertySpacing, spacing); err != nil {
		b.LogErr(err)
	}
}

// Moves child to a new position in the list of box children. The list is the
// children field of Box, and contains both widgets packed GTK_PACK_START
// as well as widgets packed GTK_PACK_END, in the order that these widgets
// were added to box . A widget's position in the box children list
// determines where the widget is packed into box . A child widget at some
// position in the list will be packed just after all other widgets of the
// same packing type that appear earlier in the list.
// Parameters:
// 	child	    the Widget to move
// 	position	the new position for child in the list of children of box,
// 	            starting from 0. If negative, indicates the end of the list
func (b *CBox) ReorderChild(child Widget, position int) {
	var children []Widget
	if position < 0 {
		position = len(b.children) - 1 + position
	}
	for idx, c := range b.children {
		if idx == position {
			children = append(children, child)
		} else if c.ObjectID() != child.ObjectID() {
			children = append(children, c)
		}
	}
	b.children = children
}

// Obtains information about how child is packed into box .
// Parameters:
// 	child	the Widget of the child to query
// 	expand	pointer to return location for “expand” child property
// 	fill	pointer to return location for “fill” child property
// 	padding	pointer to return location for “padding” child property
// 	packType	pointer to return location for “pack-type” child property
func (b *CBox) QueryChildPacking(child Widget) (expand bool, fill bool, padding int, packType PackType) {
	if cps, ok := b.property[child.ObjectID()]; ok {
		for _, cp := range cps {
			switch cp.Name() {
			case PropertyBoxChildExpand:
				if v, ok := cp.Value().(bool); ok {
					expand = v
				}
			case PropertyBoxChildFill:
				if v, ok := cp.Value().(bool); ok {
					fill = v
				}
			case PropertyBoxChildPadding:
				if v, ok := cp.Value().(int); ok {
					padding = v
				}
			case PropertyBoxChildPackType:
				if v, ok := cp.Value().(PackType); ok {
					packType = v
				}
			}
		}
	} else {
		b.LogError("%v is not a child of %v", child, b)
	}
	return
}

// Sets the way child is packed into box .
// Parameters:
// 	child	the Widget of the child to set
// 	expand	the new value of the “expand” child property
// 	fill	the new value of the “fill” child property
// 	padding	the new value of the “padding” child property
// 	packType	the new value of the “pack-type” child property
func (b *CBox) SetChildPacking(child Widget, expand bool, fill bool, padding int, packType PackType) {
	if cps, ok := b.property[child.ObjectID()]; ok {
		for _, cp := range cps {
			switch cp.Name() {
			case PropertyBoxChildExpand:
				if err := cp.Set(expand); err != nil {
					b.LogErr(err)
				}
			case PropertyBoxChildFill:
				if err := cp.Set(fill); err != nil {
					b.LogErr(err)
				}
			case PropertyBoxChildPadding:
				if err := cp.Set(padding); err != nil {
					b.LogErr(err)
				}
			case PropertyBoxChildPackType:
				if err := cp.Set(packType); err != nil {
					b.LogErr(err)
				}
			}
		}
	} else {
		b.LogError("%v is not a child of %v", child, b)
	}
}

func (b *CBox) Build(builder Builder, element *CBuilderElement) error {
	b.Freeze()
	defer b.Thaw()
	if err := b.CObject.Build(builder, element); err != nil {
		return err
	}
	for _, child := range element.Children {
		if newChild := builder.Build(child); newChild != nil {
			child.Instance = newChild
			if newChildWidget, ok := newChild.(Widget); ok {
				newChildWidget.Show()
				if len(child.Packing) > 0 {
					expand, fill, padding, packType := builder.ParsePacking(child)
					if packType == PackStart {
						b.PackStart(newChildWidget, expand, fill, padding)
					} else {
						b.PackEnd(newChildWidget, expand, fill, padding)
					}
				} else {
					b.Add(newChildWidget)
				}
				if newChildWidget.HasFlags(HAS_FOCUS) {
					newChildWidget.GrabFocus()
				}
			} else {
				b.LogError("new child object is not a Widget type: %v (%T)")
			}
		}
	}
	return nil
}

// The Container type implements a version of Widget.ShowAll() where all the
// children of the Container have their ShowAll() method called, in addition to
// calling Show() on itself first.
func (b *CBox) ShowAll() {
	b.Show()
	for _, child := range b.GetChildren() {
		child.ShowAll()
	}
}

func (b *CBox) Add(child Widget) {
	b.PackStart(child, false, true, 0)
}

func (b *CBox) GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool) {
	if b.focusChainSet {
		return b.focusChain, true
	}
	var children []interface{}
	for _, child := range b.getBoxChildren() {
		if child.packType == PackStart {
			children = append(children, child.widget)
		}
	}
	for _, child := range b.getBoxChildren() {
		if child.packType == PackEnd {
			children = append(children, child.widget)
		}
	}
	for _, child := range children {
		if cc, ok := child.(Container); ok {
			fc, _ := cc.GetFocusChain()
			for _, cChild := range fc {
				focusableWidgets = append(focusableWidgets, cChild)
			}
		} else if cw, ok := child.(Widget); ok {
			if cw.CanFocus() && cw.IsVisible() {
				focusableWidgets = append(focusableWidgets, child)
			}
		}
	}
	return
}

func (b *CBox) getBoxChildren() (children []*cBoxChild) {
	bChildren := b.GetChildren()
	nChildren := len(bChildren)
	expand := make([]bool, nChildren)
	fill := make([]bool, nChildren)
	padding := make([]int, nChildren)
	packType := make([]PackType, nChildren)
	for idx, child := range bChildren {
		expand[idx], fill[idx], padding[idx], packType[idx] = b.QueryChildPacking(child)
		if child.IsVisible() && packType[idx] == PackStart {
			children = append(children, &cBoxChild{
				widget:   child,
				expand:   expand[idx],
				fill:     fill[idx],
				padding:  padding[idx],
				packType: packType[idx],
			})
		}
	}
	for idx, child := range bChildren {
		if child.IsVisible() && packType[idx] == PackEnd {
			children = append(children, &cBoxChild{
				widget:   child,
				expand:   expand[idx],
				fill:     fill[idx],
				padding:  padding[idx],
				packType: packType[idx],
			})
		}
	}
	return
}

// func (b *CBox) GetSizeRequest() (width, height int) {
// 	children := b.getBoxChildren()
// 	nChildren := len(children)
// 	if nChildren <= 0 {
// 		return
// 	}
// 	orientation := b.GetOrientation()
// 	isVertical := orientation == cdk.ORIENTATION_VERTICAL
// 	spacing := b.GetSpacing()
// 	var w, h int
// 	if b.GetHomogeneous() {
// 		// get the size of the largest child and request that for all children
// 		for _, child := range children {
// 			req := cdk.MakeRectangle(child.widget.GetSizeRequest())
// 			if w < req.W {
// 				w = req.W
// 				if !isVertical && child.padding > 0 {
// 					w += child.padding * 2
// 				}
// 			}
// 			if h < req.H {
// 				h = req.H
// 				if isVertical && child.padding > 0 {
// 					h += child.padding * 2
// 				}
// 			}
// 		}
// 		if isVertical {
// 			width = w
// 			height = (nChildren * h) + utils.FloorI((nChildren-1)*spacing, 0)
// 		} else {
// 			width = (nChildren * w) + utils.FloorI((nChildren-1)*spacing, 0)
// 			height = h
// 		}
// 		return
// 	}
// 	// add up the sizes of all children, including spacing and child padding
// 	sizes := make([]*cdk.Rectangle, nChildren)
// 	tally := cdk.NewRectangle(0, 0)
// 	for idx, child := range children {
// 		sizes[idx] = cdk.NewRectangle(child.widget.GetSizeRequest())
// 		sizes[idx].Floor(0, 0)
// 		if w < sizes[idx].W {
// 			w = sizes[idx].W
// 			if !isVertical && child.padding > 0 {
// 				w += child.padding * 2
// 			}
// 		}
// 		tally.W += sizes[idx].W
// 		if !isVertical && child.padding > 0 {
// 			tally.W += child.padding * 2
// 		}
// 		if h < sizes[idx].H {
// 			h = sizes[idx].H
// 			if isVertical && child.padding > 0 {
// 				h += child.padding * 2
// 			}
// 		}
// 		tally.H += sizes[idx].H
// 		if !isVertical && child.padding > 0 {
// 			tally.H += child.padding * 2
// 		}
// 	}
// 	if isVertical {
// 		width = w
// 		height = tally.H
// 	} else {
// 		width = tally.W
// 		height = h
// 	}
// 	return
// }

func (b *CBox) Resize() cdk.EventFlag {
	children := b.getBoxChildren()
	numChildren := len(children)
	if numChildren == 0 {
		b.Emit(SignalResize, b)
		return cdk.EVENT_STOP
	}
	spacing := b.GetSpacing()
	origin := b.GetOrigin().NewClone()
	alloc := b.GetAllocation().NewClone()
	orientation := b.GetOrientation()
	isVertical := orientation == cdk.ORIENTATION_VERTICAL
	homogeneous := b.GetHomogeneous()
	// intermediaries
	var increment int
	var gaps []int
	if isVertical {
		increment, gaps = utils.SolveSpaceAlloc(numChildren, alloc.H, spacing)
	} else {
		increment, gaps = utils.SolveSpaceAlloc(numChildren, alloc.W, spacing)
	}
	nextPoint := origin.NewClone()
	if homogeneous {
		return b.resizeHomogeneous(isVertical, gaps, increment, numChildren, origin, nextPoint, alloc, children)
	}
	return b.resizeDynamicAlloc(isVertical, gaps, increment, spacing, numChildren, origin, nextPoint, alloc, children)
}

func (b *CBox) resizeHomogeneous(isVertical bool, gaps []int, increment, numChildren int, origin, nextPoint *cdk.Point2I, alloc *cdk.Rectangle, children []*cBoxChild) cdk.EventFlag {
	// assume child.expand == true
	var consumed int
	tracking := make([]struct {
		x, y, w, h int
		rw, rh     int
		aw, ah     int
		extra      int
		overflow   int
	}, numChildren)

	// first: build up tracking dataset

	for idx, child := range children {
		req := cdk.NewRectangle(child.widget.GetSizeRequest())
		if child.fill {
			if isVertical {
				tracking[idx].w = alloc.W
				tracking[idx].h = increment
			} else { // horizontal
				tracking[idx].w = increment
				tracking[idx].h = alloc.H
			}
			tracking[idx].aw = tracking[idx].w
			tracking[idx].ah = tracking[idx].h
		} else { // !child.fill
			if isVertical {
				tracking[idx].w = alloc.W
				tracking[idx].aw = alloc.W
				tracking[idx].ah = increment
				if req.H <= -1 && req.H > increment {
					tracking[idx].h = increment
					req.H = increment
				} else {
					tracking[idx].h = req.H
				}
			} else { // horizontal
				tracking[idx].h = alloc.H
				tracking[idx].aw = increment
				tracking[idx].ah = alloc.H
				if req.W <= -1 && req.W > increment {
					tracking[idx].w = increment
					req.W = increment
				} else {
					tracking[idx].w = req.W
				}
			} // if isVertical
		} // if child.fill
		req.Floor(0, 0)
		tracking[idx].rw = req.W
		tracking[idx].rh = req.H
	}
	for idx, _ := range children {
		if isVertical {
			consumed += tracking[idx].ah
		} else {
			consumed += tracking[idx].aw
		}
	} // for each child

	// solve positions

	for idx, child := range children {
		if isVertical {
			if tracking[idx].h < tracking[idx].ah {
				delta := tracking[idx].ah - tracking[idx].h
				if ca, ok := child.widget.(Alignable); ok {
					_, yAlign := ca.GetAlignment()
					pad := int(float64(delta) * yAlign)
					tracking[idx].y += pad
					tracking[idx].overflow += delta - pad
				} else {
					tracking[idx].overflow += delta
				}
			}
		} else { // horizontal
			if tracking[idx].w < tracking[idx].aw {
				delta := tracking[idx].aw - tracking[idx].w
				if ca, ok := child.widget.(Alignable); ok {
					xAlign, _ := ca.GetAlignment()
					pad := int(float64(delta) * xAlign)
					tracking[idx].x += pad
					tracking[idx].overflow += delta - pad
				} else {
					tracking[idx].overflow += delta
				}
			}
		} // if isVertical
	} // for each child

	// allocate children and update canvas

	for idx, child := range children {
		local := cdk.NewPoint2I(tracking[idx].x, tracking[idx].y)
		childSize := cdk.NewRectangle(tracking[idx].w, tracking[idx].h)
		childAlloc := cdk.NewRectangle(tracking[idx].aw, tracking[idx].ah)
		nextPoint.Add(local.X, local.Y)
		if canvas, ok := b.canvases[child.widget.ObjectID()]; ok {
			x := nextPoint.X - origin.X
			y := nextPoint.Y - origin.Y
			canvas.SetOrigin(cdk.MakePoint2I(x, y))
			canvas.Resize(cdk.MakeRectangle(childAlloc.W, childAlloc.H), b.GetTheme().Content.Normal)
		} else {
			b.LogError("missing canvas for child: %v", child)
		}
		child.widget.SetOrigin(nextPoint.X, nextPoint.Y)
		child.widget.SetAllocation(*childSize)
		child.widget.Resize()
		if isVertical {
			nextPoint.Y += tracking[idx].h + tracking[idx].overflow
		} else {
			nextPoint.X += tracking[idx].w + tracking[idx].overflow
		}
		if len(gaps) > idx {
			if isVertical {
				nextPoint.Y += gaps[idx]
			} else {
				nextPoint.X += gaps[idx]
			}
		}
	}

	b.Emit(SignalResize, b)
	return cdk.EVENT_STOP
}

func (b *CBox) resizeDynamicAlloc(isVertical bool, gaps []int, increment, spacing, numChildren int, origin, nextPoint *cdk.Point2I, alloc *cdk.Rectangle, children []*cBoxChild) cdk.EventFlag {
	var (
		extraSpace   int
		numExpanding int
		consumed     int
	)
	tracking := make([]struct {
		x, y, w, h int
		rw, rh     int
		aw, ah     int
		extra      int
		overflow   int
	}, numChildren)

	// first: build up tracking dataset

	for idx, child := range children {
		req := cdk.NewRectangle(child.widget.GetSizeRequest())
		if child.expand { // expand
			numExpanding += 1
			if child.fill { // expand && fill
				if isVertical {
					tracking[idx].w = alloc.W
					tracking[idx].h = increment
				} else {
					tracking[idx].w = increment
					tracking[idx].h = alloc.H
				}
				tracking[idx].aw = tracking[idx].w
				tracking[idx].ah = tracking[idx].h
			} else { // expand && !fill
				if isVertical { // expand && !fill && vertical
					if req.H <= -1 || req.H > increment {
						req.H = increment
					}
					tracking[idx].w = alloc.W
					tracking[idx].h = req.H
					tracking[idx].aw = alloc.W
					tracking[idx].ah = increment
				} else { // expand && !fill && horizontal
					if req.W <= -1 || req.W > increment {
						req.W = increment
					}
					tracking[idx].w = req.W
					tracking[idx].h = alloc.H
					tracking[idx].aw = increment
					tracking[idx].ah = alloc.H
				}
			} // else expand, !fill
		} else { // if !expand (assume !fill)
			if isVertical { // !expand, !fill, vertical
				req.W = alloc.W
				if req.H <= -1 || req.H > increment {
					req.H = increment
				}
				if req.H < increment {
					delta := increment - req.H
					extraSpace += delta
				}
			} else { // !expand, !fill, horizontal
				req.H = alloc.H
				if req.W <= -1 || req.W > increment {
					req.W = increment
				}
				if req.W < increment {
					delta := increment - req.W
					extraSpace += delta
				}
			}
			tracking[idx].w = req.W
			tracking[idx].h = req.H
			tracking[idx].aw = req.W
			tracking[idx].ah = req.H
		} // else expand
		tracking[idx].rw = req.W
		tracking[idx].rh = req.H
	}

	for idx, _ := range children {
		if isVertical {
			consumed += tracking[idx].ah
		} else {
			consumed += tracking[idx].aw
		}
	}
	var total = 0
	if isVertical {
		total = alloc.H - consumed
	} else {
		total = alloc.W - consumed
	}
	if extraSpace > 0 {
		var (
			dist []int
			err  error
		)
		if dist, gaps, err = utils.Distribute(total, extraSpace, numExpanding, numChildren, spacing); err != nil {
			b.LogErr(err)
		} else {
			di, nDist := 0, len(dist)
			for idx, child := range children {
				if child.expand {
					if di < nDist {
						if child.fill {
							if isVertical {
								tracking[idx].ah += dist[di]
								tracking[idx].h = tracking[idx].ah
							} else {
								tracking[idx].aw += dist[di]
								tracking[idx].w = tracking[idx].aw
							}
						} else { // !fill
							if isVertical {
								tracking[idx].ah += dist[di]
							} else {
								tracking[idx].aw += dist[di]
							}
						}
						di += 1
					}
				}
			}
		}
	}

	// solve positions

	for idx, child := range children {
		if ca, ok := child.widget.(Alignable); ok {
			xAlign, yAlign := ca.GetAlignment()
			if isVertical {
				if tracking[idx].h < tracking[idx].ah {
					delta := tracking[idx].ah - tracking[idx].h
					inc := int(float64(delta) * yAlign)
					tracking[idx].y += inc
					tracking[idx].overflow += delta - inc
				}
			} else {
				if tracking[idx].w < tracking[idx].aw {
					delta := tracking[idx].aw - tracking[idx].w
					inc := int(float64(delta) * xAlign)
					tracking[idx].x += inc
					tracking[idx].overflow += delta - inc
				}
			}
		}
	}

	// allocate space and update canvas

	for idx, child := range children {
		track := tracking[idx]
		nextPoint.X += track.x
		nextPoint.Y += track.y
		if canvas, ok := b.canvases[child.widget.ObjectID()]; ok {
			canvas.SetOrigin(cdk.MakePoint2I(nextPoint.X-origin.X, nextPoint.Y-origin.Y))
			canvas.Resize(cdk.MakeRectangle(track.w, track.h), b.GetTheme().Content.Normal)
		} else {
			b.LogError("missing canvas for child: %v", child)
		}
		child.widget.SetOrigin(nextPoint.X, nextPoint.Y)
		child.widget.SetAllocation(cdk.MakeRectangle(track.w, track.h))
		child.widget.Resize()
		if isVertical {
			nextPoint.Y += track.h + track.overflow
		} else {
			nextPoint.X += track.w + track.overflow
		}
		if len(gaps) > idx {
			if isVertical {
				nextPoint.Y += gaps[idx]
			} else {
				nextPoint.X += gaps[idx]
			}
		}
	}

	b.Emit(SignalResize, b)
	return cdk.EVENT_STOP
}

func (b *CBox) Draw(canvas cdk.Canvas) cdk.EventFlag {
	debug, _ := b.GetBoolProperty(cdk.PropertyDebug)
	debugChildren, _ := b.GetBoolProperty(PropertyDebugChildren)
	orientation := b.GetOrientation()
	children := b.getBoxChildren()
	canvas.Fill(b.GetTheme())
	for _, child := range children {
		if child.widget.IsVisible() {
			v, ok := b.canvases[child.widget.ObjectID()]
			if !ok {
				b.LogError("missing child canvas: %v", child.widget.ObjectID())
				continue
			}
			child.widget.Draw(v)
			if debugChildren && orientation == cdk.ORIENTATION_VERTICAL {
				v.DebugBox(cdk.ColorPink, child.widget.ObjectInfo()+" ["+b.ObjectInfo()+"]")
			} else if debugChildren {
				v.DebugBox(cdk.ColorPurple, child.widget.ObjectInfo()+" ["+b.ObjectInfo()+"]")
			}
			if err := canvas.Composite(v); err != nil {
				b.LogError("composite error: %v", err)
			}
		}
	}
	if debug && orientation == cdk.ORIENTATION_VERTICAL {
		canvas.DebugBox(cdk.ColorPink, b.ObjectInfo())
	} else if debug {
		canvas.DebugBox(cdk.ColorPurple, b.ObjectInfo())
	}
	return cdk.EVENT_STOP
}

// Whether the children should all be the same size.
// Flags: Read / Write
// Default value: FALSE
const PropertyHomogeneous cdk.Property = "homogeneous"

// The amount of space between children.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 0
const PropertySpacing cdk.Property = "spacing"

const PropertyDebugChildren cdk.Property = "debug-children"

const PropertyBoxChildPackType cdk.Property = "box-child--pack-type"

const PropertyBoxChildExpand cdk.Property = "box-child--expand"

const PropertyBoxChildFill cdk.Property = "box-child--fill"

const PropertyBoxChildPadding cdk.Property = "box-child--padding"
