## Docker Swarm cluster

AMP requires a Docker Swarm cluster to operate.
If you use an [automated way to deploy AMP in a cluster](https://github.com/appcelerator/amp-swarm-deploy), this step is already done, you can skip this page.
If you installed AMP locally, here are the instructions to create a local Docker Swarm cluster:

    $ docker swarm init
    
It should provide the following output:

    Swarm initialized: current node (ej9yivb39rrq2iyk3itdqvcq1) is now a manager.
      
    To add a worker to this swarm, run the following command:
      
        docker swarm join \
        --token SWMTKN-1-08xe2j6h2y812exq4rw5cj7j98112gn2ar88s9kkniimmn4i74-1bkwl472uc7llf4divn7k3bkv \
        10.128.27.12:2377
      
    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
 
Make sure you Docker Swarm cluster is up and ready by typing the following command:

    $ docker node ls
    
