package gohpygossh_test

import (
	"fmt"
	"gohpygossh"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/larstobi/go-multipass/multipass"
)

func TestHello(t *testing.T) {
	got := gohpygossh.Hello()
	want := "Hello, world"
	if got != want {
		t.Error("Unexpected value")
	}
}

func TestGenerateKeyPairAndCloudInit(t *testing.T) {
	// Create a temporary directory for multipass operations
	tmpDir, err := os.MkdirTemp("", "g-opg-test")

	if err != nil {
		t.Fatal(err)
	}

	r := gohpygossh.GenerateKeyPairAndCloudInit(tmpDir, "cloud-user")

	if r.Err != "<nil>" {
		t.Fatal(err)
	}

	result, _ := os.ReadDir(tmpDir)
	gotFiles := len(result)
	// Directory should include key pair (2 files) and 1 cloud init file
	if gotFiles != 3 {
		t.Error("Unexpected value")
	}

	want := "id_rsa"
	if !strings.Contains(r.SshKeyPath, want) {
		t.Error("Unexpected value for 'ssh_key_abs_path")
	}

	want_cloud_init := "cloud-init.yaml"
	if !strings.Contains(r.CloudInitPath, want_cloud_init) {
		t.Error("Unexpected value for 'cloud_init_abs_path's")
	}

	defer os.RemoveAll(tmpDir)
}

func TestRunWithMultipass(t *testing.T) {
	// Install go-multipass if not already installed
	if _, err := exec.LookPath("multipass"); err != nil {
		t.Skip("multipass not installed")
	}

	// Create a temporary directory for multipass operations
	tmpDir, err := os.MkdirTemp("", "multipass-test2")

	if err != nil {
		t.Fatal(err)
	}

	// cloud_init_file, identity_file, err := gohpygossh.GenerateKeyPairAndCloudInit(tmpDir)
	r := gohpygossh.GenerateKeyPairAndCloudInit(tmpDir, "cloud-user")
	if r.Err != "<nil>" {
		t.Fatal(err)
	}

	// FIXME: Use dynamic VM name
	// short_id, _ := gohpygossh.GenerateShortUUID(4)
	// dyn_vm_name := fmt.Sprintf("testvm%s", short_id)
	dyn_vm_name := "myvm"

	instance, err := multipass.LaunchV2(&multipass.LaunchReqV2{
		Image:         "lts",
		CPUS:          "2",
		Disk:          "10g",
		Name:          dyn_vm_name,
		Memory:        "3g",
		CloudInitFile: r.CloudInitPath,
	})

	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	// Get the hostname of the multipass instance
	// hostname, err := instance.GetHostname()
	fmt.Println(instance.IP)

	info, err := multipass.Info(&multipass.InfoRequest{Name: dyn_vm_name})
	// hostname, err := multipass.Info(&multipass.InfoRequest{Name: "multipass-test"})
	hostname := info.IP
	if err != nil {
		t.Fatal(err)
	}

	// Execute the 'ls' command on the remote server
	// FIXME: 'cloud-user' is a magic string, we should pass into
	// the 'GenerateKeyPaiirAndCloudInit()' function
	output, err := gohpygossh.Run(hostname, "cloud-user", r.SshKeyPath, "echo 'Connection success'")
	if err != nil {
		t.Fatal(err)
	}

	// Verify the output contains expected files or directories
	if !strings.Contains(output, "Connection success") {
		t.Error("Expected 'Connection success' in output but found:", output)
	}

	// TODO: What exactly does this do?
	defer multipass.Delete(&multipass.DeleteRequest{Name: dyn_vm_name})
	defer os.RemoveAll(tmpDir)

}
