package clipboard

// clipboard writer, reader program names
const (
	WriterOnMac = "pbcopy"
	WriterOnWin = "clip"
	WriterOnLin = "xsel"

	ReaderOnMac = "pbpaste"
	ReaderOnWin = "clip"
	ReaderOnLin = "xclip"
)

// std instance
var std = New()

// Reset clipboard data
func Reset() error {
	return std.Reset()
}

// ReadString contents from clipboard
func ReadString() (string, error) {
	return std.ReadString()
}

// WriteString contents to clipboard and flush
func WriteString(s string) error {
	_, err := std.WriteString(s)
	if err != nil {
		return err
	}

	return std.Flush()
}