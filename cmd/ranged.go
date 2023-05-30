/*
Copyright © 2023 Adam Worley adam.worley@netwealth.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	DirectoryPath string
)

// rangedCmd represents the ranged command
var rangedCmd = &cobra.Command{
	Use:   "ranged",
	Short: "Update nuget references to use ranged versioning",
	Long: `Convert a static version reference in a .csproj file (e.g. <PackageReference Include="Azure.Data.Tables" Version="12.8.0" />)
	to use a ranged version instead (e.g. ...Version="[12.8.0,13.0)").
	
	The version will use the current static version as the minimum and update the max version to go up to but not include
	the next major version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("updating .csproj files")

		// Read the directory path
		// directoryPath := "C:/repos/Netwealth.MessageService"

		err := updateCSProjFiles(DirectoryPath)
		if err != nil {
			log.Fatalf("Failed to update .csproj files: %v", err)
		}

		fmt.Println("Conversion complete.")
	},
}

func init() {
	nugetCmd.AddCommand(rangedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rangedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	rangedCmd.Flags().StringVarP(&DirectoryPath, "directory", "d", "", "The path to the directory (e.g. /Netwealth.Project)")
	rangedCmd.MarkFlagRequired("directory")
}

// Updates all .csproj files in the given directory and its subdirectories
func updateCSProjFiles(directoryPath string) error {
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".csproj" {
			err = updateCSProjFile(path)
			if err != nil {
				log.Printf("Failed to update .csproj file: %s (%v)", path, err)
			}
		}

		return nil
	})

	return err
}

// Updates a single .csproj file
func updateCSProjFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	outputLines := make([]string, 0)

	// Regex pattern to match NuGet package version numbers
	versionPattern := regexp.MustCompile(`<PackageReference Include="[^"]+" Version="([^"]+)"`)

	for scanner.Scan() {
		line := scanner.Text()

		if versionPattern.MatchString(line) {
			// Extract the package version number from the line
			matches := versionPattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				originalVersion := matches[1]

				// Ignore already ranged versions
				if strings.HasPrefix(originalVersion, "[") && strings.HasSuffix(originalVersion, ")") {
					outputLines = append(outputLines, line)
					continue
				}

				rangedVersion := convertToRangedVersion(originalVersion)
				newLine := strings.Replace(line, originalVersion, rangedVersion, 1)
				outputLines = append(outputLines, newLine)
				continue
			}
		}

		outputLines = append(outputLines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Write the modified content back to the .csproj file
	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	for _, line := range outputLines {
		fmt.Fprintln(outputFile, line)
	}

	return nil
}

// Converts a NuGet package version number to a ranged version up to the next major release
func convertToRangedVersion(version string) string {
	segments := strings.Split(version, ".")

	if len(segments) > 0 {
		majorVersion, err := strconv.Atoi(segments[0])
		if err == nil {
			// Increment the major version by 1 to get the next major release version
			segments[0] = strconv.Itoa(majorVersion + 1)
		}
	}

	// Format the upper version as '2.0'
	upperVersion := segments[0] + ".0"

	// Join the segments back to form the ranged version
	return "[" + version + "," + upperVersion + ")"
}
