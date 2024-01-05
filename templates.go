package main

import (
	"github.com/magiconair/properties"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type TemplateInfo struct {
	Name       string
	Container  string
	User       string
	WorkingDir string `yaml:"working_dir"`
}

type TemplateActions struct {
	Adduser string
	Install string
	Update  string
	Start   string
	Stop    string
}

type Template struct {
	Info    TemplateInfo
	Actions TemplateActions
}

func parseTemplateName(templateName string) string {
	if !strings.HasSuffix(templateName, ".yml") {
		return templateName + ".yml"
	}
	return templateName
}

func getTemplateNames() ([]string, error) {
	repoTree, err := getGithubRepoTree("https://api.github.com/repos/mrnavastar/conductor/git/trees/master?recursive=1")
	if err != nil {
		return nil, err
	}

	var url string
	for i := range repoTree.Tree {
		if repoTree.Tree[i].Path == "templates" {
			url = repoTree.Tree[i].Url
			break
		}
	}

	repoTree, err = getGithubRepoTree(url)
	if err != nil {
		return nil, err
	}

	var names []string
	for i := range repoTree.Tree {
		names = append(names, repoTree.Tree[i].Path)
	}
	return names, nil
}

func getTemplateRaw(templateName string) ([]byte, error) {
	templateName = parseTemplateName(templateName)
	err := os.MkdirAll(getAppDir()+"/templates", os.ModePerm)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(getAppDir() + "/templates/" + templateName)
	if err != nil {
		bytes, err = downloadFile("https://raw.githubusercontent.com/MrNavaStar/Conductor/master/templates/"+templateName, getAppDir()+"/templates/"+templateName)
	}
	return bytes, nil
}

func parseTemplate(templateName string) (Template, error) {
	var template Template

	templateName = parseTemplateName(templateName)
	bytes, err := getTemplateRaw(templateName)
	if err != nil {
		return template, err
	}

	if err := yaml.Unmarshal(bytes, &template); err != nil {
		return template, err
	}
	return template, nil
}

func getTemplateVars(templateName string) (map[string]string, error) {
	templateMap := make(map[string]interface{})

	data, err := getTemplateRaw(templateName)
	if err != nil {
		return mapToStringMap(templateMap), err
	}

	if err := yaml.Unmarshal(data, templateMap); err != nil {
		return mapToStringMap(templateMap), err
	}

	delete(templateMap, "info")
	delete(templateMap, "actions")
	return mapToStringMap(templateMap), nil
}

func overrideTemplateVars(templateName string, cliArgs []string) (map[string]string, error) {
	templateVars, err := getTemplateVars(templateName)
	if err != nil {
		return nil, err
	}

	for i, s := range cliArgs {
		if i == 0 {
			continue
		}

		arg := strings.Split(s, "=")
		if len(arg) != 2 {
			continue
		}

		_, ok := templateVars[arg[0]]
		if ok {
			templateVars[arg[0]] = arg[1]
		}
	}
	return templateVars, nil
}

func parseTemplateVars(templateVars map[string]string) string {
	var cmd = ""
	for key := range templateVars {
		if len(templateVars[key]) == 0 {
			continue
		}
		cmd = key + "=" + templateVars[key] + " && " + cmd
	}
	return cmd
}

func parseServerTemplateVars(serverName string) (string, error) {
	vars, err := properties.LoadFile(getAppDir()+"/servers/"+serverName+"/.conductor.properties", properties.UTF8)
	if err != nil {
		return "", err
	}
	return parseTemplateVars(vars.Map()), nil
}

func saveTemplateVarsCmd(templateVars map[string]string) string {
	var cmd = "echo -e \"# This file is auto generated, DO NOT MODIFY!!!\n# Modifying incorrectly may break this server.\n"
	for key := range templateVars {
		cmd += key + "=${" + key + "}\n"
	}
	return cmd + "\" > .conductor.properties"
}
