package grafio

type RGBA struct {
	R, G, B, A uint8
}

type TextOpts struct {
	Size       int32
	XCof, YCof float32
	Color      RGBA
}

type Drawer interface {
	Text(txt string, opts TextOpts) error
	Background(r, g, b, a uint8) error
	Rectangle(x, y, w, h int32, rgba RGBA) error

	Present(f func() error) error
	LoadResources(fontsPath, texturesPath string) (func() error, error)
	//SetMainFont(fontFileName string) error

	ScreenHeight() int32
	ScreenWidth() int32
}

func sizeCal(size int32, cof float32) int32 {
	return int32(float32(size) * (float32(cof)))
}
