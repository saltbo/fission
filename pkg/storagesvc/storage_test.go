package storagesvc

import (
	"bytes"
	"fmt"
	"github.com/gw123/glog"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewS3Storage(t *testing.T) {
	input := map[string]string{
		"bucketName":      "tmpBucket",
		"subDir":          "a/b/c",
		"accessKeyID":     "tmpAccessKeyID",
		"secretAccessKey": "tmpSecretAccessKey",
		"region":          "ap-south-1",
	}

	os.Setenv("STORAGE_S3_BUCKET_NAME", input["bucketName"])
	os.Setenv("STORAGE_S3_SUB_DIR", input["subDir"])
	os.Setenv("STORAGE_S3_ACCESS_KEY_ID", input["accessKeyID"])
	os.Setenv("STORAGE_S3_SECRET_ACCESS_KEY", input["secretAccessKey"])
	os.Setenv("STORAGE_S3_REGION", input["region"])

	storage := NewS3Storage().(s3Storage)

	for k, v := range input {
		valueInStruct := reflect.Indirect(reflect.ValueOf(storage)).FieldByName(k).String()

		if valueInStruct != v {
			t.Errorf("Incorrect s3Storage field. Got: %s, Want %s", valueInStruct, v)
		}
	}

	if storage.storageType != StorageTypeS3 {
		t.Errorf("Incorrect storageType field. Got: %s, Want %s", storage.storageType, StorageTypeS3)
	}

	// TestGetStorageType
	if storage.getStorageType() != storage.storageType {
		t.Errorf("Incorrect getStorateType() method implementation. Got: %s, Want %s", storage.getStorageType(), storage.storageType)
	}

}

func TestNewLocalStorage(t *testing.T) {
	storage := NewLocalStorage("/fission").(localStorage)

	// // When SUBDIR env is not set, expect a default "fission-functions" value.
	// if storage.subDir != "fission-functions" {
	// 	t.Errorf("Incorrect subDir field. Got: %s, Want %s", storage.subDir, "fission-functions")
	// }
	if storage.storageType != StorageTypeLocal {
		t.Errorf("Incorrect storageType field. Got: %s, Want %s", storage.storageType, StorageTypeLocal)
	}
}

func TestNewS3StorageQiniu(t *testing.T) {
	//input := map[string]string{
	//	"bucketName":      "func-fsrv",
	//	"subDir":          "loki",
	//	"accessKeyID":     "BMEpw0S-VbwPTAYEiZFmVXRutAJFSNH68UTI92jI",
	//	"secretAccessKey": "KJp_bcCaURH-EIa9qzOWG-7OGKDEcgkMiI7_ndHV",
	//	"region":          "cn-north-1",
	//	"endpoint":        "s3-cn-north-1.qiniucs.com",
	//}

	/**
	s3endpoint    = "s3.bj.bcebos.com"
	s3region      = "bj"
	s3accessKeyId = "fb483d4a4912494c8ef4e4c7d422e1c4"
	s3SecretKeyId = "8312e297a7644bf486e293cebb4ddb1b"
	*/

	input := map[string]string{
		"bucketName":      "prod-faas-func",
		"subDir":          "dev",
		"accessKeyID":     "fb483d4a4912494c8ef4e4c7d422e1c4",
		"secretAccessKey": "8312e297a7644bf486e293cebb4ddb1b",
		"region":          "bj",
		"endpoint":        "s3.bj.bcebos.com",
	}

	os.Setenv("STORAGE_S3_BUCKET_NAME", input["bucketName"])
	os.Setenv("STORAGE_S3_SUB_DIR", input["subDir"])
	os.Setenv("STORAGE_S3_ACCESS_KEY_ID", input["accessKeyID"])
	os.Setenv("STORAGE_S3_SECRET_ACCESS_KEY", input["secretAccessKey"])
	os.Setenv("STORAGE_S3_REGION", input["region"])
	//STORAGE_S3_ENDPOINT
	os.Setenv("STORAGE_S3_ENDPOINT", input["endpoint"])

	storage := NewS3Storage().(s3Storage)

	for k, v := range input {
		valueInStruct := reflect.Indirect(reflect.ValueOf(storage)).FieldByName(k).String()

		if valueInStruct != v {
			t.Errorf("Incorrect s3Storage field. Got: %s, Want %s", valueInStruct, v)
		}
	}

	if storage.storageType != StorageTypeS3 {
		t.Errorf("Incorrect storageType field. Got: %s, Want %s", storage.storageType, StorageTypeS3)
	}

	// TestGetStorageType
	if storage.getStorageType() != storage.storageType {
		t.Errorf("Incorrect getStorateType() method implementation. Got: %s, Want %s", storage.getStorageType(), storage.storageType)
	}

	l, err := storage.dial()
	if err != nil {
		glog.WithErr(err).Error("storage dial")
		return
	}

	c, err := l.Container(input["bucketName"])
	if err != nil {
		glog.WithErr(err).Error("storage dial")
		return
	}
	item, err := c.Put("abc", bytes.NewReader([]byte("xxxxx")), 5, nil)
	if err != nil {
		glog.WithErr(err).Error("storage Put")
		return
	}

	glog.WithField("item", item).Infof("success")

	item, err = c.Item("abc")
	if err != nil {
		glog.WithErr(err).Error("storage Put")
		return
	}

	reader, err := item.Open()
	if err != nil {
		glog.WithErr(err).Error("storage Put")
		return
	}

	data, _ := ioutil.ReadAll(reader)
	fmt.Println("read data: ", string(data))

}
