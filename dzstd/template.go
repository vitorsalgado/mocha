package dzstd

import "io"

// TemplateEngine initializes templates and creates TemplateRenderer instances.
type TemplateEngine interface {
	// Load is executed on server initialization.
	Load() error

	// Parse pre-compiles the source template and returns a renderer for it.
	// It is usually run once during mock setup.
	Parse(string) (TemplateRenderer, error)
}

// TemplateRenderer defines a renderer for templates.
// Each template will have a renderer associated to it.
type TemplateRenderer interface {
	// Render renders the previously parsed template to the given io.Writer.
	// The second parameter is template data and can be nil.
	Render(io.Writer, any) error
}
