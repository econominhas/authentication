package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/econominhas/authentication/internal/utils"
)

type EnvVar struct {
	Name          string
	Numeric       bool
	Required      bool
	AllowedValues []string
}

var requiredEnvVars = []EnvVar{
	{
		Name:     "ENV",
		Required: true,
		AllowedValues: []string{
			"dev",
			"prod",
		},
	},
	{
		Name:     "PORT",
		Required: true,
		Numeric:  true,
	},
	{
		Name:     "GOOGLE_CLIENT_ID",
		Required: true,
	},
	{
		Name:     "GOOGLE_CLIENT_SECRET",
		Required: true,
	},
	{
		Name:     "FACEBOOK_CLIENT_ID",
		Required: true,
	},
	{
		Name:     "FACEBOOK_CLIENT_SECRET",
		Required: true,
	},
	{
		Name:     "PASETO_PRIVATE_KEY",
		Required: true,
	},
	{
		Name:     "DATABASE_URL",
		Required: true,
	},
}

func validateEnvs() {
	for _, v := range requiredEnvVars {
		envVar := os.Getenv(v.Name)

		if v.Required && envVar == "" {
			log.Fatalf("Missing \"%s\"", v.Name)
			panic(1)
		}

		if v.Numeric {
			match, err := regexp.MatchString("^[0-9]*$", envVar)

			if err != nil || !match {
				log.Fatalf("\"%s\" must be numeric", v.Name)
				panic(1)
			}
		}

		if len(v.AllowedValues) > 0 && !utils.InArray(envVar, v.AllowedValues) {
			log.Fatalf("Variable \"%s\" don't match allowed values: %s", v.Name, strings.Join(v.AllowedValues, ", "))
			panic(1)
		}
	}
}
