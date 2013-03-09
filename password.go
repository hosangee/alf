package main

// Password type required by the ssh packages
type password string

func (p password) Password(user string) (string, error) {
	return string(p), nil
}
