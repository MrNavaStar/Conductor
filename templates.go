package main

import (
	"github.com/pterm/pterm"
	urfave "github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
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

func parseTemplate(filename string) (Template, error) {
	var template Template

	data, err := os.ReadFile(filename)
	if err != nil {
		return template, err
	}

	if err := yaml.Unmarshal(data, &template); err != nil {
		return template, err
	}
	return template, nil
}

func getTemplateVars(filename string) (map[string]string, error) {
	templateMap := make(map[string]interface{})

	data, err := os.ReadFile(filename)
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

func parseScript(templateCmd string) string {
	re := regexp.MustCompile(`(\s+\n|\n+)`)
	var cmd, _ = strings.CutSuffix(re.ReplaceAllString(templateCmd, " && "), " && ")
	return cmd
}

func cliGetTemplateVars(c *urfave.Context) urfave.ExitCoder {
	templateName := c.Args().Get(0)

	if len(templateName) == 0 {
		return nil
	}

	if !strings.HasSuffix(templateName, ".yml") {
		templateName = templateName + ".yml"
	}

	vars, err := getTemplateVars("templates/" + templateName)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	for key := range vars {
		pterm.NewRGB(252, 140, 3).Print(key)
		pterm.NewRGB(255, 255, 255).Print(":")
		pterm.NewRGB(3, 252, 90).Println(vars[key])
	}
	return nil
}
