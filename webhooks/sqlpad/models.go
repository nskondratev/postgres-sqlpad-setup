package sqlpad

import "time"

// ACL ...
type ACL struct {
	ID        int       `json:"id"`
	QueryID   string    `json:"queryId"`
	UserID    *string   `json:"userId"`
	GroupID   string    `json:"groupId"`
	Write     bool      `json:"write"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Query ...
type Query struct {
	ID   string `json:"id"`
	ACL  []ACL  `json:"acl"`
	User struct {
		Email string `json:"email"`
	} `json:"createdByUser"`
}

// User ...
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"createdAt"`
}
