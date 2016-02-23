package main

import (
	"io/ioutil"
	"log"
)

const DataFileName = "data.json"
const InitalData = `[]`

type Data struct {
	jsonData  string
	Resources []map[string]string
}

func (d *Data) saveFile() error {
	err := ioutil.WriteFile(DataFileName, []byte(d.jsonData), 0600)
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) getFileData() (string, error) {
	filename := DataFileName
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("data.json could not be read. Creating. Err: %s", err.Error())
		d.jsonData = InitalData

		err = d.saveFile()
		if err != nil {
			log.Panicf("data.json could not be read and unable to create new. Err: %s", err.Error())
		}

		return d.jsonData, nil
	}

	return string(b), nil
}
