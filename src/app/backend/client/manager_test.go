// Copyright 2017 The Kubernetes Authors.
// Copyright 2020 Authors of Arktos - file modified.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"crypto/tls"
	"net/http"
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/kubernetes/dashboard/src/app/backend/args"
	"github.com/kubernetes/dashboard/src/app/backend/errors"
	"k8s.io/client-go/rest"
)

func TestNewClientManager(t *testing.T) {
	cases := []struct {
		kubeConfigPath, apiserverHost string
	}{
		{"", "test"},
	}

	for _, c := range cases {
		manager := NewClientManager(c.kubeConfigPath, c.apiserverHost)

		if manager == nil {
			t.Fatalf("NewClientManager(%s, %s): Expected manager not to be nil",
				c.kubeConfigPath, c.apiserverHost)
		}
	}
}

func TestClient(t *testing.T) {
	cases := []struct {
		request *restful.Request
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{}),
				},
			},
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "http://localhost:8080")
		_, err := manager.Client(c.request)

		if err != nil {
			t.Fatalf("Client(%v): Expected client to be created but error was thrown:"+
				" %s", c.request, err.Error())
		}
	}
}

func TestSecureClient(t *testing.T) {
	cases := []struct {
		request       *restful.Request
		expectedError bool
		err           error
	}{
		{
			request: &restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{}),
					TLS:    &tls.ConnectionState{},
				},
			},
			expectedError: true,
			err:           errors.NewUnauthorized(errors.MsgLoginUnauthorizedError),
		},
		{
			request: &restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{"Authorization": {"Bearer asd"}}),
					TLS:    &tls.ConnectionState{},
				},
			},
			expectedError: false,
			err:           nil,
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "http://localhost:8080")
		_, err := manager.Client(c.request)

		if err == nil && !c.expectedError {
			continue
		}

		if !c.expectedError && err != nil {
			t.Fatalf("Client(%v): Expected client to be created but error was thrown:"+
				" %s", c.request, err.Error())
		}

		if c.expectedError && err == nil {
			t.Fatalf("Expected error %v but got nil", c.err)
		}

		if c.err.Error() != err.Error() {
			t.Fatalf("Expected error %v but got %v", c.err, err)
		}
	}
}

func TestAPIExtensionsClient(t *testing.T) {
	cases := []struct {
		request *restful.Request
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{}),
				},
			},
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "http://localhost:8080")
		_, err := manager.APIExtensionsClient(c.request)

		if err != nil {
			t.Fatalf("APIExtensionsClient(%v): Expected API extensions client to be created"+
				" but error was thrown: %s", c.request, err.Error())
		}
	}
}

func TestSecureAPIExtensionsClient(t *testing.T) {
	cases := []struct {
		request       *restful.Request
		expectedError bool
		err           error
	}{
		{
			request: &restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{}),
					TLS:    &tls.ConnectionState{},
				},
			},
			expectedError: true,
			err:           errors.NewUnauthorized(errors.MsgLoginUnauthorizedError),
		},
		{
			request: &restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{"Authorization": {"Bearer asd"}}),
					TLS:    &tls.ConnectionState{},
				},
			},
			expectedError: false,
			err:           nil,
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "http://localhost:8080")
		_, err := manager.APIExtensionsClient(c.request)

		if err == nil && !c.expectedError {
			continue
		}

		if !c.expectedError && err != nil {
			t.Fatalf("APIExtensions(%v): Expected client to be created but error"+
				" was thrown: %s", c.request, err.Error())
		}

		if c.expectedError && err == nil {
			t.Fatalf("Expected error %v but got nil", c.err)
		}

		if c.err.Error() != err.Error() {
			t.Fatalf("Expected error %v but got %v", c.err, err)
		}
	}
}

func TestCSRFKey(t *testing.T) {
	manager := NewClientManager("", "http://localhost:8080")
	key := manager.CSRFKey()

	if len(key) == 0 {
		t.Fatal("CSRFKey(): Expected csrf key to be autogenerated.")
	}
}

func TestConfig(t *testing.T) {
	cases := []struct {
		request  *restful.Request
		expected string
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{
						"Authorization": {"Bearer test-token"},
					}),
					TLS: &tls.ConnectionState{},
				},
			},
			"test-token",
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "https://localhost:8080")
		cfg, err := manager.Config(c.request)

		if err != nil {
			t.Fatalf("Config(%v): Expected config to be created but error was thrown:"+
				" %s",
				c.request, err.Error())
		}

		if cfg.GetConfig().BearerToken != c.expected {
			t.Fatalf("Config(%v): Expected token to be %s but got %s",
				c.request, c.expected, cfg.GetConfig().BearerToken)
		}
	}
}

func TestClientCmdConfig(t *testing.T) {
	args.GetHolderBuilder().SetEnableSkipLogin(true)
	cases := []struct {
		request  *restful.Request
		expected string
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{
						"Authorization": {"Bearer test-token"},
					}),
					TLS: &tls.ConnectionState{},
				},
			},
			"test-token",
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "https://localhost:8080")
		cmdCfg, err := manager.ClientCmdConfig(c.request)

		if err != nil {
			t.Fatalf("Config(%v): Expected client config to be created but error was thrown:"+
				" %s",
				c.request, err.Error())
		}

		var bearerToken string
		if cmdCfg != nil {
			cfg, err := cmdCfg.ClientConfig()
			if err != nil {
				t.Fatalf("Config(%v): Expected config to be created but error was thrown:"+
					" %s",
					c.request, err.Error())
			}
			bearerToken = cfg.GetConfig().BearerToken
		}

		if bearerToken != c.expected {
			t.Fatalf("Config(%v): Expected token to be %s but got %s",
				c.request, c.expected, bearerToken)
		}
	}
}

func TestVerberClient(t *testing.T) {
	manager := NewClientManager("", "http://localhost:8080")
	_, err := manager.VerberClient(&restful.Request{Request: &http.Request{TLS: &tls.ConnectionState{}}}, &rest.Config{})

	if err != nil {
		t.Fatalf("VerberClient(): Expected verber client to be created but got error: %s",
			err.Error())
	}
}

func TestClientManager_InsecureClients(t *testing.T) {
	manager := NewClientManager("", "http://localhost:8080")
	if manager.InsecureClient() == nil {
		t.Fatalf("InsecureClient(): Expected insecure client not to be nil")
	}
}

func TestClientManager_InsecureAPIExtensionsClient(t *testing.T) {
	manager := NewClientManager("", "http://localhost:8080")
	if manager.InsecureAPIExtensionsClient() == nil {
		t.Fatalf("InsecureClient(): Expected insecure client not to be nil")
	}
}

func TestImpersonationUserClient(t *testing.T) {
	args.GetHolderBuilder().SetEnableSkipLogin(true)
	cases := []struct {
		request                   *restful.Request
		expected                  string
		expectedImpersonationUser string
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{
						"Authorization":    {"Bearer test-token"},
						"Impersonate-User": {"impersonatedUser"},
					}),
					TLS: &tls.ConnectionState{},
				},
			},
			"test-token",
			"impersonatedUser",
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "https://localhost:8080")
		cfg, err := manager.Config(c.request)
		//authInfo := manager.extractAuthInfo(c.request)
		if err != nil {
			t.Fatalf("Config(%v): Expected config to be created but error was thrown:"+
				" %s",
				c.request, err.Error())
		}

		if cfg.GetConfig().BearerToken != c.expected {
			t.Fatalf("Config(%v): Expected token to be %s but got %s",
				c.request, c.expected, cfg.GetConfig().BearerToken)
		}

		if cfg.GetConfig().Impersonate.UserName != c.expectedImpersonationUser {
			t.Fatalf("Config(%v): Expected impersonated user to be %s but got %s",
				c.request, c.expectedImpersonationUser, cfg.GetConfig().Impersonate.UserName)
		}

	}
}

func TestNoImpersonationUserWithNoBearerClient(t *testing.T) {
	args.GetHolderBuilder().SetEnableSkipLogin(true)
	cases := []struct {
		request *restful.Request
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{}),
					TLS:    &tls.ConnectionState{},
				},
			},
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "https://localhost:8080")
		cfg, err := manager.Config(c.request)
		//authInfo := manager.extractAuthInfo(c.request)
		if err != nil {
			t.Fatalf("Config(%v): Expected config to be created but error was thrown:"+
				" %s",
				c.request, err.Error())
		}

		if len(cfg.GetConfig().BearerToken) > 0 {
			t.Fatalf("Config(%v): Expected no token but got %s",
				c.request, cfg.GetConfig().BearerToken)
		}

		if len(cfg.GetConfig().Impersonate.UserName) > 0 {
			t.Fatalf("Config(%v): Expected no impersonated user but got %s",
				c.request, cfg.GetConfig().Impersonate.UserName)
		}

	}
}

func TestImpersonationOneGroupClient(t *testing.T) {
	args.GetHolderBuilder().SetEnableSkipLogin(true)
	cases := []struct {
		request                     *restful.Request
		expected                    string
		expectedImpersonationUser   string
		expectedImpersonationGroups []string
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{
						"Authorization":     {"Bearer test-token"},
						"Impersonate-User":  {"impersonatedUser"},
						"Impersonate-Group": {"group1"},
					}),
					TLS: &tls.ConnectionState{},
				},
			},
			"test-token",
			"impersonatedUser",
			[]string{"group1"},
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "https://localhost:8080")
		cfg, err := manager.Config(c.request)
		//authInfo := manager.extractAuthInfo(c.request)
		if err != nil {
			t.Fatalf("Config(%v): Expected config to be created but error was thrown:"+
				" %s",
				c.request, err.Error())
		}

		if cfg.GetConfig().BearerToken != c.expected {
			t.Fatalf("Config(%v): Expected token to be %s but got %s",
				c.request, c.expected, cfg.GetConfig().BearerToken)
		}

		if cfg.GetConfig().Impersonate.UserName != c.expectedImpersonationUser {
			t.Fatalf("Config(%v): Expected impersonated user to be %s but got %s",
				c.request, c.expectedImpersonationUser, cfg.GetConfig().Impersonate.UserName)
		}

		if len(cfg.GetConfig().Impersonate.Groups) != 1 {
			t.Fatalf("Config(%v): Expected one impersonated group but got %d",
				c.request, len(cfg.GetConfig().Impersonate.Groups))
		}

		if cfg.GetConfig().Impersonate.Groups[0] != c.expectedImpersonationGroups[0] {
			t.Fatalf("Config(%v): Expected impersonated group to be %s but got %s",
				c.request, cfg.GetConfig().Impersonate.Groups[0], c.expectedImpersonationGroups[0])
		}
	}
}

func TestImpersonationTwoGroupClient(t *testing.T) {
	args.GetHolderBuilder().SetEnableSkipLogin(true)
	cases := []struct {
		request                     *restful.Request
		expected                    string
		expectedImpersonationUser   string
		expectedImpersonationGroups []string
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{
						"Authorization":     {"Bearer test-token"},
						"Impersonate-User":  {"impersonatedUser"},
						"Impersonate-Group": {"group1", "groups2"},
					}),
					TLS: &tls.ConnectionState{},
				},
			},
			"test-token",
			"impersonatedUser",
			[]string{"group1", "groups2"},
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "https://localhost:8080")
		cfg, err := manager.Config(c.request)
		//authInfo := manager.extractAuthInfo(c.request)
		if err != nil {
			t.Fatalf("Config(%v): Expected config to be created but error was thrown:"+
				" %s",
				c.request, err.Error())
		}

		if cfg.GetConfig().BearerToken != c.expected {
			t.Fatalf("Config(%v): Expected token to be %s but got %s",
				c.request, c.expected, cfg.GetConfig().BearerToken)
		}

		if cfg.GetConfig().Impersonate.UserName != c.expectedImpersonationUser {
			t.Fatalf("Config(%v): Expected impersonated user to be %s but got %s",
				c.request, c.expectedImpersonationUser, cfg.GetConfig().Impersonate.UserName)
		}

		if len(cfg.GetConfig().Impersonate.Groups) != 2 {
			t.Fatalf("Config(%v): Expected two impersonated group but got %d",
				c.request, len(cfg.GetConfig().Impersonate.Groups))
		}

		if cfg.GetConfig().Impersonate.Groups[0] != c.expectedImpersonationGroups[0] {
			t.Fatalf("Config(%v): Expected impersonated group to be %s but got %s",
				c.request, cfg.GetConfig().Impersonate.Groups[0], c.expectedImpersonationGroups[0])
		}

		if cfg.GetConfig().Impersonate.Groups[1] != c.expectedImpersonationGroups[1] {
			t.Fatalf("Config(%v): Expected impersonated group to be %s but got %s",
				c.request, cfg.GetConfig().Impersonate.Groups[1], c.expectedImpersonationGroups[1])
		}
	}
}

func TestImpersonationExtrasClient(t *testing.T) {
	args.GetHolderBuilder().SetEnableSkipLogin(true)
	cases := []struct {
		request                    *restful.Request
		expected                   string
		expectedImpersonationUser  string
		expectedImpersonationExtra map[string][]string
	}{
		{
			&restful.Request{
				Request: &http.Request{
					Header: http.Header(map[string][]string{
						"Authorization":             {"Bearer test-token"},
						"Impersonate-User":          {"impersonatedUser"},
						"Impersonate-Extra-scope":   {"views", "writes"},
						"Impersonate-Extra-service": {"iguess"},
					}),
					TLS: &tls.ConnectionState{},
				},
			},
			"test-token",
			"impersonatedUser",
			map[string][]string{"scope": {"views", "writes"},
				"service": {"iguess"}},
		},
	}

	for _, c := range cases {
		manager := NewClientManager("", "https://localhost:8080")
		cfg, err := manager.Config(c.request)
		//authInfo := manager.extractAuthInfo(c.request)
		if err != nil {
			t.Fatalf("Config(%v): Expected config to be created but error was thrown:"+
				" %s",
				c.request, err.Error())
		}

		if cfg.GetConfig().BearerToken != c.expected {
			t.Fatalf("Config(%v): Expected token to be %s but got %s",
				c.request, c.expected, cfg.GetConfig().BearerToken)
		}

		if cfg.GetConfig().Impersonate.UserName != c.expectedImpersonationUser {
			t.Fatalf("Config(%v): Expected impersonated user to be %s but got %s",
				c.request, c.expectedImpersonationUser, cfg.GetConfig().Impersonate.UserName)
		}

		if len(cfg.GetConfig().Impersonate.Extra) != 2 {
			t.Fatalf("Config(%v): Expected two impersonated extra but got %d",
				c.request, len(cfg.GetConfig().Impersonate.Extra))
		}

		if cfg.GetConfig().Impersonate.Extra["service"][0] != c.expectedImpersonationExtra["service"][0] {
			t.Fatalf("Config(%v): Expected service extra to be %s but got %s",
				c.request, cfg.GetConfig().Impersonate.Extra["service"][0], c.expectedImpersonationExtra["service"][0])

		}

		//check multi value scope

		if len(cfg.GetConfig().Impersonate.Extra["scope"]) != 2 {
			t.Fatalf("Config(%v): Expected two scope impersonated extra but got %d",
				c.request, len(cfg.GetConfig().Impersonate.Extra["scope"]))
		}

		if cfg.GetConfig().Impersonate.Extra["scope"][0] != c.expectedImpersonationExtra["scope"][0] {
			t.Fatalf("Config(%v): Expected scope extra to be %s but got %s",
				c.request, c.expectedImpersonationExtra["scope"][0], cfg.GetConfig().Impersonate.Extra["scope"][0])

		}

		if cfg.GetConfig().Impersonate.Extra["scope"][1] != c.expectedImpersonationExtra["scope"][1] {
			t.Fatalf("Config(%v): Expected scope extra to be %s but got %s",
				c.request, c.expectedImpersonationExtra["scope"][1], cfg.GetConfig().Impersonate.Extra["scope"][1])

		}

		if len(cfg.GetConfig().Impersonate.Extra["scope"]) != 2 {
			t.Fatalf("Config(%v): Expected two scope impersonated extra but got %d",
				c.request, len(cfg.GetConfig().Impersonate.Extra["scope"]))
		}
	}
}
