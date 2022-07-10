package api

import (
	_ "embed"
	"html/template"
	"net/http"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/api/types"
	"beryju.io/ddet/web/dist"
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
		r.m.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/ui/", http.StatusFound)
		})
		r.m.Path("/ui/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := t.Execute(w, nil)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		})
		var handler http.Handler
		if extconfig.Get().Debug {
			handler = http.StripPrefix("/ui/static", http.FileServer(http.Dir("./web/dist")))
		} else {
			handler = http.StripPrefix("/ui/static", http.FileServer(http.FS(dist.Static)))
		}
		r.m.PathPrefix("/ui/static/").Handler(handler)
	})
}
