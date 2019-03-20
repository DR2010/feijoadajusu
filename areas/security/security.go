package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"feijoadajusu/areas/commonstruct"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// LoginUserV2 but not
func LoginUserV2(sysid string, redisclient *redis.Client, userid string, password string) commonstruct.Credentials {

	mongodbvar := new(commonstruct.DatabaseX)

	// Updated on 008-Oct-2018
	// mongodbvar.APIServer, _ = redisclient.Get("MSAPImainIPAddress").Result()
	mongodbvar.APIServer, _ = redisclient.Get(sysid + "SecurityMicroserviceURL").Result()

	apiURL := mongodbvar.APIServer
	resource := "/securitylogin"

	data := url.Values{}
	data.Add("userid", userid)
	data.Add("password", password)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()

	var response commonstruct.Credentials

	if err := json.NewDecoder(resp2.Body).Decode(&response); err != nil {
		log.Println(err)
	}

	if resp2.Status == "200 OK" {
		return response
	}

	response.ApplicationID = "None"
	response.JWT = "Error"

	return response
}

// LoginUser something
func LoginUser(sysid string, redisclient *redis.Client, userid string, password string) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	// mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	// Updated on 02-Sep-2018
	// mongodbvar.APIServer, _ = redisclient.Get("MSAPImainIPAddress").Result()

	// Updated on 10-Oct-2018
	mongodbvar.APIServer, _ = redisclient.Get(sysid + "SecurityMicroserviceURL").Result()

	apiURL := mongodbvar.APIServer
	resource := "/securitylogin"

	data := url.Values{}
	data.Add("userid", userid)
	data.Add("password", password)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()

	var response string

	if err := json.NewDecoder(resp2.Body).Decode(&response); err != nil {
		log.Println(err)
	}

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		emptydisplay.ErrorCode = "200 OK"
		emptydisplay.ErrorDescription = "200 OK"
		emptydisplay.ReturnedValue = response

	} else {
		emptydisplay.IsSuccessful = "N"
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "404 Shit happens!... and it happened!"

	}

	return emptydisplay
}

// SignUp function
func SignUp(sysid string, redisclient *redis.Client, userid string, preferredname string, password string, passwordvalidate string, applicationid string) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	// mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	// Updated on 02-Sep-2018
	// mongodbvar.APIServer, _ = redisclient.Get("MSAPImainIPAddress").Result()

	// Updated on 10-Oct-2018
	mongodbvar.APIServer, _ = redisclient.Get(sysid + "SecurityMicroserviceURL").Result()

	var emptydisplay commonstruct.Resultado

	apiURL := mongodbvar.APIServer
	resource := "/securitysignup"

	if userid == "" {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "user id not suppplied"
		return emptydisplay
	}

	if password == "" {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "password not suppplied"
		return emptydisplay
	}

	if password != passwordvalidate {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "Password mismatch"
		return emptydisplay
	}

	var passwordhashed = Hashstring(password)
	var passwordvalidatehashed = Hashstring(passwordvalidate)

	// passwordhashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// passwordvalidatehashed, _ := bcrypt.GenerateFromPassword([]byte(passwordvalidate), bcrypt.DefaultCost)

	// passwordhasheds := string(passwordhashed)
	// passwordvalidatehasheds := string(passwordvalidatehashed)

	data := url.Values{}
	data.Add("userid", userid)
	data.Add("preferredname", preferredname)
	data.Add("password", passwordhashed)
	data.Add("passwordvalidate", passwordvalidatehashed)
	data.Add("applicationid", applicationid)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())

	// Call method here
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay

}

// GetResetPasswordCode function
func GetResetPasswordCode(sysid string, redisclient *redis.Client, bodybyte []byte) commonstruct.Resultado {

	envirvar := new(commonstruct.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get(sysid + "SecurityMicroserviceURL").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/requestcode"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		var resultado = resp2.Body
		log.Println(resultado)

	}

	var hr commonstruct.Resultado
	hr.ErrorCode = "00"
	hr.ErrorDescription = "All good"

	return hr

}

// ApplyNewPassword function
func ApplyNewPassword(sysid string, redisclient *redis.Client, bodybyte []byte) commonstruct.Resultado {

	envirvar := new(commonstruct.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get(sysid + "SecurityMicroserviceURL").Result()

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/changepassword"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()

	resultado := ""

	if err := json.NewDecoder(resp2.Body).Decode(&resultado); err != nil {
		log.Println(err)
	}

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	log.Println(resultado)

	var hr commonstruct.Resultado
	hr.ErrorCode = "00"
	hr.ErrorDescription = resultado

	return hr

}

// UpdateUserRoles function
func UpdateUserRoles(sysid string, redisclient *redis.Client, bodybyte []byte) commonstruct.Resultado {

	envirvar := new(commonstruct.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get(sysid + "SecurityMicroserviceURL").Result()

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/updateuserdetails"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()

	resultado := ""

	if err := json.NewDecoder(resp2.Body).Decode(&resultado); err != nil {
		log.Println(err)
	}

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	log.Println(resultado)

	var hr commonstruct.Resultado
	hr.ErrorCode = "00"
	hr.ErrorDescription = resultado

	return hr

}

// ResetPassword function
func ResetPassword(sysid string, redisclient *redis.Client, userid string, preferredname string, password string, passwordvalidate string, applicationid string) commonstruct.Resultado {

	mongodbvar := new(commonstruct.DatabaseX)

	// Updated on 10-Oct-2018
	mongodbvar.APIServer, _ = redisclient.Get(sysid + "SecurityMicroserviceURL").Result()

	var emptydisplay commonstruct.Resultado

	apiURL := mongodbvar.APIServer
	resource := "/securityresetpassword"

	if userid == "" {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "user id not suppplied"
		return emptydisplay
	}

	if password == "" {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "password not suppplied"
		return emptydisplay
	}

	if password != passwordvalidate {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "Password mismatch"
		return emptydisplay
	}

	var passwordhashed = Hashstring(password)
	var passwordvalidatehashed = Hashstring(passwordvalidate)

	data := url.Values{}
	data.Add("userid", userid)
	data.Add("password", passwordhashed)
	data.Add("passwordvalidate", passwordvalidatehashed)
	data.Add("applicationid", applicationid)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())

	// Call method here
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay

}

// ValidateToken is half way
func ValidateToken(redisclient *redis.Client, httprequest *http.Request) string {

	var credtemp commonstruct.Credentials

	cookie, _ := httprequest.Cookie("DanBTCjwt")
	if cookie == nil {
		return "NotOkToLogin"
	}

	cookieinbytes := []byte(cookie.Value)
	_ = json.Unmarshal(cookieinbytes, &credtemp)

	var key = credtemp.KeyJWT

	tokenstored, _ := redisclient.Get(key).Result()

	var ret = "NotOkToLogin"
	if tokenstored == credtemp.JWT {
		ret = "OkToLogin"
	}

	return ret
}

// ValidateTokenV2 will get info from cache
func ValidateTokenV2(redisclient *redis.Client, httprequest *http.Request) (string, commonstruct.Credentials) {
	var credentials commonstruct.Credentials

	credentials.ApplicationID = "Restaurante"
	credentials.UserID = "Anonymous"
	credentials.Name = "Anonymous"
	credentials.IsAdmin = "No"
	credentials.IsAnonymous = "Yes"

	// The system will store an object in cache and the key must be the used ID
	// The same user can logon in 2 places, I think
	// Users can't be mixed, I can't trust the variables since it is completely stateless - each request is stateless

	// Machine credentials
	//
	clientsecret := httprequest.FormValue("macdantoken")

	if clientsecret != "" {
		// Issue keys - should be stored in the database API Key or Secret I think
		//
		if clientsecret == "BypassSecurity" {
			var credentialsmachine commonstruct.Credentials
			credentialsmachine.ApplicationID = "Restaurante"
			credentialsmachine.UserID = "Machine"
			credentialsmachine.JWT = clientsecret
			return "OkToLogin", credentialsmachine
		}
	}

	credentials.JWT = "Error"

	jwtincookie := ""
	useridincookie := ""

	cookiekeyJWT := "DanBTCjwt"
	cookiekeyUSERID := "DanBTCuserid"

	cookieJWT, err := httprequest.Cookie(cookiekeyJWT)

	if err != nil {
		log.Println(err)
		log.Println("Not found Cookie: " + cookiekeyJWT)
	}

	if cookieJWT == nil {
		return "NotOkToLogin", credentials
	}

	cookieUSERID, err2 := httprequest.Cookie(cookiekeyUSERID)
	if err2 != nil {
		log.Println(err2)
		log.Println("Not found Cookie: " + cookiekeyUSERID)
	}

	if cookieUSERID == nil {
		return "NotOkToLogin", credentials
	}

	jwtincookie = cookieJWT.Value
	useridincookie = cookieUSERID.Value

	var keyredis = cookiekeyJWT + useridincookie

	tokenstored, _ := redisclient.Get(keyredis).Result()
	tokenstoredbytes := []byte(tokenstored)

	_ = json.Unmarshal(tokenstoredbytes, &credentials)

	var ret = "NotOkToLogin"
	if credentials.JWT == jwtincookie {
		credentials.IsAnonymous = "No"
		ret = "OkToLogin"
	} else {
		credentials.ApplicationID = "Restaurante"
		credentials.UserID = "Anonymous"
		credentials.Name = "Anonymous"
		credentials.IsAdmin = "No"
	}

	return ret, credentials
}

// Hashstring is just for hashing - only reference key
func Hashstring(str string) string {

	s := str
	h := sha1.New()
	h.Write([]byte(s))

	sha1hash := hex.EncodeToString(h.Sum(nil))

	return sha1hash
}

// GetUserDetails function
func GetUserDetails(sysid string, redisclient *redis.Client, bodybyte []byte) commonstruct.Credentials {

	envirvar := new(commonstruct.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get(sysid + "SecurityMicroserviceURL").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/getuserdetails"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay commonstruct.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		var resultado = resp2.Body
		log.Println(resultado)

	}

	var usercred commonstruct.Credentials
	if err := json.NewDecoder(resp2.Body).Decode(&usercred); err != nil {
		log.Println(err)
	} else {
		log.Println(usercred.ApplicationID)
		log.Println(usercred.IsAdmin)
		log.Println(usercred.UserID)
		log.Println(usercred.Status)
	}

	return usercred

}

func keyfortheday(day int) string {

	var key = "De tudo, ao meu amor serei atento antes" +
		"E com tal zelo, e sempre, e tanto" +
		"Que mesmo em face do maior encanto" +
		"Dele se encante mais meu pensamento" +
		"Quero vivê-lo em cada vão momento" +
		"E em seu louvor hei de espalhar meu canto" +
		"E rir meu riso e derramar meu pranto" +
		"Ao seu pesar ou seu contentamento" +
		"E assim quando mais tarde me procure" +
		"Quem sabe a morte, angústia de quem vive" +
		"Quem sabe a solidão, fim de quem ama" +
		"Eu possa lhe dizer do amor (que tive):" +
		"Que não seja imortal, posto que é chama" +
		"Mas que seja infinito enquanto dure"

	stringSlice := strings.Split(key, " ")
	var stringSliceFinal []string

	x := 0
	for i := 0; i < len(stringSlice); i++ {
		if len(stringSlice[0]) > 3 {
			stringSliceFinal[x] = stringSlice[i]
			x++
		}
	}

	return stringSliceFinal[day]
}

func getjwtfortoday() string {

	_, _, day := time.Now().Date()

	s := keyfortheday(day)
	h := sha1.New()
	h.Write([]byte(s))

	sha1hash := hex.EncodeToString(h.Sum(nil))

	return sha1hash
}

// Encrypt string to base64 crypto using AES
func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// Decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
