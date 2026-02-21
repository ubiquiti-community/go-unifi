package unifi

import "fmt"

type LoginRequiredError struct{}

func (err *LoginRequiredError) Error() string {
	return "login required"
}

type NotFoundError struct {
	Type  string
	Attr  string
	Value string
}

func (err *NotFoundError) Error() string {
	if err.Attr != "" && err.Value != "" {
		return fmt.Sprintf("not found: type=%s, attr=%s, value=%s", err.Type, err.Attr, err.Value)
	} else {
		return fmt.Sprintf("not found: type=%s", err.Type)
	}
}

type APIError struct {
	RC      string
	Message string
}

func (err *APIError) Error() string {
	return err.Message
}
