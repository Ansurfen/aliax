// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package template

import (
	"aliax/internal/aio"
	"html/template"
	"io"
	"os"

	"github.com/caarlos0/log"
)

func Execute(w io.Writer, s string, data map[string]string) error {
	tmpl, err := template.New("").Funcs(aliaxFuncs).Parse(s)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

var aliaxFuncs = map[string]any{
	"env": os.Getenv,
	"aliax_env": func(key string) string {
		if env == nil {
			log.Info("parsing aliax.env.yaml")
			err := aio.ReadYAML("aliax.env.yaml", &env)
			if err != nil {
				log.WithError(err).Fatal("fail to parse env")
			}
		}
		return env[key]
	},
}

var env map[string]string
