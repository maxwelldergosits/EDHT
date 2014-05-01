#EDHT
====

##Elastic Distributed Hash Table

our repo for the CS5412 project


###Prerequisites:


  golang distributon

  packages:

    github.com/mad293/mlog
    github.com/peterbourgon/diskv

###Building it:


  go install EDHT/coordinator
  go install EDHT/daemon


###Running it:
  coordinator [options]

  no options needed by default but import ones is the data directory that stores where you want to put the saved state, by default it wants to make a new directory with the time.  You can specify which port you want to bind to, your address that you put is where other servers in the system will try and access you, so if you want to be somewhere other than the localhost you will have to put your hostname/ip-address.

  if you are not connecting to an already established system you will probably need to specify the number of shards you want to have in the system and the number of faults you want to handle.
