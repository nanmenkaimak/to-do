package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nanmenkaimak/to-do-list/internal/taskstore"
	"mime"
	"net/http"
	"strconv"
	"time"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

func (ts *taskServer) CreateTask(w http.ResponseWriter, r *http.Request) {

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type ResponseID struct {
		ID int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expected application/json Content-type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	renderJSON(w, ResponseID{ID: id})
}

func (ts *taskServer) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	alltasks := ts.store.GetAllTasks()
	renderJSON(w, alltasks)
}

func (ts *taskServer) GetTask(w http.ResponseWriter, r *http.Request) {
	// Here and elsewhere, not checking error of Atoi because the router only
	// matches the [0-9]+ regex.
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *taskServer) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *taskServer) DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	ts.store.DeleteAllTasks()
}

func (ts *taskServer) GetTasksByTag(w http.ResponseWriter, r *http.Request) {
	tag := mux.Vars(r)["tag"]
	tasks := ts.store.GetTasksByTag(tag)
	renderJSON(w, tasks)
}

func (ts *taskServer) GetTasksByDue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	if month < int(time.January) || month > int(time.December) {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", r.URL.Path), http.StatusBadRequest)
		return
	}
	day, _ := strconv.Atoi(vars["day"])

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)

	renderJSON(w, tasks)
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
