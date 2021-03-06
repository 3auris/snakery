package grafio

var (
	// ColorGreen rgba green color
	ColorGreen = RGBA{R: 34, G: 139, B: 34, A: 10}

	// ColorWhite rgba white color
	ColorWhite = RGBA{R: 255, G: 255, B: 255, A: 255}

	// ColorBlack rgba black color
	ColorBlack = RGBA{R: 0, G: 0, B: 0, A: 0}
)

// RGBA have rgba color values
type RGBA struct {
	R, G, B, A uint8
}

// TextAlign type of text alignment side
type TextAlign int

// Right align text to stick in the right
const Right TextAlign = 2

// TextOpts options of text
type TextOpts struct {
	Size       int32
	XCof, YCof float32
	Color      RGBA
	Align      TextAlign
}

// RectOpts options of rectangle
type RectOpts struct {
	Texture string
	Color   RGBA
}

// Drawer an engine who can draw on window
type Drawer interface {
	// Background draws the whole background to the given RGBA color
	Background(rgba RGBA) error

	// Text writes given text with given options to the window
	Text(txt string, opts TextOpts) error

	// ColorRect draw rectangle with the given color
	ColorRect(x, y, w, h int32, rgba RGBA) error

	// TextureRect draw rectangle with the given texture file name
	TextureRect(x, y, w, h int32, texture string) error

	// Present draws everything into the window
	Present(f func() error) error

	// LoadResources loads fonts and textures of the given path
	LoadResources(fontsPath, texturesPath string) (func() error, error)

	// SetMainFont Set default font, if someone will use Text function the Main font will be used
	SetMainFont(fontFileName string) error

	// ScreenHeight returns the height of screen in pixels
	ScreenHeight() int32

	// ScreenWidth returns the width of screen in pixels
	ScreenWidth() int32
}
