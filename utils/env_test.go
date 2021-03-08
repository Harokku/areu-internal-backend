package utils

import "testing"

func TestReadEnv(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Variable exist",
			args:    args{s: "SET"},
			want:    "OK",
			wantErr: false,
		},
		{
			name:    "Variable not found",
			args:    args{s: "nonExistentVar"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Varialbe not set",
			args:    args{s: "NOTSET"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadEnv(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
