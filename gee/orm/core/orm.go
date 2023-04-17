package core

import (
	"database/sql"

	"github.com/limerence-code/goproject/gee/orm/log"
	"github.com/limerence-code/goproject/gee/orm/session"
)

type Engine struct {
	db *sql.DB
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
	engine := &Engine{db: db}
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
	// return sessio.NewSession(engine.db)
	return session.New(engine.db), nil
}
