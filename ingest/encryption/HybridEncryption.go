package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"ingest/logger"
	"io/ioutil"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//var log = logger.BlockchainLog

func InitatePubicPrivateKey() (*rsa.PrivateKey, *rsa.PublicKey) {
	var log = logger.BlockchainLog

	privateKeyMem, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("Cannot Generate Key")
		log.Println(err)
		os.Exit(1)
	}
	publicKeyMem := &privateKeyMem.PublicKey

	return privateKeyMem, publicKeyMem
}

func GenerateAESKey() []byte {
	mathrand.Seed(time.Now().UnixNano())
	aeskey := make([]byte, 32)
	rand.Read(aeskey)
	return aeskey
}
func EncryptKeyByPublicKey(encryptKey *rsa.PublicKey, aeskey []byte) []byte {
	var log = logger.BlockchainLog

	EncyrptedMsg, err := rsa.EncryptPKCS1v15(rand.Reader, encryptKey, aeskey)
	if err != nil {
		log.Println("Encrypting AES Key Encounter error")
		log.Println(err)
		os.Exit(1)
	}
	// fmt.Println("Original Key:", aeskey)
	// fmt.Println("Encrypted Key:", encryptKey)
	return EncyrptedMsg
}

func GetCipherAndNonce(encryptedaeskey []byte) (cipher.Block, []byte) {
	var log = logger.BlockchainLog

	cipherblock, err := aes.NewCipher(encryptedaeskey)
	if err != nil {
		log.Println("Creating CipherBlock by AES key Encounter error")
		log.Println(err)
		os.Exit(1)
	}
	aead, err := cipher.NewGCM(cipherblock)
	if err != nil {
		log.Println("Creating AEAD from AES Key Encounter error")
		log.Println(err)
		os.Exit(1)
	}
	nonce := make([]byte, aead.NonceSize())
	// nonce := []byte("OwenOwenOwen")
	// fmt.Println("Nonce Size :", aead.NonceSize())
	// fmt.Println("Nonce :", nonce)
	return cipherblock, nonce
}

func EncryptDataByAESKey(cb cipher.Block, nonce []byte, message []byte) []byte {
	var log = logger.BlockchainLog

	aead, err := cipher.NewGCM(cb)
	if err != nil {
		log.Println("Creating AEAD During Data Encryption Encounter error")
		log.Println(err)
		os.Exit(1)
	}
	encryptedmsg := aead.Seal(nil, nonce, message, nil)
	return encryptedmsg
}

func DecryptKeyByPrivateKey(private_key *rsa.PrivateKey, encrypted_key []byte) []byte {
	var log = logger.BlockchainLog

	decyrptedkey, err := rsa.DecryptPKCS1v15(rand.Reader, private_key, encrypted_key)
	if err != nil {
		log.Println("Decrypting Key Encounter error")
		log.Println(err)
	}
	// fmt.Println("Encrypted Key:", encrypted_key)
	// fmt.Println("Decrypted Key:", decyrptedkey)
	return decyrptedkey
}

func DecryptDataByAESKey(cb cipher.Block, nonce []byte, encyrpted_msg []byte) []byte {
	var log = logger.BlockchainLog

	decrypt_aead, err := cipher.NewGCM(cb)
	if err != nil {
		log.Println("Creating AEAD During Data Decryption Encounter error")
		log.Println(err)
		os.Exit(1)
	}
	decryptmsg, err := decrypt_aead.Open(nil, nonce, encyrpted_msg, nil)
	if err != nil {
		log.Println("Decrypting Data Encounter error")
		log.Println(err)
		os.Exit(1)
	}
	return decryptmsg
}

func CreateKeyFiles(publickey *rsa.PublicKey, privatekey *rsa.PrivateKey, encryptionkey []byte, destname string) (string, string, string) {
	var log = logger.BlockchainLog

	keypath := filepath.Join("..", "key")
	err := os.MkdirAll(keypath, os.ModePerm)
	if err != nil {
		log.Println("Creating Keypath encounters error")
		log.Println(err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privatekey)
	log.Println("Private Key : ", privateKeyBytes)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	Private_Key_Name := destname + "_Private.pem"
	privatePem, err := os.OpenFile(filepath.Join(keypath, Private_Key_Name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error when create %s: %s \n", Private_Key_Name, err)
		os.Exit(1)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		log.Printf("error when encode %s: %s \n", Private_Key_Name, err)
		os.Exit(1)
	}

	privatePem.Close()

	publicKeyBytes := x509.MarshalPKCS1PublicKey(publickey)
	if err != nil {
		log.Printf("error when dumping publickey: %s \n", err)
		os.Exit(1)
	}
	// fmt.Println("Public Key : ", publicKeyBytes)

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	Public_Key_Name := destname + "_Public.pem"
	publicPem, err := os.OpenFile(filepath.Join(keypath, Public_Key_Name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error when create %s: %s \n", Public_Key_Name, err)
		os.Exit(1)
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		log.Printf("error when encode %s: %s \n", Public_Key_Name, err)
		log.Println(err)
		os.Exit(1)
	}

	publicPem.Close()

	encryption_key_name := destname + "_AES.key"

	aes_block := &pem.Block{
		Type:  "AES KEY",
		Bytes: encryptionkey,
	}
	err = ioutil.WriteFile(filepath.Join(keypath, encryption_key_name), pem.EncodeToMemory(aes_block), 0644)
	if err != nil {
		log.Println("Writing Encryption Key encounters error")
		log.Println(err)
		os.Exit(1)
	}
	return Private_Key_Name, Public_Key_Name, encryption_key_name
}

func ReadAESKey(keyFile string) []byte {
	var log = logger.BlockchainLog

	encrypt_key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Println("Opening AES file encounters error")
		log.Println(err)
	}
	log.Println("The AES key path is :")
	log.Println(keyFile)
	encrypt_key_block, _ := pem.Decode(encrypt_key)
	return encrypt_key_block.Bytes
}

func ReadPrivateKey(keyFile string) *rsa.PrivateKey {
	var log = logger.BlockchainLog

	private_key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Println("Opening RSA Private file encounters error")
		log.Println(err)
		os.Exit(1)
	}

	block, _ := pem.Decode(private_key)
	if block == nil {
		log.Println("Nil Block")
	}

	private_key_reader, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("Parsing RSA Private Key file encounters error")
		log.Println(err)
		os.Exit(1)
	}
	return private_key_reader
}

func ReadPublicKey(keyFile string) *rsa.PublicKey {
	var log = logger.BlockchainLog

	public_key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Println("Opening RSA Public file encounters error")
		log.Println(err)
		os.Exit(1)
	}

	block, _ := pem.Decode(public_key)
	if block == nil {
		log.Println("Nil Block")
	}

	public_key_reader, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		log.Println("Parsing RSA Public Key file encounters error")
		log.Println(err)
		os.Exit(1)
	}
	return public_key_reader
}

func CheckKeyExist(name string) bool {
	var log = logger.BlockchainLog
	keypath := filepath.Join("..", "key")
	err := os.MkdirAll(keypath, os.ModePerm)
	if err != nil {
		log.Println("Creating Keypath encounters error")
		log.Println(err)
	}
	dir_files, err := ioutil.ReadDir(filepath.Clean(keypath))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir_files {
		if strings.Contains(file.Name(), name) && strings.Contains(file.Name(), ".key") {
			return true
		}
	}
	return false
}

func UnitTest_HybridEncryption(message []byte, KeyName string) {
	var log = logger.BlockchainLog

	privatekey, publickey := InitatePubicPrivateKey()
	aes_key := GenerateAESKey()
	log.Println("Original Message: ", message)
	cb, nonce := GetCipherAndNonce(aes_key)
	encryptedmsg := EncryptDataByAESKey(cb, nonce, message)
	encrypted_key := EncryptKeyByPublicKey(publickey, aes_key)
	log.Println("Original Message:", message)
	log.Println("Encrypted Message:", encryptedmsg)

	decrypted_key := DecryptKeyByPrivateKey(privatekey, encrypted_key)
	decrypt_cb, decrypt_nonce := GetCipherAndNonce(decrypted_key)
	decryptedmsg := DecryptDataByAESKey(decrypt_cb, decrypt_nonce, encryptedmsg)
	log.Println("Encrypted Message: ", encryptedmsg)
	log.Println("Descrypted Message: ", decryptedmsg)

	pri_name, pub_name, encrypt_name := CreateKeyFiles(publickey, privatekey, encrypted_key, KeyName)
	log.Println("Private Key Name:", pri_name)
	log.Println("Public Key Name:", pub_name)
	log.Println("Encrypt Key Name:", encrypt_name)
}
