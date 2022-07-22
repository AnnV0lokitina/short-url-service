package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"

	"github.com/AnnV0lokitina/short-url-service/pkg/userid"
)

var (
	secretKey  = []byte("secret key")
	cookieName = "login"
)

func authorizeUserAndSetCookie(w http.ResponseWriter, r *http.Request) (uint32, error) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		return 0, err
	}
	setUserIDCookie(w, userID)
	return userID, nil
}

func getUserIDFromRequest(r *http.Request) (userID uint32, err error) {
	login, err := getLoginFromCookie(r)
	if err == nil {
		return getIDFromLogin(login)
	}
	return userid.GenerateUserID()
}

func setUserIDCookie(w http.ResponseWriter, userID uint32) {
	login := generateLogin(userID)
	setLoginToCookie(w, login)
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

func idToByte(id uint32) []byte {
	d := make([]byte, userid.IDLength)
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
	id := binary.BigEndian.Uint32(data[:userid.IDLength])
	h := hmac.New(sha256.New, secretKey)
	h.Write(data[:userid.IDLength])
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[userid.IDLength:]) {
		return id, nil
	}
	return 0, nil
}
