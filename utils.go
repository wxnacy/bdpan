package bdpan

import (
	"bdpan/common"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func httpResponseToInterface(r *http.Response, i interface{}) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bodyBytes, i); err != nil {
		return err
	}
	return nil
}

func encrypt(src []byte) ([]byte, error) {
	key, err := GetKey()
	if err != nil {
		return nil, err
	}
	return common.AesEncrypt(src, key)
}

func encryptInterfaceToHex(i interface{}) (string, error) {
	str, err := common.ToMapString(i)
	if err != nil {
		return "", err
	}
	bytes, err := encrypt([]byte(str))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func decrypt(src []byte) ([]byte, error) {
	key, err := GetKey()
	if err != nil {
		return nil, err
	}
	return common.AesDecrypt(src, key)
}

func decryptHexToInterface(src string, i interface{}) error {
	str, err := hex.DecodeString(src)
	if err != nil {
		return err
	}
	bytes, err := decrypt(str)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, i)
	if err != nil {
		return err
	}
	return nil
}

func SplitFile(path, tmpdir string, fragmentSize int64) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	var fileChunk = fragmentSize //  * (1 << 20) // 1 MB, change this to your requirement

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	// fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
	basename := filepath.Base(path)

	paths := make([]string, 0)
	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(float64(fileChunk), float64(fileSize-int64(int64(i)*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// write to disk
		// fileName := "somebigfile_" + strconv.FormatUint(i, 10)
		fileName := filepath.Join(tmpdir, basename+"_"+strconv.FormatUint(i, 10))
		_, err := os.Create(fileName)

		if err != nil {
			return nil, err
		}

		// write/save buffer to disk
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)

		// fmt.Println("Split to : ", fileName)
		paths = append(paths, fileName)
	}
	if len(paths) == 0 {
		return []string{path}, nil
	}
	return paths, nil
}
