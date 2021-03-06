package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}

func StructConv(args ...interface{}) interface{} {
	if len(args) < 2 {
		return nil
	}
	arg1 := reflect.ValueOf(args[0])
	if arg1.Kind() != reflect.Struct {
		return nil
	}

	type ret struct {
		key string
		val string
	}

	fields := reflect.TypeOf(arg1)
	values := reflect.ValueOf(arg1)
	num := fields.NumField()

	rets := []ret{}
	//out := make([]interface{}, num)
	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		var r ret
		var v string
		switch value.Kind() {
		case reflect.String:
			v = value.String()
		case reflect.Int:
			v = (strconv.FormatInt(value.Int(), 10))
		case reflect.Int32:
			v = strconv.FormatInt(value.Int(), 10)
		case reflect.Int64:
			v = strconv.FormatInt(value.Int(), 10)
		default:
			v = value.String()
		}
		r.key = field.Name
		r.val = v
		//fmt.Print("Type:", field.Type, ",", field.Name, "=", value, "\n")
		rets = append(rets, r)
	}
	return rets
}

func ExecCmd(client *ssh.Client, query string) (bytes.Buffer, error) {

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(query); err != nil {
		log.Fatal("Failed to run: " + err.Error())
		return b, err
	}
	//fmt.Println(b.String())
	return b, nil
}

func GetPasswd() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()
	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

func GetFilenameDate(fName, extn string) string {
	// Use layout string for time format.
	const layout = "20060102150405"
	// Place now in the string.
	t := time.Now()
	return fName + "-" + t.Format(layout) + "." + extn
}

func WriteCsvReport(filename, data string) {
	var err error
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		var createFile, err = os.Create(filename)
		if err != nil {
			fmt.Println("error:", err)
		}
		defer createFile.Close()
	}
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer file.Close()
	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("error:", err)
	}
	// Save file changes.
	err = file.Sync()
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("Report Written Successfully.")
}
