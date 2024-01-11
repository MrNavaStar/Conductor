package main

import (
	"github.com/creasty/defaults"
	"github.com/magiconair/properties"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type TemplateInfo struct {
	Name       string `default:""`
	Container  string `default:"mrnavastar/conductor:server"`
	User       string `default:"conductor"`
	WorkingDir string `default:"/conductor" yaml:"working-dir"`
}

type TemplateActions struct {
	RootInstall string `default:"" yaml:"root-install"`
	Install     string `default:""`
	Update      string `default:""`
	Start       string `default:""`
	Stop        string `default:""`
	Broadcast   string `default:""`
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

func getTemplateAsBytes(templateName string) ([]byte, error) {
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
	bytes, err := getTemplateAsBytes(templateName)
	if err != nil {
		return template, err
	}

	if err := defaults.Set(&template); err != nil {
		return template, err
	}

	if err := yaml.Unmarshal(bytes, &template); err != nil {
		return template, err
	}

	template.Info.Name = templateName
	return template, nil
}

func getTemplateVars(templateName string) (map[string]string, error) {
	templateMap := make(map[string]interface{})

	data, err := getTemplateAsBytes(templateName)
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
		if ok || arg[0] == "ntfy_url" {
			templateVars[arg[0]] = arg[1]
		}
	}
	return templateVars, nil
}

func templateVarsToCmd(templateVars map[string]string) string {
	var cmd = ""
	for key := range templateVars {
		if len(templateVars[key]) == 0 {
			continue
		}
		cmd = key + "=\"" + templateVars[key] + "\" && " + cmd
	}
	return cmd
}

func getInstallRootCmd(template Template, serverName string, templateVars map[string]string) string {
	templateVars["conductor_template"] = template.Info.Name
	templateVars["conductor_name"] = serverName

	return templateVarsToCmd(templateVars) +
		"mkdir -p " + template.Info.WorkingDir +
		"\ncd " + template.Info.WorkingDir +
		"\n" + template.Actions.RootInstall
}

func getInstallCmd(template Template, templateVars map[string]string) string {
	return templateVarsToCmd(templateVars) +
		"\ncd " + template.Info.WorkingDir +
		"\n" + template.Actions.Install +
		"\n" + template.Actions.Update +
		"\n" + saveServerArgsCmd(templateVars)
}

func getUpdateCmd(serverName string) (string, error) {
	serverArgs, err := readServerArgs(serverName)
	if err != nil {
		return "", err
	}

	template, err := parseTemplate(serverArgs["conductor_template"])
	if err != nil {
		return "", err
	}

	return templateVarsToCmd(serverArgs) +
			"\ncd " + template.Info.WorkingDir +
			"\n" + template.Actions.Update +
			"\n" + saveServerArgsCmd(serverArgs),
		nil
}

func getStartCmd(template Template, templateVars map[string]string) string {
	cmd := templateVarsToCmd(templateVars) + "cd " + template.Info.WorkingDir + " && " + template.Actions.Start + " &"
	println(cmd)
	_, ok := templateVars["ntfy_url"]
	if ok {
		cmd += " || curl -d \"🛑 Oh no! $conductor_template server $conductor_name has crashed with exit code $?!\" $ntfy_url"
	}

	return cmd
}

func readServerArgs(serverName string) (map[string]string, error) {
	vars, err := properties.LoadFile(getAppDir()+"/servers/"+serverName+"/.conductor.properties", properties.UTF8)
	if err != nil {
		return nil, err
	}
	return vars.Map(), nil
}

func saveServerArgsCmd(templateVars map[string]string) string {
	var cmd = "rm -f .conductor.properties && echo -e \"# This file is auto generated, DO NOT MODIFY!!!\n# Modifying incorrectly may break this server.\n"
	for key := range templateVars {
		cmd += key + "=${" + key + "}\n"
	}
	return cmd + "\" > .conductor.properties"
}
