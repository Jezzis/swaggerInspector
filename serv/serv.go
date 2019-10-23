package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"swaggerdoc/lib"
	"time"
)

type global struct {
	Url    string
	Prefix string
	Token  string
	R      int64
}

func joinComment(source, newLine string) string {
	if source != "" {
		return source + "\n" + newLine
	} else {
		return newLine
	}
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Printf("currDir: %s\n", currDir)

		wwwDir := filepath.Dir(currDir) + "/serv"
		t, _ := template.ParseFiles(wwwDir + "/editor.html")
		_ = t.Execute(w, global{R: rand.New(rand.NewSource(time.Now().UnixNano())).Int63()})

	})

	http.HandleFunc("/make", func(w http.ResponseWriter, r *http.Request) {
		var commentString string

		r.ParseForm()
		d := r.Form["request"][0]

		lib.FindRequest(d, "")
		req := lib.AllRequest[0]
		comment := lib.MakeComment(req)
		lib.AllRequest = lib.AllRequest[1:]

		commentString = joinComment(commentString, "/**")
		for _, c := range comment {
			commentString = joinComment(commentString, " *"+c)
		}

		commentString = joinComment(commentString, " */")
		commentString = joinComment(commentString, "\n\n")

		funcStruct := lib.MakeFuncStruct(req)

		for _, f := range funcStruct {
			commentString = joinComment(commentString, " "+f)
		}

		commentString = joinComment(commentString, "\n\n")

		w.Write([]byte(commentString))

	})

	http.ListenAndServe(":9999", nil)
}
