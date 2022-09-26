package crc

import (
	"github.com/tyundev/x/utils/md"

	"github.com/sigurn/crc16"
)

const (
	Reiway        = "Reiway@23031993@18021995"
	PublicKeyQRCG = "9e52c703ba0579a31a8461d27634e5d4" //yeuthuongvanvat@kientrichamchithattam
)

func CheckCmd(allNum, cmd, addr, numData string) (bool, isSetup bool) {
	if allNum != AllNumByte_Setup && allNum != AllNumByte_Status {
		return false, isSetup
	}
	if cmd != CMD_Setup && cmd != CMD_Status {
		return false, isSetup
	}
	if addr != ADDR_Setup && addr != ADDR_Status {
		return false, isSetup
	}
	if numData != NumDATA_Setup && numData != NumDATA_Status {
		return false, isSetup
	}
	switch cmd {
	case CMD_Setup:
		isSetup = true
	default:
		isSetup = false
	}
	return true, isSetup
}

func GetCrc16(data string) int {
	var datac = []rune(data)
	var c = 0
	var flag = 0
	var i = 0
	var crc16 = 0xffff
	var len = len(datac)
	for i < len {
		crc16 = crc16 ^ int(datac[i])
		for c = 0; c < 8; c++ {
			flag = (crc16 & 0x01)
			crc16 = (crc16 >> 1)
			if flag != 0 {
				crc16 = (crc16 ^ 0xa001)
			}
		}
		i++
	}
	return crc16
}

func GetCheckSumReiway(data string) int64 {
	data += Reiway
	var md5, _ = md.Encrypt([]byte(data), PublicKeyQRCG)
	if md5 == "" {
		return 0
	}
	//data, _ = auth.GererateHashedPassword(data)
	table := crc16.MakeTable(crc16.CRC16_X_25)
	checkSum := crc16.Checksum([]byte(data), table)
	return int64(checkSum)
}
