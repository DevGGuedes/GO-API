package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//Usuario
type Usuario struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

//UsuarioHandler analisa o request e delega para a função
func UsuarioHandler(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/usuarios/")
	id, _ := strconv.Atoi(sid)

	switch {
	case r.Method == "GET" && id > 0:
		UsuarioPorID(w, r, id)
	case r.Method == "GET":
		usuarioTodos(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Erro....")
	}
}

func UsuarioPorID(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("mysql", "root:@/unip_lpbd")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var u Usuario
	db.QueryRow("SELECT cd_funcionario, nm_funcionario FROM funcionario where cd_funcionario = ?", id).Scan(&u.ID, &u.Nome)

	json, _ := json.Marshal(u)

	w.Header().Set("Content-Type", "applicantion/json")
	fmt.Fprint(w, string(json))

}

func usuarioTodos(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@/unip_lpbd")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, _ := db.Query("SELECT cd_funcionario, nm_funcionario FROM funcionario")
	defer rows.Close()

	var usuarios []Usuario

	for rows.Next() {
		var usuario Usuario
		rows.Scan(&usuario.ID, &usuario.Nome)
		usuarios = append(usuarios, usuario)
	}

	json, _ := json.Marshal(usuarios)

	w.Header().Set("Content-Type", "applicantion/json")
	fmt.Fprint(w, string(json))
}
