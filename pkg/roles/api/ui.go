package api

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/api/types"
	"beryju.io/ddet/web/dist"
	"github.com/go-chi/chi/v5"
)

//go:embed ui/index.html
var IndexTemplate string

func (r *APIRole) setupUI() {
	t, err := template.New("ddet.ui").Parse(IndexTemplate)
	if err != nil {
		panic(err)
	}
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		if !extconfig.Get().Debug {
			return
		}
		r.m.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/ui/", http.StatusFound)
		})
		r.m.Get("/ui/", func(w http.ResponseWriter, r *http.Request) {
			err := t.Execute(w, nil)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		})
		r.m.Route("/ui/static/", func(r chi.Router) {
			var handler http.Handler
			if extconfig.Get().Debug {
				handler = http.StripPrefix("/ui/static", http.FileServer(http.Dir("./web/dist")))
			} else {
				handler = http.StripPrefix("/ui/static", http.FileServer(http.FS(dist.Static)))
			}
			r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("foo")
				handler.ServeHTTP(w, r)
			})
		})
	})
}
