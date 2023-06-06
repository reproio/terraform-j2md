package format

import (
	"reflect"
	"testing"
)

func Test_formatJsonChangeValue(t *testing.T) {
	type args struct {
		old interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "plain string",
			args: args{
				old: "plain string",
			},
			want:    "plain string",
			wantErr: false,
		},
		{
			name: "json string",
			args: args{
				old: `{"foo":"bar"}`,
			},
			want: `{
  "foo": "bar"
}`,
			wantErr: false,
		},
		{
			name: "map with json string",
			args: args{
				old: map[string]interface{}{
					"policy": `{"foo":"bar"}`,
					"foo":    "bar",
				},
			},
			want: map[string]interface{}{
				"policy": `{
  "foo": "bar"
}`,
				"foo": "bar",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatJsonChangeValue(tt.args.old)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatJsonChangeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatJsonChangeValue(): \ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}
