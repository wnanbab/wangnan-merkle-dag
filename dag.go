package merkledag

import "hash"

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	//将分片写到KVstore中
	//判断数据类型
	switch n := node.type() {
	case File :
		StoreFile(store,node,h)
		break
	
	case Dir:
		StoreDir(store,node,h)
		break
	}
	return nil
}

//存文件的方法
func StoreFile(store KVStore, node File, h hash.Hash) []byte {
	//1.获取文件数据
	data := node.Bytes() //获取数据
	h.write(data)  //计算哈希值
	hashValue := h.sum(nil)
	//2.将数据写入KVstore,检查错误
	err := store.Put(hashValue,data)
	if err != nil {
		return nil
	}
	return hashValue
}

func StoreDir(store KVStore, dir Dir, h hash.Hash) []byte {

	iter := dir.It(); //调用了dir目录的It方法，获取目录迭代器
	//迭代器遍历每个子节点
	for iter.Next() {
		childNode := iter.Node() ;//调用迭代器Node方法，获取每个节点信息
		childHash := Add(store,node,h) //递归调用 Add 函数存储子节点数据

		newLink := Link {
			Name : childNode.name(),
			Hash : childHash,
			Size : int(childNode.Size()),
		} 

		dir.Links = append(dir.Links,newLink)
	}

	//对整个目录节点进行序列化，计算哈希值
	data := dir.Bytes()
	h.Write(data)
	hashValue := h.Sum(nil)//计算最终的哈希
	//将数据写入KVstore
	err := store.Put(hashValue,data)
	if err != nil { //错误处理
		return nil
	}

	return hashValue  //返回计算得到的哈希值
}
