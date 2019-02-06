// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jlo

import (
	json "encoding/json"
	"sort"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson748ea4a3DecodeGithubComDcmnComJlo(in *jlexer.Lexer, out *Entry) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
	} else {
		in.Delim('{')
		if !in.IsDelim('}') {
			*out = make(Entry)
		} else {
			*out = nil
		}
		for !in.IsDelim('}') {
			key := string(in.String())
			in.WantColon()
			var v1 interface{}
			if m, ok := v1.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := v1.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				v1 = in.Interface()
			}
			(*out)[key] = v1
			in.WantComma()
		}
		in.Delim('}')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson748ea4a3EncodeGithubComDcmnComJlo(out *jwriter.Writer, in Entry) {
	if in == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		keys := make([]string, 0, len(in))
		for v2Name := range in {
			keys = append(keys, v2Name)
		}
		sort.Strings(keys)

		out.RawByte('{')
		v2First := true
		for _, v2Name := range keys {
			if v2First {
				v2First = false
			} else {
				out.RawByte(',')
			}
			out.String(v2Name)
			out.RawByte(':')
			v2Value := in[v2Name]
			if m, ok := v2Value.(easyjson.Marshaler); ok {
				m.MarshalEasyJSON(out)
			} else if m, ok := v2Value.(json.Marshaler); ok {
				out.Raw(m.MarshalJSON())
			} else {
				out.Raw(json.Marshal(v2Value))
			}
		}
		out.RawByte('}')
	}
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Entry) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson748ea4a3EncodeGithubComDcmnComJlo(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Entry) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson748ea4a3DecodeGithubComDcmnComJlo(l, v)
}