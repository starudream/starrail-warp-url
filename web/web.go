package web

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/core/v2/utils/osutil"
)

const Index = "index.html"

//go:embed *.html *.js *.css
var FS embed.FS

func Handler(data string) http.Handler {
	fs := http.FileServer(http.FS(FS))
	tpl := osutil.Must1(template.New("root").ParseFS(FS, Index))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/"+Index {
			osutil.PanicErr(tpl.ExecuteTemplate(w, Index, gh.M{"data": data}))
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
