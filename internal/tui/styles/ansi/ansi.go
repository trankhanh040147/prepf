// Package ansi provides ANSI style configuration types.
// This is a minimal stub for compatibility with glamour v0.10.0.
// The actual ansi package exists in glamour v2, but prepf uses v0.10.0.
package ansi

// StylePrimitive represents basic style properties
type StylePrimitive struct {
	BlockPrefix     *string
	BlockSuffix     *string
	Prefix          *string
	Suffix          *string
	Color           *string
	BackgroundColor *string
	Bold            *bool
	Italic          *bool
	Underline       *bool
	CrossedOut      *bool
	Format          string
}

// StyleBlock represents a block-level style
type StyleBlock struct {
	StylePrimitive StylePrimitive
	Indent         *uint
	IndentToken    *string
	Margin         *uint
}

// StyleList represents list styling
type StyleList struct {
	LevelIndent int
}

// StyleTask represents task list styling
type StyleTask struct {
	StylePrimitive StylePrimitive
	Ticked         string
	Unticked       string
}

// StyleTable represents table styling
type StyleTable struct {
	StyleBlock StyleBlock
}

// StyleCodeBlock represents code block styling
type StyleCodeBlock struct {
	StyleBlock StyleBlock
	Chroma     *Chroma
}

// Chroma represents syntax highlighting styles
type Chroma struct {
	Text                StylePrimitive
	Error               StylePrimitive
	Comment             StylePrimitive
	CommentPreproc      StylePrimitive
	Keyword             StylePrimitive
	KeywordReserved     StylePrimitive
	KeywordNamespace    StylePrimitive
	KeywordType         StylePrimitive
	Operator            StylePrimitive
	Punctuation         StylePrimitive
	Name                StylePrimitive
	NameBuiltin         StylePrimitive
	NameTag             StylePrimitive
	NameAttribute       StylePrimitive
	NameClass           StylePrimitive
	NameDecorator       StylePrimitive
	NameFunction        StylePrimitive
	LiteralNumber       StylePrimitive
	LiteralString       StylePrimitive
	LiteralStringEscape StylePrimitive
	GenericDeleted      StylePrimitive
	GenericEmph         StylePrimitive
	GenericInserted     StylePrimitive
	GenericStrong       StylePrimitive
	GenericSubheading   StylePrimitive
	Background          StylePrimitive
}

// StyleConfig represents the complete markdown style configuration
type StyleConfig struct {
	Document            StyleBlock
	BlockQuote          StyleBlock
	List                StyleList
	Heading             StyleBlock
	H1                  StyleBlock
	H2                  StyleBlock
	H3                  StyleBlock
	H4                  StyleBlock
	H5                  StyleBlock
	H6                  StyleBlock
	Strikethrough       StylePrimitive
	Emph                StylePrimitive
	Strong              StylePrimitive
	HorizontalRule      StylePrimitive
	Item                StylePrimitive
	Enumeration         StylePrimitive
	Task                StyleTask
	Link                StylePrimitive
	LinkText            StylePrimitive
	Image               StylePrimitive
	ImageText           StylePrimitive
	Code                StyleBlock
	CodeBlock           StyleCodeBlock
	Table               StyleTable
	DefinitionDescription StylePrimitive
}

