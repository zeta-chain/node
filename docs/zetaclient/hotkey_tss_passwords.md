## Hot Key and TSS key-share Passwords

### Zetaclient

During startup of the client process, a series of prompts will appear asking for passwords using stdin.
* Hot Key password
* TSS Key share password

**It's extremely important to take note of these passwords or commit them to memory.**

### Hot Key

#### File backend

* The hot key will use the existing keyring that holds your operator key. The file will be encrypted with your existing password,
make sure to use this same password when starting the client.

#### Test backend

* You will still be prompted for a password, but you need to leave it blank which indicates the test backend is being used. 

### TSS Key-Share

During key-gen, the password you enter will be used to encrypt the generated key-share file. The key data will be stored in
memory once the process is running. If the client needs to be restarted, this key-share file needs to be present on your
machine and will be decrypted using the password you've entered.