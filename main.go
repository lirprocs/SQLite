package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
)

func generateKey(key string) []byte {
	//key := "pop1234567890pop"
	return []byte(key)
}

func encryptFile(filename string, key []byte) error {
	inputFile, err := os.Open(filename)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(filename + ".enc")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}

	outputFile.Write(iv)

	stream := cipher.NewCFBEncrypter(cipherBlock, iv)
	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	_, err = io.Copy(writer, inputFile)
	if err != nil {
		return err
	}

	inputFile.Close()
	os.Remove(filename)
	return nil
}

func decryptFile(filename string, key []byte) error {
	inputFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	//defer inputFile.Close()

	outputFile, err := os.Create(filename[:len(filename)-4]) // убираем ".enc" из имени файла
	if err != nil {
		return err
	}
	defer outputFile.Close()

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = inputFile.Read(iv)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(cipherBlock, iv)
	reader := &cipher.StreamReader{S: stream, R: inputFile}

	_, err = io.Copy(outputFile, reader)
	if err != nil {
		return err
	}

	inputFile.Close()
	return nil
}

func Alter_password(folderName string, old_pass, new_pass []byte) string {
	err := decryptFile(folderName+".enc", old_pass)
	if err != nil {
		fmt.Println("Error decrypting folder:", err)
	}

	err = encryptFile(folderName, new_pass)
	if err != nil {
		fmt.Println("Error encrypting folder:", err)
	}
	return "Password alter successfully"
}

func Delete(DB *sql.DB) {
	fmt.Printf("\nDelete:\n")
	type Obor struct {
		id   int
		name string
	}

	row, _ := DB.Query("SELECT Оборудование.№, Оборудование.Название FROM Оборудование")
	o := Obor{}
	fmt.Printf("До удаления:\n")
	for row.Next() {
		row.Scan(&o.id, &o.name)
		fmt.Printf("%d %s\n", o.id, o.name)
	}

	_, err := DB.Exec("DELETE FROM Оборудование WHERE `№` > 10;")
	if err != nil {
		fmt.Println("Error DELETE:", err)
		return
	} else {
		row, _ := DB.Query("SELECT Оборудование.№, Оборудование.Название FROM Оборудование")
		o := Obor{}
		fmt.Printf("После удаления:\n")
		for row.Next() {
			row.Scan(&o.id, &o.name)
			fmt.Printf("%d %s\n", o.id, o.name)
		}
	}
}

func Insert(DB *sql.DB) {
	fmt.Printf("\nInsert:\n")
	type Obor struct {
		id   int
		name string
	}

	row, _ := DB.Query("SELECT Оборудование.№, Оборудование.Название FROM Оборудование")
	o := Obor{}
	fmt.Printf("До добавления:\n")
	for row.Next() {
		row.Scan(&o.id, &o.name)
		fmt.Printf("%d %s\n", o.id, o.name)
	}

	_, err := DB.Exec("INSERT INTO Оборудование ( №, Название, Тип, Цена, Год_выпуска) VALUES (22, 'Петли', 'Звуковое', 10000, '2022-01-01');")
	if err != nil {
		fmt.Println("Error INSERT INTO:", err)
		return
	} else {
		row, _ := DB.Query("SELECT Оборудование.№, Оборудование.Название FROM Оборудование")
		o := Obor{}
		fmt.Printf("После добавления:\n")
		for row.Next() {
			row.Scan(&o.id, &o.name)
			fmt.Printf("%d %s\n", o.id, o.name)
		}
	}
}

func Select(DB *sql.DB) {
	fmt.Printf("\nSelect:\n")
	type Obor struct {
		id   int
		name string
	}

	row, _ := DB.Query("SELECT Оборудование.№, Оборудование.Название FROM Оборудование WHERE Оборудование.№ > 8")
	o := Obor{}
	fmt.Printf("Выбока оборудования с № > 8:\n")
	for row.Next() {
		row.Scan(&o.id, &o.name)
		fmt.Printf("%d %s\n", o.id, o.name)
	}
}

func Update(DB *sql.DB) {
	fmt.Printf("\nUpdate:\n")
	type Obed struct {
		name  string
		city  string
		price int
	}

	row, _ := DB.Query("SELECT Объединения.Название, Объединения.Город, Объединения.Цена_выступления FROM Объединения WHERE Город = 'Москва';")
	o := Obed{}
	fmt.Printf("ДО обновления:\n")
	for row.Next() {
		row.Scan(&o.name, &o.city, &o.price)
		fmt.Printf("%s %s %d\n", o.name, o.city, o.price)
	}

	_, err := DB.Exec("UPDATE Объединения SET Цена_выступления = 18888 WHERE Город = 'Москва';")
	if err != nil {
		fmt.Println("Error UPDATE:", err)
		return
	} else {
		row, _ := DB.Query("SELECT Объединения.Название, Объединения.Город, Объединения.Цена_выступления FROM Объединения WHERE Город = 'Москва';")
		o := Obed{}
		fmt.Printf("После обновления:\n")
		for row.Next() {
			row.Scan(&o.name, &o.city, &o.price)
			fmt.Printf("%s %s %d\n", o.name, o.city, o.price)
		}
	}
}

func back(DB *sql.DB, folderName string, key []byte) {
	DB.Exec("INSERT INTO Оборудование ( №, Название, Тип, Цена, Год_выпуска) VALUES (11, 'LED-панель', 'Видео', 10000, '2022-01-01');")
	DB.Exec("DELETE FROM Оборудование WHERE `№` > 20;")
	DB.Exec("UPDATE Объединения SET Цена_выступления = 55555 WHERE Город = 'Москва';")

	DB.Close()
	encryptFile(folderName, key)
	fmt.Printf("Всё вернулось обратно. \n")
}

func test(DB *sql.DB) {
	err := DB.Ping()
	if err != nil {
		fmt.Println("Ваш пароль неверный\n", err)
	} else {
		fmt.Println("Ваш пароль верный\n")
	}
}

func main() {
	folderName := "REP.db"
	key := generateKey("pop1234567890pop")
	//new_key := generateKey("opo1234567890opo")

	decryptFile(folderName+".enc", key)

	DB, err := sql.Open("sqlite3", "./"+folderName)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer DB.Close()

	//test(DB)
	//Delete(DB)
	//Insert(DB)
	//Select(DB)
	//Update(DB)

	//back(DB, folderName, key) //ВЕРНУТЬ КАК БЫЛО

	//Alter_password(folderName, key, new_key)
}
