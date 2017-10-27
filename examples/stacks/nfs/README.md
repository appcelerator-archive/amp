NFS client
==========

http server mounting a NFS share, and serving the same files across the replicas.

Get the URL of your NFS server. If it's the one deployed with the AMP cluster, you have it in the outputs of the stack.

    $ export NFS_SERVER=fs-12345678.efs.us-west-2.amazonaws.com
    $ amp stack deploy -c ./nfs.yml

Optionally you can verify that the NFS endpoint has been correctly set on the service

    $ amp service inspect nfs_client | jq '.Spec.TaskTemplate.ContainerSpec.Mounts[0].VolumeOptions.DriverConfig.Options.o'

You can log in one of the containers, add a file in /html/:

    $ echo "nfs works!" > /html/nfs.html

The app will be available at [http://nfs.examples.local.appcelerator.io/nfs.html](http://nfs.examples.local.appcelerator.io/nfs.html)

Test with

    $ curl -i http://nfs.examples.local.appcelerator.io/nfs.html

or open it in a browser.
Repeating it will distribute the request on the replicas (the Server header is the container ID). The body should be the same, since the file is shared on the NFS server.

You can replace the Docker image in the stack file by the official `nginx` image, however you'll have to modify the docroot from `/html` to `/usr/share/nginx/html`, and you'll lose the Server header reflecting the container hostname.
