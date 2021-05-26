package echox

const (
	// HmacWithSHA224 对称加密算法从这里开始
	HmacWithSHA224     Algorithm = "hmac-sha224"
	HmacWithSHA256     Algorithm = "hmac-sha256"
	HmacWithSHA384     Algorithm = "hmac-sha384"
	HmacWithSHA512     Algorithm = "hmac-sha512"
	HmacWithRipemd160  Algorithm = "hmac-ripemd160"
	HmacWithSHA3224    Algorithm = "hmac-sha3-224"
	HmacWithSHA3256    Algorithm = "hmac-sha3-256"
	HmacWithSHA3384    Algorithm = "hmac-sha3-384"
	HmacWithSHA3512    Algorithm = "hmac-sha3-512"
	HmacWithSHA512224  Algorithm = "hmac-sha512-224"
	HmacWithSHA512256  Algorithm = "hmac-sha512-256"
	HmacWithBlake2s256 Algorithm = "hmac-blake2s-256"
	HmacWithBlake2b256 Algorithm = "hmac-blake2b-256"
	HmacWithBlake2b384 Algorithm = "hmac-blake2b-384"
	HmacWithBlake2b512 Algorithm = "hmac-blake2b-512"
	Blake2sWith256     Algorithm = "blake2s-256"
	Blake2bWith256     Algorithm = "blake2b-256"
	Blake2bWith384     Algorithm = "blake2b-384"
	Blake2bWith512     Algorithm = "blake2b-512"
	// RsaWithSHA224 RAS非对称加密算法从这里开始
	RsaWithSHA224    Algorithm = "rsa-sha224"
	RsaWithSHA256    Algorithm = "rsa-sha256"
	RsaWithSHA384    Algorithm = "rsa-sha384"
	RsaWithSHA512    Algorithm = "rsa-sha512"
	RsaWithRipemd160 Algorithm = "rsa-ripemd160"
	// EcdsaWithSHA224 ECDSA非对称加密算法从这里开始
	EcdsaWithSHA224    Algorithm = "ecdsa-sha224"
	EcdsaWithSHA256    Algorithm = "ecdsa-sha256"
	EcdsaWithSHA384    Algorithm = "ecdsa-sha384"
	EcdsaWithSHA512    Algorithm = "ecdsa-sha512"
	EcdsaWithRipemd160 Algorithm = "ecdsa-ripemd160"
)

// Algorithm 签名算法
type Algorithm string
