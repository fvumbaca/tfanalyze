package main

import "bytes"

type Diagnostic struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
}

type Diagnostics []Diagnostic

func (d Diagnostics) Append(diag Diagnostic) Diagnostics {
	return append(d, diag)
}

func (d Diagnostics) AppendAll(diags Diagnostics) Diagnostics {
	return append(d, diags...)
}

func (d Diagnostics) Error() string {
	var buff bytes.Buffer

	for _, diag := range d {
		buff.WriteString(diag.Summary)
		buff.WriteString("\n")
	}
	return buff.String()
}

func (d Diagnostics) IsFatal() bool {
	for _, diag := range d {
		if diag.Severity == string(SeverityError) {
			return true
		}
	}
	return false
}
