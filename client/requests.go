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

type Funcionario struct {
	Id           int    `json:"id"`
	Nome         string `json:"nome"`
	Email        string `json:"email"`
	Cpf          int    `json:"cpf"`
	Salario      string `json:"salario"`
	Idade        int    `json:"idade"`
	Departamento int    `json:"departamento"`
}

type Response struct {
	Id int `json:"id"`
}

func RequestsHandler(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/funcionarios/")
	id, _ := strconv.Atoi(sid)

	switch {
	case r.Method == "GET" && id > 0:
		FuncionarioPorId(w, r, id)
	case r.Method == "GET":
		BuscaTodosFuncionarios(w, r)
	case r.Method == "POST":
		CadastraFuncionario(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Endpoint não localizado")
	}

}

func CadastraFuncionario(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var f Funcionario
	var resp Response

	err := decoder.Decode(&f)
	if err != nil {
		panic(err)
	}

	/*log.Println(f.Nome)
	log.Println(f.Idade)
	log.Println(f.Email)*/

	db, err := sql.Open("mysql", "root:@/unip_lpbd")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO funcionario (nm_funcionario, ds_email_funcionario, cd_cpf_funcionario, vl_salario_funcionario,  idade_funcionario, cd_departamento) VALUES (?,?,?,?,?,?)")

	res, _ := stmt.Exec(f.Nome, f.Email, f.Cpf, f.Salario, f.Idade, f.Departamento)
	id, _ := res.LastInsertId()

	resp.Id = int(id)

	json, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "applicantion/json")
	fmt.Fprint(w, string(json))

}

func FuncionarioPorId(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("mysql", "root:@/unip_lpbd")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var f Funcionario
	db.QueryRow("SELECT cd_funcionario, nm_funcionario, ds_email_funcionario, cd_cpf_funcionario, vl_salario_funcionario, idade_funcionario, cd_departamento FROM funcionario where cd_funcionario = ?", id).Scan(&f.Id, &f.Nome, &f.Email, &f.Cpf, &f.Salario, &f.Idade, &f.Departamento)

	json, _ := json.Marshal(f)

	w.Header().Set("Content-Type", "applicantion/json")
	fmt.Fprint(w, string(json))

}

func BuscaTodosFuncionarios(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@/unip_lpbd")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, _ := db.Query("SELECT cd_funcionario, nm_funcionario, ds_email_funcionario, cd_cpf_funcionario, vl_salario_funcionario, idade_funcionario, cd_departamento FROM funcionario")
	defer rows.Close()

	var funcionarios []Funcionario

	for rows.Next() {
		var funcionario Funcionario
		rows.Scan(&funcionario.Id, &funcionario.Nome, &funcionario.Email, &funcionario.Cpf, &funcionario.Salario, &funcionario.Idade, &funcionario.Departamento)
		funcionarios = append(funcionarios, funcionario)
	}

	json, _ := json.Marshal(funcionarios)
	w.Header().Set("Content-Type", "applicantion/json")
	fmt.Fprint(w, string(json))

}
