package queries

import "embed"

//go:embed *.sql
var fs embed.FS

func Get(name string) string {
	data, err := fs.ReadFile(name + ".sql")
	if err != nil {
		panic(err)
	}
	return string(data)
}
