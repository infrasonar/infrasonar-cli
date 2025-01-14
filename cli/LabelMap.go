package cli

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type LabelMap struct {
	labels  map[string]*Label
	reverse map[int]string
}

func NewLabelMap() *LabelMap {
	return &LabelMap{
		labels:  map[string]*Label{},
		reverse: map[int]string{},
	}
}

var reName = regexp.MustCompile(`[a-zA-Z0-9]+`)

func (m *LabelMap) Append(label *Label) {
	if _, exists := m.reverse[label.Id]; exists {
		return
	}
	name := strings.Join(reName.FindAllString(label.Name, -1), "_")
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
	m.labels[nn] = label
	m.reverse[label.Id] = nn
}

func (m *LabelMap) Labels() map[string]*Label {
	return m.labels
}

func (m *LabelMap) LabelById(labelId int) *Label {
	if key, ok := m.reverse[labelId]; ok {
		if label, ok := m.labels[key]; ok {
			return label
		}
	}
	return nil
}

func (m *LabelMap) LabelByKey(key string) *Label {
	if label, ok := m.labels[key]; ok {
		return label
	}
	return nil
}

func (m *LabelMap) GetName(labelId int) string {
	if name, ok := m.reverse[labelId]; ok {
		return name
	}
	return fmt.Sprintf("%d", labelId)
}
