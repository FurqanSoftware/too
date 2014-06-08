// Copyright 2014 The Too Authors. All rights reserved.

package too

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type Rater struct {
	e    *Engine
	kind string
}

// Add adds a rating by user for item
func (r Rater) Add(user User, item Item) error {
	yes, err := redis.Bool(r.e.c.Do("SISMEMBER", fmt.Sprintf("%s:%s:%s", r.e.class, item, r.kind), user))
	if err != nil {
		return err
	}

	if !yes {
		_, err = r.e.c.Do("ZINCRBY", fmt.Sprintf("%s:mosts:%s", r.e.class, r.kind), 1, item)
		if err != nil {
			return err
		}
	}

	_, err = r.e.c.Do("SADD", fmt.Sprintf("%s:%s:%s", r.e.class, user, r.kind), item)
	if err != nil {
		return err
	}

	_, err = r.e.c.Do("SADD", fmt.Sprintf("%s:%s:%s", r.e.class, item, r.kind), user)
	if err != nil {
		return err
	}

	err = r.e.Similars.update(user)
	if err != nil {
		return err
	}

	err = r.e.Suggestions.update(user)
	if err != nil {
		return err
	}

	return nil
}
