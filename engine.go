// Copyright 2014 The Too Authors. All rights reserved.

package too

import "github.com/garyburd/redigo/redis"

type Engine struct {
	c     redis.Conn
	class string

	Likes    Rater
	Dislikes Rater

	Similars    Similars
	Suggestions Suggestions
}

// New returns a new engine for given class connected to Redis server at addr.
func New(url, class string) (*Engine, error) {
	c, err := redis.DialURL(url)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		c:     c,
		class: class,
	}
	e.Likes = Rater{e, "likes"}
	e.Dislikes = Rater{e, "dislikes"}
	e.Similars = Similars{e}
	e.Suggestions = Suggestions{e}
	return e, nil
}
