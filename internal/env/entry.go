package env

// Entry represents a single key-value pair parsed from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string // inline or preceding comment, if any
}
