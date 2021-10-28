/*
 * File: formaterror.go
 * Project: formaterror
 * File Created: Thursday, 28th October 2021 10:05:02 am
 * Author: Ananda Yudhistira (anandabayu12@gmail.com)
 * -----
 * Last Modified: Thursday, 28th October 2021 10:05:30 am
 * Modified By: Ananda Yudhistira (anandabayu12@gmail.com>)
 * -----
 * Copyright 2021 Ananda Yudhistira, -
 */

package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "nickname") {
		return errors.New("Nickname Already Taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}

	if strings.Contains(err, "title") {
		return errors.New("Title Already Taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}

	return errors.New("Incorrect Details")
}
