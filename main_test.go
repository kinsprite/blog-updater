package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mainRouter(t *testing.T) {
	type args struct {
		uri     string
		event   string
		payload string
	}

	serverSecret = "xxx-yyy-zzz"

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Github ping",
			args: args{uri: "/github-webhooks", event: "ping",
				payload: `{"zen":"3232aa","hook_id":123,"hook":{}`},
			want: `{"message":"Pong"}`,
		},
		{
			name: "Github push",
			args: args{uri: "/github-webhooks", event: "push",
				payload: `{}`},
			want: `{"message":"Done"}`,
		},
	}

	router := mainRouter()

	ts := httptest.NewServer(router)
	defer ts.Close()

	getResp := func(args args) string {
		client := &http.Client{}

		url := fmt.Sprintf("%s%s", ts.URL, args.uri)
		req, err := http.NewRequest("POST", url, strings.NewReader(args.payload))
		sig := generateSignature([]byte(args.payload))

		req.Header.Add("X-GitHub-Event", args.event)
		req.Header.Add("X-Hub-Signature", sig)
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := getResp(tt.args.uri); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("mainRouter() = %v, want %v", got, tt.want)
			// }

			got := getResp(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
