package template

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/adohe/kube2haproxy/proxy"
	"github.com/golang/glog"
)

type TemplateData struct {
	IPs        map[string]bool
	RouteTable map[string]*proxy.ServiceUnit
	Rservice   *proxy.ServiceUnit
}

func hasIP(ipsMap map[string]bool, ip string) bool {
	glog.V(1).Infof("hasIP:%s",ip)
	return true	
	return ipsMap[ip]
}

// Returns string content of a rendered template
func RenderTemplate(templateName, templateContent string, data interface{}) ([]byte, error) {
	tpl := template.Must(template.New(templateName).Parse(templateContent))

	strBuffer := new(bytes.Buffer)

	err := tpl.Execute(strBuffer, data)
	if err != nil {
		return nil, err
	}

	return strBuffer.Bytes(), nil
}

// Returns string content of a rendered template (with template functions)
func RenderTemplateWithFuncs(templateName, templateContent string, data interface{}) ([]byte, error) {
	funcMap := template.FuncMap{
		"hasIP":   hasIP,
		"ToLower": strings.ToLower,
	}

	tpl := template.Must(template.New(templateName).Funcs(funcMap).Parse(templateContent))

	strBuffer := new(bytes.Buffer)

	err := tpl.Execute(strBuffer, data)
	if err != nil {
		return nil, err
	}

	return strBuffer.Bytes(), nil
}
