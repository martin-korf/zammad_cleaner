////CGO_ENABLED=0 go build -o ./zammad_cleaner main.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Ticket struct {
	ID      int64  `json:"id"`
	CloseAt string `json:"close_at"`
}

type Assets struct {
	Tickets map[string]Ticket `json:"Ticket"`
}

type TicketResponse struct {
	Assets Assets `json:"assets"`
}

func main() {

	var (
		apiTokenFlag string
		baseURLFlag  string
		cutOffDate   string
	)

	flag.StringVar(&apiTokenFlag, "token", "", "API token")
	flag.StringVar(&baseURLFlag, "base-url", "", "Base URL")
	flag.StringVar(&cutOffDate, "cut-off-date", "2021-07-01", "Cut-off date for ticket deletion (format: YYYY-MM-DD)")

	flag.Parse()

	if apiTokenFlag == "" || baseURLFlag == "" {
		fmt.Println("Please provide both API token and base URL.")
		return
	}

	cutOffDateTime, err := time.Parse("2006-01-02", cutOffDate)
	if err != nil {
		fmt.Println("Error parsing cut-off date:", err)
		return
	}

	apiURL := fmt.Sprintf("%s/api/v1/tickets/search?query=state_id:4&sort_by=id&order_by=asc", baseURLFlag)
	baseURL := baseURLFlag
	apiToken := fmt.Sprintf("token=%s", apiTokenFlag)

	client := &http.Client{}
	for {
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Authorization", "Token "+apiToken)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error getting response:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var ticketResponse TicketResponse
		err = json.Unmarshal(body, &ticketResponse)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		processedAnyTicket := false
		for _, ticket := range ticketResponse.Assets.Tickets {
			lastCloseDate, err := time.Parse("2006-01-02T15:04:05.999Z", ticket.CloseAt)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				continue
			}

			if lastCloseDate.Before(cutOffDateTime) {
				processedAnyTicket = true
				deleteURL := fmt.Sprintf("%s/api/v1/tickets/%d", baseURL, ticket.ID)
				deleteReq, err := http.NewRequest("DELETE", deleteURL, nil)
				if err != nil {
					fmt.Println("Error creating delete request:", err)
					continue
				}
				deleteReq.Header.Set("Authorization", "Token "+apiToken)
				deleteResp, err := client.Do(deleteReq)
				if err != nil {
					fmt.Println("Error sending delete request:", err)
					continue
				}
				defer deleteResp.Body.Close()

				if deleteResp.StatusCode == http.StatusOK {
					fmt.Printf("Ticket with ID %d was successfully deleted at %s.\n", ticket.ID, ticket.CloseAt)
				} else {
					fmt.Printf("Error deleting ticket with ID %d, status code: %d\n", ticket.ID, deleteResp.StatusCode)
				}
			}
		}

		if !processedAnyTicket {
			break
		}
	}
}
