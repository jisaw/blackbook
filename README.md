# blackbook
A CLI for saving ssh information

To save a new ssh server name and ip:
<code>$>blackbook -new server1</code>
Fill in the address and password. All of the server addresses are AES encrpyted. 

To delete an address:
<code>blackbook -del server1</code>

To get your name and ip:
<code>blackbook server1</code>

In The Works:
-Running the actual command to connect to ssh
-Running ssh commands through blackbook
