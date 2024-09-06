package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v3"
	"log"
)

// Config All field of structs needs to be exported (ie public, ie capitalized)
type Config struct {
	Hosts    map[string]Host    `yaml:"hosts"`
	Services map[string]Service `yaml:"services"`
}

type Host struct {
	Ipv4 string `yaml:"ipv4"`
}

type Service struct {
	Description string   `yaml:"description"`
	Protocols   []string `yaml:"protocols"`
	Host        string   `yaml:"host"`
}

func CheckCommand() *cli.Command {
	return &cli.Command{
		Name:      "infra",
		HideHelp:  true,
		Usage:     "Perform Infra checks",
		Action:    CheckAction(),
		ArgsUsage: "registry [options]",
		Flags: []cli.Flag{&cli.StringFlag{
			Name: "jq-expression",
			// backtick is used as the name of the variable
			Usage: "Set the `JQ_EXPRESSION` used to parse the JSON API response ",
			Action: func(ctx context.Context, command *cli.Command, v string) error {
				if v == "" {
					return cli.Exit("jq expression should not be empty", 1)
				}
				return nil
			},
		}},
	}
}

func CheckAction() func(c context.Context, command *cli.Command) error {
	return func(c context.Context, command *cli.Command) error {

		viperSonitorConf := viper.New()
		viperSonitorConf.SetConfigName("sonitor")
		viperSonitorConf.AddConfigPath("example")
		err := viperSonitorConf.ReadInConfig()
		if err != nil {
			panic(err)
		}
		test := viperSonitorConf.Get("hosts")
		fmt.Println(test)
		var config Config

		err = viperSonitorConf.Unmarshal(&config)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}

		for serviceName, service := range config.Services {
			fmt.Println("Service:", serviceName, "Description:", service.Description)
			host := config.Hosts[service.Host]
			fmt.Println("Host Ipv4", host.Ipv4)
		}

		return nil

	}
}
