package model

import "testing"

func TestVarValue_IsOptionalWithDefault(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name string
		v    VarValue
		args args
		want bool
	}{
		{
			name: "パイプ経由の default の場合、true を返すこと",
			v:    VarValue(`{{.VARIABLE | default "default-value"}}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "パイプ経由の default であり、空白がない場合、true を返すこと",
			v:    VarValue(`{{.VARIABLE|default "default-value"}}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "パイプ経由の default であり、空白がある場合、true を返すこと",
			v:    VarValue(`{{   .VARIABLE   |   default   "default-value"  }}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "パイプ経由の default であり、デフォルト値が他の変数に依存している場合、true を返すこと",
			v:    VarValue(`{{.VARIABLE | default .ANOTHER}}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "先頭に default がある場合、true を返すこと",
			v:    VarValue(`{{default "default-value" .VARIABLE}}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "先頭に default があり、空白がある場合、true を返すこと",
			v:    VarValue(`{{   default   "default-value"   .VARIABLE  }}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "先頭に default があり、デフォルト値が他の変数に依存している場合、true を返すこと",
			v:    VarValue(`{{default .ANOTHER .VARIABLE}}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "値ありの変数の場合、false を返すこと",
			v:    VarValue(`value`),
			args: args{
				name: "VARIABLE",
			},
			want: false,
		},
		{
			name: "デフォルト値なしのオプショナル変数の場合、false を返すこと",
			v:    VarValue(`{{.VARIABLE}}`),
			args: args{
				name: "VARIABLE",
			},
			want: false,
		},
		{
			name: "デフォルト値なしのオプショナル変数であり、空白ありの場合、false を返すこと",
			v:    VarValue(`{{   .VARIABLE   }}`),
			args: args{
				name: "VARIABLE",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsOptionalWithDefault(tt.args.name); got != tt.want {
				t.Errorf("IsOptionalWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
