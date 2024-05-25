package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/econominhas/authentication/internal/models"
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

func validateEnvs(logger models.Logger) {
	for _, v := range requiredEnvVars {
		envVar := os.Getenv(v.Name)

		if v.Required && envVar == "" {
			logger.Error(
				fmt.Sprintf("Missing \"%s\"", v.Name),
			)
			panic(1)
		}

		if v.Numeric {
			exp, err := regexp.Compile("^[0-9]*$")
			if err != nil {
				logger.Error("Fail to compile regex")
				panic(1)
			}

			match := exp.Match([]byte(envVar))

			if !match {
				logger.Error(
					fmt.Sprintf("\"%s\" must be numeric", v.Name),
				)
				panic(1)
			}
		}

		if len(v.AllowedValues) > 0 && !utils.InArray(envVar, v.AllowedValues) {
			logger.Error(
				fmt.Sprintf(
					"Variable \"%s\" don't match allowed values: %s",
					v.Name,
					strings.Join(v.AllowedValues, ", "),
				),
			)
			panic(1)
		}
	}
}
