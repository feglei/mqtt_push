package rt_mod

import (
	"fmt"
)

var StockMapKeys []string

func StockAll() (data_str string) {
	for _, v := range StockMapKeys {
		data_str += v + ","
	}
	if len(data_str) == 0 {
		return ""
	}
	return data_str[:len(data_str)-1]
}

func GetSection(section int, index int) (data_str string) {
	sotck_len := len(StockMapKeys)
	section_size := int(sotck_len / section)
	section_start := section_size * index
	section_end := section_size * (index + 1)

	if index+1 >= section {
		section_end = sotck_len
	}
	fmt.Println("sotck_len", sotck_len)
	fmt.Println("section_size", section_size)
	fmt.Println("section_start", section_start)
	fmt.Println("section_end", section_end)
	fmt.Println(GetKeysIds(StockMapKeys[section_start:section_end]))
	return GetKeysIds(StockMapKeys[section_start:section_end])
}


func GetKeysIds(keys []string) string {
	data_str := ""
	for i := 0; i < len(keys); i++ {
		data_str += keys[i] + ","
	}
	if len(data_str) == 0 {
		return ""
	}
	return data_str[:len(data_str)-1]
}