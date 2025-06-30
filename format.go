package log

type Formatter interface {
	Format(entry *Entry) error
}
