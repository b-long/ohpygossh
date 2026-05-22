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
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Hello() string {
	return "Hello, world"
}

type KeysForSSH struct {
	CloudUser        string
	PrivKeyAbsPath   string
	PublicKeyAbsPath string
}

func GenerateKeysForSSH(destinationDir, cloudUser string) (KeysForSSH, error) {
	fmt.Printf("Generating keypair for user: %s", cloudUser)

	temp_file, err := os.CreateTemp(destinationDir, "id_rsa_test")
	if err != nil {
		return KeysForSSH{}, err
	}

	privateKeyFileName := temp_file.Name()
	cmd := fmt.Sprintf("ssh-keygen -b 2048 -t rsa -q -N '' -f %s <<<y >/dev/null 2>&1", privateKeyFileName)

	output, exec_cmd_err := exec.Command("bash", "-c", cmd).Output()
	if exec_cmd_err != nil {
		return KeysForSSH{}, exec_cmd_err
	}
	fmt.Println(output)

	publicKeyAbsPath, err := filepath.Abs(privateKeyFileName + ".pub")
	if err != nil {
		return KeysForSSH{}, err
	}

	privateKeyAbsPath, err := filepath.Abs(privateKeyFileName)
	if err != nil {
		return KeysForSSH{}, err
	}

	return KeysForSSH{
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

func GenerateKeyPairAndCloudInit(destinationDir, cloudUser string) (KeysAndInit, error) {
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

// newSSHClient dials an SSH connection using private-key authentication.
// FIXME: host key verification is not yet implemented; connections are
// currently vulnerable to MITM. See https://stackoverflow.com/a/63308243
func newSSHClient(hostname, username, privateKey string) (*ssh.Client, error) {
	keyBytes, err := os.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	cfg := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec // known limitation, see FIXME above
		Timeout:         60 * time.Second,
	}
	return ssh.Dial("tcp", net.JoinHostPort(hostname, "22"), cfg)
}

func Run(hostname, username, privateKey, command string) (string, error) {
	client, err := newSSHClient(hostname, username, privateKey)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(command)
	return b.String(), err
}

func Upload(hostname, username, privateKey, localPath, remotePath string) error {
	client, err := newSSHClient(hostname, username, privateKey)
	if err != nil {
		return err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(remoteFile, localFile); err != nil {
		remoteFile.Close()
		return err
	}
	return remoteFile.Close()
}

func Download(hostname, username, privateKey, remotePath, localPath string) error {
	client, err := newSSHClient(hostname, username, privateKey)
	if err != nil {
		return err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(localFile, remoteFile); err != nil {
		localFile.Close()
		return err
	}
	return localFile.Close()
}
