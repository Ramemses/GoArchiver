package shanon_fano

import (
	"testing"
	"reflect"

	"archiver/lib/compression/vlc/table"
)


func Test_bestDividePosition(t* testing.T){
	tests := []struct{
		name string
		codes []code
		want int
	}{

		{
			name : "one element",
			codes: []code{
					{Quantity: 2},
				} ,
			want: 0,
		},
		{
			name : "two elements",
			codes: []code{
					{Quantity: 2},
					{Quantity: 2},
				} ,
			want: 1,
		},
		{
			name : "three elements",
			codes: []code{
					{Quantity: 2},
					{Quantity: 1},
					{Quantity: 1},
				} ,
			want: 1,
		},
		{
			name : "important occasion (need rightmost)",
			codes: []code{
					{Quantity: 1},
					{Quantity: 1},
					{Quantity: 1},
				} ,
			want: 1,
		},
		{
			name : "many elements",
			codes: []code{
					{Quantity: 3},
					{Quantity: 2},
					{Quantity: 3},
					{Quantity: 3},
					{Quantity: 3},
					{Quantity: 1},
				} ,
			want: 3,
		},
		
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t* testing.T){
			if got := bestDividePosition(tt.codes); got != tt.want{
				t.Errorf("assingCodes() = #%v#, want #%v#", got, tt.want)
			}
		})

	}

}

		
func Test_assignCodes(t* testing.T){
	tests := []struct{
		name string
		codes []code
		want []code
	}{
		{
			name : "one element",
			codes: []code{
					{Quantity: 2},
				} ,
			want: []code{
					{Quantity: 2, Bits: 0, Size: 1},
			},
		},
		{
			name : "two elements",
			codes: []code{
					{Quantity: 2},
					{Quantity: 2},
				} ,
			want: []code{
					{Quantity: 2, Bits: 0, Size: 1},
					{Quantity: 2, Bits: 1, Size: 1},	
			},
		},
		{
			name : "three elements, certian position",
			codes: []code{
					{Quantity: 2},
					{Quantity: 1},
					{Quantity: 1},
				} ,
			want: []code{
					{Quantity: 2, Bits: 0, Size: 1},
					{Quantity: 1, Bits: 2, Size: 2},	
					{Quantity: 1, Bits: 3, Size: 2},	
			},
		},
		{
			name : "three elements, certian position",
			codes: []code{
					{Quantity: 2},
					{Quantity: 2},
					{Quantity: 2},
				} ,
			want: []code{
					{Quantity: 2, Bits: 0, Size: 1},
					{Quantity: 2, Bits: 2, Size: 2},	
					{Quantity: 2, Bits: 3, Size: 2},	
			},
		},
		
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t* testing.T){
			assignCodes(tt.codes)

			if !reflect.DeepEqual(tt.codes, tt.want){
				t.Errorf("assingCode()  = got: %v, want: %v", tt.codes, tt.want)
			}
		})

	}

}

func Test_build(t* testing.T){
	tests := []struct{
		name string
		str string
		want encodingTable
	}{
		{
			name : "base test",
			str : "abbbcc",
			want: encodingTable{
				'a': code{
					Char: 'a',
					Quantity: 1,
					Bits: 3,
					Size: 2,
				},
				'b': code{
					Char: 'b',
					Quantity: 3,
					Bits: 0,
					Size: 1,
				},
				'c': code{
					Char: 'c',
					Quantity: 2,
					Bits: 2,
					Size: 2,
				},

			},
		},
		
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t* testing.T){
			got := build(tt.str)

			if !reflect.DeepEqual(got, tt.want){
				t.Errorf("build() = got: %v, want: %v", got, tt.want)
			}
		})

	}

}

		
func Test_NewTable(t* testing.T){
	type args struct{
		g Generator
		text string
	}


	tests := []struct{
		name string
		args args
		want table.EncodingTable 
	}{
		{
			name : "base test",
			args : args{
				g : NewGenerator(),
				text : "abbbcc",
			},
			want: table.EncodingTable{
				'a': "11",
				'b': "0",
				'c': "10",
				},
		},
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t* testing.T){
			got := tt.args.g.NewTable(tt.args.text)

			if !reflect.DeepEqual(got, tt.want){
				t.Errorf("NewTable() = got: %v, want: %v", got, tt.want)
			}
		})

	}

}

		








