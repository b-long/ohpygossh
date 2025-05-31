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
	Err              string
}

func GenerateKeysForSsh(destinationDir string, cloudUser string) KeysForSsh {
	// Generate a new SSH key pair
	fmt.Printf("Generating keypair for user: %s", cloudUser)

	temp_file, err := os.CreateTemp(destinationDir, "id_rsa_test")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	privateKeyFileName := temp_file.Name()
	cmd := fmt.Sprintf("ssh-keygen -b 2048 -t rsa -q -N '' -f %s <<<y >/dev/null 2>&1", privateKeyFileName)

	output, exec_cmd_err := exec.Command("bash", "-c", cmd).Output()
	if exec_cmd_err != nil {
		fmt.Println(exec_cmd_err)
		os.Exit(1)
	}
	fmt.Println(output)

	publicKeyAbsPath, err := filepath.Abs("" + privateKeyFileName + ".pub")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	privateKeyAbsPath, err := filepath.Abs(privateKeyFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return KeysForSsh{
		CloudUser:        cloudUser,
		PublicKeyAbsPath: publicKeyAbsPath,
		PrivKeyAbsPath:   privateKeyAbsPath,
		Err:              fmt.Sprint(err),
	}
}

type KeysAndInit struct {
	CloudInitPath string
	SshKeyPath    string
	CloudUser     string
	Err           string
}

func GenerateKeyPairAndCloudInit(destinationDir string, cloudUser string) KeysAndInit {
	// Generate a new SSH key pair
	fmt.Printf("Generating keypair & cloud-init for user: %s", cloudUser)

	temp_file, err := os.CreateTemp(destinationDir, "id_rsa_test")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	privateKeyFileName := temp_file.Name()
	cmd := fmt.Sprintf("ssh-keygen -b 2048 -t rsa -q -N '' -f %s <<<y >/dev/null 2>&1", privateKeyFileName)

	output, exec_cmd_err := exec.Command("bash", "-c", cmd).Output()
	if exec_cmd_err != nil {
		fmt.Println(exec_cmd_err)
		os.Exit(1)
	}
	fmt.Println(output)

	// Get the public SSH key (file content)
	publicKey, err := os.ReadFile("" + privateKeyFileName + ".pub")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a new cloud-init file
	cloud_init_file_path := fmt.Sprintf("%s/cloud-init.yaml", destinationDir)
	cloudInitFile, err := os.Create(cloud_init_file_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write the cloud-init file
	cloudInitFile.WriteString(`
users:
- name: ` + cloudUser + `
  groups: users, admin
  sudo: ALL=(ALL) NOPASSWD:ALL
  ssh_authorized_keys:
    - ` + strings.TrimSpace(string(publicKey)) + `
`)
	defer cloudInitFile.Close()

	cloud_init_abs_path, err := filepath.Abs(cloudInitFile.Name())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// publicKeyAbsPath, err := filepath.Abs("" + privateKeyFileName + ".pub")
	privateKeyAbsPath, err := filepath.Abs(privateKeyFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(privateKeyAbsPath)

	// return cloud_init_abs_path, privateKeyAbsPath, err
	return KeysAndInit{
		CloudInitPath: cloud_init_abs_path,
		SshKeyPath:    privateKeyAbsPath,
		CloudUser:     cloudUser,
		Err:           fmt.Sprint(err),
	}
}

func GenerateShortUUID(length int) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, length/3*4) // Adjust for base64 encoding
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
	f, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer f.Close()

	// Chunk size
	const maxSize = 4

	// Create buffer
	byteBuffer := make([]byte, maxSize)

	key, err := ssh.ParsePrivateKey(byteBuffer)
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
	// signer, err := ssh.ParseRawPrivateKey([]byte(privateKey))
	// key, err := x509.ParsePKCS1PrivateKey([]byte(privateKey))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// signer, err := ssh.NewSignerFromKey(key)

	// ssh.PublicKeys(key),
	// ssh.PublicKeys(key),
	// ssh.ParseRawPrivateKey(privateKey),
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
