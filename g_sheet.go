package main

import (
	"encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"
	"bytes"
	"strings"

        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/sheets/v4"
)

//func getClient(config *oauth2.Config) *http.Client {
func getClient() *http.Client {

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
			                "authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func get_google_sheet(page string) string {
	var buf bytes.Buffer

	client := getClient()
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadSheetId := os.Getenv("GOOGLE_SHEET_ID")
	readRange := page
	resp, err := srv.Spreadsheets.Values.Get(spreadSheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		buf.WriteString("No data found.")
	} else {
		//fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			buf.WriteString(row[0].(string))
			buf.WriteString("\t")
			buf.WriteString(row[1].(string))
			buf.WriteString("\t")
			buf.WriteString(row[2].(string))
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func add_google_sheet_row(sheet_page, row_data string) {
	log.Println("add google seet")
	ctx := context.Background()

	c := getClient()

	sheetsService, err := sheets.New(c)
	if err != nil {
		log.Fatal(err)
	}

	spreadSheetId := os.Getenv("GOOGLE_SHEET_ID")

	// How the input data should be interpreted
	valueInputOption := "RAW"

	// How the input data should be inserted
	insertDataOption := "INSERT_ROWS"

	insertRow := strings.Split(row_data, " ")
	fmt.Printf("%T\n", insertRow)
	fmt.Println(insertRow)

	var rb sheets.ValueRange

	rbValue := []interface{}{insertRow[1], insertRow[2], insertRow[3]}
	rb.Values = append(rb.Values, rbValue)

	write_resp, err := sheetsService.Spreadsheets.Values.Append(spreadSheetId, sheet_page, &rb).
		ValueInputOption(valueInputOption).
		InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Print(write_resp)
}
