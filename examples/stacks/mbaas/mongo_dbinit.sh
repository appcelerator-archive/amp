#!/bin/bash

# based on https://github.com/appcelerator/arrow-cloud-docker/blob/master/swarm/scripts/mongo-configure-replica-set.sh

primary_host_ip=mongo-primary
primary_host_port=27017
admin_username=${MBAAS_MONGO_ADMIN_USERNAME:-appcadmin}
admin_password=${MBAAS_MONGO_ADMIN_PASSWORD:-cocoafish}
products="${MBAAS_MONGO_PRODUCTS:-arrowcloud arrowdb dashboard}"
arrowcloud_username=${MBAAS_MONGO_ARROWCLOUD_USERNAME:-appcelerator}
arrowcloud_password=${MBAAS_MONGO_ARROWCLOUD_PASSWORD:-cocoafish}
arrowdb_username=${MBAAS_MONGO_ARROWDB_USERNAME:-appcelerator}
arrowdb_password=${MBAAS_MONGO_ARROWDB_PASSWORD:-cocoafish}
dashboard_username=${MBAAS_MONGO_DASHBOARD_USERNAME:-appcelerator}
dashboard_password=${MBAAS_MONGO_DASHBOARD_PASSWORD:-cocoafish}

# return 0 means master is ready
check_is_master() {
	local primary_host_ip=$1
	local master=`mongo ${primary_host_ip}:${primary_host_port} --quiet --eval "db.isMaster().ismaster"`
	if [ "${master}" = "true" ]; then
		return 0
	else
		return 1
	fi
}

# $1 is the primary container service name
# $2 is the secondary container service name
# $3 is the arbiter container service name
# assign the configuration to a variable "CONFIG"
generate_rs_config() {
	# CONFIG="{"_id" : "data","members" : [{"_id" : 0,"host" : "$1:27017","priority": 100},{"_id" : 1,"host" : "$2:27017","priority": 1},{"_id" : 2,"host" : "$3:27017","arbiterOnly": true}]}"
	CONFIG="{\"_id\" : \"data\",
				\"members\" : [
					{
						\"_id\" : 0,
						\"host\" : \"$1:27017\",
						\"priority\": 100
					},
					{
						\"_id\" : 1,
						\"host\" : \"$2:27017\",
						\"priority\": 1
					},
					{
						\"_id\" : 2,
						\"host\" : \"$3:27017\",
						\"arbiterOnly\": true
					}
				]
			}"
}

# $1 is the primary host IP address
initiate_rs() {
	local primary_host_ip=$1
	mongo ${primary_host_ip}:${primary_host_port} --eval "rs.initiate(${CONFIG})"
	return $?
}

# $1 is the primary host IP address
# $2 is the username of admin user
# $3 is the password of that user
create_admin_user() {
	local primary_host_ip=$1
	local admin_username=$2
	local admin_password=$3

	echo "will add the admin user"
	mongo ${primary_host_ip}:${primary_host_port} --eval "db=db.getSiblingDB(\"admin\"); db.createUser( { user: \"${admin_username}\", pwd: \"${admin_password}\", roles: [ \"userAdminAnyDatabase\", \"clusterAdmin\", \"readAnyDatabase\", \"readWriteAnyDatabase\", \"dbAdminAnyDatabase\", { db: \"local\", role: \"readWrite\" } ] })"
	if [ $? -ne 0 ]; then
		echo "error create admin user"
		exit 1
	fi
}

create_user() {
	local primary_host_ip=$1
	local username=$2
	local password=$3
	local product=$4
	local dbs

	case $product in
		arrowcloud)
			dbs=("arrowcloud" "arrowcloud_log" "registry_auth")
			;;
		arrowdb)
			dbs=("acs_api" "acs_log" "push_notification_logs" "acs_delayed_jobs" "acs_global_apps")
			;;
		dashboard)
			dbs=("360" "360_sessions")
			;;
		*)
			;;
	esac
  echo "adding user $username for dbs $dbs"

	for db in "${dbs[@]}"; do
		mongo ${primary_host_ip}:${primary_host_port} --eval "db=db.getSiblingDB(\"${db}\"); db.createUser( { user: \"${username}\", pwd: \"${password}\", roles: [ \"readWrite\" ] } )"
		if [ $? != 0 ]; then
			echo "error create user in db ${db}"
			exit 1
		fi
	done
}

generate_rs_config "mongo-primary" "mongo-secondary" "mongo-arbiter"
initiate_rs ${primary_host_ip}
rc=$?
if [ ${rc} -ne 0 ]; then
	echo "errors when rs.initiate"
	exit 1
fi
echo "waiting for $primary_host_ip to be a master" >&2
SECONDS=0
while [[ $SECONDS -lt 20 ]]; do
  check_is_master ${primary_host_ip}
  rc=$?
  if [ ${rc} -ne 0 ]; then
    sleep 3
	  continue
  else
    echo "$primary_host_ip is a master node ($SECONDS sec after)" >&2
    echo "waiting 10 more sec before creating the users..." >&2
		sleep 10
		break
	fi
done
if [ ${rc} -ne 0 ]; then
	echo "${primary_host_ip} is not master yet"
	exit 1
fi

create_admin_user ${primary_host_ip} "${admin_username}" "${admin_password}"
create_user ${primary_host_ip} "${arrowcloud_username}" "${arrowcloud_password}" "arrowcloud"
create_user ${primary_host_ip} "${arrowdb_username}" "${arrowdb_password}" "arrowdb"
create_user ${primary_host_ip} "${dashboard_username}" "${dashboard_password}" "dashboard"

echo DONE
