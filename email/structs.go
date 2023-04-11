package email

type Templates_Dir struct {
	Dir string
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}
