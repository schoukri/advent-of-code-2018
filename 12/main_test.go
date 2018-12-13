package main

import "testing"

func Test_prepareGen(t *testing.T) {

	tests := []struct {
		gen   string
		want  string
		want1 int
	}{
		{
			gen:   "#",
			want:  ".....#.....",
			want1: 5,
		},
		{
			gen:   ".....#.....",
			want:  ".....#.....",
			want1: 0,
		},
		{
			gen:   "....#....",
			want:  ".....#.....",
			want1: 1,
		},
		{
			gen:   "#..#.#..##......###...###",
			want:  ".....#..#.#..##......###...###.....",
			want1: 5,
		},
		{
			gen:   ".....#..#.#..##......###...###.....",
			want:  ".....#..#.#..##......###...###.....",
			want1: 0,
		},
		{
			gen:   ".....#...#....#.....#..#..#..#.....",
			want:  ".....#...#....#.....#..#..#..#.....",
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.gen, func(t *testing.T) {
			got, got1 := prepareGen(tt.gen)
			if got != tt.want {
				t.Errorf("prepareGen() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("prepareGen() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
