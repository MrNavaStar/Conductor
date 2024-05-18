package main

import (
	_ "embed"
	"github.com/creasty/defaults"
	"github.com/magiconair/properties"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type TemplateInfo struct {
	Container  string `default:"mrnavastar/conductor:server"`
	User       string `default:"conductor"`
	WorkingDir string `default:"/conductor" yaml:"working-dir"`
}

type TemplateActions struct {
	Install   string `default:""`
	Update    string `default:""`
	Start     string `default:""`
	Stop      string `default:""`
	Broadcast string `default:""`
}

type TemplateMembers struct {
	Info    TemplateInfo
	Actions TemplateActions
}

type Template struct {
	Name    string
	Vars    map[string]string
	Info    TemplateInfo
	Actions TemplateActions
}

func parseTemplateName(templateName string) string {
	if !strings.HasSuffix(templateName, ".yml") {
		return templateName + ".yml"
	}
	return templateName
}

func getRepoTemplateNames() ([]string, error) {
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
		names = append(names, strings.TrimSuffix(repoTree.Tree[i].Path, ".yml"))
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
		if err != nil {
			return nil, err
		}
	}
	return bytes, nil
}

func getTemplate(templateName string) (Template, error) {
	var templateMembers TemplateMembers
	var template Template

	bytes, err := getTemplateAsBytes(parseTemplateName(templateName))
	if err != nil {
		return template, err
	}
	if err := defaults.Set(&templateMembers); err != nil {
		return template, err
	}
	if err := yaml.Unmarshal(bytes, &templateMembers); err != nil {
		return template, err
	}

	templateMap := make(map[string]interface{})
	if err := yaml.Unmarshal(bytes, templateMap); err != nil {
		return template, err
	}
	delete(templateMap, "info")
	delete(templateMap, "actions")

	template.Name = templateName
	template.Info = templateMembers.Info
	template.Actions = templateMembers.Actions
	template.Vars = mapToStringMap(templateMap)
	return template, nil
}

func (template Template) getVarsCmd() string {
	var cmdVars = ""
	for key := range template.Vars {
		cmdVars = key + "=\"" + template.Vars[key] + "\" && " + cmdVars
	}
	return cmdVars
}

func (template Template) getDeployCmd() string {
	return template.getVarsCmd() +
		"wget --output-document=/bin/gdmp https://github.com/MrNavaStar/GDMP/releases/download/1.0.0/gdmp\n" +
		"chmod +x /bin/gdmp\n" +
		"mkdir -p " + template.Info.WorkingDir + "\ncd " + template.Info.WorkingDir + "\n" +
		template.Actions.Install + "\n" +
		template.Actions.Update + "\n" +
		"chown -R " + template.Info.User + ":" + template.Info.User + " ."
}

func (template Template) getUpdateCmd() string {
	return template.getVarsCmd() + "cd " + template.Info.WorkingDir + "\n" + template.Actions.Update
}

func (template Template) getStartCmd() string {
	var cmd string
	if template.Actions.Stop == "" {
		cmd = template.Actions.Start
	} else {
		cmd = "gdmp --term --int -x -w \"" + template.Actions.Stop + "\" " + template.Actions.Start
	}
	return template.getVarsCmd() + "cd " + template.Info.WorkingDir + "\n" + cmd
}

/*func getProcessedTemplate(templateName string, overrides []string) (map[string]string, error) {
	templateVars, err := getTemplateVars(templateName)
	if err != nil {
		return nil, err
	}

	for i, s := range overrides {
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
}*/

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
