package crypto

type Cryptor interface {
	EncryptedCode(code string) (enCode string, err error)
	DecryptedCode(enCode string) (code string, err error)
}
