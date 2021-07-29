package flutter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
)

// keytool -genkey -v -keystore ~/upload-keystore.jks -keyalg RSA -keysize 2048 -validity 10000 -alias upload

// UploadKeystore create keystore and modify build.gradle
func UploadKeystore() *cobra.Command {
	keystoreCmd := &cobra.Command{
		Use:   "uploadkeystore",
		Short: "create upload keystore for android",
		Long:  "create upload keystore for android and modify build.gradle, all in one ",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			homeDir, _ := os.UserHomeDir()
			genkey := exec.Command(
				"keytool",
				// "-genkey -v -keystore ~/upload-keystore.jks -keyalg RSA -keysize 2048 -validity 10000 -alias upload",
				"-genkey", "-v",
				"-keystore", filepath.Join(homeDir, "/upload-keystore.jks"),
				"-keyalg", "RSA", "-keysize", "2048",
				"-validity", "10000", "-alias", "upload",
			)
			genkey.StderrPipe()
			genkey.Stderr = os.Stderr
			genkey.Stdin = os.Stdin
			genkey.Stdout = os.Stdout
			//TODO: handle existing keyfile
			err := genkey.Run()
			if err != nil {
				return errors.Wrap(err, "err execute")
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
	return keystoreCmd
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password from previous step: ")
	password, _ := reader.ReadString('\n')
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
