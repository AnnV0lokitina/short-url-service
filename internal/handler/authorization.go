package handler

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
)

var (
	secretKey  = []byte("secret key")
	idLength   = 4
	cookieName = "login"
)

func authorization(w http.ResponseWriter, r *http.Request) (uint32, error) {
	var userID uint32
	login, err := getLoginFromCookie(r)
	if err == nil {
		userID, err = getIDFromLogin(login)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}
	userID, err = generateUserID()
	if err != nil {
		return 0, err
	}
	login = generateLogin(userID)
	setLoginToCookie(w, login)
	return userID, nil
}

func getLoginFromCookie(request *http.Request) (string, error) {
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func setLoginToCookie(w http.ResponseWriter, login string) {
	cookie := &http.Cookie{Name: cookieName, Value: login, HttpOnly: false}
	http.SetCookie(w, cookie)
}

func generateUserID() (uint32, error) {
	b := make([]byte, idLength)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

func idToByte(id uint32) []byte {
	d := make([]byte, idLength)
	binary.BigEndian.PutUint32(d[0:], id)
	return d
}

func generateLogin(id uint32) string {
	b := idToByte(id)
	signature := createSignature(b)
	b = append(b, signature...)
	return hex.EncodeToString(b)
}

func createSignature(src []byte) []byte {
	key := secretKey //[:16]
	h := hmac.New(sha256.New, key)
	h.Write(src)
	return h.Sum(nil)
}

func getIDFromLogin(encodedLogin string) (uint32, error) {
	data, err := hex.DecodeString(encodedLogin)
	if err != nil {
		return 0, err
	}
	id := binary.BigEndian.Uint32(data[:idLength])
	h := hmac.New(sha256.New, secretKey)
	h.Write(data[:idLength])
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[idLength:]) {
		return id, nil
	}
	return 0, nil
}
