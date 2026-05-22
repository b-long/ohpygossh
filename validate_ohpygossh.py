"""
This file serves as a test of ohpygossh.
"""

import glob
import sys
import json
from shutil import copy, which
from pathlib import Path
import os
from os import listdir, chdir
from os.path import join
import subprocess
import tempfile
import time
from shlex import split

GIT_ROOT = Path(__file__).parent.resolve()
DOCKER_IN_VAGRANT = GIT_ROOT / "ssh-servers" / "docker-in-vagrant"


def run_multiple_commands(commands, description_of_command="") -> tuple[str, str]:
    process = subprocess.Popen(
        "/bin/bash", stdin=subprocess.PIPE, stdout=subprocess.PIPE, text=True
    )
    out, err = process.communicate(commands)

    print(
        f"""Start multi-command output -- {description_of_command}

STDOUT
{out}

STDERR
{err}

Finished -- {description_of_command}
    """
    )


def get_vagrant_ssh_field(field: str, vagrantfile_dir: Path):
    """Retrieves some 'ssh-confg' from the active Vagrant VM."""
    # Captures the output of vagrant ssh-config
    process = subprocess.Popen(
        ["vagrant", "ssh-config"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        text=True,
        cwd=vagrantfile_dir,
    )

    ssh_config_output, err = process.communicate()

    # Iterates over lines and searches for the line starting
    # with a given field (e.g. "Port" or "HostName").
    #
    # Note, these fields usually are prefixed with whitespace
    for line in ssh_config_output.splitlines():
        if line.lstrip().startswith(field):
            # Extracts the hostname (second word) and returns it
            return line.split()[1]

    # If no "HostName" line is found, raise an exception
    raise Exception(f"Failed to find {field=} in vagrant ssh-config")


def test_with_keys_only():
    keys_only_user = "keysOnlyUser"
    try:
        print("Python validation starting, for 'test_with_keys_only'")

        from ohpygossh.gohpygossh import GenerateKeysForSsh, KeysForSsh

        with tempfile.TemporaryDirectory() as tmpDir:
            print("'ohpygossh': Generating SSH key pair")
            keys: KeysForSsh = GenerateKeysForSsh(tmpDir, keys_only_user)

            print(
                f"""


Keypair written to disk

{keys.PrivKeyAbsPath=}

{keys.PublicKeyAbsPath=}

"""
            )

            commands = f"""

echo "Gathering SSH public key content 'ohpygossh' public SSH key has content"


# Change directory
cd {DOCKER_IN_VAGRANT}

cat {keys.PublicKeyAbsPath} > new_id_rsa.pub

# List content in the directory
ls -la

# Launch machine
vagrant up --provider=docker

echo "Checking vagrant status"
vagrant status && echo "That's the status"
"""
            run_multiple_commands(
                commands=commands, description_of_command="Vagrant up, etc."
            )

            ssh_hostname = get_vagrant_ssh_field(
                "HostName", vagrantfile_dir=DOCKER_IN_VAGRANT
            )

            ssh_port = get_vagrant_ssh_field("Port", vagrantfile_dir=DOCKER_IN_VAGRANT)
            run_multiple_commands(
                commands=f"""

echo "Testing generated keypair"
echo "Performing authentication as 'vagrant' user"

ssh -o StrictHostKeyChecking=no -i {keys.PrivKeyAbsPath} vagrant@{ssh_hostname} -p {ssh_port} 'echo "We did it!"'

""",
                description_of_command="Open connection, with ssh -i flag and run command",
            )

    except Exception as e:
        raise RuntimeError("An unexpected error occurred testing ohpygossh") from e


def test_say_hello():
    try:
        print("Python validation starting...")

        from ohpygossh.gohpygossh import Hello

        # Prints a basic str value: 'Hello, world'
        print(Hello())
    except Exception as e:
        raise RuntimeError("An unexpected error occurred testing: Hello") from e


def test_with_cloud_init():
    try:
        print("Python validation starting, for 'test_with_cloud_init'")

        from ohpygossh.gohpygossh import GenerateShortUUID

        # Prints str like 'my-generated-id-MJKB'
        print(f"my-generated-id-{GenerateShortUUID(length=4)}")

        from ohpygossh.gohpygossh import GenerateKeyPairAndCloudInit, KeysAndInit

        this_file_dir = Path(__file__).parent
        vagrant_test_dir = this_file_dir / "ssh-servers" / "cloud-init"

        with tempfile.TemporaryDirectory() as tmpDir:
            print("Created temporary directory", tmpDir)

            # Copy all files
            for filename in glob.glob(join(vagrant_test_dir, "*.*")):
                copy(filename, tmpDir)

            # Copy the Vagrantifle
            for filename in glob.glob(join(vagrant_test_dir, "Vagrantfile")):
                copy(filename, tmpDir)

            kai: KeysAndInit = GenerateKeyPairAndCloudInit(tmpDir, "validation-user")

            print(listdir(tmpDir))

            # Change the working directory
            chdir(tmpDir)

            # Run vagrant up command
            vm_process = subprocess.Popen(
                ["vagrant", "up"], stdout=subprocess.PIPE, stderr=subprocess.PIPE
            )

            # Wait for the process to finish and capture output
            stdout, stderr = vm_process.communicate()

            # Check for errors
            if stderr:
                print("Error bringing up Vagrant VM:", stderr.decode("utf-8"))
            else:
                print("Vagrant VM successfully brought up.")

            ssh_cmd = f"ssh -o StrictHostKeyChecking=no -i {kai.SshKeyPath} {kai.CloudUser}@192.168.56.10 'echo \"Connection success\"'"
            split_ssh_cmd = split(ssh_cmd)
            ssh_process = subprocess.Popen(
                split_ssh_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE
            )

            stdout, stderr = ssh_process.communicate()

            # Check for errors
            if ssh_process.returncode != 0:
                print("Error running SSH command:", stderr.decode("utf-8"))
            else:
                print(
                    f"""SSH command successfully executed.

SSH result:
{stdout.decode("utf-8")}

"""
                )

    except Exception as e:
        raise RuntimeError("An unexpected error occurred testing ohpygossh") from e


def _multipass_cmd() -> list:
    """On Linux snap installs the multipassd socket is root-owned; prefix with sudo."""
    if sys.platform == "linux" and os.geteuid() != 0 and which("sudo"):
        return ["sudo", "multipass"]
    return ["multipass"]


def test_with_multipass():
    if not which("multipass"):
        print("multipass not installed, skipping test_with_multipass")
        return

    try:
        print("Python validation starting, for 'test_with_multipass'")

        from ohpygossh.gohpygossh import (
            Download,
            GenerateKeyPairAndCloudInit,
            GenerateShortUUID,
            Run,
            Upload,
        )

        short_id = GenerateShortUUID(4)
        vm_name = f"ohpytest-{short_id}".lower().replace("_", "-")
        mp = _multipass_cmd()

        # multipass snap AppArmor profile allows @{HOME}/** but not /tmp/**;
        # create the temp dir under $HOME so the daemon can read the cloud-init file.
        with tempfile.TemporaryDirectory(dir=Path.home()) as tmpDir:
            kai = GenerateKeyPairAndCloudInit(tmpDir, "cloud-user")

            print(f"Launching multipass VM: {vm_name}")
            try:
                launch_cmd = [
                    *mp,
                    "launch",
                    "lts",
                    "--name",
                    vm_name,
                    "--cpus",
                    "1",
                    "--memory",
                    "1G",
                    "--disk",
                    "5G",
                    "--cloud-init",
                    kai.CloudInitPath,
                ]
                for attempt, wait in enumerate([0, 10, 30, 90]):
                    if wait:
                        print(
                            f"Launch attempt {attempt} failed, retrying in {wait}s..."
                        )
                        time.sleep(wait)
                    if subprocess.run(launch_cmd).returncode == 0:
                        break
                else:
                    raise RuntimeError("multipass launch failed after 4 attempts")

                info_result = subprocess.run(
                    [*mp, "info", "--format", "json", vm_name],
                    capture_output=True,
                    text=True,
                    check=True,
                )
                info_data = json.loads(info_result.stdout)
                vm_info = info_data.get("info", {}).get(vm_name, {})
                ipv4_addresses = vm_info.get("ipv4", [])
                if not ipv4_addresses:
                    raise RuntimeError(
                        f"No IPv4 address found for VM {vm_name}. Info: {info_data}"
                    )
                ip = ipv4_addresses[0]
                print(f"VM IP: {ip}")

                output = Run(
                    ip, kai.CloudUser, kai.SshKeyPath, "echo 'Connection success'"
                )

                if "Connection success" not in output:
                    raise RuntimeError(f"Unexpected SSH output: {output!r}")

                print(f"Run e2e test passed. Output: {output!r}")

                # SFTP round-trip: upload a file, verify it landed, download it back.
                unique_marker = GenerateShortUUID(8)
                upload_content = f"ohpygossh-sftp-test:{unique_marker}\n"
                local_upload = os.path.join(tmpDir, "upload.txt")
                remote_path = f"/home/{kai.CloudUser}/uploaded.txt"
                local_download = os.path.join(tmpDir, "download.txt")

                Path(local_upload).write_text(upload_content)

                Upload(ip, kai.CloudUser, kai.SshKeyPath, local_upload, remote_path)
                print(f"Upload complete: {local_upload} -> {remote_path}")

                cat_output = Run(
                    ip, kai.CloudUser, kai.SshKeyPath, f"cat {remote_path}"
                )
                if cat_output.strip() != upload_content.strip():
                    raise RuntimeError(
                        f"Remote file content mismatch after upload. Got: {cat_output!r}"
                    )
                print("Remote file content verified via Run/cat.")

                Download(ip, kai.CloudUser, kai.SshKeyPath, remote_path, local_download)
                print(f"Download complete: {remote_path} -> {local_download}")

                downloaded_content = Path(local_download).read_text()
                if downloaded_content != upload_content:
                    raise RuntimeError(
                        f"Downloaded content does not match original. Got: {downloaded_content!r}"
                    )

                print(f"SFTP round-trip e2e test passed. Marker: {unique_marker!r}")

            finally:
                print(f"Deleting multipass VM: {vm_name}")
                subprocess.run([*mp, "delete", "--purge", vm_name], check=False)

    except Exception as e:
        raise RuntimeError(
            "An unexpected error occurred testing ohpygossh with multipass"
        ) from e


if __name__ == "__main__":
    test_say_hello()

    test_with_multipass()

    if which("vagrant"):
        test_with_keys_only()
    else:
        print("vagrant not installed, skipping test_with_keys_only")
