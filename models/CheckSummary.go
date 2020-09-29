package models

//CheckSummary ...
type CheckSummary struct {
	upLinks   []string
	downLinks []string
}

//RecordUpLink ...
func (cs *CheckSummary) RecordUpLink(link string) {
	cs.upLinks = append(cs.upLinks, link)
}

//RecordDownLink ...
func (cs *CheckSummary) RecordDownLink(link string) {
	cs.downLinks = append(cs.downLinks, link)
}

//GetUpLinks ...
func (cs *CheckSummary) GetUpLinks() []string {
	return cs.upLinks
}

//GetDownLinks ...
func (cs *CheckSummary) GetDownLinks() []string {
	return cs.downLinks
}

//GetNumUpLinks ...
func (cs *CheckSummary) GetNumUpLinks() int {
	return len(cs.upLinks)
}

//GetNumDownLinks ...
func (cs *CheckSummary) GetNumDownLinks() int {
	return len(cs.downLinks)
}
