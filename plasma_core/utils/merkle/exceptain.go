package merkle

type MemberNotExistException struct {
	msg string
}

func (e *MemberNotExistException)Error() string{
	return e.msg
}
