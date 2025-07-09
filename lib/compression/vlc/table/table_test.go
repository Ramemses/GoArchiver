package table

import (
	"testing"
	"reflect"
)


func Test_BuildEncodingTree(t* testing.T){
	tests := []struct{
		name string
		ec encodingTable
		want DecodingTree
	}{
		{
			name: "base test",
			ec: encodingTable{
				'a': "11",
				'b' : "1001",
 				'z' : "0101", 
			},
			want: DecodingTree{
				Left: &DecodingTree{
					Right: &DecodingTree{
						Left: &DecodingTree{
							Right: &DecodingTree{
								Data: "z",
							},
						},

					},
				},
				Right: &DecodingTree{
					Left: &DecodingTree{
						Left: &DecodingTree{
							Right: &DecodingTree{
								Data: "b",
							},
						},
					},
					Right: &DecodingTree{
						Data: "a",
					},
				},

			},
		},
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t* testing.T){
			if got := tt.ec.DecodingTree(); !reflect.DeepEqual(got, tt.want){
				t.Errorf("DecodingTree() = #%v#, want #%v#", got, tt.want)
			}
		})

	}

}

