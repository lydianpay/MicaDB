package db

// Options holds the settings used when creating a new MicaDB
type Options struct {
	// Filename is what the file on disk will be named. The default is micadb.bin
	Filename string

	// IsTest when enabled, the file will be created in the temp directory and removed afterwards
	IsTest bool

	// CustomTypes is used to pass any custom types that may be stored in MicaDB
	CustomTypes []any

	// BackupFrequency is how often in seconds the in-memory data store should be persisted to disk
	// 0 (Default) will occur on every change, -1 will disable backups
	BackupFrequency int
}

func (options Options) init() {
	if options.Filename == "" {
		options.Filename = "micadb.bin"
	}

}
