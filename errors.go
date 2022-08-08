package Mailer

import "errors"

var ErrEmptySubject = errors.New("empty mail subject")
var ErrEmptyBody = errors.New("empty mail body")
var ErrNoRecipients = errors.New("no recipients")
var ErrNoRemitent = errors.New("no remitent")
