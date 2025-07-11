package messagesrepo

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

var _ gomock.Matcher = CursorMatcher{}

// CursorMatcher is intended to be used only in tests.
type CursorMatcher struct {
	c Cursor
}

func NewCursorMatcher(c Cursor) CursorMatcher {
	return CursorMatcher{c: c}
}

func (cm CursorMatcher) Matches(x any) bool {
	v, ok := x.(*Cursor)
	if !ok {
		return false
	}

	// Фикс: Два курсора равны, если равны их PageSize и LastCreatedAt
	if cm.c.PageSize != v.PageSize {
		return false
	}
	if cm.c.LastCreatedAt.UnixNano() != v.LastCreatedAt.UnixNano() {
		return false
	}
	return true
}

func (cm CursorMatcher) String() string {
	return fmt.Sprintf("{ps=%d, last_created_at=%d}", cm.c.PageSize, cm.c.LastCreatedAt.UnixNano())
}
