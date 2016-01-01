// Copyright 2014 The Too Authors. All rights reserved.

package too

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type Rater struct {
	e           *Engine
	kind        string
	memberships map[User][]Item
}

type BatchRaterOp struct {
	User  User
	Items []Item
}

func (r Rater) Batch(ops []BatchRaterOp, updateSimilarsAndSuggestions bool) error {
	// Disable Auto Update, Suggestions and Similiars will be updated later
	r.e.DisableAutoUpdateSimilarsAndSuggestions()

	// Cache memberships to be used into redis transaction
	r.memberships = make(map[User][]Item, 0)
	r.cacheMemberships(ops)

	// Start a transaction
	r.e.c.Send("MULTI")
	for _, op := range ops {
		for _, item := range op.Items {
			err := r.Add(op.User, item)
			if err != nil {
				// Rollback if found  error
				r.e.c.Send("DISCARD")
				return err
			}
		}
	}
	// Commit the transaction
	r.e.c.Do("EXEC")

	// After finished, update Suggestions and Similiars
	if updateSimilarsAndSuggestions {
		for _, op := range ops {
			r.e.Update(op.User)
		}
	}

	r.e.EnableAutoUpdateSimilarsAndSuggestions()

	return nil
}

func (r Rater) cacheMemberships(ops []BatchRaterOp) {
	for _, op := range ops {
		for _, item := range op.Items {
			r.e.c.Send("WATCH", r.memberKey(item))
			yes, err := r.userIsMember(op.User, item, false)
			if yes && err != nil {
				r.memberships[op.User] = append(r.memberships[op.User], item)
			}
		}
	}
}

func (r Rater) isCachedInMembership(user User, item Item) bool {
	for _, _item := range r.memberships[user] {
		if _item == item {
			return true
		}
	}
	return false
}

func (r Rater) userIsMember(user User, item Item, useCache bool) (bool, error) {
	yes, err := redis.Bool(r.e.c.Do("SISMEMBER", r.memberKey(item), user))
	if err != nil && useCache {
		return r.isCachedInMembership(user, item), nil
	}
	return yes, err
}

func (r Rater) memberKey(item Item) string {
	return fmt.Sprintf("%s:%s:%s", r.e.class, item, r.kind)
}

// Add adds a rating by user for item
func (r Rater) Add(user User, item Item) error {
	yes, err := r.userIsMember(user, item, true)
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

	if r.e.autoUpdateSimilarsAndSuggestions {
		err = r.e.Update(user)
	}

	if err != nil {
		return err
	}

	return nil
}
