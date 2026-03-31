package unifi

import "context"

func Ptr[T any](v T) *T {
	return &v
}

// 1. Define an interface for the "Page" structure.
// This ensures the generic function knows the Page has a TotalCount.
type Paginatable[T any] interface {
	GetTotalCount() int64
	GetData() []T
}

// 2. The Universal Paginator
func FetchAll[T any, P Paginatable[T]](
	ctx context.Context,
	fetcher func(offset int32) (P, error),
) ([]T, error) {
	var allItems []T
	var currentOffset int32 = 0
	const pageSize int32 = 50

	for {
		page, err := fetcher(currentOffset)
		if err != nil {
			return nil, err
		}

		items := page.GetData()
		allItems = append(allItems, items...)

		// Exit if we have all items or the page is empty
		if int64(len(allItems)) >= page.GetTotalCount() || len(items) == 0 {
			break
		}

		currentOffset += pageSize
	}

	return allItems, nil
}
