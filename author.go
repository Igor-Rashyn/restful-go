package main

// Author model
type Author struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"first_name,omitempty"`
	Lastname  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}
