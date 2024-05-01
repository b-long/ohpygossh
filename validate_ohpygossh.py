"""
This file serves as a test of ohpygossh.
"""

import glob
from shutil import copy
from pathlib import Path
from os import listdir, chdir
from os.path import join
import subprocess
import tempfile
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
        ["vagrant", "ssh-config"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, text=True, cwd=vagrantfile_dir
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

            # commands = f"""cat {keys.PublicKeyAbsPath} | vagrant ssh --command "cat >> ~/.ssh/authorized_keys && echo 'SSH key copied' " """
            # run_multiple_commands(
            #     commands=commands,
            #     description_of_command="cat public key to authorized_keys",
            # )

            ssh_hostname = get_vagrant_ssh_field("HostName", vagrantfile_dir=DOCKER_IN_VAGRANT)
            # print(f"{ssh_hostname}")

            ssh_port = get_vagrant_ssh_field("Port", vagrantfile_dir=DOCKER_IN_VAGRANT)
            # print(f"{ssh_port}")
            run_multiple_commands(
                commands=f"""

echo "Testing generated keypair"
echo "Performing authenticaton as 'vagrant' user"

ssh -o StrictHostKeyChecking=no -i {keys.PrivKeyAbsPath} vagrant@{ssh_hostname} -p {ssh_port} 'echo "We did it!"'

""",
                description_of_command="Open connection, with ssh -i flag and run command",
            )

#             run_multiple_commands(
#                 commands=f"""

# echo "Testing generated keypair"
# echo "Performing authenticaton as '{keys_only_user}' user"

# ssh -i {keys.PrivKeyAbsPath} {keys_only_user}@{ssh_hostname} -p {ssh_port} -c 'echo "We did it!"'

# """,
#                 description_of_command="Open connection, with ssh -i flag and run command",
#             )

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
        vagrant_test_dir = this_file_dir / "ssh-servers" / "new"

        with tempfile.TemporaryDirectory() as tmpDir:
            print("Created temporary directory", tmpDir)

            td_path = Path(tmpDir)
            copied_vagrant = td_path / "vm"

            # Copy all files
            for filename in glob.glob(join(vagrant_test_dir, "*.*")):
                copy(filename, tmpDir)

            # Copy the Vagrantifle
            for filename in glob.glob(join(vagrant_test_dir, "Vagrantfile")):
                copy(filename, tmpDir)

            kai: KeysAndInit = GenerateKeyPairAndCloudInit(tmpDir, "validation-user")
            # Prints the entire struct, like:
            # print(kai)
            #   gohpygossh.KeysAndInit{CloudInitPath=/var/folders/r7/srtk3z715s1bqzq2xy0mlsk80000gn/T/tmp56kflw4t/cloud-init.yaml, Err=<nil>, SshKeyPath=/var/folders/r7/srtk3z715s1bqzq2xy0mlsk80000gn/T/tmp56kflw4t/id_rsa_test2220445842, handle=1}

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

            # At this point, we can run
            # ssh -o StrictHostKeyChecking=no -i id_rsa_test2669721341 cloud-user@192.168.56.10 'echo "Connection success"'
            ssh_cmd = f"ssh -o StrictHostKeyChecking=no -i {kai.SshKeyPath} cloud-user@192.168.56.10 'echo \"Connection success\"'"
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

        # breakpoint()

    except Exception as e:
        raise RuntimeError("An unexpected error occurred testing ohpygossh") from e


if __name__ == "__main__":
    test_say_hello()

    # test_with_cloud_init()

    test_with_keys_only()
    # ssh_port = get_vagrant_ssh_field("Port", vagrantfile_dir=DOCKER_IN_VAGRANT)
    # print(f"{ssh_port}")
