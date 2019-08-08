package beeplus

var (
	Options map[string]string
)

func init() {
	Options = map[string]string{
		"autoIndex": "true", //在router中是否启用自动index
	}
}
func SetOptions(key, value string) {
	Options[key] = value
}
func GetOptions(key string) string {
	return Options[key]
}
