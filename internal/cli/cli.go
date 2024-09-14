package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "soft-crusher",
		Usage: "API generation and deployment tool",
		Commands: []*cli.Command{
			{
				Name:    "analyze",
				Aliases: []string{"a"},
				Usage:   "Analyze Go files in the current directory",
				Action: func(c *cli.Context) error {
					analyzer := NewFunctionAnalyzer()
					err := analyzer.AnalyzeDirectory("./")
					if err != nil {
						return fmt.Errorf("error analyzing directory: %v", err)
					}
					fmt.Println("Analysis completed successfully!")
					return nil
				},
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate API code, documentation, and tests",
				Action: func(c *cli.Context) error {
					analyzer := NewFunctionAnalyzer()
					err := analyzer.AnalyzeDirectory("./")
					if err != nil {
						return fmt.Errorf("error analyzing directory: %v", err)
					}

					designer := NewAPIDesigner()
					designer.DesignAPI(analyzer.Functions)

					generator := NewCodeGenerator(designer)
					err = generator.GenerateAPICode()
					if err != nil {
						return fmt.Errorf("error generating API code: %v", err)
					}

					docGenerator := NewDocumentationGenerator(designer)
					err = docGenerator.GenerateSwaggerDoc()
					if err != nil {
						return fmt.Errorf("error generating Swagger documentation: %v", err)
					}

					testGenerator := NewTestingSuiteGenerator(designer)
					err = testGenerator.GenerateTests()
					if err != nil {
						return fmt.Errorf("error generating test suite: %v", err)
					}

					err = testGenerator.UpdateGoModFile()
					if err != nil {
						return fmt.Errorf("error updating go.mod file: %v", err)
					}

					fmt.Println("API generation completed successfully!")
					return nil
				},
			},
			{
				Name:    "deploy",
				Aliases: []string{"d"},
				Usage:   "Generate deployment configurations",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "soft-crusher-api",
						Usage:   "Name of the API",
					},
					&cli.StringFlag{
						Name:    "version",
						Aliases: []string{"v"},
						Value:   "1.0.0",
						Usage:   "Version of the API",
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   8080,
						Usage:   "Port number for the API",
					},
				},
				Action: func(c *cli.Context) error {
					helper := NewDeploymentHelper(
						c.String("name"),
						c.String("version"),
						c.Int("port"),
					)

					err := helper.GenerateDockerfile()
					if err != nil {
						return fmt.Errorf("error generating Dockerfile: %v", err)
					}

					err = helper.GenerateKubernetesManifests()
					if err != nil {
						return fmt.Errorf("error generating Kubernetes manifests: %v", err)
					}

					err = helper.GenerateDockerComposeFile()
					if err != nil {
						return fmt.Errorf("error generating Docker Compose file: %v", err)
					}

					fmt.Println("Deployment configurations generated successfully!")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}