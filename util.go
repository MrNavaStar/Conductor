package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
)

type RepoFile struct {
	Path string
	Url  string
}

type RepoTree struct {
	Tree []RepoFile
}

func getGithubRepoTree(url string) (RepoTree, error) {
	var repoTree RepoTree
	resp, err := http.Get(url)
	if err != nil {
		return repoTree, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return repoTree, err
	}

	err = json.Unmarshal(data, &repoTree)
	if err != nil {
		return repoTree, err
	}
	return repoTree, nil
}

func downloadFile(url string, fullPath string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	_, err = file.Write(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func mapToStringMap(currentMap map[string]interface{}) map[string]string {
	newMap := map[string]string{}
	for key := range currentMap {
		switch currentMap[key].(type) {
		case string:
			newMap[key] = currentMap[key].(string)
		case int:
			newMap[key] = strconv.Itoa(currentMap[key].(int))
		case bool:
			newMap[key] = strconv.FormatBool(currentMap[key].(bool))
		}
	}
	return newMap
}

func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + "/.cache/conductor", nil
}
