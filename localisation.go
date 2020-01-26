package alexa

import (
	"reflect"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Localisation maps locale codes to according Translations
type Localisation map[string]Translation

// Translation is a set of translated strings. The value can be either a string or a string array
type Translation map[string]interface{}

// Translator actively translates keys to values
type Translator struct {
	Phrases Translation
	Format  func(string, ...interface{}) string
}

// R (for Replace) is a shortcut for map[string]interface{}, while the value must be int (any), float (any) or string
type R map[string]interface{}

// GetTranslator gets a translator for the given language
func (loc Localisation) GetTranslator(locale string) *Translator {
	if ph, ok := loc[locale]; ok {
		lang, err := language.Parse(locale)
		if err != nil {
			lang, err = language.Parse("en-US")
		}
		if err == nil {
			prnt := message.NewPrinter(lang)
			return &Translator{
				Phrases: ph,
				Format: func(str string, stuff ...interface{}) string {
					return prnt.Sprintf(str, stuff...)
				},
			}
		}
	}
	return nil
}

// GetString gets a string from the value according to the given key
func (tr *Translator) GetString(key string) string {
	if val, exists := tr.Phrases[key]; exists {
		switch val.(type) {
		case string:
			return val.(string)
		case []string:
			arr := val.([]string)
			if count := len(arr); count > 0 {
				return arr[random.Intn(len(arr))]
			}
		}
	}
	return ""
}

// GetArray gets an array from the value according to the given key
func (tr *Translator) GetArray(key string) []string {
	if val, exists := tr.Phrases[key]; exists {
		switch val.(type) {
		case string:
			str := val.(string)
			return []string{str}
		case []string:
			return val.([]string)
		}
	}
	return []string{}
}

// GetStringAndReplace gets a translated string and replaces given keys with given values.
// Place key in {brackets} to be replaced here!
func (tr *Translator) GetStringAndReplace(key string, replace R) string {
	str := tr.GetString(key)

	for k, v := range replace {
		str = strings.Replace(str, "{"+k+"}", tr.Format("%v", v), -1)
	}

	return str
}

// GetStringWithVariables gets a translated string, where all placeholders are filled with the values from the given struct.
// Place key in {brackets} to be replaced here! As tag name use alexa.
func (tr *Translator) GetStringWithVariables(key string, data interface{}) string {
	str := tr.GetString(key)

	tr.setFields(&str, data, "")

	return str
}

func (tr *Translator) setFields(str *string, data interface{}, prefix string) {
	t, v := reflect.TypeOf(data), reflect.ValueOf(data)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		tf, vf := t.Field(i), v.Field(i)

		name := tf.Name
		if tag := tf.Tag.Get("alexa"); tag != "" {
			name = tag
		}
		if prefix != "" {
			name = prefix + "." + name
		}

		if vf.Kind() == reflect.Struct || (vf.Kind() == reflect.Ptr && vf.Elem().Kind() == reflect.Struct) {
			tr.setFields(str, vf.Interface(), name)
		} else {
			var value string
			if tf.Type.String() == "string" {
				value = vf.String()
			} else if strings.HasPrefix(tf.Type.String(), "int") {
				value = tr.Format("%d", vf.Int())
			} else if strings.HasPrefix(tf.Type.String(), "float") {
				value = tr.Format("%v", vf.Float())
			}
			*str = strings.Replace(*str, "{"+name+"}", value, -1)
		}
	}
}
