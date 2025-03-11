package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/go-ldap/ldap/v3"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ldapURL := os.Getenv("LDAP_URL")
	userDN := os.Getenv("LDAP_USER_DN")
	userFilter := os.Getenv("LDAP_USER_FILTER")
	groupDN := os.Getenv("LDAP_GROUP_DN")
	groupFilter := os.Getenv("LDAP_GROUP_FILTER")
	groupAttribute := os.Getenv("LDAP_GROUP_ATTRIBUTE")
	usernameFromEnv := os.Getenv("LDAP_USERNAME")
	passwordFromEnv := os.Getenv("LDAP_PASSWORD")

	var username, password string

	if ldapURL == "" || userDN == "" || userFilter == "" {
		log.Fatal("Missing LDAP configuration in .env file")
	}

	rootCmd := &cobra.Command{
		Use:     "ldap-auth",
		Short:   "Authenticate against LDAP server",
		Version: fmt.Sprintf("%s (built at %s from %s)", version, date, commit),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Enter username (default: '%s'): ", usernameFromEnv)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			username = strings.TrimSpace(input)

			fmt.Printf("Enter password (default: %q): ", strings.Repeat("*", len(passwordFromEnv)))
			bytePassword, _ := term.ReadPassword(syscall.Stdin)
			fmt.Println()
			password = string(bytePassword)
		},
	}

	rootCmd.Flags().StringVarP(&username, "username", "u", "", "LDAP username")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "LDAP password")

	help := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		help(cmd, args)
		os.Exit(0)
	})

	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Println(rootCmd.Version)
			os.Exit(0)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	if username == "" {
		username = usernameFromEnv
	}

	if password == "" {
		password = passwordFromEnv
	}

	conn, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatalf("Failed to connect to LDAP server: %v", err)
	}
	defer conn.Close()

	if username == "" || password == "" {
		log.Fatal("Username and password are required")
	}

	userDN = fmt.Sprintf(userDN, username)
	if err := conn.Bind(userDN, password); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	fmt.Println("Authentication successful!")

	if groupDN == "" || groupFilter == "" || groupAttribute == "" {
		log.Println("Skipping group search due to missing configuration")
		return
	}

	searchGroupFilter := fmt.Sprintf(groupFilter, userDN)
	searchRequest := ldap.NewSearchRequest(
		groupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchGroupFilter,
		[]string{groupAttribute},
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		log.Fatalf("Failed to search groups: %v", err)
	}

	fmt.Println("Groups:")
	for _, entry := range searchResult.Entries {
		fmt.Println(" -", entry.GetAttributeValue(groupAttribute))
	}
}
