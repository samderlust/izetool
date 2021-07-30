package flutter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
				log.Println("create in homedir")
				homeDir, _ := os.UserHomeDir()
				extractedPath := strings.Split(keyFromFLag, "~")[1]
				fmt.Printf("extractedPath: %v\n", extractedPath)
				keystorePath = filepath.Join(homeDir, extractedPath)
			}
			fmt.Printf("keystorePath: %v\n", keystorePath)

			_, err := os.Stat(keystorePath)
			if os.IsNotExist(err) {
				log.Println("keystore file does not exist,")
			} else {
				log.Println("keystore file already exist, skip creating")
			}

			// genkey := exec.Command(
			// 	"keytool",
			// 	// "-genkey -v -keystore ~/atestkey.jks -keyalg RSA -keysize 2048 -validity 10000 -alias upload",
			// 	"-genkey", "-v",
			// 	"-keystore", keystorePath,
			// 	"-keyalg", "RSA", "-keysize", "2048",
			// 	"-validity", "10000", "-alias", "upload",
			// )
			// genkey.StderrPipe()
			// genkey.Stderr = os.Stderr
			// genkey.Stdin = os.Stdin
			// genkey.Stdout = os.Stdout
			// //TODO: handle existing keyfile
			// err = genkey.Run()
			// if err != nil {
			// 	return errors.Wrap(err, "err execute")
			// }

			// if err := createKeyFile(); err != nil {
			// 	return errors.Wrap(err, "err create key file")
			// }

			// if err := modifyBuildGradle(); err != nil {
			// 	return err
			// }

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

	fmt.Print("Enter your new password: ")
	password, _ = reader.ReadString('\n')
	fmt.Print("Re-enter your new password: ")
	verifyPassword, _ := reader.ReadString('\n')
	if password != verifyPassword {
		return nil, errors.New("Passwords don't match")
	}

	fmt.Print("Enter your name: ")
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

	genvalues := fmt.Sprintf("-storepass %s -dname 'CN=%s OU=%s O=%s L=%s S=%s C=%s'",
		password, name, organizationUnit, organizationName, city, state, country,
	)

	return &genvalues, nil
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

	// reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password for your keystore file: ")
	// password, _ := reader.ReadString('\n')
	trimmedPassword := strings.Trim(password, "\n")

	keyFile.WriteString(
		fmt.Sprintf("storePassword=%s\nkeyPassword=%s\nkeyAlias=upload\nstoreFile=%s",
			trimmedPassword,
			trimmedPassword,
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
