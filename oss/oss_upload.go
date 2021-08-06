package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	slog "github.com/jau1jz/cornus/commons/log"
	"io"
	"io/ioutil"
	"time"
)

type Client interface {
	UploadAndSignUrl(fileReader io.Reader, objectName string, expiredInSec int64) (string, error)
	DeleteByObjectName(objectName string)
	UploadByReader(fileReader io.Reader, fileName string) (err error)
	DownloadFile(fileName string) (data []byte, err error)
	IsFileExist(fileName string) (isExist bool, err error)
	GetFileURL(fileName string, expireTime time.Duration) (url string, err error)
}

type ClientImp struct {
	ossBucket       string
	accessKeyID     string
	accessKeySecret string
	ossEndPoint     string
}

func (slf *ClientImp) GetFileURL(fileName string, expireTime time.Duration) (url string, err error) {
	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Slog.ErrorF("ClientImp IsFileExist Error:%s", err)
		return "", err
	}
	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Slog.ErrorF("ClientImp IsFileExist  Error:%s", err)
		return "", err
	}
	url, err = bucket.SignURL(fileName, oss.HTTPGet, int64(expireTime))
	if err != nil {
		return "", err
	}
	return

}

func ClientInstance(ossBucket, accessKeyID, accessKeySecret, ossEndPoint string) Client {
	return &ClientImp{
		ossBucket:       ossBucket,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		ossEndPoint:     ossEndPoint,
	}
}

func (slf *ClientImp) IsFileExist(fileName string) (isExist bool, err error) {
	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Slog.ErrorF("ClientImp IsFileExist Error:%s", err)
		return false, err
	}
	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Slog.ErrorF("ClientImp IsFileExist  Error:%s", err)
		return false, err
	}
	return bucket.IsObjectExist(fileName)
}

func (slf *ClientImp) UploadAndSignUrl(fileReader io.Reader, objectName string, expiredInSec int64) (string, error) {
	// 创建OSSClient实例。
	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return "", err
	}
	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return "", err
	}
	err = bucket.PutObject(objectName, fileReader)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return "", err
	}
	//oss.Process("image/format,png")
	signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, expiredInSec)
	if err != nil {
		bucket.DeleteObject(objectName)
		return "", err
	}
	return signedURL, nil
}

func (slf *ClientImp) DeleteByObjectName(objectName string) {
	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return
	}
	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return
	}
	err = bucket.DeleteObject(objectName)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
	}
}

func (slf *ClientImp) UploadByReader(fileReader io.Reader, fileName string) (err error) {

	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return
	}

	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return
	} else {
		fmt.Println("bucket ok")
	}

	err = bucket.PutObject(fileName, fileReader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return
}

func (slf *ClientImp) DownloadFile(fileName string) (data []byte, err error) {

	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return
	}

	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Slog.ErrorF("Error:%s", err)
		return
	} else {
		fmt.Println("bucket ok")
	}

	body, err := bucket.GetObject(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	data, err = ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return
}
