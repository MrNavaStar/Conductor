package main

import (
	"encoding/json"
	"fmt"
	"github.com/MadAppGang/httplog"
	"github.com/docker/docker/client"
	"net/http"
)

func getMux(cli *client.Client) *http.ServeMux {
	m := http.NewServeMux()
	m.Handle("/api/templates/repo", httplog.Logger(http.HandlerFunc(getRepoTemplates)))
	//m.Handle("/api/templates", httplog.Logger(http.HandlerFunc(getRepoTemplates)))

	m.Handle("/api/servers", httplog.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		servers(w, r, cli)
	})))

	//m.Handle("GET /api/plugins")
	//m.Handle("POST /api/plugins")
	//m.Handle("POST /api/plugins/{name}")
	//m.Handle("DELETE /api/plugins")
	return m
}

func getRepoTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	names, err := getRepoTemplateNames()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(names)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func processServer(server Server, r *http.Request, cli *client.Client) error {
	switch r.Method {
	case "POST":
		return server.deploy(cli)
	case "DELETE":
		return server.delete(cli)
	case "PATCH":
		if server.Running {
			return server.start()
		} else {
			return server.stop()
		}
	}
	return nil
}

func servers(w http.ResponseWriter, r *http.Request, cli *client.Client) {
	var server Server
	err := json.NewDecoder(r.Body).Decode(&server)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = processServer(server, r, cli)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}
