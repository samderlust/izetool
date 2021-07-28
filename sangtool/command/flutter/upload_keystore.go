package flutter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const ()

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
			err := genkey.Run()
			if err != nil {
				return errors.Wrap(err, "err execute")
			}

			if err := createKeyFile(); err != nil {
				return errors.Wrap(err, "err create key file")
			}
			return nil
		},
	}
	return keystoreCmd
}

func createKeyFile() error {
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

	fmt.Print("Enter password from previous step")
	var input string
	fmt.Sprintln(&input)

	keyFile.WriteString(
		fmt.Sprintf("storePassword=%s\nkeyPassword=%s\nkeyAlias=upload\nstoreFile=~/upload-keystore.jks", input, input))

	return nil
}
