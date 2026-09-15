package model

import (
	"reflect"
	"testing"
)

func FuzzScriptedProjection(f *testing.F) {
	for _, s := range []string{"", `{\k20}one{\K30}two`, `{\p1}m 0 0 l 1 1{\p0}text`, `😀{\t(0,100,\bord2)}é`, "bad{nested{brace}", `a\N\n\hb`} {
		f.Add(s, uint8(0))
	}
	f.Fuzz(func(t *testing.T, text string, wrap uint8) {
		if len(text) > 4096 {
			t.Skip()
		}
		a, e := ProjectScriptedText(text, int(wrap%4), 1000, 2000)
		b, e2 := ProjectScriptedText(text, int(wrap%4), 1000, 2000)
		if (e == nil) != (e2 == nil) || !reflect.DeepEqual(a, b) {
			t.Fatal("nondeterministic projection")
		}
		if e != nil {
			return
		}
		r := []rune(text)
		position := 0
		for _, s := range a.Spans {
			if s.StartScalar != position || s.EndScalar < s.StartScalar || s.EndScalar > len(r) || s.Raw != string(r[s.StartScalar:s.EndScalar]) {
				t.Fatalf("invalid source span: %#v", s)
			}
			position = s.EndScalar
		}
		if position != len(r) {
			t.Fatal("native Text span content disappeared")
		}
		for _, tag := range a.Tags {
			if tag.StartScalar < 0 || tag.EndScalar < tag.StartScalar || tag.EndScalar > len(r) || tag.Raw != string(r[tag.StartScalar:tag.EndScalar]) {
				t.Fatalf("invalid tag offsets: %#v", tag)
			}
		}
		end := int64(1000)
		for _, token := range a.Tokens {
			if token.StartMilliseconds < end || token.EndMilliseconds < token.StartMilliseconds || token.EndMilliseconds > 2000 {
				t.Fatalf("invalid derived token timing: %#v", token)
			}
			end = token.EndMilliseconds
		}
	})
}
