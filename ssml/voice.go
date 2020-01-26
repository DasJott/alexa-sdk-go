package ssml

// NewVoice returns a function to return a given string, surrounded by certain voice tags
func NewVoice(name string) func(string) string {
	nme := name
	return func(txt string) string {
		return Voice(nme, txt)
	}
}

// Voice returns text, surrounded by voice tags using name
func Voice(name, text string) string {
	return "<voice name=\"" + name + "\">" + text + "</voice>"
}
