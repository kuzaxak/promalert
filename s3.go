package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/globalsign/mgo/bson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadFile(bucket, region string, plot io.WriterTo) (string, error) {
	// get the file size and read
	// the file content into a buffer
	s := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	_, err := s.Config.Credentials.Get()

	f, err := ioutil.TempFile("", "promplot-*.png")
	if err != nil {
		return "", fmt.Errorf("failed to create tmp file: %v", err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close tmp file: %v", err))
		}
		err := os.Remove(f.Name())
		if err != nil {
			panic(fmt.Errorf("failed to delete tmp file: %v", err))
		}
	}()
	_, err = plot.WriteTo(f)
	if err != nil {
		return "", fmt.Errorf("failed to write plot to file: %v", err)
	}

	fileInfo, _ := f.Stat()

	size := fileInfo.Size()
	buffer := make([]byte, size)
	_, err = f.Seek(0, io.SeekStart)
	_, err = f.Read(buffer)

	// create a unique file name for the file
	tempFileName := "pictures/" + bson.NewObjectId().Hex() + ".png"

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(tempFileName),
		ACL:           aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(int64(size)),
		ContentType:   aws.String(http.DetectContentType(buffer)),
		//ContentDisposition: aws.String("inline"),
		//ServerSideEncryption: aws.String("AES256"),
		//StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucket, region, tempFileName), err
}
