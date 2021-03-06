package db

import (
	"akita/common"
	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"time"
)

func save(ctx echo.Context) error {
	if !Sever.IsMaster() {
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
	_, err = Sever.Insert(key, src, length)
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
	value, err := Sever.Seek(key)
	if err != nil {
		return ctx.JSON(http.StatusOK, "seek key: "+key+" fail. ")
	}
	return ctx.Blob(http.StatusOK, "image/jpeg", value)
}

func del(ctx echo.Context) error {
	if !Sever.IsMaster() {
		return ctx.JSON(http.StatusBadRequest, "sorry this akita node isn't master node! ")
	}
	key := ctx.QueryParam("key")
	if key == "" {
		return ctx.JSON(http.StatusOK, "key can not be empty!  ")
	}
	_, delOffset, err := Sever.Delete(key)
	if err != nil {
		return ctx.JSON(http.StatusOK, "delete key: "+key+"fail. ")
	}
	return ctx.JSON(http.StatusOK, delOffset)
}

func syn(ctx echo.Context) error { // deal with slaves sync request
	if !Sever.IsMaster() {
		http.Error(ctx.Response(), "sorry, slaves server can not sync data", http.StatusBadRequest)
		return nil
	}
	reqBody := ctx.Request().Body
	defer reqBody.Close()

	offsetBuf, err := ioutil.ReadAll(reqBody)
	if err != nil {
		http.Error(ctx.Response(), err.Error(), http.StatusBadRequest)
		return err
	}
	syncOffset := &SyncOffset{}
	err = proto.Unmarshal(offsetBuf, syncOffset)
	if err != nil {
		common.Error.Printf("proto data unmarshal error: %s \n", err)
		http.Error(ctx.Response(), err.Error(), http.StatusInternalServerError)
		return err
	}
	common.Info.Printf("request offset is %d\n", syncOffset.Offset)

	complete := make(chan error)
	dataCh := make(chan []byte)
	go func() {
		data, err := Sever.db.GetDataByOffset(syncOffset.Offset)
		dataCh <-data
		complete <-err
	}()
	data := <-dataCh
	err = <-complete

	syncData := &SyncData{}
	if err != nil {
		if err == common.ErrNoDataUpdate {
			notifier := make(chan struct{})
			Sever.register(ctx.Request().Host, notifier)
			select {
			case <-time.After(1000 * time.Millisecond):
				syncData.Code = 0
				syncData.Data = nil
			case <-notifier:
				go func() {
					data, err := Sever.db.GetDataByOffset(syncOffset.Offset)
					dataCh <-data
					complete <-err
				}()
				data = <-dataCh
				err = <-complete

				common.Info.Printf("the data length is %d\n", len(data))
				if err != nil {
					common.Error.Printf("get data by offset error :%s\n", err)
					syncData.Code = 0
					syncData.Data = nil
				}
				syncData.Code = 1
				syncData.Data = data
			}
		} else {
			common.Error.Printf("get data by offset error :%s\n", err)
			syncData.Code = 0
			syncData.Data = nil
		}
	} else {
		syncData.Code = 1
		syncData.Data = data
		common.Info.Printf("the data length is %d\n", len(data))
	}
	protoData, _ := proto.Marshal(syncData)
	ctx.Response().Header().Set("content-type", "application/protobuf") // use protobuf format to transport data
	if _, err = ctx.Response().Write(protoData); err != nil {           // if response error, reply 500 error
		http.Error(ctx.Response(), err.Error(), http.StatusInternalServerError)
	}
	return nil
}
