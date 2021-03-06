package grafio

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Sdl2Draw adapter for sdl2 go bindings
type Sdl2Draw struct {
	r *sdl.Renderer

	fonts    map[string]*ttf.Font
	textures map[string]*sdl.Texture

	mainFont string

	w, h int32
}

// NewSdl2Draw creates new sdl2draw drawer
func NewSdl2Draw(font string, w, h int32) (sdl2draw *Sdl2Draw, destroy func() error, erro error) {
	r, destroy, err := prepareSdl2(w, h)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not prepare sdl2 engine")
	}

	return &Sdl2Draw{
		mainFont: font,

		fonts:    map[string]*ttf.Font{},
		textures: map[string]*sdl.Texture{},

		r: r,
		w: w,
		h: h,
	}, destroy, nil
}

// SetMainFont sets the default font for all text
func (s *Sdl2Draw) SetMainFont(fontFileName string) error {
	if _, ok := s.fonts[fontFileName]; !ok {
		return fmt.Errorf("font %s is not loaded", fontFileName)
	}

	s.mainFont = fontFileName

	return nil
}

// ColorRect draws into sdl2 window given color rectangle in the given coordinates
func (s *Sdl2Draw) ColorRect(x, y, w, h int32, rgba RGBA) error {
	if err := s.r.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A); err != nil {
		return errors.Wrap(err, "could not set draw color")
	}

	if err := s.r.FillRect(&sdl.Rect{X: x, Y: y, W: w, H: h}); err != nil {
		return errors.Wrap(err, "could not fill rect")
	}

	return nil
}

// TextureRect draws into sdl2 window with given texture rectangle in the given coordinates
func (s *Sdl2Draw) TextureRect(x, y, w, h int32, texture string) error {
	if _, ok := s.textures[texture]; !ok {
		return fmt.Errorf("texture %s is not found", texture)
	}

	rect := &sdl.Rect{X: x, Y: y, W: w, H: h}
	if err := s.r.Copy(s.textures[texture], nil, rect); err != nil {
		return errors.Wrap(err, "could not copy texture")
	}

	return nil
}

// ScreenHeight returns the height in pixels of the sdl2 created screen
func (s Sdl2Draw) ScreenHeight() int32 {
	return s.h
}

// ScreenWidth returns the width in pixels of the sdl2 created screen
func (s Sdl2Draw) ScreenWidth() int32 {
	return s.w
}

// Background draws whole background of sdl2 window to the given rgba color
func (s *Sdl2Draw) Background(rgba RGBA) error {
	if err := s.r.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A); err != nil {
		return errors.Wrap(err, "couldn't set color")
	}

	if err := s.r.FillRect(nil); err != nil {
		return errors.Wrap(err, "couldn't fill rect")
	}

	return nil
}

// Text draws given text with options to the sdl2 window
func (s *Sdl2Draw) Text(txt string, opts TextOpts) (erro error) {
	c := sdl.Color{R: opts.Color.R, G: opts.Color.G, B: opts.Color.B, A: opts.Color.A}
	surface, err := s.fonts[s.mainFont].RenderUTF8Solid(txt, c)
	if err != nil {
		return errors.Wrap(err, "could not render title")
	}
	defer surface.Free()

	texture, err := s.r.CreateTextureFromSurface(surface)
	if err != nil {
		return errors.Wrap(err, "could not create texture")
	}

	defer func() {
		if err = texture.Destroy(); err != nil {
			erro = errors.Wrap(err, "could not destroy texture")
		}
	}()

	shift := 0
	if opts.Align == Right && len(txt) > 1 {
		shift = (len(txt) * int(opts.Size)) - int(opts.Size)
	}

	rect := &sdl.Rect{
		X: sizeCal(s.w, opts.XCof) - int32(shift),
		Y: sizeCal(s.h, opts.YCof),
		W: opts.Size * int32(len(txt)),
		H: opts.Size + 20,
	}

	if err := s.r.Copy(texture, nil, rect); err != nil {
		return errors.Wrap(err, "could not copy texture")
	}

	return nil
}

// Present makes white background of the window calls user's function
func (s *Sdl2Draw) Present(f func() error) error {
	if err := s.r.Clear(); err != nil {
		return errors.Wrap(err, "could not clear the renderer")
	}

	if err := s.Background(ColorWhite); err != nil {
		return errors.Wrap(err, "could not set the background")
	}

	if err := f(); err != nil {
		return errors.Wrap(err, "could not execute user given function")
	}

	s.r.Present()

	return nil
}

// LoadResources load resources of fonts and textures given path
func (s *Sdl2Draw) LoadResources(fontsPath, texturesPath string) (func() error, error) {
	textures, err := ioutil.ReadDir(texturesPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read dir")
	}

	for _, f := range textures {
		if f.IsDir() {
			continue
		}

		image, err := img.Load(texturesPath + "/" + f.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to create texture: %v\n", err)
		}

		texture, err := s.r.CreateTextureFromSurface(image)
		if err != nil {
			return nil, fmt.Errorf("failed to create texture: %v\n", err)
		}
		image.Free()

		s.textures[f.Name()] = texture
	}

	fonts, err := ioutil.ReadDir(fontsPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read dir")
	}

	for _, f := range fonts {
		font, err := ttf.OpenFont(fontsPath+"/"+f.Name(), 124)
		if err != nil {
			return nil, fmt.Errorf("could not load font: %v", err)
		}

		s.fonts[f.Name()] = font
	}

	return func() error { return s.destroy() }, nil
}

func (s *Sdl2Draw) destroy() error {
	for _, texture := range s.textures {
		if err := texture.Destroy(); err != nil {
			return errors.Wrap(err, "could not destroy texture")
		}
	}

	for _, font := range s.fonts {
		font.Close()
	}
	return nil
}

func sizeCal(size int32, cof float32) int32 {
	return int32(float32(size) * (float32(cof)))
}
