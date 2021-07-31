package flutter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	replacement1 = `
def keystoreProperties = new Properties()
def keystorePropertiesFile = rootProject.file('key.properties')
if (keystorePropertiesFile.exists()) {
	keystoreProperties.load(new FileInputStream(keystorePropertiesFile))
}

android {
	`
	replacement2 = `
	signingConfigs {
		release {
			keyAlias keystoreProperties['keyAlias']
			keyPassword keystoreProperties['keyPassword']
			storeFile keystoreProperties['storeFile'] ? file(keystoreProperties['storeFile']) : null
			storePassword keystoreProperties['storePassword']
		}
	}
	buildTypes {`
	replacement3 = `signingConfig signingConfigs.release`

	keystoreFlag = "keystoreflag"
)

var (
	password string
)

// keytool -genkey -v -keystore ~/upload-keystore.jks -keyalg RSA -keysize 2048 -validity 10000 -alias upload

/*
-storepass
CN=commonName
OU=organizationUnit
O=organizationName
L=localityName
S=stateName
C=country */

// UploadKeystore create keystore and modify build.gradle
func UploadKeystore() *cobra.Command {
	keystoreCmd := &cobra.Command{
		Use:   "uploadkeystore",
		Short: "create upload keystore for android",
		Long:  "create upload keystore for android and modify build.gradle, all in one",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyFromFLag, _ := cmd.Flags().GetString(keystoreFlag)

			var keystorePath string
			if strings.Contains("~/upload-keystore.jks", "~") {
				homeDir, _ := os.UserHomeDir()
				extractedPath := strings.Split(keyFromFLag, "~")[1]
				fmt.Printf("extractedPath: %v\n", extractedPath)
				keystorePath = filepath.Join(homeDir, extractedPath)
			}

			_, err := os.Stat(keystorePath)
			if !os.IsNotExist(err) {
				log.Println("keystore file already exist, skip creating")

			} else {
				if err := prompt2Passwords(); err != nil {
					return err
				}

				out, err := promptKeygenInput()
				if err != nil {
					return err
				}

				genkey := exec.Command(
					"keytool",
					"-genkey", "-v",
					"-keystore", keystorePath,
					"-keyalg", "RSA", "-keysize", "2048",
					"-validity", "10000", "-alias", "upload",
					"-storepass", password,
					"-dname", *out,
				)

				genkey.Stderr = os.Stderr
				genkey.Stdin = os.Stdin
				genkey.Stdout = os.Stdout

				err = genkey.Run()
				if err != nil {
					return errors.Wrap(err, "err execute")
				}

			}
			if err := createKeyFile(); err != nil {
				return errors.Wrap(err, "err create key file")
			}

			if err := modifyBuildGradle(); err != nil {
				return err
			}
			return nil
		},
	}
	keystoreCmd.Flags().StringP(keystoreFlag, "s", "~/upload-keystore.jks", "place to store key")

	return keystoreCmd
}

/*
-storepass
CN=commonName
OU=organizationUnit
O=organizationName
L=localityName
S=stateName
C=country */
func promptKeygenInput() (*string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nEnter your name: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("Enter your organization name: ")
	organizationName, _ := reader.ReadString('\n')
	fmt.Print("Enter your organization unit: ")
	organizationUnit, _ := reader.ReadString('\n')
	fmt.Print("Enter your city or locality: ")
	city, _ := reader.ReadString('\n')
	fmt.Print("Enter your state or province: ")
	state, _ := reader.ReadString('\n')
	fmt.Print("Enter two-letter country code: ")
	country, _ := reader.ReadString('\n')

	genvalues := fmt.Sprintf(
		"CN=%s OU=%s O=%s L=%s S=%s C=%s",
		strings.TrimSuffix(name, "\n"), strings.TrimSuffix(organizationUnit, "\n"),
		strings.TrimSuffix(organizationName, "\n"), strings.TrimSuffix(city, "\n"),
		strings.TrimSuffix(state, "\n"), strings.TrimSuffix(country, "\n"),
	)

	return &genvalues, nil
}

func prompt2Passwords() error {
	attempt := 1
	for {
		if attempt > 3 {
			log.Println("Too many failures - try later")
			os.Exit(1)
			break
		}

		promptAndVerifyPassword()

		fmt.Print("\nRe-enter your new password: ")
		input2, _ := term.ReadPassword(0)
		verifyPassword := string(input2)

		if password != verifyPassword {
			log.Println("\nPasswords don't match")
		} else if len(password) < 6 {
			log.Println("\nPassword must have at least 6 letters")
		} else {
			break
		}
		attempt++
	}

	return nil
}

func promptAndVerifyPassword() error {
	attempt := 1
	for {
		if attempt > 3 {
			log.Println("Too many failures - try later")
			os.Exit(1)
			break
		}

		fmt.Print("Enter your new password: ")
		input1, _ := term.ReadPassword(0)
		if len(input1) < 6 {
			log.Println("\nPassword must have at least 6 letters")
		} else {
			password = string(input1)
			break
		}

		attempt++
	}
	return nil
}

func createKeyFile() error {
	homeDir, _ := os.UserHomeDir()
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to Getwd")
	}

	keyPath := filepath.Join(cwd, "android/key.properties")
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}

	defer keyFile.Close()

	if len(password) == 0 {

		fmt.Print("Enter password for your keystore file: ")
		input1, _ := term.ReadPassword(0)
		password = string(input1)
	}

	keyFile.WriteString(
		fmt.Sprintf("storePassword=%s\nkeyPassword=%s\nkeyAlias=upload\nstoreFile=%s",
			password,
			password,
			filepath.Join(homeDir, "/upload-keystore.jks"),
		),
	)

	return nil
}

func modifyBuildGradle() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to Getwd")
	}

	buildGradlePath := filepath.Join(cwd, "/android/app/build.gradle")

	data, _ := ioutil.ReadFile(buildGradlePath)

	updatedContent := strings.Replace(string(data), "android {", replacement1, 1)
	updatedContent = strings.Replace(updatedContent, "buildTypes {", replacement2, 1)
	updatedContent = strings.Replace(updatedContent, "signingConfig signingConfigs.debug", replacement3, 1)

	err = ioutil.WriteFile(buildGradlePath, []byte(updatedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
