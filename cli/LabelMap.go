package cli

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type LabelMap struct {
	labels  map[string]int
	reverse map[int]string
}

func NewLabelMap() *LabelMap {
	return &LabelMap{
		labels:  map[string]int{},
		reverse: map[int]string{},
	}
}

var reName = regexp.MustCompile(`[a-zA-Z0-9]+`)

func (m *LabelMap) Append(labelId int, name string) {
	if _, exists := m.reverse[labelId]; exists {
		return
	}
	name = strings.Join(reName.FindAllString(name, -1), "_")
	if len(name) == 0 || !unicode.IsLetter([]rune(name)[0]) {
		name = "_" + name
	}
	nn := name
	i := 1
	_, exists := m.labels[nn]
	for exists {
		nn = fmt.Sprintf("%s_%d", name, i)
		_, exists = m.labels[nn]
	}
	m.labels[nn] = labelId
	m.reverse[labelId] = nn
}

func (m *LabelMap) Labels() map[string]int {
	return m.labels
}

func (m *LabelMap) GetName(labelId int) string {
	if name, ok := m.reverse[labelId]; ok {
		return name
	}
	return fmt.Sprintf("%d", labelId)
}

func (m *LabelMap) GetId(name string) int {
	if labalId, ok := m.labels[name]; ok {
		return labalId
	}
	return 0
}
