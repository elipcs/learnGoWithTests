package reflection

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {

	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"struct with one string field",
			struct{ Name string }{"Eli"},
			[]string{"Eli"},
		},
		{
			"struct with two string fields",
			struct {
				Name string
				City string
			}{"Eli", "Campina Grande"},
			[]string{"Eli", "Campina Grande"},
		},
		{
			"struct with non string field",
			struct {
				Name string
				Age  int
			}{"Eli", 33},
			[]string{"Eli"},
		},
		{
			"nested fields",
			Person{
				"Eli",
				Profile{33, "Campina Grande"},
			},
			[]string{"Eli", "Campina Grande"},
		},
		{
			"pointers to things",
			&Person{
				"Eli",
				Profile{33, "Campina Grande"},
			},
			[]string{"Eli", "Campina Grande"},
		},
		{
			"slices",
			[]Profile{
				{33, "Campina Grande"},
				{34, "Joao Pessoa"},
			},
			[]string{"Campina Grande", "Joao Pessoa"},
		},
		{
			"arrays",
			[2]Profile{
				{33, "Campina Grande"},
				{34, "Joao Pessoa"},
			},
			[]string{"Campina Grande", "Joao Pessoa"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})
	}

	t.Run("with maps", func(t *testing.T) {
		aMap := map[string]string{
			"Foo": "Bar",
			"Baz": "Boz",
		}

		var got []string
		walk(aMap, func(input string) {
			got = append(got, input)
		})

		assertContains(t, got, "Bar")
		assertContains(t, got, "Boz")
	})

	t.Run("with channels", func(t *testing.T) {
		aChannel := make(chan Profile)

		go func() {
			aChannel <- Profile{33, "Sao Paulo"}
			aChannel <- Profile{34, "Rio de Janeiro"}
			close(aChannel)
		}()
	
		var got []string
		want := []string{"Sao Paulo", "Rio de Janeiro"}

		walk(aChannel, func(input string) {
			got = append(got, input)
		})
	
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("with function", func(t *testing.T) {
		aFunction := func() (Profile, Profile) {
			return Profile{33, "Berlin"}, Profile{34, "Katowice"}
		}

		var got []string
		want := []string{"Berlin", "Katowice"}

		walk(aFunction, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func assertContains(t testing.TB, haystack []string, needle string) {
	t.Helper()
	contains := false
	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}
	if !contains {
		t.Errorf("expected %v to contain %q but it didn't", haystack, needle)
	}
}
