// Package translate implements translation on Google Translate.
package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/StalkR/goircbot/lib/tls"
)

type lresult struct {
	Data ldata
}

type ldata struct {
	Languages []Language
}

// Language represents a supported language with language code and name.
type Language struct {
	Language string
	Name     string
}

// Languages returns the list of supported Google Translate languages for a
// given target language or empty string for all supported languages.
func Languages(target string, key string) ([]Language, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tls.Config("www.googleapis.com"),
		},
	}
	base := "https://www.googleapis.com/language/translate/v2/languages"
	params := url.Values{}
	params.Set("key", key)
	if target != "" {
		params.Set("target", target)
	}
	resp, err := client.Get(fmt.Sprintf("%s?%s", base, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := &lresult{}
	err = json.Unmarshal(contents, r)
	if err != nil {
		return nil, err
	}
	return r.Data.Languages, nil
}

type tresult struct {
	Data tdata
}
type tdata struct {
	Translations []Translation
}

// Translation represents a result with text and detected source language.
type Translation struct {
	TranslatedText         string
	DetectedSourceLanguage string
}

// Translate translates a text on Google Translate from a language to another.
// It requires a Google API Key (key), valid source and target languages.
// For automatic source language detection, use empty string.
func Translate(source, target, text, key string) (*Translation, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tls.Config("www.googleapis.com"),
		},
	}
	base := "https://www.googleapis.com/language/translate/v2"
	params := url.Values{}
	params.Set("key", key)
	if source != "" {
		params.Set("source", source)
	}
	params.Set("target", target)
	params.Set("q", text)
	resp, err := client.Get(fmt.Sprintf("%s?%s", base, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := &tresult{}
	err = json.Unmarshal(contents, r)
	if err != nil {
		return nil, err
	}
	if len(r.Data.Translations) == 0 {
		return nil, errors.New("no translation")
	}
	t := &r.Data.Translations[0]
	t.TranslatedText = html.UnescapeString(t.TranslatedText)
	return t, nil
}
