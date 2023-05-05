package session

import (
	"testing"

	"github.com/limerence-code/goproject/gee/orm/log"
)

type Account struct {
	Id       int64 `geeorm:"Primary Key"`
	Password string
}

func (a *Account) BeforeInsert(s *Session) error {
	log.Info("Before insert", a)
	a.Id += 1000
	return nil
}

func (a *Account) AfterQuery(s *Session) error {
	log.Info("Afert query", a)
	a.Password = "******"
	return nil
}

func TestHookCallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	s.DropTable()
	s.CreateTable()
	s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})
	account := &Account{}
	err := s.First(account)
	if err != nil || account.Id != 1001 || account.Password != "******" {
		t.Fatal("Failed to call hooks after query, got", account)
	}
}
