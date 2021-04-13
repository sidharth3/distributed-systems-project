package structs

import (
	"strings"
)

type DirTree struct {
	parent   string
	children []*DirTree
}

func InitDir(parent string) *DirTree {
	return &DirTree{parent, make([]*DirTree, 0)}
}

func (d *DirTree) Insert(pathString string) {
	pathList := strings.Split(pathString, "/")
	d.insert(pathList[1:])
}

func (d *DirTree) insert(fullPath []string) {
	if len(fullPath) == 0 {
		return
	}
	exists := false
	for _, treeItem := range d.children {
		if treeItem.parent == fullPath[0] {
			treeItem.insert(fullPath[1:])
			exists = true
		}
	}
	if !exists {
		newChild := &DirTree{fullPath[0], make([]*DirTree, 0)}
		newChild.insert(fullPath[1:])
		d.children = append(d.children, newChild)
	}
}

func (d *DirTree) FormatString() string {
	dir := d.parent + "\n"
	dir = formatString(dir, d.children, 1)
	return dir
}

func formatString(dir string, children []*DirTree, count int) string {
	if len(children) == 0 {
		return dir
	}
	next_count := count + 1
	for _, treeItem := range children {
		dir += strings.Repeat("-", count) + treeItem.parent + "\n"
		dir = formatString(dir, treeItem.children, next_count)
	}
	return dir
}
