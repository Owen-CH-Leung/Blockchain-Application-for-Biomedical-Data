package image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"upload-view/backend/logger"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"time"
)

type Metadata_jpg2json struct {
	FileName         string
	Size             int64
	ModificationTime time.Time
	Image            []byte
}

type Metadata_json2jpg struct {
	FileName         string
	Size             int64
	ModificationTime time.Time
	Image            string
}

// OrgFile must end with ".jpg", DestFile must end with ".json"
func JpgToMetaData(OrgFile string) Metadata_jpg2json {
	var log = logger.BlockchainLog

	f_handler, err := os.Open(OrgFile)
	if err != nil {
		log.Println("Opening Image encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	defer f_handler.Close()
	fileInfo, err := f_handler.Stat()
	if err != nil {
		log.Println("Getting File stats encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	size := fileInfo.Size()
	filename := fileInfo.Name()
	modtime := fileInfo.ModTime()
	log.Println("size : ", size)
	log.Println("File Name : ", filename)
	log.Println("Modification Time : ", modtime)

	img_bytes, err := ioutil.ReadAll(f_handler)
	if err != nil {
		log.Println("Reading File as Byte encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Byte Array Size : ", len(img_bytes))
	//log.Println("Image Byte :", img_bytes)

	return Metadata_jpg2json{
		FileName:         filename,
		Size:             size,
		ModificationTime: modtime,
		Image:            img_bytes,
	}
}

func MetaDataToJson(json_mapping Metadata_jpg2json, DestFile string) {
	var log = logger.BlockchainLog

	jsonByte, err := json.MarshalIndent(json_mapping, "", " ")
	if err != nil {
		log.Println("JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}

	jsonpath := filepath.Join(".", "json")
	err = os.MkdirAll(jsonpath, os.ModePerm)
    if err != nil {
        log.Println("Creating json directory encounters error")
        log.Println(err)
    }
	_ = ioutil.WriteFile(filepath.Join(jsonpath, DestFile), jsonByte, 0644)
}

// OrgFile must end with ".json", DestFile must end with ".jpg"
func JsonToMetaData(OrgFile string) Metadata_jpg2json {
	var log = logger.BlockchainLog

	jsonfile, err := os.Open(OrgFile)
	if err != nil {
		log.Println("Opening json file encounters Error")
		log.Println(err)
	}
	defer jsonfile.Close()

	jsonByte, err := ioutil.ReadAll(jsonfile)
	if err != nil {
		log.Println("Read json as Byte encounters Error")
		log.Println(err)
		os.Exit(1)
	}

	var json_map Metadata_json2jpg

	err = json.Unmarshal(jsonByte, &json_map)
	if err != nil {
		log.Println("Unmarshal JSON File stats encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	// log.Println("MetaData Image String:", json_map.Image)
	// log.Println("MetaData Image Byte:", []byte(json_map.Image))

	byte_arr, err := base64.StdEncoding.DecodeString(json_map.Image)
	if err != nil {
		log.Println("Convert Base64 to Byte encounters Error")
		log.Println(err)
		os.Exit(1)
	}

	return Metadata_jpg2json{
		FileName:         json_map.FileName,
		Size:             json_map.Size,
		ModificationTime: json_map.ModificationTime,
		Image:            byte_arr,
	}
}

func MetaDataToJpg(json_map Metadata_jpg2json, DestFile string) {
	var log = logger.BlockchainLog

	buffer := bytes.NewReader(json_map.Image)
	img, _, err := image.Decode(buffer)
	if err != nil {
		log.Println("Decode to Image encounters Error")
		log.Println(err)
		os.Exit(1)
	}

	jpgpath := filepath.Join(".", "jpg")
	err = os.MkdirAll(jpgpath, os.ModePerm)
    if err != nil {
        log.Println("Creating jpg directory encounters error")
        log.Println(err)
    }
	output, err := os.Create(filepath.Join(jpgpath, DestFile))
	if err != nil {
		log.Println("Creating new file encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	defer output.Close()

	var opts jpeg.Options
	opts.Quality = 100

	_ = jpeg.Encode(output, img, &opts)
}

func UnitTest_JpgToJson(OrgFile string, DestFile string) {
	mapping := JpgToMetaData(OrgFile)
	MetaDataToJson(mapping, DestFile)
}

func UnitTest_JsonToJpg(OrgFile string, DestFile string) {
	mapping := JsonToMetaData(OrgFile)
	MetaDataToJpg(mapping, DestFile)
}
