package main

import (
	_ "github.com/btcsuite/btcd/btcec"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var pubkeyCurve = elliptic.P256()
var ipfsUrl = "localhost:5001"
var txtDB = "IPNS.txt"

var (
	requestFlag = flag.String("request-type", "", "record-set|record-get|sign")
	uidFlag = flag.String("uid", "", "user identity: <username>:<pubkey>")
	ipfsFlag = flag.String("ipfs", "", "ipfs-link")
	signFlag = flag.String("sign", "", "signature for new ipfs link")
	nameFlag = flag.String("username", "", "username")
)

func generateKeys () (*ecdsa.PrivateKey, error) { //сгенерировать новую пару ключей

	keyPair := new(ecdsa.PrivateKey)
	keyPair, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return keyPair, nil
}

func signMessage(userName string, msg string) { //сгенерировать подпись для заданной строки новой парой ключей

	keyPair, err := generateKeys()
	if err != nil {
		fmt.Println("Error: failed to generate keys")
		return
	}
	publicKey := elliptic.MarshalCompressed(pubkeyCurve, keyPair.PublicKey.X, keyPair.PublicKey.Y)

	signHash := sha256.Sum256([]byte(msg))
	sign, serr := ecdsa.SignASN1(rand.Reader, keyPair, signHash[:])
	if serr != nil {
		fmt.Println("Error: failed to generate signature")
		return
	}
	signature := hex.EncodeToString(sign)

	fmt.Printf("UserID = %s:%x \n", userName, publicKey) //uid пользователя с новым pubkey
	fmt.Println("Signature =", signature)
}

func verifySignature(msg string, publicKeyStr string, signatureStr string) (bool, error) {

	publicKey, parseErr := hex.DecodeString(publicKeyStr)
	if parseErr!=nil {
		return false, parseErr
	}
	signature, parseErr := hex.DecodeString(signatureStr)
	if parseErr!=nil {
		return false, parseErr
	}

	X, Y := elliptic.UnmarshalCompressed(pubkeyCurve, publicKey)
	if (X==nil || Y==nil) {
		return false, errors.New("failed to parse Public Key")
	}
	publicKeyStruct := new(ecdsa.PublicKey)
	publicKeyStruct.Curve = pubkeyCurve
	publicKeyStruct.X = X
	publicKeyStruct.Y = Y //перевели строку с publicKey в необходимый формат

	signHash := sha256.Sum256([]byte(msg))
	valid := ecdsa.VerifyASN1(publicKeyStruct, signHash[:], signature)

	return valid, nil
}

func getStringsFromFile(fileName string) ([]string, error) {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	fileData := strings.Split(string(fileContent), "\n")
	fileData = fileData[:len(fileData)-1]
	return fileData, nil
}

func writeToFile (data *[]string) (error) {
	file, err := os.Create(txtDB) //обновить существующий файл или создать новый
	if err != nil {
		return err
	}
	defer file.Close()
	for _, val := range *data {
		_, _ = file.WriteString(val + "\n")
	}
	return nil
}

func updateIPNSFile (fileName string, userName string, newEntry string) (string, error) { //добавить запись в файл

	dbStrings := make([]string, 0)
	change := "added"

	_, fileErr := os.Stat(fileName)
	if !os.IsNotExist(fileErr) { //если файл существует, ищем старую запись пользователя
		dbStrings, fileErr = getStringsFromFile(txtDB)
		if fileErr!=nil {
			return "", fileErr
		}
		dbStringsLen := len(dbStrings)
		i := sort.Search(dbStringsLen, func(i int) bool { return (dbStrings[i]>=userName+":")})
		if i<dbStringsLen && strings.HasPrefix(dbStrings[i], userName+":") {
			if dbStrings[i]==newEntry { //если старая запись найдена и совпадает с новой, оставляем как есть
				return "", errors.New("Entry already exists")
			}
			dbStringsLen--
			dbStrings[i] = dbStrings[dbStringsLen] //если не совпадает, удаляем старую запись из файла
			dbStrings = dbStrings[:dbStringsLen]
			change = "updated"
		}
	}
	dbStrings = append(dbStrings, newEntry) //добавляем новую запись
	sort.SliceStable(dbStrings, func(i, j int) bool {return dbStrings[i]<dbStrings[j]})

	fileErr = writeToFile(&dbStrings)
	if fileErr!=nil {
		return "", fileErr
	}
	return change, nil
}

func getLinkByUsername (userName string) (string, error) { //найти в файле ipfs-строку для указанного uid

	dbStrings, fileErr := getStringsFromFile(txtDB)
	if fileErr != nil {
		return "", fileErr
	}
	dbStringsLen := len(dbStrings)
	i := sort.Search(dbStringsLen, func(i int) bool { return (dbStrings[i] >= userName+"|") })
	if i >= dbStringsLen || !strings.HasPrefix(dbStrings[i], userName+"|") {
		return "", errors.New("Entry not found")
	}

	words := strings.Split(dbStrings[i], "|")
	if len(words)!=2 || words[0]=="" || words[1]=="" {
		return "", errors.New("Error: incorrect entry")
	}
	link := words[1]
	return link, nil
}


func setRecord (uid string, ipfs string, sign string) {

	words := strings.Split(uid, ":") //uid = userName + publicKey
	if len(words)!=2 || words[0]=="" || words[1]=="" {
		fmt.Println("Error: incorrect uid")
		return
	}
	userName := words[0]
	publicKey := words[1]

	verified, err := verifySignature(ipfs, publicKey, sign) //проверка подписи
	if err!=nil {
		fmt.Println("Error: failed to verify the signature: ", err)
		return
	}
	if !verified {
		fmt.Println("Wrong signature")
		return
	}

	newEntry := uid+"|"+ipfs //обновление записи в файле
	change, fileErr := updateIPNSFile(txtDB, userName, newEntry)
	if fileErr!=nil {
		fmt.Println(fileErr)
		return
	}
	fmt.Printf("Entry %s successfully\n", change)
}

func getRecord (uid string) {

	ipfsLink, err := getLinkByUsername (uid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("link =", ipfsLink)

	sh := shell.NewShell(ipfsUrl)
	shID, shellErr := sh.ID()
	if shellErr!=nil || shID == nil {
		fmt.Printf("Connection failed\nMake sure your local node is running on %s", ipfsUrl)
		return
	}
	data, err := sh.BlockGet(ipfsLink)
	if err != nil {
		fmt.Println("Failed to get file: ", err)
		return
	}
	myStr := string(data[6:len(data)-2]) //убираем ненужные символы в начале данных
	fmt.Println(myStr)
}


func main() {

	flag.Parse()

	if *requestFlag=="record-set" && (*uidFlag!="" && *ipfsFlag!="" && *signFlag!="" ) {
		setRecord(*uidFlag, *ipfsFlag, *signFlag)
		return
	}
	if *requestFlag=="record-get" && *uidFlag!="" {
		getRecord(*uidFlag)
		return
	}
	if *requestFlag=="sign" && (*nameFlag!="" && *ipfsFlag!="") {
		signMessage(*nameFlag, *ipfsFlag)
		return
	}
	fmt.Println("Flags error!")
	return
}
