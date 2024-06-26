package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// 大端字节序（big-endian）: 将多字节数据类型的最高有效字节存储在内存的最低地址处，最低有效字节存储在最高地址处(符合人类的阅读习惯)

//  123456789
// 1 最高有效字节 先放进内存(内存最低地址处)

func main() {
	var a uint32 = 12373642
	fmt.Printf("%x\n", a) // bcce8a

	dataBuff := bytes.NewBuffer([]byte{})
	binary.Write(dataBuff, binary.LittleEndian, a)

	fmt.Println(dataBuff.Bytes())                             // [8a ce bc 0] byte是uint8，每个uint8是8bit，每个16进制数等于4个二进制数位
	fmt.Printf("%x\n", dataBuff.Bytes())                      // 8acebc00
	fmt.Println(binary.LittleEndian.Uint32(dataBuff.Bytes())) // 12373642

	// c8400066020000000001
	dataByte := []byte{0xc8, 0x40, 0x00, 0x66, 0x02, 0x00, 0x00, 0x00, 0x00, 0x01}
	t := binary.LittleEndian.Uint32(dataByte[0:4]) // type
	fmt.Println(t)

	sys_type := t >> 24
	fmt.Println(sys_type)

	fmt.Println(t & 0x3FF)

	fmt.Println(binary.LittleEndian.Uint32(dataByte[4:8])) // length

	str := "hello"
	strDataBuff := bytes.NewBuffer([]byte{})
	err := binary.Write(strDataBuff, binary.LittleEndian, []byte(str))
	fmt.Println(err)

	o := make([]byte, 1)
	err = binary.Read(bytes.NewBuffer(strDataBuff.Bytes()[0:1]), binary.LittleEndian, o)
	fmt.Println(err)
	fmt.Println("o", string(o))

}
