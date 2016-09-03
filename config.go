/***
  This file is part of config.

  Copyright (c) 2015 Peter Sztan <sztanpet@gmail.com>

  config is free software; you can redistribute it and/or modify it
  under the terms of the GNU Lesser General Public License as published by
  the Free Software Foundation; either version 3 of the License, or
  (at your option) any later version.

  config is distributed in the hope that it will be useful, but
  WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
  Lesser General Public License for more details.

  You should have received a copy of the GNU Lesser General Public License
  along with config; If not, see <http://www.gnu.org/licenses/>.
***/

package config

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/influxdata/toml"
)

func Init(config interface{}, example, path string) error {

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	// empty? initialize it
	if info, err := f.Stat(); err == nil && info.Size() == 0 {
		io.WriteString(f, example)
		f.Seek(0, 0)
	}

	if err := ReadConfig(f, config); err != nil {
		return err
	}

	return nil
}

func ReadConfig(r io.Reader, d interface{}) error {
	dec := toml.NewDecoder(r)
	return dec.Decode(d)
}

func WriteConfig(w io.Writer, d interface{}) error {
	enc := toml.NewEncoder(w)
	return enc.Encode(d)
}

func Save(config interface{}, path string) error {
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile(dir, "tmpconf-")
	if err != nil {
		return err
	}

	err = WriteConfig(f, config)
	if err != nil {
		return err
	}
	_ = f.Close()

	return os.Rename(f.Name(), path)
}
