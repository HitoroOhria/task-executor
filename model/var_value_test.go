package model

import "testing"

func TestVarValue_IsOptional(t *testing.T) {
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
			name: "オプショナル変数の場合、true を返すこと",
			v:    VarValue(`{{.VARIABLE}}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "オプショナル変数であり、空白がある場合、true を返すこと",
			v:    VarValue(`{{   .VARIABLE   }}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "オプショナル変数であり、空白がない場合、true を返すこと",
			v:    VarValue(`{{.VARIABLE}}`),
			args: args{
				name: "VARIABLE",
			},
			want: true,
		},
		{
			name: "デフォルト値ありのオプショナル変数（パイプ）の場合、false を返すこと",
			v:    VarValue(`{{.VARIABLE | default "default-value"}}`),
			args: args{
				name: "VARIABLE",
			},
			want: false,
		},
		{
			name: "デフォルト値ありのオプショナル変数（先頭）の場合、false を返すこと",
			v:    VarValue(`{{default "default-value" .VARIABLE}}`),
			args: args{
				name: "VARIABLE",
			},
			want: false,
		},
		{
			name: "特殊文字を含む変数名の場合、true を返すこと",
			v:    VarValue(`{{.VAR_NAME-WITH.SPECIAL}}`),
			args: args{
				name: "VAR_NAME-WITH.SPECIAL",
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
			name: "部分的に一致する変数名の場合、false を返すこと",
			v:    VarValue(`{{.VARIABLE_OTHER}}`),
			args: args{
				name: "VARIABLE",
			},
			want: false,
		},
		{
			name: "異なる変数名の場合、false を返すこと",
			v:    VarValue(`{{.OTHER_VARIABLE}}`),
			args: args{
				name: "VARIABLE",
			},
			want: false,
		},
		{
			name: "空文字列の場合、false を返すこと",
			v:    VarValue(``),
			args: args{
				name: "VARIABLE",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsOptional(tt.args.name); got != tt.want {
				t.Errorf("IsOptional() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
