# wsl2安装goalng

```sh
sudo apt update
sudo apt install curl git
wget https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
go version
```

```shell
#!/bin/bash

# Update package list
sudo apt update

# Install required packages
sudo apt install curl git -y

# Download Go archive
wget https://go.dev/dl/go1.20.2.linux-amd64.tar.gz

# Extract archive and replace any existing installation
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz

# Add Go binary directory to PATH environment variable
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Reload .bashrc file to apply changes to the current shell session
source ~/.bashrc

# Verify Go installation
go version
```

