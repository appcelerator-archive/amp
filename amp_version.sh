# tries to get
_amp_version(){
  local _version
  local _default_version=latest
  if [ -n "${FORCE_AMP_VERSION-}" ]; then
    echo $FORCE_AMP_VERSION
    return 0
  fi
  git version >/dev/null || true
  if [ $? -ne 0 ]; then
    # git is not installed
    return 0
  fi
  if [ ! -d .git ]; then
    # not a git repo
    echo $_default_version
    return
  fi
  _version=$(git rev-parse --abbrev-ref HEAD)
  if [ "x$_version" = "xHEAD" ]; then
    # probably on a tag
    _version=$(git name-rev --name-only HEAD | sed 's/tags\/\(.*\)../\1/')
  fi

  if [ -z "$_version" ]; then
    echo $_default_value
  elif [ "$_version" = "master" ]; then
    echo $_default_value
  else
    echo $_version
  fi
}
