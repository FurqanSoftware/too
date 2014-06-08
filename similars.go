// Copyright 2014 The Too Authors. All rights reserved.

package too

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type Similars struct {
	e *Engine
}

// Of returns n users similar to user
func (s Similars) Of(user User, n int) ([]User, error) {
	results, err := redis.Strings(s.e.c.Do("ZREVRANGE", fmt.Sprintf("%s:%s:%s", s.e.class, user, "similars"), 0, n-1))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	users := []User{}
	for _, user := range results {
		users = append(users, User(user))
	}
	return users, nil
}

// Jaccard returns the Jaccard coefficient between user and other
func (s Similars) Jaccard(user, other User) (float64, error) {
	likes, err := redis.Strings(s.e.c.Do("SINTER", fmt.Sprintf("%s:%s:%s", s.e.class, user, s.e.Likes.kind), fmt.Sprintf("%s:%s:%s", s.e.class, other, s.e.Likes.kind)))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}

	dislikes, err := redis.Strings(s.e.c.Do("SINTER", fmt.Sprintf("%s:%s:%s", s.e.class, user, s.e.Dislikes.kind), fmt.Sprintf("%s:%s:%s", s.e.class, other, s.e.Dislikes.kind)))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}

	antiLikes, err := redis.Strings(s.e.c.Do("SINTER", fmt.Sprintf("%s:%s:%s", s.e.class, user, s.e.Likes.kind), fmt.Sprintf("%s:%s:%s", s.e.class, other, s.e.Dislikes.kind)))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}

	antiDislikes, err := redis.Strings(s.e.c.Do("SINTER", fmt.Sprintf("%s:%s:%s", s.e.class, user, s.e.Dislikes.kind), fmt.Sprintf("%s:%s:%s", s.e.class, other, s.e.Likes.kind)))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}

	return float64(len(likes)+len(dislikes)-len(antiLikes)-len(antiDislikes)) / float64(len(likes)+len(dislikes)+len(antiLikes)+len(antiDislikes)), nil
}

func (s Similars) update(user User) error {
	items, err := redis.Strings(s.e.c.Do("SUNION", fmt.Sprintf("%s:%s:%s", s.e.class, user, s.e.Likes.kind), fmt.Sprintf("%s:%s:%s", s.e.class, user, s.e.Dislikes.kind)))
	if err != nil && err != redis.ErrNil {
		return err
	}

	args := []interface{}{}
	for _, item := range items {
		args = append(args, fmt.Sprintf("%s:%s:%s", s.e.class, item, s.e.Likes.kind))
		args = append(args, fmt.Sprintf("%s:%s:%s", s.e.class, item, s.e.Dislikes.kind))
	}
	users, err := redis.Strings(s.e.c.Do("SUNION", args...))
	if err != nil && err != redis.ErrNil {
		return err
	}

	for _, other := range users {
		if other != string(user) {
			v, err := s.Jaccard(user, User(other))
			if err != nil {
				return err
			}

			_, err = s.e.c.Do("ZADD", fmt.Sprintf("%s:%s:similars", s.e.class, user), v, other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
