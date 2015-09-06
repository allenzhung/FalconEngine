/*****************************************************************************
 *  file name : StringIndexDic.go
 *  author : Wu Yinghao
 *  email  : wyh817@gmail.com
 *
 *  file description : 字符串到ID的hash映射文件，可以为每个term生成唯一的ID，用于后续的
 *   				   倒排索引
 *
******************************************************************************/

package utils

import (
	"fmt"
	"os"
	"bytes"
	"syscall"
	"encoding/binary"
)

type StringIdxDic struct {
	Lens        int64
	StringMap	map[string]int64
	Index    	int64
	Name		string
}







func NewStringIdxDic(name string) *StringIdxDic {

	this := &StringIdxDic{StringMap:make(map[string]int64),Lens:0,Index:1,Name:name}
	return this
}

func (this *StringIdxDic) Display() {
	
	fmt.Printf("========================= Lens : %v  Count :%v =========================\n", this.Lens, this.Index)
	for k,v := range this.StringMap {
		fmt.Printf("Key : %v \t\t--- Value : %v  \n", k, v)
	}
	fmt.Printf("===============================================================================\n")
	
}

/*****************************************************************************
*  function name : PutKeyForInt
*  params : 输入的key
*  return :
*
*  description : 在hash表中添加一个key，产生一个key的唯一id
*
******************************************************************************/
func (this *StringIdxDic) Put(key string) int64 {
	
	id:= this.Find(key)
	if id!=-1{
		return id
	}
	
	this.StringMap[key] = this.Index
	this.Index ++ 
	this.Lens ++
	
	return this.StringMap[key]
	
}

/*****************************************************************************
*  function name : Length
*  params : nil
*  return : int64
*
*  description : 返回哈希表长度
*
******************************************************************************/
func (this *StringIdxDic) Length() int64 {
	
	return this.Lens

}

func (this *StringIdxDic) Find(key string) int64 {
	
	value,has_key:=this.StringMap[key]
	if has_key{
		return value
	}
	return -1
}


func (this *StringIdxDic) WriteToFile() error {
	
	fmt.Printf("Writing to File [%v]...\n", this.Name)
	file_name := fmt.Sprintf("./index/%v.dic",this.Name)
	fout, err := os.Create(file_name)
	defer fout.Close()
	if err != nil {
		return err
	}
	
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, this.Lens)
	err = binary.Write(buf, binary.LittleEndian, this.Index)
	if err != nil {
			fmt.Printf("Lens ERROR :%v \n",err)
		}
	for k,v := range this.StringMap {
		err := binary.Write(buf, binary.LittleEndian, int64(len(k)))
		if err != nil {
			fmt.Printf("Write Lens Error :%v \n",err)
		}
		err = binary.Write(buf, binary.LittleEndian, []byte(k))
		if err != nil {
			fmt.Printf("Write Key Error :%v \n",err)
		}
		err = binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			fmt.Printf("Write Value Error :%v \n",err)
		}
	}
	fout.Write(buf.Bytes())
	return nil
}


func (this *StringIdxDic) ReadFromFile() error {
	
	file_name := fmt.Sprintf("./index/%v.dic",this.Name)
	f, err := os.Open(file_name)
	defer f.Close()
	if err != nil {
		return  err
	}
	
	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("ERR:%v", err)
	}
	
	MmapBytes, err := syscall.Mmap(int(f.Fd()), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)

	if err != nil {
		fmt.Printf("MAPPING ERROR  %v \n", err)
		return nil
	}

	defer syscall.Munmap(MmapBytes)
	
	
	this.Lens=int64(binary.LittleEndian.Uint64(MmapBytes[:8]))
	this.Index=int64(binary.LittleEndian.Uint64(MmapBytes[8:16]))
	var start int64 = 16
	var i int64 = 0
	for i=0;i<this.Lens;i++{
		lens:=int64(binary.LittleEndian.Uint64(MmapBytes[start:start+8]))
		start+=8
		key:=string(MmapBytes[start:start+lens])
		start+=lens
		value:=int64(binary.LittleEndian.Uint64(MmapBytes[start:start+8]))
		start+=8
		this.StringMap[key]=value
	}
	

	return nil
	
	
}