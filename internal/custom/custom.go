package custom

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/types"
)

func findPluginWithName(c *config.Config, name string) *config.Plugin {
	for _, plugin := range c.Plugins {
		if plugin.Name == name {
			return &plugin
		}
	}
	return nil
}

func Run(c *config.Config, name string) ([]types.Check, error) {
	plugin := findPluginWithName(c, name)
	checks, err := ExecuteCommand(c, plugin)
	return checks, err

}

func ExecuteCommand(c *config.Config, plugin *config.Plugin) ([]types.Check, error) {
	checks := []types.Check{}
	check := types.Check{}
	check.Name = plugin.Name
	check.Description = plugin.Description
	check.Status = "OK"

	cmd := exec.Command(plugin.Command, plugin.Args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	result := types.Result{}
	if strings.TrimRight(outb.String(), "\n") == plugin.ExpectedOutput {
		result.Message = fmt.Sprint("Output matched: ", plugin.ExpectedOutput)
		result.Status = "OK"
	} else {
		result.Message = fmt.Sprint("Output did not match: ", plugin.ExpectedOutput, " instead got: ", outb.String())
		result.Status = "FAIL"
	}
	check.Results = append(check.Results, result)
	checks = append(checks, check)
	return checks, nil

}