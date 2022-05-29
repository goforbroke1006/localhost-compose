package domain

import (
	"reflect"
	"testing"
)

func TestServiceSpec_CommandList(t *testing.T) {
	type fields struct {
		WorkingDir string
		Build      BuildSpec
		Command    string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "empty list",
			fields: fields{
				WorkingDir: "",
				Build:      BuildSpec{},
				Command:    ``,
			},
			want: nil,
		},
		{
			name: "one command",
			fields: fields{
				WorkingDir: "",
				Build:      BuildSpec{},
				Command: `
					ls /
				`,
			},
			want: []string{
				"ls /",
			},
		},
		{
			name: "multiple commands",
			fields: fields{
				WorkingDir: "",
				Build:      BuildSpec{},
				Command: `
					ls /
					echo 'Hello World'
				`,
			},
			want: []string{
				"ls /",
				"echo 'Hello World'",
			},
		},
		{
			name: "dont split semicolons",
			fields: fields{
				WorkingDir: "",
				Build:      BuildSpec{},
				Command: `
					ls /; echo 'Hello World'
					ping -c 1 google.com
				`,
			},
			want: []string{
				"ls /; echo 'Hello World'",
				"ping -c 1 google.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := ServiceSpec{
				WorkingDir: tt.fields.WorkingDir,
				Build:      tt.fields.Build,
				Command:    tt.fields.Command,
			}
			if got := ss.CommandList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandList() = %v, want %v", got, tt.want)
			}
		})
	}
}
