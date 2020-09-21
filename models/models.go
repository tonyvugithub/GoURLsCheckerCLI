package models

//LinkStatus "class"
type LinkStatus struct {
	url  string
	live bool
}

//SetURL ..
func (pt *LinkStatus) SetURL(url string) {
	pt.url = url
}

//SetLiveStatus ...
func (pt *LinkStatus) SetLiveStatus(status bool) {
	pt.live = status
}

//GetURL ..
func (pt *LinkStatus) GetURL() string {
	return pt.url
}

//GetLiveStatus ..
func (pt *LinkStatus) GetLiveStatus() bool {
	return pt.live
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
