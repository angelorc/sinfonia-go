#===============================================================================
# STEP 1: Create node servers
#===============================================================================
resource "digitalocean_droplet" "bitsong-indexer-1" {
  image  = "ubuntu-22-04-x64"
  name   = "bitsong-indexer-1"
  region = "fra1"
  size   = "s-2vcpu-2gb"
  graceful_shutdown = true
  monitoring = true
  ssh_keys = [
    data.digitalocean_ssh_key.angelo.id
  ]

  connection {
    host        = self.ipv4_address
    user        = "root"
    type        = "ssh"
    private_key = file(var.pvt_key)
    timeout     = "2m"
  }

  provisioner "file" {
    destination = "/etc/apt/apt.conf.d/00auto-conf"
    content     = <<EOF
Dpkg::Options {
    "--force-confdef";
    "--force-confold";
}
    EOF
  }

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      "export DEBIAN_FRONTEND=noninteractive",
      "while ps aux | grep -q [a]pt; do sleep 1; done",
      "while fuser /var/lib/apt/lists/lock >/dev/null 2>&1; do sleep 1; done",
      "while fuser /var/lib/dpkg/lock >/dev/null 2>&1; do sleep 1; done",
      "apt-get -q update",
      "apt-get -q upgrade -y",
      "apt-get -q install -y apt-transport-https ca-certificates build-essential make ufw curl jq git",
      "sed -i 's/#\\?PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config",
      "systemctl restart ssh",
      "echo Installing golang...",
      "wget -q -O - https://git.io/vQhTU | bash -s -- --version 1.18.3",
      "source ~/.bashrc",
      "go version",
      "echo Installing sinfonia-go...",
      "cd $HOME",
      "git clone https://github.com/angelorc/sinfonia-go.git",
      "cd sinfonia-go/bitsong",
      "make install",
      "sinfonia-bitsong"
    ]
  }
}