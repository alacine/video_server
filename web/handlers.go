package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		pg := &HomePage{Name: "ruian"}
		t, err := template.ParseFiles("./templates_example/home.html")
		if err != nil {
			log.Printf("Parsing template home.html error: %v", err)
			return
		}
		t.Execute(w, pg)
		return
	}
	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
	}
}

func userHomeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fname := r.FormValue("username")
	var pg *UserPage
	if len(cname.Value) != 0 {
		pg = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		pg = &UserPage{Name: fname}
	}
	t, err := template.ParseFiles("./templates_example/userhome.html")
	if err != nil {
		log.Printf("Parsing userhome.html error %v", err)
		return
	}
	t.Execute(w, pg)
}

func apiHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != http.MethodPost {
		re, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}
	res, _ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}
	if err := json.Unmarshal(res, apibody); err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(re))
		return
	}
	resquest(apibody, w, r)
	defer r.Body.Close()
}
