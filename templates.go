package main

import (
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
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cacheDir+"/templates", os.ModePerm)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cacheDir + "/templates/" + templateName)
	if err != nil {
		data, err = downloadFile("https://raw.githubusercontent.com/MrNavaStar/Conductor/master/templates/"+templateName, cacheDir+"/templates/"+templateName)
	}
	return data, nil
}

func parseTemplate(templateName string) (Template, error) {
	var template Template

	templateName = parseTemplateName(templateName)
	data, err := getTemplateRaw(templateName)
	if err != nil {
		return template, err
	}

	if err := yaml.Unmarshal(data, &template); err != nil {
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

/*func parseScript(templateCmd string) string {
	re := regexp.MustCompile(`(\s+\n|\n+)`)
	var cmd, _ = strings.CutSuffix(re.ReplaceAllString(templateCmd, " && "), " && ")
	return cmd
}*/
