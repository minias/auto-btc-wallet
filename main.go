package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cosmos/btcutil/base58"
)

const version = byte(0x00)
const walletFile = "btc-wallet.json"

// Wallet
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// Wallets stores bunch of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// Get wallet address
func (w Wallet) GetAddress() string {
	publicKeyHash := HashPublicKey(w.PublicKey)

	return base58.CheckEncode(publicKeyHash, version)
}

// Hash public key
func HashPublicKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	Hasher := sha256.New()
	_, err := Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	public := Hasher.Sum(nil)
	return public
}

// Generate New Wallet
func NewWallet() *Wallet {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return &Wallet{*privateKey, publicKey}
}

// CreateWallet adds a Wallet to Wallets
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()
	ws.Wallets[address] = wallet
	return address
}

// Creates wallets and fills it from a file if it exists
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	if err != nil {
		return &wallets, err
	}
	return &wallets, err
}

// LoadFromFile loads wallets from the file
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(fileContent, ws)
	if err != nil {
		log.Panic(err)
	}
	return nil
}

// SaveToFile saves wallets to a file
func (ws Wallets) SaveToFile() {
	jsonData, err := json.Marshal(ws)
	if err != nil {
		log.Panic(err)
	}
	err = os.WriteFile(walletFile, jsonData, 0666)
	if err != nil {
		log.Panic(err)
	}
}

// wallet.go
func (w Wallet) MarshalJSON() ([]byte, error) {
	mapStringAny := map[string]any{
		"PrivateKey": map[string]any{
			"D": w.PrivateKey.D,
			"PublicKey": map[string]any{
				"X": w.PrivateKey.PublicKey.X,
				"Y": w.PrivateKey.PublicKey.Y,
			},
			"X": w.PrivateKey.X,
			"Y": w.PrivateKey.Y,
		},
		"PublicKey": w.PublicKey,
	}
	return json.Marshal(mapStringAny)
}

// Returns addresses stored at wallet file
func (ws *Wallets) GetAddresses() []string {
	var addrs []string
	for address := range ws.Wallets {
		addrs = append(addrs, address)
	}
	return addrs
}

// Returns a Wallet by address
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func main() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	Wallet := wallets.GetWallet(address)
	wallets.SaveToFile()
	fmt.Printf("PrivateKey:%x\n", Wallet.PublicKey)
	fmt.Printf("Adress:    %s\n", Wallet.GetAddress())
	fmt.Printf("PublicKey: %x\n", Wallet.PublicKey)
}
