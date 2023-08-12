package main

import (
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

type Rules struct {
	Resources     []ResourceRule     `hcl:"resource,block"`
	ListResources []ListResourceRule `hcl:"list_resource,block"`
}

type ResourceRule struct {
	Type      string      `hcl:"type,label"`
	Name      string      `hcl:"name,label"`
	OnCreate  *HookParams `hcl:"on_create,block"`
	OnDestroy *HookParams `hcl:"on_destroy,block"`
}

func (r *ResourceRule) ValidateInputs() error {
	if r.OnCreate != nil {
		if !(r.OnCreate.Severity == SeverityInform || r.OnCreate.Severity == SeverityWarn || r.OnCreate.Severity == SeverityError) {
			return fmt.Errorf("'%s' must be a supported severity of 'inform', 'warn', or 'error'", r.OnCreate.Severity)
		}
	}

	if (r.OnDestroy != nil && r.OnCreate != nil) || (r.OnDestroy == nil && r.OnCreate == nil) {
		return fmt.Errorf("either on_create or on_destroy must be specified")
	}
	return nil
}

func (r *ResourceRule) Eval(resource *tfjson.ResourceChange) Diagnostics {
	var diags Diagnostics
	for _, action := range resource.Change.Actions {
		if action == "create" && r.OnCreate != nil {
			diags = diags.AppendAll(r.OnCreate.Eval())
		}
	}
	return diags
}

type HookParams struct {
	Severity Severity `hcl:"severity"`
	Message  string   `hcl:"message"`
}

func (h HookParams) Eval() Diagnostics {
	return Diagnostics{{Severity: string(h.Severity), Summary: h.Message}}
}

type ListResourceRule struct {
	Type      string      `hcl:"type,label"`
	Name      string      `hcl:"name,label"`
	Index     string      `hcl:"index,label"`
	OnCreate  *HookParams `hcl:"on_create,block"`
	OnDestroy *HookParams `hcl:"on_destroy,block"`
}

func (r *ListResourceRule) ValidateInputs() error {
	if r.OnCreate != nil {
		if !(r.OnCreate.Severity == SeverityInform || r.OnCreate.Severity == SeverityWarn || r.OnCreate.Severity == SeverityError) {
			return fmt.Errorf("'%s' must be a supported severity of 'inform', 'warn', or 'error'", r.OnCreate.Severity)
		}
	}

	if (r.OnDestroy != nil && r.OnCreate != nil) || (r.OnDestroy == nil && r.OnCreate == nil) {
		return fmt.Errorf("either on_create or on_destroy must be specified")
	}
	return nil
}

func (r *ListResourceRule) Eval(resource *tfjson.ResourceChange) Diagnostics {
	var diags Diagnostics
	for _, action := range resource.Change.Actions {
		if action == "create" && r.OnCreate != nil {
			diags = diags.AppendAll(r.OnCreate.Eval())
		}
	}
	return diags
}
