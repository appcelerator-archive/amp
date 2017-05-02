provider=""
curl -m 3 -f 169.254.169.254/latest/meta-data/ 2>/dev/null && provider=aws || true
awk -F/ '$2 == "docker"' /proc/self/cgroup | read && provider=docker || true
