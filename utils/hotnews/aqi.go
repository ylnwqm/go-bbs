package hotnews

var QueryName = "query.hotnews"

type News interface {
	Get(url string) map[string]interface{}
}

