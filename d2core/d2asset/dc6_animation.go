package d2asset

import (
	"errors"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	d2iface "github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

// DC6Animation is an animation made from a DC6 file
type DC6Animation struct {
	animation
	dc6Path  string
	palette  d2iface.Palette
	renderer d2iface.Renderer
}

// CreateDC6Animation creates an Animation from d2dc6.DC6 and d2dat.DATPalette
func CreateDC6Animation(renderer d2iface.Renderer, dc6Path string,
	palette d2iface.Palette) (d2iface.Animation, error) {
	dc6, err := loadDC6(dc6Path)
	if err != nil {
		return nil, err
	}

	animation := animation{
		directions:     make([]animationDirection, dc6.Directions),
		playLength:     defaultPlayLength,
		playLoop:       true,
		originAtBottom: true,
	}

	anim := DC6Animation{
		animation: animation,
		dc6Path:   dc6Path,
		palette:   palette,
		renderer:  renderer,
	}

	err = anim.SetDirection(0)

	return &anim, err
}

// SetDirection decodes and sets the direction
func (a *DC6Animation) SetDirection(directionIndex int) error {
	const smallestInvalidDirectionIndex = 64
	if directionIndex >= smallestInvalidDirectionIndex {
		return errors.New("invalid direction index")
	}

	direction := d2dcc.Dir64ToDcc(directionIndex, len(a.directions))
	if !a.directions[direction].decoded {
		err := a.decodeDirection(direction)
		if err != nil {
			return err
		}
	}

	a.directionIndex = direction
	a.frameIndex = 0

	return nil
}

func (a *DC6Animation) decodeDirection(directionIndex int) error {
	dc6, err := loadDC6(a.dc6Path)
	if err != nil {
		return err
	}

	startFrame := directionIndex * int(dc6.FramesPerDirection)

	for i := 0; i < int(dc6.FramesPerDirection); i++ {
		dc6Frame := dc6.Frames[startFrame+i]

		sfc, err := a.renderer.NewSurface(int(dc6Frame.Width), int(dc6Frame.Height),
			d2enum.FilterNearest)
		if err != nil {
			return err
		}

		indexData := make([]int, dc6Frame.Width*dc6Frame.Height)
		for i := range indexData {
			indexData[i] = -1
		}

		x := 0
		y := int(dc6Frame.Height) - 1
		offset := 0

		for {
			b := int(dc6Frame.FrameData[offset])
			offset++

			if b == 0x80 {
				if y == 0 {
					break
				}
				y--
				x = 0
			} else if b&0x80 > 0 {
				transparentPixels := b & 0x7f
				for i := 0; i < transparentPixels; i++ {
					indexData[x+y*int(dc6Frame.Width)+i] = -1
				}
				x += transparentPixels
			} else {
				for i := 0; i < b; i++ {
					indexData[x+y*int(dc6Frame.Width)+i] = int(dc6Frame.FrameData[offset])
					offset++
				}
				x += b
			}
		}

		bytesPerPixel := 4
		colorData := make([]byte, int(dc6Frame.Width)*int(dc6Frame.Height)*bytesPerPixel)

		for i := 0; i < int(dc6Frame.Width*dc6Frame.Height); i++ {
			// TODO: Is this == -1 or < 1?
			if indexData[i] < 1 {
				continue
			}

			c := a.palette.GetColors()[indexData[i]]
			colorData[i*bytesPerPixel] = c.R()
			colorData[i*bytesPerPixel+1] = c.G()
			colorData[i*bytesPerPixel+2] = c.B()
			colorData[i*bytesPerPixel+3] = c.A()
		}

		if err := sfc.ReplacePixels(colorData); err != nil {
			return err
		}

		a.directions[directionIndex].decoded = true
		a.directions[directionIndex].frames = append(a.directions[directionIndex].frames, &animationFrame{
			width:   int(dc6Frame.Width),
			height:  int(dc6Frame.Height),
			offsetX: int(dc6Frame.OffsetX),
			offsetY: int(dc6Frame.OffsetY),
			image:   sfc,
		})
	}

	return nil
}

// Clone creates a copy of the animation
func (a *DC6Animation) Clone() d2iface.Animation {
	animation := *a
	return &animation
}
