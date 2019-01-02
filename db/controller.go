package db

import (
	"akita/common"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

func save(ctx echo.Context) error {
	if ! Sev.IsMaster() {
		return ctx.JSON(http.StatusBadRequest, "sorry this akita node isn't master node! ")
	}
	key := ctx.FormValue("key")
	if key == "" {
		return ctx.JSON(http.StatusOK, "key can not be empty! ")
	}
	if len(common.StringToByteSlice(key)) > 10*common.K {
		return ctx.JSON(http.StatusOK, common.ErrKeySize)
	}
	file, err := ctx.FormFile("file")
	if file == nil {
		return ctx.JSON(http.StatusOK, "file can not be empty! ")
	}
	if err != nil {
		common.Error.Printf("Get form file fail: %s\n", err)
		return ctx.JSON(http.StatusOK, "file upload fail. Please try again. ")
	}
	var length int64
	if length = file.Size; length > 10*common.M {
		return ctx.JSON(http.StatusOK, "file is too large to save. ")
	}
	src, err := file.Open()
	defer src.Close()
	if err != nil {
		common.Error.Printf("File open fail: %s\n", err)
		return ctx.JSON(http.StatusOK, err)
	}
	_, err = Sev.Insert(key, src, length)
	if err != nil {
		return ctx.JSON(http.StatusOK, "save key: "+key+" fail")
	}
	return ctx.JSON(http.StatusOK, "save  key: "+key+" success! ")
}

func search(ctx echo.Context) error {
	key := ctx.QueryParam("key")
	if key == "" {
		return ctx.JSON(http.StatusOK, "key can not be empty!  ")
	}
	value, err := Sev.Seek(key)
	if err != nil {
		return ctx.JSON(http.StatusOK, "seek key: "+key+" fail. ")
	}
	return ctx.JSON(http.StatusOK, value)
}

func del(ctx echo.Context) error {
	if ! Sev.IsMaster() {
		return ctx.JSON(http.StatusBadRequest, "sorry this akita node isn't master node! ")
	}
	key := ctx.QueryParam("key")
	if key == "" {
		return ctx.JSON(http.StatusOK, "key can not be empty!  ")
	}
	_, delOffset, err := Sev.Delete(key)
	if err != nil {
		return ctx.JSON(http.StatusOK, "delete key: "+key+"fail. ")
	}
	return ctx.JSON(http.StatusOK, delOffset)
}

func syn(ctx echo.Context) error {
	if ! Sev.IsMaster() {
		return ctx.JSON(http.StatusBadRequest, nil)
	}
	offsetStr := ctx.QueryParam("offset")
	offset, _ := strconv.Atoi(offsetStr)
	timeout := time.After(1000 * time.Millisecond)
	select {
	case <-timeout:
		return ctx.JSON(http.StatusOK, nil)
	case size := <-Sev.synChan:
		if size > int64(offset) {
			length := size - int64(offset)
			buf, err := common.ReadFileToByte(Sev.dB.dataFile, int64(offset), length)
			if err != nil {
				return ctx.JSON(http.StatusOK, nil)
			}
			return ctx.JSON(http.StatusOK, buf)
 		}else {
 			return ctx.JSON(http.StatusOK, nil)
		}
	}
}
