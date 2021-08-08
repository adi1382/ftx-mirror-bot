//+build windows

package wmic

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"golang.org/x/sys/windows/registry"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// machineID returns the key MachineGuid in registry `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Cryptography`.
// If there is an error running the commad an empty string is returned.

type diskDrive struct {
	Caption      string
	DeviceID     string
	SerialNumber string
	Model        string
	Partitions   uint
	Size         uint64
}

// build a command line for wmic command and format as csv output
func buildCommand(cmd interface{}) ([]string, error) {
	cmdString := make([]string, 0, 0)
	s := reflect.Indirect(reflect.ValueOf(cmd))
	t := s.Type()
	if s.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("unknown interface")
	}

	cmdString = append(cmdString, t.Name())
	cmdString = append(cmdString, "get")

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	cmdString = append(cmdString, strings.Join(fields, ","))
	cmdString = append(cmdString, "/format:csv")

	return cmdString, nil
}

//Execute the wmic command and return the stdout/stderr
func runCmd(dst interface{}) (string, error) {
	cmdLineOpt, _ := buildCommand(dst)

	run := exec.Command("wmic", cmdLineOpt...)

	var stdout, stderr bytes.Buffer
	run.Stdout = &stdout
	run.Stderr = &stderr

	err := run.Run()
	if err != nil {
		return string(stderr.Bytes()), err
	} else {
		return string(stdout.Bytes()), err
	}
}

//Parse the csv format output of the runCmd
func parseResult(stdout string, dst interface{}) error {
	dv := reflect.ValueOf(dst).Elem()
	t := dv.Type().Elem()

	dv.Set(reflect.MakeSlice(dv.Type(), 0, 0))

	lines := strings.Split(stdout, "\n")
	var header []int = nil

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			v := reflect.New(t)
			r := csv.NewReader(strings.NewReader(line))

			r.FieldsPerRecord = t.NumField() + 1

			records, err := r.ReadAll()
			if err != nil {
				return err
			}
			//Find the field number of the record
			if header == nil {
				header = make([]int, len(records[0]), len(records[0]))
				for i, record := range records[0] {
					for j := 0; j < t.NumField(); j++ {
						if record == t.Field(j).Name {
							header[i] = j
						}
					}
				}
				continue
			} else {
				for i, record := range records[0] {
					f := reflect.Indirect(v).Field(header[i])
					switch t.Field(header[i]).Type.Kind() {
					case reflect.String:
						f.SetString(record)
					case reflect.Uint, reflect.Uint64:
						uintVal, err := strconv.ParseUint(record, 10, 64)
						if err != nil {
							return err
						}
						f.SetUint(uintVal)
					case reflect.Bool:
						bVal, err := strconv.ParseBool(record)
						if err != nil {
							return err
						}
						f.SetBool(bVal)
					default:
						return errors.New("unknown data type")
					}
				}

			}
			dv.Set(reflect.Append(dv, reflect.Indirect(v)))
		}
	}

	return nil
}

func getDiskDriveInfo() ([]diskDrive, error) {
	var disk []diskDrive

	output, err := runCmd(disk)

	if err != nil {
		return nil, err
	}

	err = parseResult(output, &disk)

	return disk, err
}

func machineID() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = k.Close()
	}()

	s, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "", err
	}
	return s, nil
}

func protect(appID, id string) string {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appID))
	hash := sha256.Sum256(mac.Sum(nil))
	return hex.EncodeToString(hash[:])
}

func GetHashedKey() string {
	//privateKey := ""

	diskSerialNumbers := make([]string, 0, 2)

	drives, err := getDiskDriveInfo()
	if err != nil {
		panic(err)
	}

	for i := range drives {
		diskSerialNumbers = append(diskSerialNumbers, strings.TrimSpace(drives[i].SerialNumber))
	}
	sort.Strings(diskSerialNumbers)
	privateKey, err := machineID()

	if err != nil {
		panic(err)
	}

	for i := range diskSerialNumbers {
		privateKey += diskSerialNumbers[i]
	}

	return protect("ef393ae49fb4a0426955a0ed47ab352c20d8c355cc4c00be30e15a1e102e7f40", protect("MirrorTrader", privateKey))
}
