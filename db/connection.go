package db

import (
	"akita/common"
	"sync/atomic"
)

// 所有的数据库操作，都需要获取 Connection
type Connection struct {

}

func (conn *Connection) Insert(key string, fileName string) (bool, error) {		// 插入数据
	keyBuf := common.StringToByteSlice(key)
	bufLen, err := common.GetFileSize(fileName)
	if err != nil {
		return false, err
	}
	valueBuf, err := common.ReadFileToByte(fileName, 0, bufLen)
	if err != nil {
		return false, err
	}
	ks := len(keyBuf)
	vs := len(valueBuf)
	if ks > common.K {
		return false, common.ErrKeySize
	}
	if vs > 10 * common.M {
		return false, common.ErrFileSize
	}
	dataRecord := &DataRecord{
		dateHeader: &DataHeader{
			Ks: int32(ks),
			Vs: int32(vs),
			Flag: common.WriteFlag,
		    },
		key: keyBuf,
		value: valueBuf,
	}
	akDb := getSingletonAkitaDb()
	offset := akDb.size
	errorChan  := make(chan error)
	lengthChan := make(chan int64)
	go func(record *DataRecord) {
		length, err := akDb.WriteRecord(record)
		errorChan  <-err
		lengthChan <-length
	}(dataRecord)

	if err = <-errorChan; err != nil {
		return false, err
	}
	akMap  := getSingletonAkitaMap()
	ir := &indexRecord{offset:offset, size: <-lengthChan}
	akMap.set(key, ir)												 	     // 设置 map 索引
	return true, nil
}

func (conn *Connection) Seek(key string) ([]byte, error) {					 // 查找数据
	akMap := getSingletonAkitaMap()
	ir, err := akMap.get(key)										 	     // 获取该记录的起始 offset
	if err != nil {
		return nil, err
	}
	valueChan := make(chan []byte)
	errChan   := make(chan error)

	go func() {
		akDb := getSingletonAkitaDb()
		value, err := akDb.ReadRecord(ir.offset, ir.size)
		valueChan <- value
		errChan <- err
	}()

	if err = <- errChan; err != nil {
		return nil, err
	}
	value := <- valueChan
	return value, nil
}

func (conn *Connection) Delete(key string) (bool, []byte, error)  {				 				// 删除数据, 返回删除的 []byte

	akMap := getSingletonAkitaMap()
	ok, ir, err := akMap.del(key)
	if !ok {
		return false, nil, err
	}
	keyBuf := common.StringToByteSlice(key)
	ks := len(keyBuf)
	dataRecord := &DataRecord{
		dateHeader: &DataHeader{
			Ks: int32(ks),
			Vs: int32(0),
			Flag: common.DeleteFlag,
	        },
		key: keyBuf,
		value: nil,
	}
	offset := akMap.CurOffset
	errChan    := make(chan error, 2)
	offsetChan := make(chan int64)
	valueChan  := make(chan []byte)
	akDb := getSingletonAkitaDb()
	// 读取数据
	go func() {
		value, err := akDb.ReadRecord(ir.offset, ir.size)
		valueChan <- value
		errChan   <- err
	}()

	if err = <- errChan; err != nil {
		return false, nil, err
	}
	value := <- valueChan

	go func(filePath string, from int64, record *DataRecord) {
		curOffset, err := akDb.WriteRecord(record)
		offsetChan <- curOffset
		errChan    <- err
	}(common.DefaultDataFile, offset, dataRecord)

	if err = <- errChan; err != nil {
		return false, nil, err
	}

	atomic.CompareAndSwapInt64(&akMap.CurOffset, offset, <-offsetChan)		 // 设置当前 offset
	return true, value, nil
}
