package core

import (
	"database/sql"
	"fmt"

	"github.com/limerence-code/goproject/gee/orm/dialect"
	"github.com/limerence-code/goproject/gee/orm/log"
	"github.com/limerence-code/goproject/gee/orm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	dail, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return nil, fmt.Errorf("dialect %s Not Found", driver)
	}

	engine := &Engine{db: db, dialect: dail}
	log.Info("connection db success")
	return engine, nil
}

func (engine *Engine) Close() error {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close connection db err", err)
		return err
	}
	log.Info("close connection db success")
	return nil
}

func (engine *Engine) NewSession() (*session.Session, error) {
	return session.New(engine.db, engine.dialect), nil
}

//----------------------事务----------------------
type TxFunc func(s *session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (res interface{}, err error) {
	s, _ := engine.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			s.RollBack()
			panic(p)
		} else if err != nil {
			s.RollBack()
		} else {
			s.Commit()
		}
	}()
	return f(s)
}
