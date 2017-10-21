package haproxy

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/adohe/kube2haproxy/util/config"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
)

// Configuration object for constructing Haproxy.
type HaproxyConfig struct {
	ConfigPath       string
	TemplatePath     string
	ReloadScriptPath string
	ConfiguePath     string
	PidPath          string
	ExecFile         string
	ReloadInterval   time.Duration
}

// Haproxy represents a real HAProxy instance.
type Haproxy struct {
	config     HaproxyConfig
	configurer *config.Configurer
}

func NewInstance(cfg HaproxyConfig) (*Haproxy, error) {
	configurer, err := config.NewConfigurer(cfg.ConfigPath)
	if err != nil {
		return nil, err
	}
	return &Haproxy{
		config:     cfg,
		configurer: configurer,
	}, nil
}

func (p *Haproxy) Reload(cfgBytes []byte) error {
	// First update haproxy.cfg
	err := p.configurer.WriteConfig(cfgBytes)
	if err != nil {
		return err
	}

	start := time.Now()
	defer func() {
		glog.V(4).Infof("ReloadHaproxy took %v", time.Since(start))
	}()

	// Reload
	glog.V(4).Infof("Ready to reload haproxy")
	cmd := exec.Command(p.config.ReloadScriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error reloading haproxy: %v\n%s", err, out)
	}
	return nil
}

func (p *Haproxy) Reload2(filename string, cfgBytes []byte) error {
	glog.V(1).Infof("Reload2:filename:%s,ConfiguePath:%s", filename, p.config.ConfiguePath)
	tmpConfigFile := fmt.Sprintf("%s/%s.cfg", p.config.ConfiguePath, filename)
	glog.V(1).Infof("tmpConfigFile:%s", tmpConfigFile)
	err := ioutil.WriteFile(tmpConfigFile, cfgBytes, os.FileMode(0644))
	if err != nil {
		glog.Errorf("failed to write tmp config file: %v", err)
		return err
	}

	start := time.Now()
	defer func() {
		glog.V(4).Infof("ReloadHaproxy took %v", time.Since(start))
	}()

	pidfile := fmt.Sprintf("%s/%s.pid", p.config.PidPath, filename)
	// Reload
	glog.V(4).Infof("Ready to reload haproxy pidfile:%s", pidfile)
	cmdline := fmt.Sprintf("%s -f %s -p %s -D -sf $(cat %s)", p.config.ExecFile, tmpConfigFile, pidfile, pidfile)
	glog.V(1).Infof("cmdline:%s", cmdline)
	cmd := exec.Command("sh", "-c", cmdline)
	//cmdline := fmt.Sprintf("%s/%s.pid", p.config.PidPath, filename)
	//cmd := exec.Command(p.config.ExecFile, "-f", tmpConfigFile, "-p", cmdline, "-D", "-sf")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error reloading haproxy: %v\n%s", err, out)
	}
	return nil
}
