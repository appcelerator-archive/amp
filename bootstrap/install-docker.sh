_install_docker() {
  wget -qO- https://get.docker.com/ | sh
  usermod -G docker ubuntu
  systemctl enable docker.service
}
