package main

type ChoiceSource interface {
	download() (string, error)
	install() error
}

type ChoiceGithub struct {
}

func (c ChoiceGithub) download() (string, error) {
	return "hogehoge", nil
}
