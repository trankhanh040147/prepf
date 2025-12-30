package diffview

import (
	"fmt"
	"image/color"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/prepf/internal/ansiext"
)

var _ chroma.Formatter = chromaFormatter{}

// chromaFormatter is a custom formatter for Chroma that uses Lip Gloss for
// foreground styling, while keeping a forced background color.
type chromaFormatter struct {
	bgColor color.Color
}

// colorToHex converts a color.Color to hex string
func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

// Format implements the chroma.Formatter interface.
func (c chromaFormatter) Format(w io.Writer, style *chroma.Style, it chroma.Iterator) error {
	for token := it(); token != chroma.EOF; token = it() {
		value := strings.TrimRight(token.Value, "\n")
		value = ansiext.Escape(value)

		entry := style.Get(token.Type)
		if entry.IsZero() {
			if _, err := fmt.Fprint(w, value); err != nil {
				return err
			}
			continue
		}

		s := lipgloss.NewStyle().
			Background(lipgloss.Color(colorToHex(c.bgColor)))

		if entry.Bold == chroma.Yes {
			s = s.Bold(true)
		}
		if entry.Underline == chroma.Yes {
			s = s.Underline(true)
		}
		if entry.Italic == chroma.Yes {
			s = s.Italic(true)
		}
		if entry.Colour.IsSet() {
			s = s.Foreground(lipgloss.Color(entry.Colour.String()))
		}

		if _, err := fmt.Fprint(w, s.Render(value)); err != nil {
			return err
		}
	}
	return nil
}
