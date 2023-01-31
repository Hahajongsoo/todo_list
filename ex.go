package main

import (
	"ch31/db_utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render

type Todo struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	Completed bool   `json:"completed,omitempty"`
}

var todoMap map[int]Todo

func MakeWebHandler() http.Handler {
	todoMap = make(map[int]Todo)
	mux := mux.NewRouter()
	mux.Handle("/", http.FileServer(http.Dir("public")))
	mux.HandleFunc("/todos", GetTodoListHandler).Methods("GET")
	mux.HandleFunc("/todos", PostTodoHandler).Methods("POST")
	mux.HandleFunc("/todos/{id:[0-9]+}", RemoveTodoHandler).Methods("DELETE")
	mux.HandleFunc("/todos/{id:[0-9]+}", UpdateTodoHandler).Methods("PUT")
	return mux
}

type Todos []Todo

func (t Todos) Len() int {
	return len(t)
}

func (t Todos) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Todos) Less(i, j int) bool {
	return t[i].ID > t[j].ID
}

func GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
	todos := make(Todos, 0)
	db := db_utils.UseDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM todo_list")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	err = db_utils.RowsToStructs(rows, &todos)
	if err != nil {
		log.Fatal(err)
	}
	rd.JSON(w, http.StatusOK, todos)
}

func PostTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db := db_utils.UseDB()
	defer db.Close()

	_, err = db.Exec("INSERT INTO todo_list VALUES (NULL, ?, ?)", todo.Name, todo.Completed)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var id int
	err = db.QueryRow("SELECT MAX(id) FROM todo_list").Scan(&id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	todo.ID = id
	rd.JSON(w, http.StatusCreated, todo)
}

type Success struct {
	Success bool `json:"success"`
}

func RemoveTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	db := db_utils.UseDB()
	defer db.Close()
	_, err := db.Exec("DELETE FROM todo_list WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
		rd.JSON(w, http.StatusNotFound, Success{false})
		return
	} else {
		rd.JSON(w, http.StatusOK, Success{true})
		_, err = db.Exec("SET @COUNT=0")
		if err != nil {
			log.Fatal(err)
			rd.JSON(w, http.StatusNotFound, Success{false})
			return
		}
		_, err = db.Exec("UPDATE todo_list SET todo_list.id = @COUNT:=@COUNT+1")
		if err != nil {
			log.Fatal(err)
			rd.JSON(w, http.StatusNotFound, Success{false})
			return
		}
		_, err = db.Exec("ALTER TABLE todo_list AUTO_INCREMENT=1")
		if err != nil {
			log.Fatal(err)
			rd.JSON(w, http.StatusNotFound, Success{false})
			return
		}
	}
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	db := db_utils.UseDB()
	defer db.Close()
	_, err = db.Exec("UPDATE todo_list SET completed = ? WHERE id = ?", newTodo.Completed, id)
	if err != nil {
		log.Fatal(err)
		rd.JSON(w, http.StatusNotFound, Success{false})
		return
	} else {
		rd.JSON(w, http.StatusOK, Success{true})
	}
}

func main() {
	rd = render.New()
	m := MakeWebHandler()
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Started App")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
