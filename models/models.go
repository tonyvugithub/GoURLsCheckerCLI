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
