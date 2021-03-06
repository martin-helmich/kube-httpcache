package controller

import (
	"github.com/martin-helmich/kube-httpcache/watcher"
	"io"
	"io/ioutil"
	"text/template"
)

type TemplateData struct {
	Backends       watcher.BackendList
	PrimaryBackend *watcher.Backend
}

type VarnishController struct {
	SecretFile   string
	Storage      string
	FrontendAddr string
	FrontendPort int
	AdminAddr    string
	AdminPort    int

	vclTemplate    *template.Template
	updates        chan *watcher.BackendConfig
	configFile     string
	secret         []byte
	localAdminAddr string
}

func NewVarnishController(
	secretFile string,
	storage string,
	frontendAddr string,
	frontendPort int,
	adminAddr string,
	adminPort int,
	updates chan *watcher.BackendConfig,
	vclTemplateFile string,
) (*VarnishController, error) {
	contents, err := ioutil.ReadFile(vclTemplateFile)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("vcl").Parse(string(contents))
	if err != nil {
		return nil, err
	}

	secret, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return nil, err
	}

	return &VarnishController{
		SecretFile:   secretFile,
		Storage:      storage,
		FrontendAddr: frontendAddr,
		FrontendPort: frontendPort,
		AdminAddr:    adminAddr,
		AdminPort:    adminPort,
		vclTemplate:  tmpl,
		updates:      updates,
		configFile:   "/tmp/vcl",
		secret:       secret,
	}, nil
}

func (v *VarnishController) renderVCL(target io.Writer, backendList watcher.BackendList, primary *watcher.Backend) (error) {
	err := v.vclTemplate.Execute(target, &TemplateData{
		Backends:       backendList,
		PrimaryBackend: primary,
	})

	return err
}
