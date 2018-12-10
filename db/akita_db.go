package db

import (
	"akita/common"
	"akita/utils"
	"sync"
)

// 数据文件对应结构体
type AkitaDB struct {
	mutex 		sync.Mutex			// 互斥锁
	dataFile    string			    // 数据文件
	size        int64				// 记录文件大小
}

var dbInstance *AkitaDB

// 全局单例的 AkitaDB
func getSingletonAkitaDb() *AkitaDB{
	if dbInstance == nil {
		dbInstance = &AkitaDB{dataFile: common.DefaultDataFile, size: 0}
	}
	return dbInstance
}

func (db *AkitaDB) Reload() (bool, error) { 											// 数据重新载入
	return false, nil
}


// 向数据文件中写入一条记录
func (db *AkitaDB)WriteRecord (offset int64, record * DataRecord) (int64, error) {	// 将记录写入
	ksBuf, err := utils.Int32ToByteSlice(record.dateHeader.Ks)
	if err != nil {
		return 0, err
	}
	vsBuf, err := utils.Int32ToByteSlice(record.dateHeader.Vs)
	if err != nil {
		return 0, err
	}
	flagBuf, err := utils.Int32ToByteSlice(record.dateHeader.Flag)
	if err != nil {
		return 0, err
	}
	recordBuf := utils.AppendByteSlice(ksBuf, vsBuf, flagBuf, record.key, record.value)
	crc32 := utils.CreateCrc32(recordBuf)
	crcBuf, err := utils.UintToByteSlice(crc32)
	if err != nil {
		return 0, err
	}
	recordBuf = append(recordBuf, crcBuf...)
	db.mutex.Lock()																		// 互斥锁上锁
	curOffset, err := common.WriteFileWithByte(db.dataFile, offset, recordBuf)
	if err != nil {
		return 0, err
	}
	db.size = curOffset
	defer db.mutex.Unlock()																// 解锁
	return curOffset, nil
}

func (db *AkitaDB)ReadRecord(offset int64) ([]byte, error) {
	kvsBuf, err := common.ReadFileToByte(db.dataFile, offset, common.KvsByteLength)
	if err != nil {
		return nil, err
	}
	ksBuf := kvsBuf[0:common.KsByteLength:common.KsByteLength]
	vsBuf := kvsBuf[common.KsByteLength:len(kvsBuf):common.VsByteLength]
	ks, err := utils.ByteSliceToInt32(ksBuf)
	if err != nil {
		return nil, err
	}
	vs, err := utils.ByteSliceToInt32(vsBuf)
	if err != nil {
		return nil, err
	}
	fkvLength := common.FlagByteLength + int64(ks) + int64(vs)
	recordWithoutKvsBuf, err := common.ReadFileToByte(db.dataFile, offset + common.KvsByteLength, fkvLength + common.CrcByteLength)
	if err != nil {
		return nil, err
	}
	flagKeyValBuf := recordWithoutKvsBuf[0:fkvLength]
	valueBuf := recordWithoutKvsBuf[fkvLength + int64(ks) - 1:fkvLength]
	crc32Buf := recordWithoutKvsBuf[fkvLength:]
	recordWithoutCrc32Buf := utils.AppendByteSlice(kvsBuf, flagKeyValBuf)
	recordCrc32, err := utils.ByteSliceToUint(crc32Buf)
	if err != nil {
		return nil, err
	}
	checkCrc32 := utils.CreateCrc32(recordWithoutCrc32Buf)
	if err != nil {
		return nil, err
	}
	if recordCrc32 != checkCrc32 {
		return nil, common.ErrDataHasBeenModified
	}
	return valueBuf, nil
}

func (conn *Connection) Close() error {									// 关闭连接, 使 Connection 实现 io.Closer
	return nil
}



