package set

import (
	"reflect"
	"testing"
)

func TestMakeSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []string
		want map[string]struct{}
	}{
		{name: "empty", in: nil, want: map[string]struct{}{}},
		{name: "single", in: []string{"a"}, want: map[string]struct{}{"a": {}}},
		{name: "duplicates", in: []string{"a", "a", "b"}, want: map[string]struct{}{"a": {}, "b": {}}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := MakeSet(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("MakeSet() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	t.Parallel()

	a := MakeSet([]string{"a", "b"})
	b := MakeSet([]string{"b", "c"})
	empty := MakeSet[string](nil)

	if got := Intersect(a, b); !reflect.DeepEqual(got, map[string]struct{}{"b": {}}) {
		t.Fatalf("Intersect(a,b) = %v", got)
	}

	if got := Intersect(a, empty); len(got) != 0 {
		t.Fatalf("Intersect with empty should be empty, got %v", got)
	}

	if got := Intersect(empty, empty); len(got) != 0 {
		t.Fatalf("Intersect(empty,empty) should be empty, got %v", got)
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()

	a := MakeSet([]string{"a", "b"})
	b := MakeSet([]string{"b", "c"})
	empty := MakeSet[string](nil)

	if got := Union(a, b); !reflect.DeepEqual(got, map[string]struct{}{"a": {}, "b": {}, "c": {}}) {
		t.Fatalf("Union(a,b) = %v", got)
	}

	if got := Union(a, empty); !reflect.DeepEqual(got, a) {
		t.Fatalf("Union with empty should be a, got %v", got)
	}

	if got := Union(empty, empty); len(got) != 0 {
		t.Fatalf("Union(empty,empty) should be empty, got %v", got)
	}
}
