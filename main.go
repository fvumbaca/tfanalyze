package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/hashicorp/hcl/v2/hclsimple"
	tfjson "github.com/hashicorp/terraform-json"
)

type Severity string

const (
	SeverityInform Severity = "inform"
	SeverityWarn   Severity = "warn"
	SeverityError  Severity = "error"
)

func main() {
	var projectDir string
	ctx := context.Background()

	if len(os.Args) > 1 {
		projectDir = os.Args[1]
	}

	rulesFilename := filepath.Join(projectDir, "rules.hcl")

	var rules Rules

	err := hclsimple.DecodeFile(rulesFilename, nil, &rules)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var validationErrors bool
	for _, rule := range rules.Resources {
		err = rule.ValidateInputs()
		if err != nil {
			fmt.Println("Invalid rule:", err)
			validationErrors = true
		}
	}

	if validationErrors {
		os.Exit(1)
	}

	plan, err := getPlan(ctx, projectDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	diags := evalPlan(plan, rules)
	for _, diag := range diags {
		fmt.Println(diag.Summary)
	}

	if diags.IsFatal() {
		os.Exit(1)
	}
}

func evalPlan(plan *tfjson.Plan, rules Rules) Diagnostics {
	var diags Diagnostics
	for _, resource := range plan.ResourceChanges {
		for _, rule := range rules.Resources {
			if resource.Type == rule.Type && resource.Name == rule.Name {
				diags = diags.AppendAll(rule.Eval(resource))
			}
		}

		for _, rule := range rules.ListResources {
			if resource.Type == rule.Type && resource.Name == rule.Name && (fmt.Sprint(resource.Index) == rule.Index || rule.Index == "*") {
				diags = diags.AppendAll(rule.Eval(resource))
			}
		}
	}
	return diags
}

func getPlan(ctx context.Context, projectDir string) (*tfjson.Plan, error) {
	id := uuid.New().String()
	tmpDir := os.TempDir()

	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return nil, err
	}

	planFilename := filepath.Join(tmpDir, id+".tfplan")

	cmd := exec.CommandContext(ctx, "terraform", "plan", "-out", planFilename)
	if projectDir != "" {
		cmd.Dir = projectDir
	}

	cmd.Stdout = io.Discard
	cmd.Stderr = os.Stderr // TODO: capture stderr in case of error

	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	defer os.Remove(planFilename)

	cmd = exec.CommandContext(ctx, "terraform", "show", "-json", planFilename)

	if projectDir != "" {
		cmd.Dir = projectDir
	}

	var buff bytes.Buffer

	cmd.Stdout = &buff
	cmd.Stderr = os.Stderr // TODO: capture stderr in case of error

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	var plan tfjson.Plan

	err = plan.UnmarshalJSON(buff.Bytes())
	if err != nil {
		return nil, err
	}

	return &plan, nil
}
