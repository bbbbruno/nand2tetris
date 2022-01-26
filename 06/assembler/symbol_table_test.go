package main

import "testing"

func TestAddSymbol(t *testing.T) {
	st := NewSymbolTableWithOpts(20, make(map[symbol]addr))
	st.AddSymbol("NEW")
	if addr := st.nextAddr; addr != 21 {
		t.Errorf("Expected 21, got %#v", addr)
	}
	if addr := st.symbols["NEW"]; addr != 20 {
		t.Errorf("Expected 20, got %#v", addr)
	}
}

func TestAddr(t *testing.T) {
	st := NewSymbolTableWithOpts(17, map[symbol]addr{"ARG": 1, "NEW": 16})
	testCases := []struct {
		in    symbol
		want  addr
		found bool
	}{
		{in: "NEW", want: 16, found: true},
		{in: "100", want: 100, found: true},
		{in: "OLD", want: 0, found: false},
	}
	for _, test := range testCases {
		if addr, ok := st.Addr(test.in); addr != test.want || ok != test.found {
			t.Errorf("st.Addr(%#v) got (%#v, %#v), want (%#v, %#v)", test.in, addr, ok, test.want, test.found)
		}
	}
}
