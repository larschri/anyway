package anyway

import (
	"fmt"
	"testing"
)

func TestLookup(t *testing.T) {
	var tests = []struct {
		keys   []any
		out    any
		errMsg string
	}{
		{[]any{"foo", "banan"}, 42, ""},
		{[]any{"myslice", 1, "bogo"}, 1, ""},
		{[]any{"myslice", 2, "bogo"}, nil, "out of bounds at 1:2"},
	}

	thing := map[string]any{
		"foo": map[string]any{
			"banan": 42,
		},
		"myslice": []any{
			"boh",
			map[string]any{
				"bogo": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.keys), func(t *testing.T) {
			r, err := Lookup[any](thing, tt.keys...)
			if err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("unexpected error %v", err.Error())
				}
				return
			}

			if r != tt.out {
				t.Errorf("unexpected result %v", r)
			}
		})
	}
}

func ExampleMust() {
	mp := map[string]any{
		"foo": []any{
			42,
		},
	}
	Must(Lookup[[]any](mp, "foo"))[0] = 43
	fmt.Println(mp)
	// Output: map[foo:[43]]
}

func TestSkeleton(t *testing.T) {
	thing := map[string]any{
		"foo": map[string]any{
			"banan": "hello",
		},
		"myslice": []any{
			map[string]any{
				"bogo": "hello",
			},
		},
	}
	skeleton := Skeleton(thing, "")

	lu, err := Lookup[any](skeleton, "foo", "banan")
	if err != nil {
		t.Fatal(err)
	}
	if lu != ".foo.banan" {
		t.Errorf("unexpected %v", lu)
	}

	lu, err = Lookup[any](skeleton, "myslice", 0, "bogo")
	if err != nil {
		t.Fatal(err)
	}
	if lu != ".myslice[0].bogo" {
		t.Errorf("unexpected %v", lu)
	}
}
