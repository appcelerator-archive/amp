_install_docker() {
  wget -qO- https://get.docker.com/ | sh
  grep -qw docker /etc/group && grep -qw ubuntu /etc/passwd && usermod -G docker ubuntu || true
  systemctl enable docker.service || true
}
