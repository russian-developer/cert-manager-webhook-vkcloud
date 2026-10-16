package main

import (
	"encoding/json"
	"testing"
)

func TestParseTXTRecords(t *testing.T) {
	arrayBody := []byte(`[{"uuid":"u1","name":"_acme-challenge.dev","content":"token"}]`)
	records, err := parseTXTRecords(arrayBody)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Name != "_acme-challenge.dev" || records[0].Content != "token" {
		t.Fatalf("unexpected records: %#v", records)
	}

	wrappedBody, _ := json.Marshal(map[string]any{
		"txt_records": []map[string]string{
			{"uuid": "u2", "name": "_acme-challenge", "content": "other"},
		},
	})
	records, err = parseTXTRecords(wrappedBody)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].UUID != "u2" {
		t.Fatalf("unexpected wrapped records: %#v", records)
	}
}

func TestRecordNamesMatch(t *testing.T) {
	zone := "kitcrm.ru."
	cases := []struct {
		name    string
		apiName string
		fqdn    string
		want    bool
	}{
		{
			name:    "relative api name vs fqdn",
			apiName: "_acme-challenge.dev",
			fqdn:    "_acme-challenge.dev.kitcrm.ru.",
			want:    true,
		},
		{
			name:    "both fqdn with trailing dots",
			apiName: "_acme-challenge.dev.kitcrm.ru.",
			fqdn:    "_acme-challenge.dev.kitcrm.ru.",
			want:    true,
		},
		{
			name:    "case and trailing dot",
			apiName: "_ACME-Challenge.dev",
			fqdn:    "_acme-challenge.dev.kitcrm.ru",
			want:    true,
		},
		{
			name:    "different labels",
			apiName: "_acme-challenge.prod",
			fqdn:    "_acme-challenge.dev.kitcrm.ru.",
			want:    false,
		},
		{
			name:    "apex relative",
			apiName: "@",
			fqdn:    "kitcrm.ru.",
			want:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := recordNamesMatch(tc.apiName, tc.fqdn, zone)
			if got != tc.want {
				t.Fatalf("recordNamesMatch(%q, %q, %q) = %v, want %v",
					tc.apiName, tc.fqdn, zone, got, tc.want)
			}
		})
	}
}

func TestRelativeRecordName(t *testing.T) {
	if got := relativeRecordName("_acme-challenge.dev.kitcrm.ru.", "kitcrm.ru."); got != "_acme-challenge.dev" {
		t.Fatalf("got %q", got)
	}
	if got := relativeRecordName("_acme-challenge.dev", "kitcrm.ru"); got != "_acme-challenge.dev" {
		t.Fatalf("got %q", got)
	}
	if got := relativeRecordName("kitcrm.ru.", "kitcrm.ru."); got != "@" {
		t.Fatalf("got %q", got)
	}
}
