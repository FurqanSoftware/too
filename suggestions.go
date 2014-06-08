// Copyright 2014 The Too Authors. All rights reserved.

package too

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type Suggestions struct {
	e *Engine
}

// For returns n suggested items for user
func (s Suggestions) For(user User, n int) ([]Item, error) {
	results, err := redis.Strings(s.e.c.Do("ZREVRANGE", fmt.Sprintf("%s:%s:%s", s.e.class, user, "suggestions"), 0, n-1))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	items := []Item{}
	for _, item := range results {
		items = append(items, Item(item))
	}
	return items, nil
}

func (s Suggestions) update(user User) error {
	similars, err := s.e.Similars.Of(user, 8)
	if err != nil {
		return err
	}
	if len(similars) == 0 {
		return nil
	}

	args := []interface{}{fmt.Sprintf("%s:%s:_tmp", s.e.class, user)}
	for _, similar := range similars {
		args = append(args, fmt.Sprintf("%s:%s:%s", s.e.class, similar, s.e.Likes.kind))
	}
	_, err = s.e.c.Do("SUNIONSTORE", args...)
	if err != nil {
		return err
	}
	defer s.e.c.Do("DEL", fmt.Sprintf("%s:%s:_tmp", s.e.class, user))

	items, err := redis.Strings(s.e.c.Do("SDIFF", fmt.Sprintf("%s:%s:_tmp", s.e.class, user), fmt.Sprintf("%s:%s:%s", s.e.class, user, s.e.Likes.kind), fmt.Sprintf("%s:%s:%s", s.e.class, user, s.e.Dislikes.kind)))
	if err != nil && err != redis.ErrNil {
		return err
	}

	scores := map[string]float64{}
	for _, item := range items {
		likers, err := redis.Strings(s.e.c.Do("SMEMBERS", fmt.Sprintf("%s:%s:%s", s.e.class, item, s.e.Likes.kind)))
		if err != nil && err != redis.ErrNil {
			return err
		}

		for _, liker := range likers {
			score, err := redis.Float64(s.e.c.Do("ZSCORE", fmt.Sprintf("%s:%s:similars", s.e.class, user), liker))
			if err != nil && err != redis.ErrNil {
				return err
			}

			scores[item] += score
		}

		dislikers, err := redis.Strings(s.e.c.Do("SMEMBERS", fmt.Sprintf("%s:%s:%s", s.e.class, item, s.e.Likes.kind)))
		if err != nil && err != redis.ErrNil {
			return err
		}

		for _, disliker := range dislikers {
			score, err := redis.Float64(s.e.c.Do("ZSCORE", fmt.Sprintf("%s:%s:similars", s.e.class, user), disliker))
			if err != nil && err != redis.ErrNil {
				return err
			}

			scores[item] -= score
		}

		total := len(likers) + len(dislikers)
		if total > 0 {
			scores[item] /= float64(total)
		}
	}

	_, err = s.e.c.Do("DEL", fmt.Sprintf("%s:%s:suggestions", s.e.class, user))
	if err != nil {
		return err
	}

	for _, item := range items {
		s.e.c.Do("ZADD", fmt.Sprintf("%s:%s:suggestions", s.e.class, user), scores[item], item)
	}
	return nil
}
