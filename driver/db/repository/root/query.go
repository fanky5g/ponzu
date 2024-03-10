package root

import (
	"bytes"
	"github.com/boltdb/bolt"
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"github.com/fanky5g/ponzu/internal/domain/enum"
)

// FindAllWithOptions retrieves a set of content from the db based on options
// and returns the total number of content in the namespace and the content
func (repo *repository) FindAllWithOptions(namespace string,
	order enum.SortOrder,
	pagination *entities.Pagination) (int, []interface{}, error) {
	var posts []interface{}
	var total int

	// correct bad input rather than return nil or error
	// similar to default case for opts.Order switch below
	if pagination.Count < 0 {
		pagination.Count = -1
	}

	if pagination.Offset < 0 {
		pagination.Offset = 0
	}

	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(namespace))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		c := b.Cursor()
		n := b.Stats().KeyN
		total = n

		// return nil if no content
		if n == 0 {
			return nil
		}

		var start, end int
		switch pagination.Count {
		case -1:
			start = 0
			end = n

		default:
			start = pagination.Count * pagination.Offset
			end = start + pagination.Count
		}

		// bounds check on posts given the start & end count
		if start > n {
			start = n - pagination.Count
		}
		if end > n {
			end = n
		}

		i := 0   // count of num posts added
		cur := 0 // count of num cursor moves
		switch order {
		case "desc", "":
			for k, v := c.Last(); k != nil; k, v = c.Prev() {
				if cur < start {
					cur++
					continue
				}

				if cur >= end {
					break
				}

				entity, err := repo.MarshalEntity(namespace, bytes.NewBuffer(v))
				if err != nil {
					return err
				}

				posts = append(posts, entity)
				i++
				cur++
			}

		case "asc":
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if cur < start {
					cur++
					continue
				}

				if cur >= end {
					break
				}

				entity, err := repo.MarshalEntity(namespace, bytes.NewBuffer(v))
				if err != nil {
					return err
				}

				posts = append(posts, entity)
				i++
				cur++
			}

		default:
			// results for DESC order
			for k, v := c.Last(); k != nil; k, v = c.Prev() {
				if cur < start {
					cur++
					continue
				}

				if cur >= end {
					break
				}

				entity, err := repo.MarshalEntity(namespace, bytes.NewBuffer(v))
				if err != nil {
					return err
				}

				posts = append(posts, entity)
				i++
				cur++
			}
		}

		return nil
	})

	if err != nil {
		return 0, nil, err
	}

	return total, posts, nil
}
