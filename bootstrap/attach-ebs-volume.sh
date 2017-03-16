# _attach_ebs_volume
# attach and format an ebs volume
# if mountpoint = none, the device won't be mounted
# IAM profile should contain ec2:AttachVolume, ec2:CreateVolume, ec2:DescribeVolumes, ec2:CreateTags and ec2:ModifyInstanceAttribute
_attach_ebs_volume () {
    local device=$1
    local mountpoint=${2:-none}
    local desc=${3:-none}
    local size=${4:-8}
    local type=${5:-ext4}
    local delete_on_termination=${6:-true}

    local devicelocalname
    local volumetype=gp2
    local az=$(curl 169.254.169.254/latest/meta-data/placement/availability-zone 2>/dev/null)
    local region=${az%?}
    local instanceid=$(curl 169.254.169.254/latest/meta-data/instance-id 2>/dev/null)
    local mounttimeout=300
    local IOPSARGS=""

    which jq || return 1
    which aws || return 1

    if [ "x$delete_on_termination" != "xtrue" ]; then
        delete_on_termination=false
    fi
    devicelocalname=$(echo $device | sed 's/\/sd/\/xvd/')

    # if device is already used, abort
    ls $device 2>/dev/null || ls $devicelocalname 2>/dev/null && (echo "device $device is already used"; return 1)

    if [[ -b $devicelocalname ]]; then
        echo "Device $devicelocalname already exists, won't attach a new volume" >&2
        return 1
    fi

    echo "Creating EBS volume at ${mountpoint}..." >&2
    local result        # local always return an exit 0
    result=$(aws ec2 create-volume --region $region --size $size --availability-zone $az --volume-type $volumetype $IOPSARGS)
    if [[ $? -ne 0 ]]; then
        echo "Unable to create volume" >&2
        return 1
    fi
    echo "Getting volume ID..." >&2
    local volumeid      # local always return an exit 0
    volumeid=$(echo $result | jq -r '.VolumeId')
    if [[ $? -ne 0 ]]; then
        "Unable to get volume ID" >&2
        return 1
    fi

    echo "Wait for volume $volumeid ($device) to be available... " >&2
    local state=unknown
    local startat=$(date +%s)
    local enddate=$(( startat + mounttimeout ))
    while [ "$state" != "available" -a $(date +%s) -le $enddate ] ; do
        state=$(aws ec2 describe-volumes --region $region --volume-ids $volumeid | jq -r '.Volumes[0].State')
        if [[ $? -ne 0 ]]; then
            echo "Unable to get volume $volumeid state" >&2
            return 1
        fi
        if [ "$state" != "available" ] ; then
            sleep 1
        fi
    done

    if [ "$state" != "available" ]; then
        echo "Volume $volumeid ($device) is not available after $mounttimeout sec" >&2
        return 1
    fi
    echo "Volume state: $state" >&2

    if [[ "x$desc" != "xnone" ]]; then
        aws ec2 create-tags --region $region --resources $volumeid --tags "Key=Name,Value=$desc"
        if [[ $? -ne 0 ]]; then
            echo "Unable to tag the volume $volumeid, abort" >&2
            return 3
        fi
    fi

    echo "Attaching volume $volumeid to $device... " >&2
    aws ec2 attach-volume --region $region --volume-id $volumeid --instance-id $instanceid --device "$device"
    if [[ $? -ne 0 ]]; then
        "Unable to attach volume $volumeid" >&2
        return 1
    fi
    echo "Wait for volume $volumeid to be attached... " >&2
    local state=unknown
    local startat=$(date +%s)
    local enddate=$(( startat + mounttimeout ))
    while [ "$state" != "in-use" -a $(date +%s) -le $enddate ] ; do
        state=$(aws ec2 describe-volumes --region $region --volume-ids $volumeid | jq -r '.Volumes[0].State')
        if [[ $? -ne 0 ]]; then
            echo "Unable to get volume $volumeid state"
            return 1
        fi
        if [ "$state" != "in-use" ] ; then
            sleep 1
        fi
    done

    if [ "$state" != "in-use" ]; then
        echo "Unable to attach volume $volumeid: State=$state"
        return 1
    fi

    echo "Volume state: $state"

    echo "Wait for volume $volumeid to be visible from the box... " >&2
    startat=$(date +%s)
    enddate=$(( startat + mounttimeout ))

    echo "--------------------------------------------------------" >&2
    while [ $(date +%s) -le $enddate ] ; do
        echo "Test presence of $devicelocalname..."
        fdisk -l  2>&1 | sed "s/^/    > /" # Add spaces in front of each line to be more readable in log
        if ( fdisk -l 2>/dev/null | grep -q "$devicelocalname" ) ; then
            break
        fi
        sleep 10
    done
    fdisk -l  2>&1 | sed "s/^/    > /" # Add spaces in front of each line to be more readable in log
    if ! ( fdisk -l 2>/dev/null | grep -q "$devicelocalname" ) ; then
        echo "Unable to see volume $devicelocalname from box" >&2
        return 1
    fi
    echo "[ OK ]" >&2
    echo "--------------------------------------------------------" >&2

    echo "Formatting partition $devicelocalname... " >&2
    /sbin/mkfs.$type $devicelocalname | sed "s/^/    > /"
    if [[ $? -ne 0 ]]; then
        echo "Unable to format partition $devicelocalname" >&2
        return 1
    fi
    echo "[ OK ]" >&2
    if [[ "x$mountpoint" != "xnone" ]]; then
      echo "Create directory $mountpoint... " >&2
      mkdir -p $mountpoint
      if [[ $? -ne 0 ]]; then
          echo "Unable to create directory $mountpoint" >&2
          return 1
      fi
      echo "[ OK ]"

      echo "Adding $mountpoint to /etc/fstab... " >&2
      cat >> /etc/fstab << EOF
$devicelocalname $mountpoint $type defaults 0 2
EOF

      if [[ $? -ne 0 ]]; then
          echo "Failed to add $mountpoint in fstab" >&2
          return 1
      fi
      echo "Mounting $devicelocalname on $mountpoint... " >&2
      mount $mountpoint
      if [[ $? -ne 0 ]]; then
          echo "Unable to mount $devicelocalname in directory $mountpoint"
          return 1
      fi
      # set the deleteOnTermination flag
      aws ec2 --region $region modify-instance-attribute --block-device-mappings "[{\"DeviceName\": \"$device\", \"Ebs\": {\"DeleteOnTermination\": $delete_on_termination}}]" --instance-id=$instanceid
      echo "[ OK ]"
    fi
}
