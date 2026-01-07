package service

import "context"

func (s *svc) getOffset(page int) (offset int) {
	offset = (page - 1) * s.limit
	return
}

func (s *svc) HasNextPage(ctx context.Context, total, page int) bool {
	return total > (page+1)*s.limit
}
