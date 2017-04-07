### Create a main user
    amp user signup
    [you should receive an email with a verification token]
    amp user verify
    amp login

### Create a second user (different name, same email)
    amp user signup
    amp user verify

### Organization management
    amp org create
    amp org ls
    [the newly created org will show up]

    amp org get
    [see details about the specified org]

    amp org member add
    [add the second user tou created since the logged in user is automatically added as the organization owner]

    amp org member ls
    [see newly added member]

    amp org member rm

    amp org member ls
    [no longer find the removed user]

### Team management
    amp team create

    amp team ls
    [the newly created team will show up]

    amp team get
    [see details about the specified team]

    amp team member add
    [add the second user]

    amp team member ls
    [see newly added member]

    amp team member rm

    amp team member ls
    [no longer find the removed user]

    amp team rm

    amp team ls
    [no longer find the removed team]

    amp org rm
    [delete the organization]

    amp org ls
    [no longer find the removed org]

### Retrieve account name 
    amp user forgot-login
    [an email containing the username will be sent to the registered email address]
