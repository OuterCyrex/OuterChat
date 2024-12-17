package SError

const (
	RePasswordError = iota + 101
	InValidIdError
	InValidEmailError
	NameHasBeenUsedError
	NameNotExistError
	PasswordWrongError
)

const (
	AlreadyFriendError = iota + 201
	AlreadySendRequestError
	FriendWithSelfError
	NotEvenFriendError
)

const (
	IntervalError = iota + 501
	InvalidTokenError
)
