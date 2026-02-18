//
// rubrik-exporter
//
// Exports metrics from rubrik backup for prometheus
//
// License: Apache License Version 2.0,
// Organization: Claranet GmbH
// Author: Martin Weber <martin.weber@de.clara.net>
//

package rubrik

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Session struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organizationId"`
	Token          string `json:"token"`
	UserId         string `json:"userId"`
}

type OAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (r *Rubrik) Login() error {
	// Check if service account authentication is being used
	if r.serviceAccountClientID != "" && r.serviceAccountClientSecret != "" {
		log.Print("Using service account authentication")
		return r.loginWithServiceAccount()
	}

	// Fall back to username/password authentication
	log.Print("Using username/password authentication")
	return r.loginWithUsernamePassword()
}

func (r *Rubrik) loginWithServiceAccount() error {
	// Try OAuth2 client credentials flow first
	if err := r.tryOAuth2ClientCredentials(); err == nil {
		return nil
	}

	log.Print("OAuth2 client credentials failed, trying basic auth with service account credentials")
	// Fall back to basic auth with client_id/client_secret
	return r.tryServiceAccountBasicAuth()
}

func (r *Rubrik) tryOAuth2ClientCredentials() error {
	_url := r.url + "/api/client_token"

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	var netClient = http.Client{Transport: tr}

	// Prepare Rubrik service account request
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     r.serviceAccountClientID,
		"client_secret": r.serviceAccountClientSecret,
	}

	// Create form-encoded body
	values := make(url.Values)
	for key, value := range data {
		values.Set(key, value)
	}

	req, err := http.NewRequest("POST", _url, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := netClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Rubrik service account authentication failed: HTTP %d", resp.StatusCode)
	}

	var tokenResp OAuth2TokenResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tokenResp)
	if err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	r.sessionToken = tokenResp.AccessToken
	r.isLoggedIn = true

	log.Printf("Successfully authenticated with Rubrik service account (/api/client_token)")
	return nil
}

func (r *Rubrik) tryServiceAccountBasicAuth() error {
	_url := r.url + "/api/v1/session"

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	var netClient = http.Client{Transport: tr}

	req, err := http.NewRequest("POST", _url, nil)
	if err != nil {
		return err
	}

	// Use client_id as username and client_secret as password
	req.SetBasicAuth(r.serviceAccountClientID, r.serviceAccountClientSecret)

	resp, err := netClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("service account basic auth failed: HTTP %d", resp.StatusCode)
	}

	data := json.NewDecoder(resp.Body)
	var s Session
	err = data.Decode(&s)
	if err != nil {
		return fmt.Errorf("failed to decode session response: %v", err)
	}

	r.sessionToken = s.Token
	r.isLoggedIn = true

	log.Printf("Successfully authenticated with service account basic auth")
	return nil
}

func (r *Rubrik) loginWithUsernamePassword() error {
	_url := r.url + "/api/v1/session"

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	var netClient = http.Client{Transport: tr}
	req, err := http.NewRequest("POST", _url, nil)

	if err != nil {
		log.Fatal(err)
		return err
	}
	req.SetBasicAuth(r.username, r.password)

	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	var s Session
	err = data.Decode(&s)

	r.sessionToken = s.Token
	r.isLoggedIn = true

	return nil
}

func (r *Rubrik) Logout() {
	resp, _ := r.makeRequest("DELETE", "/api/v1/session", RequestParams{})
	if resp != nil {
		resp.Body.Close()
	}
}
