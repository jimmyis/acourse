package app

import (
	"database/sql"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/acourse/pkg/appctx"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/flash"
	"github.com/acoshift/header"
	"github.com/acoshift/httprouter"
	"github.com/lib/pq"
)

func getCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	link := httprouter.GetParam(ctx, "courseID")

	// if id can parse to int64 get course from id
	id, err := strconv.ParseInt(link, 10, 64)
	if err != nil {
		// link can not parse to int64 get course id from url
		id, err = model.GetCourseIDFromURL(link)
		if err == model.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	x, err := model.GetCourse(id)
	if err == model.ErrNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		http.Redirect(w, r, "/course/"+x.URL.String, http.StatusFound)
		return
	}

	enrolled := false
	if user != nil {
		enrolled, err = model.IsEnrolled(user.ID, x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	owned := user.ID == x.UserID

	// if user enrolled or user is owner fetch course contents
	if enrolled || owned {
		x.Contents, err = model.GetCourseContents(x.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if owned {
		x.Owner = user
	} else {
		x.Owner, err = model.GetUser(x.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	page := defaultPage
	page.Title = x.Title + " | " + page.Title
	page.Desc = x.ShortDesc
	page.Image = x.Image
	page.URL = baseURL + "/course/" + url.PathEscape(x.Link())
	view.Course(w, r, &view.CourseData{
		Page:     &page,
		Course:   x,
		Enrolled: enrolled,
		Owned:    owned,
	})
}

func getCourseCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := flash.Get(ctx)

	page := defaultPage
	page.Title = "Create new Course | " + page.Title
	view.CourseCreate(w, r, &view.CourseCreateData{
		Page:  &page,
		Flash: f,
	})
}

func postCourseCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := flash.Get(ctx)
	user := appctx.GetUser(ctx)

	if !verifyXSRF(r.FormValue("X"), user.ID, "editor/create") {
		f.Add("Errors", "invalid xsrf token")
		back(w, r)
		return
	}

	var (
		title         = r.FormValue("Title")
		shortDesc     = r.FormValue("ShortDesc")
		desc          = r.FormValue("Desc")
		imageURL      string
		start         pq.NullTime
		assignment, _ = strconv.ParseBool(r.FormValue("Assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		back(w, r)
		return
	}

	if v := r.FormValue("Start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	image, info, err := r.FormFile("Image")
	if err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			back(w, r)
			return
		}

		imageURL, err = UploadProfileImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRow(`
		insert into courses
			(user_id, title, short_desc, long_desc, image, start)
		values
			($1, $2, $3, $4, $5, $6)
		returning id
	`, user.ID, title, shortDesc, desc, imageURL, start).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec(`
		insert into course_options
			(course_id, assignment)
		values
			($1, $2)
	`, id, assignment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var link sql.NullString
	db.QueryRow(`select url from courses where id = $1`, id).Scan(&link)
	if !link.Valid {
		http.Redirect(w, r, "/course/"+strconv.FormatInt(id, 10), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/course/"+link.String, http.StatusFound)
}

func getCourseEdit(w http.ResponseWriter, r *http.Request) {
	page := defaultPage
	page.Title = "Edit Course | " + page.Title

	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	course, err := model.GetCourse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.CourseEdit(w, r, &view.CourseEditData{
		Page:   &page,
		Course: course,
	})
}

func postCourseEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

	f := flash.Get(ctx)
	user := appctx.GetUser(ctx)

	if !verifyXSRF(r.FormValue("X"), user.ID, "editor/course") {
		f.Add("Errors", "invalid xsrf token")
		back(w, r)
		return
	}

	var (
		title         = r.FormValue("Title")
		shortDesc     = r.FormValue("ShortDesc")
		desc          = r.FormValue("Desc")
		imageURL      string
		start         pq.NullTime
		assignment, _ = strconv.ParseBool(r.FormValue("Assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		back(w, r)
		return
	}

	if v := r.FormValue("Start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	image, info, err := r.FormFile("Image")
	if err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			back(w, r)
			return
		}

		imageURL, err = UploadProfileImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		update courses
		set
			title = $2,
			short_desc = $3,
			long_desc = $4,
			start = $5,
			updated_at = now()
		where id = $1
	`, id, title, shortDesc, desc, start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(imageURL) > 0 {
		_, err = tx.Exec(`
			update courses
			set
				image = $2
			where id = $1
		`, id, imageURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	_, err = tx.Exec(`
		upsert into course_options
			(course_id, assignment)
		values
			($1, $2)
	`, id, assignment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var link sql.NullString
	db.QueryRow(`select url from courses where id = $1`, id).Scan(&link)
	if !link.Valid {
		http.Redirect(w, r, "/course/"+strconv.FormatInt(id, 10), http.StatusFound)
		return
	}
	http.Redirect(w, r, "/course/"+link.String, http.StatusFound)
}

func getCourseContentEdit(w http.ResponseWriter, r *http.Request) {
	page := defaultPage
	page.Title = "Edit Course | " + page.Title

	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	course, err := model.GetCourse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.CourseContentEdit(w, r, &view.CourseEditData{
		Page:   &page,
		Course: course,
	})
}