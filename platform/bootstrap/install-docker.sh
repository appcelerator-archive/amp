_install_docker() {
  local _release=$(lsb_release -is)
  wget -qO- https://get.docker.com/ | sh
  if [[ $(grep -w docker /etc/group) ]]; then
    [[ "x$_release" = "xUbuntu" ]] && usermod -G docker ubuntu || true
    [[ "x$_release" = "xDebian" ]] && usermod -G docker admin  || true
  fi
  systemctl enable docker.service || true
}
