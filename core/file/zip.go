package CoreFile

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

//ZipDir 压缩文件夹
//param src string 源文件
//param zipSrc string 目标压缩包
//return error 错误信息
func ZipDir(src string, zipSrc string) error {
	//构建ZIP文件
	d, err := os.Create(zipSrc)
	if err != nil{
		return err
	}
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	//读取目录
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	//遍历文件
	for _, fileSrc := range dir {
		if err := zipDirC(src + Sep + fileSrc.Name(), "", w); err != nil {
			return err
		}
	}
	return nil
}

//压缩目录子操作结构
//param file *os.File 文件句柄
//param prefix string
//param zw *zip.Writer 写入ZIP句柄
//return error 错误代码
func zipDirC(fileSrc string, prefix string, zw *zip.Writer) error {
	cFile, err := os.Open(fileSrc)
	if err != nil{
		return err
	}
	defer cFile.Close()
	info, err := cFile.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		info, err := cFile.Readdir(-1)
		if err != nil {
			return err
		}
		for _, vInfo := range info {
			f, err := os.Open(cFile.Name() + "/" + vInfo.Name())
			if err != nil {
				return err
			}
			err = zipDirC(fileSrc, prefix, zw)
			if err != nil {
				if err2 := f.Close(); err2 != nil{
					return errors.New(err.Error() + ", " + err2.Error())
				}
				return err
			}
			if err2 := f.Close(); err2 != nil{
				return err2
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if err != nil{
			return err
		}
		header.Name = prefix + "/" + header.Name
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, cFile)
		if err != nil {
			return err
		}
	}
	return nil
}

//UnZip 解压文件
//param zipSrc string 目标压缩包
//param dest string 解压到... eg : /dir/
//return error 错误信息
func UnZip(zipSrc string, dest string) error {
	reader, err := zip.OpenReader(zipSrc)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		if err := unZipC(file, dest); err != nil{
			return err
		}
	}
	return nil
}

//unZipC 解压子处理函数
func unZipC(file *zip.File, dest string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func(){
		_ = rc.Close()
	}()
	filename := dest + file.Name
	err = os.MkdirAll(GetDir(filename), 0755)
	if err != nil {
		return err
	}
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = w.Close()
	}()
	_, err = io.Copy(w, rc)
	if err != nil {
		return err
	}
	return nil
}