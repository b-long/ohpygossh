package gohpygossh

/*
Main based on
https://gist.github.com/iamralch/b7f56afc966a6b6ac2fc
*/
import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func Hello() string {
	return "Hello, world"
}

type KeysForSsh struct {
	CloudUser        string
	PrivKeyAbsPath   string
	PublicKeyAbsPath string
}

func GenerateKeysForSsh(destinationDir string, cloudUser string) (KeysForSsh, error) {
	fmt.Printf("Generating keypair for user: %s", cloudUser)

	temp_file, err := os.CreateTemp(destinationDir, "id_rsa_test")
	if err != nil {
		return KeysForSsh{}, err
	}

	privateKeyFileName := temp_file.Name()
	cmd := fmt.Sprintf("ssh-keygen -b 2048 -t rsa -q -N '' -f %s <<<y >/dev/null 2>&1", privateKeyFileName)

	output, exec_cmd_err := exec.Command("bash", "-c", cmd).Output()
	if exec_cmd_err != nil {
		return KeysForSsh{}, exec_cmd_err
	}
	fmt.Println(output)

	publicKeyAbsPath, err := filepath.Abs(privateKeyFileName + ".pub")
	if err != nil {
		return KeysForSsh{}, err
	}

	privateKeyAbsPath, err := filepath.Abs(privateKeyFileName)
	if err != nil {
		return KeysForSsh{}, err
	}

	return KeysForSsh{
		CloudUser:        cloudUser,
		PublicKeyAbsPath: publicKeyAbsPath,
		PrivKeyAbsPath:   privateKeyAbsPath,
	}, nil
}

type KeysAndInit struct {
	CloudInitPath string
	SshKeyPath    string
	CloudUser     string
}

func GenerateKeyPairAndCloudInit(destinationDir string, cloudUser string) (KeysAndInit, error) {
	fmt.Printf("Generating keypair & cloud-init for user: %s", cloudUser)

	temp_file, err := os.CreateTemp(destinationDir, "id_rsa_test")
	if err != nil {
		return KeysAndInit{}, err
	}

	privateKeyFileName := temp_file.Name()
	cmd := fmt.Sprintf("ssh-keygen -b 2048 -t rsa -q -N '' -f %s <<<y >/dev/null 2>&1", privateKeyFileName)

	output, exec_cmd_err := exec.Command("bash", "-c", cmd).Output()
	if exec_cmd_err != nil {
		return KeysAndInit{}, exec_cmd_err
	}
	fmt.Println(output)

	publicKey, err := os.ReadFile(privateKeyFileName + ".pub")
	if err != nil {
		return KeysAndInit{}, err
	}

	cloud_init_file_path := fmt.Sprintf("%s/cloud-init.yaml", destinationDir)
	cloudInitFile, err := os.Create(cloud_init_file_path)
	if err != nil {
		return KeysAndInit{}, err
	}
	defer cloudInitFile.Close()

	cloudInitFile.WriteString(`
users:
- name: ` + cloudUser + `
  groups: users, admin
  sudo: ALL=(ALL) NOPASSWD:ALL
  ssh_authorized_keys:
    - ` + strings.TrimSpace(string(publicKey)) + `
`)

	cloud_init_abs_path, err := filepath.Abs(cloudInitFile.Name())
	if err != nil {
		return KeysAndInit{}, err
	}

	privateKeyAbsPath, err := filepath.Abs(privateKeyFileName)
	if err != nil {
		return KeysAndInit{}, err
	}

	return KeysAndInit{
		CloudInitPath: cloud_init_abs_path,
		SshKeyPath:    privateKeyAbsPath,
		CloudUser:     cloudUser,
	}, nil
}

func GenerateShortUUID(length int) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, (length*3+3)/4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Base64 encode without padding
	encodedString := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Truncate to desired length and remove trailing '=' padding
	return encodedString[:length], nil
}

func PublicKeyFile(file string) ssh.AuthMethod {
	keyBytes, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func Run(hostname string, username string, privateKey string, command string) (string, error) {
	// envVars map[string]string) string {
	// Establish an SSH connection to the remote server
	// key1, err := ssh.ParsePrivateKey([]byte(privateKey))
	keyBytes, err := os.ReadFile(privateKey)
	if err != nil {
		return "", err
	}
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return "", err
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		/*

			FIXME: Work on supporting host key verification.

			See: https://stackoverflow.com/a/63308243

		*/
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// Env: envVars,
		// env: envVars,
		Timeout: time.Duration(time.Duration.Seconds(60)),
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(hostname, "22"), sshConfig)
	if err != nil {
		return "", err
	}
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer  // import "bytes"
	session.Stdout = &b // get output
	// you can also pass what gets input to the stdin, allowing you to pipe
	// content from client to server
	//      session.Stdin = bytes.NewBufferString("My input")

	// Finally, run the command
	err = session.Run(command)
	return b.String(), err
}
