package main

import (
	"fmt"
	"testing"
  "time"
)

func TestExtractIdAndCatigorie(t *testing.T) {

  tests := []struct {
    name string
    url string
    want [2]int
  }{
    {
      name: "valid url",
      url: "categorie/309/Offres-emploi/annonce/9402890/Stagiaire-content-manager.html",
      want: [2]int{309, 9402890},
    },
    {
      name: "valid url",
      url: "categorie/309/Offres-emploi/annonce/9294635/Responsable-juridique.html",
      want: [2]int{309, 9294635},
    },
    {
      name: "unvalid url",
      url: "www.google.com",
      want: [2]int{0, 0},
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      id, cat, _ := extractIdAndCatigorie(tt.url)
      if id != tt.want[0] || cat != tt.want[1] {
        t.Errorf("expected [%d, %d], got [%d, %d]", tt.want[0], tt.want[1], id, cat)
      }
    })
  }
}

func TestToInt(t *testing.T) {
  tests := []struct {
    name string
    input []byte
    want int
  }{
    {
      name: "3 digit input",
      input: []byte{byte('3'), byte('9'), byte('8')},
      want: 398,
    },
    {
      name: "5 digit input",
      input: []byte{byte('3'), byte('9'), byte('8'), byte('3'), byte('2')},
      want: 39832,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      result := toInt(tt.input)
      if result != tt.want {
        t.Errorf("expected %d, got %d", tt.want, result)
      }
    })
  }
}

