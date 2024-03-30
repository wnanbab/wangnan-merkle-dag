package merkledag

import (
	"encoding/json"
	"strings"
)
// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 根据hash和path， 返回对应的文件, hash对应的类型是tree
	flag,_ := store.Has(hash); //调用has方法查看存储中是否存在指定的对象
	if flag {
		objBinary,_ := store.Get(hash); //获取二进制数据
		obj := binaryToObj(objBinary); //解析二进制数据
        pathArr:=strings.Split(path,"\\");//将路径字符串按照分隔符分割成数组pathArr
		cur := 1
		//传入解析得到的对象、路径数组、当前路径索引和存储对象,该函数会查找指定文件，返回文件内容
		return getFileByDir(obj,pathArr,cur,store)
	}
	return nil
}

func getFileByDir(obj* Object,pathArr []string,cur int,store KVStore) []byte {
	//判断当前处理的路径索引是否超出了路径数组的长度，则说明已经遍历完路径，直接返回 nil
	if cur >= len(path) {
		return nil
	}
	index := 0
	//遍历当前目录对象中的所有链接
	for i := range obj.Links {
		//从 obj.Data 中获取当前链接的类型，并将其转换为字符串。index 到 index+STEP 是当前链接类型在 obj.Data 中的索引范围
		objType := string(obj.Data[index : index+STEP])
		index += STEP
		//获取当前链接的信息
		objInfo := obj.Links[i]
		//如果当前链接的名称与路径数组中当前索引对应的路径不匹配，则跳过当前链接，继续下一次迭代
		if objInfo.Name != pathArr[cur] {
			continue
		}
		switch objType {
		case TREE:
			objDirBinary, _ := store.Get(objInfo.Hash)
			objDir := binaryToObj(objDirBinary)
			ans := getFileByDir(objDir, pathArr, cur+1, store)
			if ans != nil {
				return ans
			}
		//文件类型
		case BLOB:
			ans, _ := store.Get(objInfo.Hash)
			return ans
		//列表类型
		case LIST:
			objLinkBinary, _ := store.Get(objInfo.Hash)
			objList := binaryToObj(objLinkBinary)
			ans := getFileByList(objList, store)
			return ans
		}
	}
	return nil
}
//列表对象中递归查找文件，并将所有找到的文件内容拼接成一个大的[]byte 切片返回
func getFileByList(obj *Object, store KVStore) []byte {
	ans := make([]byte, 0)
	index := 0
	for i := range obj.Links {
		curObjType := string(obj.Data[index : index+STEP])
		index += STEP
		curObjLink := obj.Links[i]
		curObjBinary, _ := store.Get(curObjLink.Hash)
		curObj := binaryToObj(curObjBinary)
		if curObjType == BLOB {
			ans = append(ans, curObjBinary...)
		} else { //List
			tmp := getFileByList(curObj, store)
			ans = append(ans, tmp...)
		}
	}
	return ans
}
//将二进制数据解析成对象
func binaryToObj(objBinary []byte) *Object {
	var res Object
	//使用 json.Unmarshal 函数将二进制数据 objBinary 解析成 Object 结构体，并存储在 res 变量中
	json.Unmarshal(objBinary, &res)
	return &res  //返回地址
}


