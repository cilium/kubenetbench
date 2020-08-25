package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// SessionCtx is the context for a session run
type Session struct {
	id  string // id identifies the run
	dir string // directory to store results/etc.
}

// NewRunCtx creates a new RunCtx
func NewSession(
	sessId string,
	sessDirBase string,
) (*Session, error) {

	sess := &Session{
		id:  sessId,
		dir: fmt.Sprintf("%s/%s", sessDirBase, sessId),
	}

	info, err_stat := os.Stat(sess.dir)
	if err_stat == nil && info.IsDir() {
		// directory exists, good to go
		return sess, nil
	} else if os.IsNotExist(err_stat) {
		// otherwise, create directory if it does not exist
		err_mkdir := os.Mkdir(sess.dir, 0755)
		if err_mkdir != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w\n", sess.dir, err_mkdir)
		}
		sess.writeScript(sessId, sessDirBase)
		return sess, nil
	} else {
		return nil, fmt.Errorf("failed to initialize session using directory %d", sess.dir)
	}
}

func InitSession(
	sessId string,
	sessDirBase string,
) (*Session, error) {

	sess := &Session{
		id:  sessId,
		dir: fmt.Sprintf("%s/%s", sessDirBase, sessId),
	}

	info, err_stat := os.Stat(sess.dir)
	if err_stat == nil && info.IsDir() {
		return nil, fmt.Errorf("session directory (%s) already exists\n", sess.dir)
	} else if os.IsNotExist(err_stat) {
		// otherwise, create directory if it does not exist
		err_mkdir := os.Mkdir(sess.dir, 0755)
		if err_mkdir != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w\n", sess.dir, err_mkdir)
		}
		sess.writeScript(sessId, sessDirBase)
		return sess, nil
	} else {
		return nil, fmt.Errorf("failed to initialize session using directory %s", sess.dir)
	}
}

func (s *Session) getSessionLabel() string {
	return fmt.Sprintf("%s=%s", sessIdLabel, s.id)
}

func (s *Session) writeScript(sid, sdbase string) {
	fname := fmt.Sprintf("%s/knb", s.dir)
	f, err := os.Create(fname)
	if err != nil {
		return
	}
	defer f.Close()

	prog, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}

	fmt.Fprintln(f, "#!/bin/sh")
	fmt.Fprintln(f, "# wrapper script for kubenetbench")
	fmt.Fprintf(f, "%s --session-id %s --session-base-dir %s \"$@\"\n", prog, sid, sdbase)

	err = os.Chmod(fname, 0755)
	if err != nil {
		return
	}

	log.Printf("wrote wrapper script: you can now use wrapper %s\n", fname)
}

func (s *Session) OpenLog() (*os.File, error) {
	fname := fmt.Sprintf("%s/log", s.dir)
	return os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}
