/*
 * @Author: gaoh
 * @Date: 2023-05-15 23:42:28
 * @LastEditTime: 2023-05-15 23:48:20
 */
package session

import "github.com/limerence-code/goproject/gee/orm/log"

func (s *Session) Begin() (err error) {
	log.Info("transaction Begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transaction Commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) RollBack() (err error) {
	log.Info("transaction RollBack")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
