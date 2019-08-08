package beeplus

type M map[string]interface{}

var (
	G         M
	LibLoader LibLoaderClass
)

func init() {
	G = make(M)
	LibLoader = LibLoaderClass{}
}
