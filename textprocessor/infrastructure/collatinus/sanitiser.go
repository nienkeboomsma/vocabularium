package collatinus

import "strings"

func sanitise(data []byte) string {
	string := string(data)

	replacer := strings.NewReplacer(
		"\x00", "",
		"\n", " ",
		"\r", " ",
	)
	safe := replacer.Replace(string)

	return strings.TrimSpace(safe)
}
