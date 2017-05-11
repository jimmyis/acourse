package model

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

// User model
type User struct {
	id        string
	role      *UserRole
	Username  string
	Name      string
	AboutMe   string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRole type
type UserRole struct {
	Admin      bool
	Instructor bool
}

// SetID sets user id
func (x *User) SetID(id string) {
	x.id = id
}

// ID returns user id
func (x *User) ID() string {
	return x.id
}

// Role returns user role
func (x *User) Role() *UserRole {
	if x.role == nil {
		x.role = &UserRole{}
	}
	return x.role
}

// Save saves user
func (x *User) Save(c redis.Conn) error {
	if len(x.id) == 0 {
		return fmt.Errorf("invalid id")
	}

	c.Send("MULTI")
	c.Send("SADD", key("u", "all"), x.id)

	x.UpdatedAt = time.Now()
	if x.CreatedAt.IsZero() {
		x.CreatedAt = x.UpdatedAt
		c.Send("ZADD", key("u", "t0"), x.CreatedAt.UnixNano(), x.id)
	} else {
		c.Send("ZADD", key("u", "t0"), x.CreatedAt.UnixNano(), x.id) // TODO: remove after migrate
	}

	c.Send("ZADD", key("u", "t1"), x.UpdatedAt.UnixNano(), x.id)
	c.Send("HSET", key("u"), x.id, enc(x))

	if x.role != nil {
		if x.role.Admin {
			c.Send("SADD", key("u", "admin"), x.id)
		} else {
			c.Send("SREM", key("u", "admin"), x.id)
		}
		if x.role.Instructor {
			c.Send("SADD", key("u", "instructor"), x.id)
		} else {
			c.Send("SREM", key("u", "instructor"), x.id)
		}
	}

	_, err := c.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// GetUsers gets users
func GetUsers(c redis.Conn, userIDs []string) ([]*User, error) {
	xs := make([]*User, len(userIDs))
	for _, userID := range userIDs {
		c.Send("SISMEMBER", key("u", "all"), userID)
		c.Send("HGET", key("u"), userID)
		c.Send("SISMEMBER", key("u", "admin"), userID)
		c.Send("SISMEMBER", key("u", "instructor"), userID)
	}
	c.Flush()
	for i, userID := range userIDs {
		exists, _ := redis.Bool(c.Receive())
		if !exists {
			c.Receive()
			c.Receive()
			c.Receive()
			continue
		}
		var x User
		b, err := redis.Bytes(c.Receive())
		if err != nil {
			return nil, err
		}
		err = dec(b, &x)
		if err != nil {
			return nil, err
		}
		x.role = &UserRole{}
		x.role.Admin, _ = redis.Bool(c.Receive())
		x.role.Instructor, _ = redis.Bool(c.Receive())
		x.id = userID
		xs[i] = &x
	}
	return xs, nil
}

// GetUser gets user from id
func GetUser(c redis.Conn, userID string) (*User, error) {
	xs, err := GetUsers(c, []string{userID})
	if err != nil {
		return nil, err
	}
	return xs[0], nil
}

// GetUserFromUsername gets user from username
func GetUserFromUsername(c redis.Conn, username string) (*User, error) {
	userID, err := redis.String(c.Do("HGET", key("u", "username"), username))
	if err == redis.ErrNil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return GetUser(c, userID)
}

// ListUsers lists users
// TODO: pagination
func ListUsers(c redis.Conn) ([]*User, error) {
	userIDs, err := redis.Strings(c.Do("ZREVRANGE", key("u", "t0"), 0, -1))
	if err != nil {
		return nil, err
	}
	return GetUsers(c, userIDs)
}