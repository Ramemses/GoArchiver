package haffman

import (
	"archiver/lib/compression/vlc/table"
	"reflect"
	"sort"
	"testing"

	"container/heap"
)

func TestNewGenerator_NewTable(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  table.EncodingTable
	}{
		{
			name:  "empty string",
			input: "",
			want:  table.EncodingTable{},
		},
		{
			name:  "single character",
			input: "a",
			want:  table.EncodingTable{'a': "0"},
		},
		{
			name:  "same characters",
			input: "aaaa",
			want:  table.EncodingTable{'a': "0"},
		},
		{
			name:  "two characters with equal frequency",
			input: "abab",
			want: table.EncodingTable{
				'a': "0",
				'b': "1",
			},
		},
		{
			name:  "three characters with equal frequency",
			input: "abcabcabc",
			want: table.EncodingTable{
				'a': "00",
				'b': "01",
				'c': "1",
			},
		},
		{
			name:  "different frequencies",
			input: "aaabbbccd",
			want: table.EncodingTable{
				'a': "1",
				'b': "01",
				'c': "001",
				'd': "000",
			},
		},
		{
			name:  "unicode characters",
			input: "世界世界和平!",
			want: table.EncodingTable{
				'世': "10",
				'界': "11",
				'和': "010",
				'平': "011",
				'!': "00",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGenerator()
			got := g.NewTable(tt.input)

			// Нормализуем таблицы для стабильного сравнения
			normalizedWant := normalizeTable(tt.want)
			normalizedGot := normalizeTable(got)

			// Проверяем соответствие ожидаемых кодов
			if !reflect.DeepEqual(normalizedWant, normalizedGot) {
				t.Errorf("NewTable() codes mismatch\nGot:  %v\nWant: %v", normalizedGot, normalizedWant)
			}

			// Проверяем префиксность кодов
			if !isPrefixCode(got) {
				t.Errorf("Generated codes are not prefix codes: %v", got)
			}

			// Проверяем оптимальность кодирования
			if len(tt.input) > 0 {
				expectedLength := calculateExpectedLength(tt.input, got)
				actualLength := calculateActualLength(tt.input, got)
				if actualLength > expectedLength {
					t.Errorf("Suboptimal coding: actual %d bits, expected %d bits",
						actualLength, expectedLength)
				}
			}
		})
	}
}

// Нормализует таблицу для стабильного сравнения
func normalizeTable(et table.EncodingTable) map[rune]string {
	codes := make([]struct {
		char rune
		code string
	}, 0, len(et))

	for char, code := range et {
		codes = append(codes, struct {
			char rune
			code string
		}{char, code})
	}

	// Сортируем по символам
	sort.Slice(codes, func(i, j int) bool {
		return codes[i].char < codes[j].char
	})

	result := make(map[rune]string)
	for _, item := range codes {
		result[item.char] = item.code
	}
	return result
}

// Вспомогательная функция для проверки префиксных кодов
func isPrefixCode(et table.EncodingTable) bool {
	codes := make([]string, 0, len(et))
	for _, code := range et {
		codes = append(codes, code)
	}

	sort.Strings(codes)

	for i := 0; i < len(codes)-1; i++ {
		// Проверяем, является ли текущий код префиксом следующего
		if len(codes[i]) <= len(codes[i+1]) &&
			codes[i] == codes[i+1][:len(codes[i])] {
			return false
		}
	}
	return true
}

// Вспомогательная функция для расчета ожидаемой длины
func calculateExpectedLength(input string, et table.EncodingTable) int {
	freq := make(map[rune]int)
	for _, ch := range input {
		freq[ch]++
	}

	total := 0
	for ch, count := range freq {
		total += len(et[ch]) * count
	}
	return total
}

// Вспомогательная функция для расчета фактической длины
func calculateActualLength(input string, et table.EncodingTable) int {
	total := 0
	for _, ch := range input {
		total += len(et[ch])
	}
	return total
}

func Test_buildHaffmanTree(t *testing.T) {
	tests := []struct {
		name  string
		queue *Queue
		want  *Node
	}{
		{
			name: "single node tree",
			queue: func() *Queue {
				q := &Queue{
					{Char: 'A', Quantite: 10},
				}
				heap.Init(q)
				return q
			}(),
			want: &Node{
				Char:     'A',
				Quantite: 10,
			},
		},
		{
			name: "two nodes",
			queue: func() *Queue {
				q := &Queue{
					{Char: 'B', Quantite: 5},
					{Char: 'A', Quantite: 10},
				}
				heap.Init(q)
				return q
			}(),
			want: &Node{
				Quantite: 15,
				Left: &Node{
					Char:     'B',
					Quantite: 5,
				},
				Right: &Node{
					Char:     'A',
					Quantite: 10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildHaffmanTree(tt.queue)

			if !compareTrees(got, tt.want) {
				t.Errorf("buildHaffmanTree() tree structure mismatch want = %v, got = %v", tt.want, got)
			}
		})
	}
}

// Вспомогательная функция для сравнения деревьев
func compareTrees(a, b *Node) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.Quantite != b.Quantite {
		return false
	}

	// Сравниваем символы только для листовых узлов
	if a.Left == nil && a.Right == nil {
		if b.Left != nil || b.Right != nil {
			return false
		}
		return a.Char == b.Char
	}

	return compareTrees(a.Left, b.Left) && compareTrees(a.Right, b.Right)
}

func Test_assignCodes(t *testing.T) {
	tests := []struct {
		name  string
		tree  *Node
		want  encodingTable
	}{
		{
			name: "single node",
			tree: &Node{Char: 'A', Quantite: 10},
			want: encodingTable{'A': {Char: 'A', Bits: 0, Size: 1}},
		},
		{
			name: "simple tree",
			tree: &Node{
				Quantite: 15,
				Left:     &Node{Char: 'B', Quantite: 5},
				Right:    &Node{Char: 'A', Quantite: 10},
			},
			want: encodingTable{
				'B': {Char: 'B', Bits: 0, Size: 1},
				'A': {Char: 'A', Bits: 1, Size: 1},
			},
		},
		{
			name: "three level tree",
			tree: &Node{
				Quantite: 6,
				Left: &Node{
					Quantite: 3,
					Left:     &Node{Char: 'C', Quantite: 1},
					Right:    &Node{Char: 'B', Quantite: 2},
				},
				Right: &Node{Char: 'A', Quantite: 3},
			},
			want: encodingTable{
				'C': {Char: 'C', Bits: 0b00, Size: 2},
				'B': {Char: 'B', Bits: 0b01, Size: 2},
				'A': {Char: 'A', Bits: 0b1, Size: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make(encodingTable)
			assignCodes(tt.tree, got)

			for char, wantNode := range tt.want {
				gotNode, ok := got[char]
				if !ok {
					t.Errorf("Character %c missing in encoding table", char)
					continue
				}

				if gotNode.Bits != wantNode.Bits || gotNode.Size != wantNode.Size {
					t.Errorf("Code for %c mismatch: got %d/%d, want %d/%d",
						char, gotNode.Bits, gotNode.Size, wantNode.Bits, wantNode.Size)
				}
			}
		})
	}
}

func Test_encodingTable_Export(t *testing.T) {
	tests := []struct {
		name string
		et   encodingTable
		want table.EncodingTable
	}{
		{
			name: "simple codes",
			et: encodingTable{
				'a': {Bits: 0b0, Size: 1},
				'b': {Bits: 0b1, Size: 1},
			},
			want: table.EncodingTable{
				'a': "0",
				'b': "1",
			},
		},
		{
			name: "multi-bit codes",
			et: encodingTable{
				'a': {Bits: 0b00, Size: 2},
				'b': {Bits: 0b01, Size: 2},
				'c': {Bits: 0b1, Size: 1},
			},
			want: table.EncodingTable{
				'a': "00",
				'b': "01",
				'c': "1",
			},
		},
		{
			name: "leading zeros",
			et: encodingTable{
				'a': {Bits: 0b0, Size: 3},
				'b': {Bits: 0b1, Size: 3},
			},
			want: table.EncodingTable{
				'a': "000",
				'b': "001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.et.Export()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Export() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCharStat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  charStat
	}{
		{
			name:  "empty string",
			input: "",
			want:  charStat{},
		},
		{
			name:  "single character",
			input: "a",
			want:  charStat{'a': 1},
		},
		{
			name:  "multiple characters",
			input: "hello world",
			want: charStat{
				'h': 1,
				'e': 1,
				'l': 3,
				'o': 2,
				' ': 1,
				'w': 1,
				'r': 1,
				'd': 1,
			},
		},
		{
			name:  "unicode characters",
			input: "こんにちは世界",
			want: charStat{
				'こ': 1,
				'ん': 1,
				'に': 1,
				'ち': 1,
				'は': 1,
				'世': 1,
				'界': 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCharStat(tt.input)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCharStat() = %v, want %v", got, tt.want)
			}
		})
	}
}
