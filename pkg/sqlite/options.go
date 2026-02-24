package sqlite

type Option func(*Sqlite)

func Driver(name string) Option {
	return func(s *Sqlite) {
		s.driverName = name
	}
}

func SourceName(name string) Option {
	return func(s *Sqlite) {
		s.dataSourceName = name
	}
}
