package view

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/flash"
	"github.com/acoshift/header"
	"github.com/dustin/go-humanize"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"golang.org/x/net/xsrftoken"
)

const templateDir = "template"

var (
	m         = minify.New()
	muExecute = &sync.Mutex{}
	templates = make(map[interface{}]*templateStruct)
	loc       *time.Location
)

type templateStruct struct {
	*template.Template
	set []string
}

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal(err)
	}

	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)

	parseTemplate(keyIndex, []string{"index.tmpl", "app.tmpl", "layout.tmpl", "component/course-card.tmpl"})
	parseTemplate(keySignIn, []string{"signin.tmpl", "auth.tmpl", "layout.tmpl"})
	parseTemplate(keySignUp, []string{"signup.tmpl", "auth.tmpl", "layout.tmpl"})
	parseTemplate(keyProfile, []string{
		"profile.tmpl", "app.tmpl", "layout.tmpl",
		"component/user-profile.tmpl",
		"component/own-course-card.tmpl",
		"component/enrolled-course-card.tmpl",
	})
	parseTemplate(keyCourse, []string{"course.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyAdminUsers, []string{"admin/users.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyAdminCourses, []string{"admin/courses.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyAdminPayments, []string{"admin/payments.tmpl", "app.tmpl", "layout.tmpl"})
}

func joinTemplateDir(files []string) []string {
	r := make([]string, len(files))
	for i, f := range files {
		r[i] = filepath.Join(templateDir, f)
	}
	return r
}

func parseTemplate(key interface{}, set []string) {
	templateName := strings.TrimSuffix(set[0], ".tmpl")
	t := template.New("")
	t.Funcs(template.FuncMap{
		"templateName": func() string {
			return templateName
		},
		"xsrf": func(action string) string {
			return xsrftoken.Generate(internal.GetXSRFSecret(), "", action)
		},
		"currency": func(v float64) string {
			return humanize.FormatFloat("#,###.##", v)
		},
		"me": func() interface{} {
			return nil
		},
		"courseType": func(v int) string {
			switch v {
			case model.Live:
				return "Live"
			case model.Video:
				return "Video"
			case model.EBook:
				return "eBook"
			default:
				return ""
			}
		},
		"date": func(v time.Time) string {
			return v.In(loc).Format("02/01/2006")
		},
		"dateTime": func(v time.Time) string {
			return v.In(loc).Format("02/01/2006 15:04:05")
		},
		"live": func() int {
			return model.Live
		},
		"video": func() int {
			return model.Video
		},
		"eBook": func() int {
			return model.EBook
		},
	})
	_, err := t.ParseFiles(joinTemplateDir(set)...)
	if err != nil {
		log.Fatalf("internal: parse template %s error; %v", templateName, err)
	}
	t = t.Lookup("root")
	if t == nil {
		log.Fatalf("internal: root template not found in %s", templateName)
	}
	templates[key] = &templateStruct{
		Template: t,
		set:      set,
	}
}

func render(w http.ResponseWriter, r *http.Request, key, data interface{}) {
	t := templates[key]
	if t == nil {
		http.Error(w, fmt.Sprintf("template not found"), http.StatusInternalServerError)
		return
	}
	if dev {
		muExecute.Lock()
		defer muExecute.Unlock()
		parseTemplate(key, t.set)
		t = templates[key]
	}

	ctx := r.Context()

	// inject template funcs
	tp := t.Funcs(template.FuncMap{
		"me": func() interface{} {
			return internal.GetUser(ctx)
		},
	})

	w.Header().Set(header.ContentType, "text/html; charset=utf-8")
	w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate, max-age=0")
	pipe := &bytes.Buffer{}
	err := tp.Execute(pipe, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = m.Minify("text/html", w, pipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	f := flash.Get(r.Context())
	f.Clear()
}