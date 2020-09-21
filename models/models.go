package models

//LinkStatus "class"
type LinkStatus struct {
	Url  string
	Live bool
}

//Set up link count info "class"
type linkCountInfo struct {
	upLink    int
	downLink  int
	totalLink int
}

func (pt *linkCountInfo) incrementUpLink() {
	(*pt).upLink++
}

func (pt *linkCountInfo) incrementDownLink() {
	(*pt).downLink++
}

func (pt *linkCountInfo) setTotalLink(total int) {
	(*pt).totalLink = total
}
