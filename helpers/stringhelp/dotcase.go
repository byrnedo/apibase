package stringhelp
import (
	"unicode/utf8"
	"bytes"
	"unicode"
)

//
// Copyright (c) 2015 Huan Du
// Modifications: 2015 Donal Byrne
// ToDotSnakeCase can convert all upper case characters in a string to
// dot separated format.
//
// Some samples.
//     "FirstName"  => "first.name"
//     "HTTPServer" => "http.server"
//     "NoHTTPS"    => "no.https"
//     "GO_PATH"    => "go.path"
//     "GO PATH"    => "go.path"      // space is converted to underscore.
//     "GO-PATH"    => "go.path"      // hyphen is converted to underscore.
func ToDotSnakeCase(str string) string {
	if len(str) == 0 {
		return ""
	}

	buf := &bytes.Buffer{}
	var prev, r0, r1 rune
	var size int

	r0 = '_'

	for len(str) > 0 {
		prev = r0
		r0, size = utf8.DecodeRuneInString(str)
		str = str[size:]

		switch {
		case r0 == utf8.RuneError:
			buf.WriteByte(byte(str[0]))

		case unicode.IsUpper(r0):
			if prev != '_' {
				buf.WriteRune('_')
			}

			buf.WriteRune(unicode.ToLower(r0))

			if len(str) == 0 {
				break
			}

			r0, size = utf8.DecodeRuneInString(str)
			str = str[size:]

			if !unicode.IsUpper(r0) {
				buf.WriteRune(r0)
				break
			}

			// find next non-upper-case character and insert `_` properly.
			// it's designed to convert `HTTPServer` to `http_server`.
			// if there are more than 2 adjacent upper case characters in a word,
			// treat them as an abbreviation plus a normal word.
			for len(str) > 0 {
				r1 = r0
				r0, size = utf8.DecodeRuneInString(str)
				str = str[size:]

				if r0 == utf8.RuneError {
					buf.WriteRune(unicode.ToLower(r1))
					buf.WriteByte(byte(str[0]))
					break
				}

				if !unicode.IsUpper(r0) {
					if r0 == '_' || r0 == ' ' || r0 == '-' {
						r0 = '_'

						buf.WriteRune(unicode.ToLower(r1))
					} else {
						buf.WriteRune('_')
						buf.WriteRune(unicode.ToLower(r1))
						buf.WriteRune(r0)
					}

					break
				}

				buf.WriteRune(unicode.ToLower(r1))
			}

			if len(str) == 0 || r0 == '_' {
				buf.WriteRune(unicode.ToLower(r0))
				break
			}

		default:
			if r0 == ' ' || r0 == '-' {
				r0 = '_'
			}

			buf.WriteRune(r0)
		}
	}

	return buf.String()
}
