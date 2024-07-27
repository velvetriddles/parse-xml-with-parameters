package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// тут писать SearchServer

var FileDataset = "dataset.xml"

type Dataset struct {
	Rows []Row `xml:"row"`
}

type Row struct {
	ID        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type (
	Users  []User
	Values url.Values
)

const (
	OrderFieldID    = "id"
	OrderFieldAge   = "age"
	OrderFieldName  = "name"
	OrderFieldEmpty = ""

	ErrorBadLimit   = "limit invalid"
	ErrorBadOffset  = "offset invalid"
	ErrorBadOrderBy = "order_by invalid"
)

func SearchServer(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("AccessToken") != "token" {
		unauthorized(w)
		return
	}

	users, err := loadUsers()
	if err != nil {
		internalServerError(w, err.Error())
		return
	}

	limit, err := parseLimitParam(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	offset, err := parseOffsetParam(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	orderField, err := parseOrderFieldParam(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	orderBy, err := parseOrderByParam(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	query := parseQueryParam(r)

	users = sortUsers(users, orderBy, orderField)
	users = queryUsers(users, query)
	users = limitOffsetUsers(users, limit, offset)

	ok(w, users)
}

func queryUsers(users Users, query string) Users {
	unique := make(map[int]User)
	for _, user := range users {
		if len(query) > 0 {
			if len(query) > 20 {
				time.Sleep(time.Second / 10)
			}
			if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
				unique[user.ID] = user
				continue
			}
		} else {
			unique[user.ID] = user
		}
	}

	result := make(Users, 0, len(unique))
	for _, user := range users {
		u, ok := unique[user.ID]
		if !ok {
			continue
		}
		result = append(result, u)
	}

	return result
}

func sortUsers(users Users, orderBy int, orderField string) Users {
	sort.SliceStable(users, func(i, j int) bool {
		switch orderBy {
		case OrderByAsc:
			switch orderField {
			case OrderFieldID:
				return users[i].ID < users[j].ID
			case OrderFieldAge:
				return users[i].Age < users[j].Age
			case OrderFieldName:
				return users[i].Name < users[j].Name
			}
		case OrderByDesc:
			switch orderField {
			case OrderFieldID:
				return users[i].ID > users[j].ID
			case OrderFieldAge:
				return users[i].Age > users[j].Age
			case OrderFieldName:
				return users[i].Name > users[j].Name
			}
		}

		return false
	})

	return users
}

func limitOffsetUsers(users Users, limit, offset int) Users {
	if offset >= len(users) {
		offset = len(users) - 1
		if offset < 0 {
			offset = 0
		}
	}
	realLimit := offset + limit
	if realLimit > len(users) {
		realLimit = len(users)
	}

	return users[offset:realLimit]
}

func parseLimitParam(r *http.Request) (int, error) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		return 0, fmt.Errorf(ErrorBadLimit)
	}

	return limit, nil
}

func parseOffsetParam(r *http.Request) (int, error) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		return 0, fmt.Errorf(ErrorBadOffset)
	}

	return offset, nil
}

func parseOrderFieldParam(r *http.Request) (string, error) {
	orderField := strings.ToLower(r.URL.Query().Get("order_field"))
	switch orderField {
	case OrderFieldID, OrderFieldAge, OrderFieldName:
	case OrderFieldEmpty:
		orderField = OrderFieldName
	default:
		return "", fmt.Errorf(ErrorBadOrderField)
	}

	return orderField, nil
}

func parseOrderByParam(r *http.Request) (int, error) {
	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil || orderBy < -1 || orderBy > 1 {
		return 0, fmt.Errorf(ErrorBadOrderBy)
	}

	return orderBy, nil
}

func parseQueryParam(r *http.Request) string {
	return r.URL.Query().Get("query")
}

func loadUsers() (Users, error) {
	file, err := os.Open(FileDataset)
	if err != nil {
		return nil, fmt.Errorf("open file error")
	}
	defer file.Close()

	var dataset Dataset
	if err := xml.NewDecoder(file).Decode(&dataset); err != nil {
		return nil, err
	}

	users := make(Users, 0, len(dataset.Rows))
	for _, row := range dataset.Rows {
		users = append(users, User{
			ID:     row.ID,
			Name:   fmt.Sprintf("%s %s", row.FirstName, row.LastName),
			Age:    row.Age,
			About:  row.About,
			Gender: row.Gender,
		})
	}

	return users, nil
}

func internalServerError(w http.ResponseWriter, desc string) {
	resp, err := json.Marshal(SearchErrorResponse{Error: desc})
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err = w.Write(resp); err != nil {
		internalServerError(w, err.Error())
		return
	}
}

func badRequest(w http.ResponseWriter, desc string) {
	resp, err := json.Marshal(SearchErrorResponse{Error: desc})
	if err != nil {
		internalServerError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if _, err = w.Write(resp); err != nil {
		internalServerError(w, err.Error())
		return
	}
}

func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func ok(w http.ResponseWriter, data interface{}) {
	resp, err := json.Marshal(data)
	if err != nil {
		internalServerError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(resp); err != nil {
		internalServerError(w, err.Error())
		return
	}
}
