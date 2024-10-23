package myerrors

import "errors"

var EmailCreateError = errors.New("email já existente")

func GetCreateCustomerErrors(e error) error {

	if e == nil {
		return nil
	}

	if e.Error() == "UNIQUE constraint failed: customers.desc_email; invalid transaction" {
		return EmailCreateError
	}

	return e
}
