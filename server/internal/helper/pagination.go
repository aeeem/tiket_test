package helper

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func CursorToDateAndID(cursor string) (lastDate string, lastID int, err error) {
	//take base64
	if cursor == "" {
		return
	}
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return
	}
	decodedString := string(b)
	strArr := strings.Split(decodedString, "=")
	if len(strArr) != 2 {
		err = errors.New("cursor not valid")
	}
	lastDate = strArr[0]
	lastID, err = strconv.Atoi(strArr[1])
	if err != nil {
		err = errors.New("Cursor format is not valid")
	}
	return
}

func DateAndIDToCursor(lastDate string, lastID int) (cursor string, err error) {
	combStr := fmt.Sprintf("%s=%d", lastDate, lastID)
	cursor = base64.StdEncoding.EncodeToString([]byte(combStr))
	return
}
