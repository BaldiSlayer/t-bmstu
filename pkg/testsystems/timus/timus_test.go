package timus

import "testing"

func Test_getNearest(t *testing.T) {
	type args struct {
		language  string
		languages []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{
				language:  "bicycle",
				languages: []string{"car", "train", "bus", "bicycle"},
			},
			want: "bicycle",
		},
		{
			name: "Test 2",
			args: args{
				language:  "tree",
				languages: []string{"flower", "grass", "tree", "bush"},
			},
			want: "tree",
		},
		{
			name: "Test 3",
			args: args{
				language:  "Python 3.8",
				languages: []string{"C++ Min GW", "Golang 1.20", "Python 3.12", "Algol"},
			},
			want: "Python 3.12",
		},
		{
			name: "Test 3",
			args: args{
				language:  "C#",
				languages: []string{"C++ Min GW", "Golang 1.20", "Python 3.12", "Algol"},
			},
			want: "Algol",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNearest(tt.args.language, tt.args.languages); got != tt.want {
				t.Errorf("getNearest() = %v, want %v", got, tt.want)
			}
		})
	}
}
