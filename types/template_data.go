package types

// TODO: make this better to house more things as we need them, like functions
type TemplateData struct {
	PageTitle string
	BoolMap   map[string]bool
	StringMap map[string]string
	IntMap    map[string]int
	ObjectMap map[string]interface{}
}
