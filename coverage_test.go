package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestSearchServerParseLimitParam(t *testing.T) {
	params := url.Values{}
	params.Add("offset", strconv.Itoa(1))
	params.Add("limit", strconv.Itoa(-1))

	req, err := http.NewRequest("GET", "/?"+params.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("AccessToken", "token")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchServer)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSearchServerParseOffsetParam(t *testing.T) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(1))
	params.Add("offset", strconv.Itoa(-1))

	req, err := http.NewRequest("GET", "/?"+params.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("AccessToken", "token")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchServer)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestFindUsers(t *testing.T) {
	cases := map[string]struct {
		DatasetName string
		AccessToken string
		URL         string
		Request     SearchRequest
		Response    *SearchResponse
		IsError     bool
	}{
		"test-1: with correct request": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "cillum",
				OrderField: OrderFieldID,
				OrderBy:    OrderByAsc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
					{
						ID:     2,
						Name:   "Brooks Aguilar",
						Age:    25,
						About:  "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			IsError: false,
		},
		"test-2: without next page": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      5,
				Offset:     1,
				Query:      "Adipisicing",
				OrderField: OrderFieldAge,
				OrderBy:    OrderByAsc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     20,
						Name:   "Lowery York",
						Age:    27,
						About:  "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n",
						Gender: "male",
					},
					{
						ID:     17,
						Name:   "Dillard Mccoy",
						Age:    36,
						About:  "Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			IsError: false,
		},
		"test-3: with empty response": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:  100,
				Offset: 50,
				Query:  "test",
			},
			Response: &SearchResponse{
				Users:    []User{},
				NextPage: false,
			},
			IsError: false,
		},
		"test-4: with wrong dataset file": {
			DatasetName: "data.xml",
			AccessToken: "token",
			Request:     SearchRequest{},
			Response:    nil,
			IsError:     true,
		},
		"test-5: without token": {
			DatasetName: FileDataset,
			AccessToken: "",
			Request:     SearchRequest{},
			Response:    nil,
			IsError:     true,
		},
		"test-6: with wrong offset": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:  5,
				Offset: -1,
			},
			Response: nil,
			IsError:  true,
		},
		"test-7: with wrong limit": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:  -1,
				Offset: 0,
			},
			Response: nil,
			IsError:  true,
		},
		"test-8: with wrong order_by": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:   1,
				Offset:  0,
				OrderBy: -2,
			},
			Response: nil,
			IsError:  true,
		},
		"test-9: with wrong order_field": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      5,
				Offset:     1,
				OrderField: "field",
			},
			Response: nil,
			IsError:  true,
		},
		"test-10: with order_by asc": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				OrderField: OrderFieldName,
				OrderBy:    OrderByAsc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     15,
						Name:   "Allison Valdez",
						Age:    21,
						About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
						Gender: "male",
					},
					{
						ID:     16,
						Name:   "Annie Osborn",
						Age:    35,
						About:  "Consequat fugiat veniam commodo nisi nostrud culpa pariatur. Aliquip velit adipisicing dolor et nostrud. Eu nostrud officia velit eiusmod ullamco duis eiusmod ad non do quis.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
			IsError: false,
		},
		"test-11: with order_by desc": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				OrderField: OrderFieldAge,
				OrderBy:    OrderByDesc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     13,
						Name:   "Whitley Davidson",
						Age:    40,
						About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
						Gender: "male",
					},
					{
						ID:     32,
						Name:   "Christy Knapp",
						Age:    40,
						About:  "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
			IsError: false,
		},
		"test-12: with order_by asis": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				OrderField: OrderFieldAge,
				OrderBy:    OrderByAsIs,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
					{
						ID:     1,
						Name:   "Hilda Mayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
			IsError: false,
		},
		"test-13: with order_field id": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				OrderField: OrderFieldID,
				OrderBy:    OrderByDesc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     34,
						Name:   "Kane Sharp",
						Age:    34,
						About:  "Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n",
						Gender: "male",
					},
					{
						ID:     33,
						Name:   "Twila Snow",
						Age:    36,
						About:  "Sint non sunt adipisicing sit laborum cillum magna nisi exercitation. Dolore officia esse dolore officia ea adipisicing amet ea nostrud elit cupidatat laboris. Proident culpa ullamco aute incididunt aute. Laboris et nulla incididunt consequat pariatur enim dolor incididunt adipisicing enim fugiat tempor ullamco. Amet est ullamco officia consectetur cupidatat non sunt laborum nisi in ex. Quis labore quis ipsum est nisi ex officia reprehenderit ad adipisicing fugiat. Labore fugiat ea dolore exercitation sint duis aliqua.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
			IsError: false,
		},
		"test-14: with order_field name": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				OrderField: OrderFieldName,
				OrderBy:    OrderByDesc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     13,
						Name:   "Whitley Davidson",
						Age:    40,
						About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
						Gender: "male",
					},
					{
						ID:     33,
						Name:   "Twila Snow",
						Age:    36,
						About:  "Sint non sunt adipisicing sit laborum cillum magna nisi exercitation. Dolore officia esse dolore officia ea adipisicing amet ea nostrud elit cupidatat laboris. Proident culpa ullamco aute incididunt aute. Laboris et nulla incididunt consequat pariatur enim dolor incididunt adipisicing enim fugiat tempor ullamco. Amet est ullamco officia consectetur cupidatat non sunt laborum nisi in ex. Quis labore quis ipsum est nisi ex officia reprehenderit ad adipisicing fugiat. Labore fugiat ea dolore exercitation sint duis aliqua.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
			IsError: false,
		},
		"test-15: with order_field age": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				OrderField: OrderFieldAge,
				OrderBy:    OrderByDesc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     13,
						Name:   "Whitley Davidson",
						Age:    40,
						About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
						Gender: "male",
					},
					{
						ID:     32,
						Name:   "Christy Knapp",
						Age:    40,
						About:  "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
			IsError: false,
		},
		"test-16: with order_field empty": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				OrderField: OrderFieldEmpty,
				OrderBy:    OrderByDesc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     13,
						Name:   "Whitley Davidson",
						Age:    40,
						About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
						Gender: "male",
					},
					{
						ID:     33,
						Name:   "Twila Snow",
						Age:    36,
						About:  "Sint non sunt adipisicing sit laborum cillum magna nisi exercitation. Dolore officia esse dolore officia ea adipisicing amet ea nostrud elit cupidatat laboris. Proident culpa ullamco aute incididunt aute. Laboris et nulla incididunt consequat pariatur enim dolor incididunt adipisicing enim fugiat tempor ullamco. Amet est ullamco officia consectetur cupidatat non sunt laborum nisi in ex. Quis labore quis ipsum est nisi ex officia reprehenderit ad adipisicing fugiat. Labore fugiat ea dolore exercitation sint duis aliqua.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
			IsError: false,
		},
		"test-17: with bad url": {
			URL:         "localhost",
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:  2,
				Offset: 0,
			},
			IsError: true,
		},
		"test-18: with timeout": {
			DatasetName: FileDataset,
			AccessToken: "token",
			Request: SearchRequest{
				Limit:  2,
				Offset: 0,
				Query:  "qwertyuiopasdfghjklzxcvbnm",
			},
			IsError: true,
		},
	}

	for name, item := range cases {
		server := httptest.NewServer(http.HandlerFunc(SearchServer))
		defer server.Close()

		url := item.URL
		if len(url) == 0 {
			url = server.URL
		}

		client := &SearchClient{
			AccessToken: item.AccessToken,
			URL:         url,
		}

		FileDataset = item.DatasetName

		response, err := client.FindUsers(item.Request)
		if err != nil && !item.IsError {
			t.Errorf("[%s] unexpected error: %v", name, err)
		}
		if err == nil && item.IsError {
			t.Errorf("[%s] expected error, got nil", name)
		}
		if !reflect.DeepEqual(item.Response, response) {
			t.Errorf("[%s] wrong result, expected %#v, got %#v", name, item.Response, response)
		}
	}
}
