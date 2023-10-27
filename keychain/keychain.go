package keychain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func NewKeychain(filePath string) (kc Keychain, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return kc, err
	}

	if kc.Header, err = unmarshalHeader(data[:20]); err != nil {
		return kc, err
	}

	tableListByte := data[20:]
	if kc.TableList, err = unmarshalTableList(tableListByte); err != nil {
		return kc, err
	}

	return kc, err
}

func unmarshalHeader(headerData []byte) (header keychainHeader, err error) {
	r := bytes.NewReader(headerData)

	if err = binary.Read(r, binary.BigEndian, &header); err != nil {
		return header, err
	}

	if string(header.Signature[:]) != signature {
		return header, fmt.Errorf("keychain signature does not match")
	}
	return header, nil
}

func unmarshalTableList(tableListDataBytes []byte) (tableList keychainTableList, err error) {
	r := bytes.NewReader(tableListDataBytes)

	var header struct {
		DataSize       int32
		NumberOfTables int32
	}
	if err = binary.Read(r, binary.BigEndian, &header); err != nil {
		return tableList, err
	}

	tableList.DataSize = header.DataSize
	tableList.NumberOfTables = header.NumberOfTables

	tableOffsets := make([]int32, header.NumberOfTables)
	if err := binary.Read(r, binary.BigEndian, &tableOffsets); err != nil {
		return tableList, err
	}

	if tableOffsets[header.NumberOfTables-1] > header.DataSize {
		return tableList, fmt.Errorf("table list does not match")
	}

	table := make([]keychainTable, tableList.NumberOfTables)
	for i, offset := range tableOffsets {
		var tableHeader keychainTableHeader
		tableDataBytes := tableListDataBytes[offset:header.DataSize]
		if tableHeader, err = unmarshalTableHeader(tableDataBytes); err != nil {
			log.Println(err.Error())
		}

		table[i].offset = offset
		table[i].Header = tableHeader

		recordList := make([]Record, tableHeader.NumberOfRecord)
		for recordIndex, recordOffset := range tableHeader.RecordOffset {
			record, err := unmarshalRecord(tableDataBytes[recordOffset:tableHeader.DataSize])
			if err != nil {
				continue
			}
			recordList[recordIndex] = record
		}
		println(tableHeader.RecordType.String())
		table[i].RecordList = recordList
	}

	tableList.Table = table
	return tableList, err
}

func unmarshalTableHeader(tableDataByte []byte) (tableHeader keychainTableHeader, err error) {
	r := bytes.NewReader(tableDataByte)

	var header struct {
		DataSize             int32
		RecordType           RecordTypes
		NumberOfRecord       int32
		RecordListOffset     int32
		Unknown1             int32
		Unknown2             int32
		NumberOfRecordOffset int32
	}

	if err = binary.Read(r, binary.BigEndian, &header); err != nil {
		return tableHeader, err
	}
	tableHeader.DataSize = header.DataSize
	tableHeader.RecordType = header.RecordType
	tableHeader.NumberOfRecord = header.NumberOfRecord
	tableHeader.RecordListOffset = header.RecordListOffset
	tableHeader.Unknown1 = header.Unknown1
	tableHeader.Unknown2 = header.Unknown2
	tableHeader.NumberOfRecordOffset = header.NumberOfRecordOffset

	recordOfOffset := make([]int32, header.NumberOfRecord)

	if err = binary.Read(r, binary.BigEndian, &recordOfOffset); err != nil {
		return tableHeader, err
	}

	if len(recordOfOffset) <= 0 {
		return tableHeader, err
	}

	if recordOfOffset[header.NumberOfRecord-1] > header.DataSize {
		return tableHeader, fmt.Errorf("record offset does not match")
	}

	tableHeader.RecordOffset = recordOfOffset
	return tableHeader, err
}

func unmarshalRecord(recordDataBytes []byte) (record Record, err error) {

	var header RecordHeader
	if header, err = unmarshalRecordHeader(recordDataBytes); err != nil {
		return record, err
	}

	record.Header = header

	return record, err
}

func unmarshalRecordHeader(recordDataBytes []byte) (recordHeader RecordHeader, err error) {
	r := bytes.NewReader(recordDataBytes)

	if err = binary.Read(r, binary.BigEndian, &recordHeader); err != nil {
		return recordHeader, err
	}

	return recordHeader, err
}
