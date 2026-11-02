package service

import (
	"os"
	"os/exec"
	"strings"
)

// Platform is the manager the app uses when the caller does not inject one.
// Install never runs in CI.
func Platform() Manager {
	return &platform{run: defaultRun}
}

// PlatformWith uses run instead of exec. Tests pass a recorder.
func PlatformWith(run CommandRunner) Manager {
	if run == nil {
		run = defaultRun
	}
	return &platform{run: run}
}

type platform struct {
	run CommandRunner
}

func (p *platform) Status() (Status, error) {
	if p == nil || Blocked() {
		return Status{}, nil
	}
	_, err := os.Stat(PlistPath)
	if err != nil {
		return Status{}, nil
	}
	return Status{Installed: true}, nil
}

func (p *platform) Install() error {
	if Blocked() {
		return ErrCI
	}
	return install(p.run)
}

func (p *platform) Uninstall() error {
	if Blocked() {
		return ErrCI
	}
	return uninstall(p.run)
}

func (p *platform) StartCore(dir string) error {
	if Blocked() {
		return ErrCI
	}
	st, err := p.Status()
	if err != nil {
		return err
	}
	if !st.Installed {
		return ErrUnsupported
	}
	token, err := os.ReadFile(userTokenPath())
	if err != nil {
		return err
	}
	_, err = Call(SocketPath, strings.TrimSpace(string(token)), Request{Op: "start", Dir: dir})
	return err
}

func defaultRun(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
