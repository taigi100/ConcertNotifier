package utils

import "testing"

func TestLocation_DistanceTo(t *testing.T) {
	type args struct {
		other Location
	}
	tests := []struct {
		name string
		loc  Location
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "",
			loc:  Location{lat: 45.75372, long: 21.22571},            //Temeshvar
			args: args{Location{lat: 37.4900318, long: 136.4664008}}, //Japan
			want: 8725308,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.loc.DistanceTo(tt.args.other); got != tt.want {
				t.Errorf("Location.DistanceTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
