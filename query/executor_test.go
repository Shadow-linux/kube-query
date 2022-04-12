package query

import "testing"

func TestCmdRunWithFile(t *testing.T) {
	type args struct {
		cmd        string
		autoOutput bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "xxx",
			args: args{
				cmd:        "echo A | grep -i a",
				autoOutput: false,
			},
			want: "A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CmdRunWithFile(tt.args.cmd, tt.args.autoOutput); got != tt.want {
				t.Errorf("CmdRunWithFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
