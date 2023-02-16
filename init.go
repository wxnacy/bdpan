package main

func initAll() {
	err = initConfigDir()
	if err != nil {
		panic(err)
	}
	err = initCryptoKey()
	if err != nil {
		panic(err)
	}
}
