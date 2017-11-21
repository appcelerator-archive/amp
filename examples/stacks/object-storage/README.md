Object store example
====================

Demonstrates how to use the object storage feature of amp.

One service will run a single task to upload a file to the S3 bucket.

Another service will run a task to check the existence of the file. It will restart until successfull.

### Cluster creation & sign up
    amp cluster create # --provider aws is the only implementation as of this writting
    amp -k user signup

### Create an object store

    amp object-store create OBJECT_STORE_NAME

If the bucket has already been created, you can retrieve the object store name with the CLI:

   amp object-store ls

### Deploying the stack

Run in this directory:

    export BUCKET_NAME=OBJECT_STORE_NAME
    amp -k stack deploy -c ./object-store.yml

### Checking the logs

   $ amp -k logs object-store_writer
   object-store_writer.1.ni6ptw8mtgw5@ip-192-168-0-186    | upload: '/etc/hostname' -> 's3://XXX/objstore-test.txt' (13 bytes in 0.0 seconds, 266.17 B/s) [1 of 1]
   amp -k logs object-store_reader
   object-store_reader.1.i5hb3salmijf@ip-192-168-16-30    | 2017-11-21 01:07        13   s3://XXX/objstore-test.txt

### cleanup

    amp object-store rm OBJECT_STORE_NAME

You should get an error because the s3 bucket is not empty. Use a flag to force the removal (make sure you don't make any mistake in the bucket name, you'll lose the data):

    amp object-store rm --force OBJECT_STORE_NAME
